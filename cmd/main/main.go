package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"

    "github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	//"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

var (
    theme *string
    version *bool
    theme_list *bool
    interactive_mod *bool
    display_mod *bool
    speed *int
    helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

//enum like struct to activate different mode in the model
type mode int

const (
    displayMode mode = iota
    typingMode
)

//go:embed themes/*
var themeFS embed.FS


// tickCmd sends a "tick" every 50ms
type tickMsg struct{}

func tickCmd(user_speed int) tea.Cmd {
    base_speed := 1000
    effective_speed := (base_speed)/(user_speed)
	// Ensure effectiveSpeed is at least 1ms to prevent invalid durations
	if effective_speed < 1 {
		effective_speed = 1
	}
	return tea.Tick(time.Millisecond*time.Duration(effective_speed), func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

//Cursor tick speed
type cursorTick struct{}

func cursorCmd(num time.Duration) tea.Cmd {
    return tea.Tick(time.Millisecond*num, func (t time.Time) tea.Msg {
        return cursorTick{}
    })
}


//TODO: Auto typing mode implementation using bubble tea
type model struct {
	content       string        // Full content of the file
	currentIndex  int           // Current index in the file content
	displayBuffer *strings.Builder // Displayed content buffer
	lexer         chroma.Lexer  // Lexer for syntax highlighting
	formatter     chroma.Formatter // Formatter for terminal output
	style         *chroma.Style // Style for syntax highlighting
	done          bool          // Whether typing simulation is complete
    cursorVisible bool          // Simulate cursor blinking behavior
    progress      progress.Model// Progress bar
    mode          mode         //model mode
    user_speed    int// Tick speed for updates
}

//Initialize model
func (m model) Init() tea.Cmd {
    // Use the -s flag to control speed
	// return concurrent commands in Batch(tick and cursor)
	return tea.Batch(
        tickCmd(m.user_speed),
        cursorCmd(250),
        )
}

// Update logic for Bubbletea
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Exit on Ctrl+C or Esc
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}

        if m.mode == typingMode {
            // If all content is displayed, finish
            if m.currentIndex >= len(m.content) {
                m.done = true
                return m, nil
            }
            // On every keystroke, add one character to displayBuffer
            if m.currentIndex < len(m.content) {
                m.displayBuffer.WriteByte(m.content[m.currentIndex])
                m.currentIndex++
            }
            // percentage value explicitly, too.
            cmd := m.progress.SetPercent(float64(m.currentIndex) / float64(len(m.content)))
            return m, cmd
        }

	case tickMsg:
        // Ignore tick messages in typing mode
        if m.mode == typingMode {
            return m, nil
        }

		// If all content is displayed, finish
		if m.currentIndex >= len(m.content) {
			m.done = true
			return m, nil
		}

		// Add the next character to the display buffer
		m.displayBuffer.WriteByte(m.content[m.currentIndex])
		m.currentIndex++

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
        cmd := m.progress.SetPercent(float64(m.currentIndex) / float64(len(m.content)))

		// Continue ticking if not done
		return m, tea.Batch(
                    tickCmd(m.user_speed),
                    cmd,)

    case cursorTick:
        //Toggle the cursor visible
        m.cursorVisible = !m.cursorVisible
        var numtick time.Duration
        if (m.mode == displayMode) {
            numtick = 250
        } else {numtick = 500}

        return m, cursorCmd(numtick)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	}

	return m, nil
}

// View logic for Bubbletea
func (m model) View() string {
    if m.mode == typingMode {
        // **Typing Mode Logic**:
        // Show only the content typed so far
        bufferContent := m.displayBuffer.String()


        // Tokenize the displayed buffer
        iterator, err := m.lexer.Tokenise(nil, bufferContent)
        if err != nil {
            return fmt.Sprintf("Error: %v", err)
        }

        // Format the highlighted content
        var highlighted strings.Builder
        m.formatter.Format(&highlighted, m.style, iterator)

        highlightedContent := highlighted.String()
        // Add the cursor at the end of the typed content
        cursor := "█"
        if m.cursorVisible {
            highlightedContent = highlightedContent + cursor
        }
        // Add fixed padding at the top and left
        const tabDown = 4 // Number of newlines at the top
        const tabIn = 3   // Number of tabs on the left
        const bottomPadding = 2 //Number of lines at the bottom

        // **New Logic**: Limit the number of lines visible
        const maxVisibleLines = 35 // Adjust this value for smoother scrolling
        lines := strings.Split(highlightedContent, "\n")
        start := 0
        if len(lines) > maxVisibleLines {
            start = len(lines) - maxVisibleLines
        }
        visibleLines := lines[start:] // Only show the last `maxVisibleLines` lines

        var output strings.Builder
        output.WriteString(strings.Repeat("\n", tabDown)) // Add top padding

        for _, line := range visibleLines {
            output.WriteString(strings.Repeat("\t", tabIn)) // Add left padding
            output.WriteString(line)
            output.WriteString("\n")
        }

        output.WriteString(strings.Repeat("\n", bottomPadding)) //Add bottom padding

        return  m.progress.View() + helpStyle.Render("\n • Esc: exit\n") + output.String()
    }

    // Extract the buffer's content
    bufferContent := m.displayBuffer.String()


	// Tokenize the typed buffer
	iterator, err := m.lexer.Tokenise(nil, bufferContent)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	// Format the highlighted content for syntax highlighting
	var highlighted strings.Builder
    m.formatter.Format(&highlighted, m.style, iterator)

    highlightedContent := highlighted.String()
	// Add the cursor block
	cursor := "█"
	if m.cursorVisible {
            highlightedContent = highlightedContent + cursor
	}

	// Add fixed padding at the top and left
	const tabDown = 4 // Number of newlines at the top
	const tabIn = 3   // Number of tabs on the left
    const bottomPadding = 2 //Number of lines at the bottom

    // **New Logic**: Limit the number of lines visible
    const maxVisibleLines = 30 // Adjust this value for smoother scrolling
    lines := strings.Split(highlightedContent, "\n")
    start := 0
    if len(lines) > maxVisibleLines {
        start = len(lines) - maxVisibleLines
    }
    visibleLines := lines[start:] // Only show the last `maxVisibleLines` lines

	var output strings.Builder
	output.WriteString(strings.Repeat("\n", tabDown)) // Add top padding

	for _, line := range visibleLines {
		output.WriteString(strings.Repeat("\t", tabIn)) // Add left padding
		output.WriteString(line)
		output.WriteString("\n")
	}

    output.WriteString(strings.Repeat("\n", bottomPadding)) //Add bottom padding

	return  m.progress.View() + helpStyle.Render("\n • Esc: exit\n") + output.String()
}


//Function to look into theme directory to print available theme
func list_theme() {
     //Open embedded directory
    files, err := fs.ReadDir(themeFS, "themes")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: cannot read directory %v\n", err)
        os.Exit(1)
    }

    for _, file := range files {
        if file.IsDir() {
            continue
        }
        if strings.HasSuffix(file.Name(), ".xml") {
            theme_name := strings.TrimSuffix(file.Name(), ".xml")
            fmt.Printf("- %v\n", theme_name)
        }
    }

    //READ form user_theme_dir
    homedir, err := os.UserHomeDir()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Cannot find user home directory\n")
        os.Exit(1)
    }

    user_theme_dir := homedir + "/.config/codeflow/themes"

    userfiles, err := os.ReadDir(user_theme_dir)
    if err == nil {//Only process if user_theme_dir exists
       for _, userfile := range userfiles {
            if userfile.IsDir() {
                continue
            }
            if strings.HasSuffix(userfile.Name(), ".xml") {
                theme_name := strings.TrimSuffix(userfile.Name(), ".xml")
                fmt.Printf("- %v\n", theme_name)
            }
        } 
    }
 
}

//function to print usage message
func print_usage() {
    fmt.Fprintf(os.Stderr, "\nUsage: %s [option] <file-path>\n\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "Options:\n")
    flag.PrintDefaults()
    fmt.Fprintf(os.Stderr, "\nExample:\n")
    fmt.Fprintf(os.Stderr, " codeview -t tokyonight-night file.go\n")
}

//initalize utility flag
func init() {
    //default use tokyonight-night
    theme = flag.String("t", "tokyonight-night", "Specify a theme name")
    version = flag.Bool("version", false, "Check program version")
    theme_list = flag.Bool("listtheme", false, "List all available themes")
    interactive_mod = flag.Bool("i", false, "Interactive mode: Showing word by word on keystroke")
    display_mod = flag.Bool("d", false, "Display mode: Print out content automatically")
    speed = flag.Int("s", 20, "Set the speed for auto-typing (1-1000)")


    //Override usage message
    flag.Usage = print_usage;
}


func main() {

    //Call the flag and it will handle args from here (usually os)
    flag.Parse()
    fmt.Printf("Speed value after parsing: %d\n", *speed)

    if *version {
        fmt.Printf("codeflow version 1.0.0\n")
        os.Exit(0)
    }

    if *theme_list {
        list_theme()
        os.Exit(0)
    }

    // Validate speed
    if *speed <= 0 {
        fmt.Fprintf(os.Stderr, "Error: Speed must be greater than 0\n")
        os.Exit(1)
    }

    //Check if user provide a file to read
    if (len(flag.Args()) == 0) {
        fmt.Fprintf(os.Stderr, "Error: No file provided\n")
        flag.Usage()
        os.Exit(1)
    }
    
    //Read through files list provided 
    for i := 0; i < len(flag.Args()); i++ {
        //Flag will ignore flag arguments in the Arg array
        filepath := flag.Arg(i)


        // Read the file content
        content, err := os.ReadFile(filepath)
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            return
        }

        fmt.Printf("\n------------ %v ---------------------\n\n", filepath)

        // Determine the lexer based on file extension
        // Language detection
        lexer := lexers.Match(filepath)
        if lexer == nil {
            lexer = lexers.Fallback
        }


        //NOTE: not supporting user themes for now
        // Load a built-in style
        // Use the built-in Catppuccin Mocha style
        var theme_name = *theme

        //validate theme_name
        themeExist := false
        for _, name := range styles.Names() {
            if theme_name == name {
                themeExist = true
                break
            }
        }
        if (!themeExist) {
            fmt.Println("Theme not found. Use '-themelist' to see available themes")
            os.Exit(1)
        }

        style := styles.Get(theme_name)


        // Create a terminal formatter with line numbers
        formatter := formatters.Get("terminal16m")
        if formatter == nil {
            fmt.Println("Terminal formatter not found")
            return
        }

        
        //NOTE: Rendering bubbletea model
        if *display_mod{
            // Initialize the Bubbletea model
            m := model{
                content:       string(content),
                lexer:         lexer,
                formatter:     formatter,
                style:         style,
                currentIndex:  0,
                displayBuffer: &strings.Builder{},
                done:          false,
                cursorVisible: true,
                progress:      progress.New(progress.WithDefaultGradient()),
                mode:          displayMode, //display mode
                user_speed:    *speed, //Set custom speed for display mode
            }

            // Initialize program and model
            // Model run in Alt Screen
            p := tea.NewProgram(m, tea.WithAltScreen())

            //Run program
            _, err := p.Run()
            if err != nil {
                fmt.Printf("Error running program: %v\n", err)
                os.Exit(1)
            }

        } else if *interactive_mod {
            //TODO: typing mode here
            // Initialize the Bubbletea model
            m := model{
                content:       string(content),
                lexer:         lexer,
                formatter:     formatter,
                style:         style,
                currentIndex:  0,
                displayBuffer: &strings.Builder{},
                done:          false,
                cursorVisible: true,
                progress:      progress.New(progress.WithDefaultGradient()),
                mode:          typingMode, //typing mode
                user_speed:    *speed, //Set custom speed for display mode
            }

            // Initialize program and model
            // Model run in Alt Screen
            p := tea.NewProgram(m, tea.WithAltScreen())

            //Run program
            _, err := p.Run()
            if err != nil {
                fmt.Printf("Error running program: %v\n", err)
                os.Exit(1)
            }
        } else {
            // Render the highlighted output line by line
            lines := strings.Split(string(content), "\n")
            for i, line := range lines {
                // Print line number in grey
                fmt.Printf("\033[90m%5d\033[0m ", i+1) // Grey ANSI escape sequence

                // Tokenize and format the current line
                lineIterator, err := lexer.Tokenise(nil, line)
                if err != nil {
                    fmt.Printf("Error tokenizing line: %v\n", err)
                    continue
                }
                formatter.Format(os.Stdout, style, lineIterator)

                fmt.Println()
            }
        }
    }
}
