# CodeFlow CLI

CodeFlow CLI is a versatile tool for displaying and interacting with text-based content dynamically. Designed with customization and ease of use in mind, it offers multiple themes and display modes to suit various user needs.


## Demo
![](assets/demo.gif)

## Features

- **Custom Themes**: Select your preferred theme with the `-t` flag.
- **Theme Listing**: View all available themes using the `-listtheme` flag.
- **Interactive Mode**: Display text word by word, advancing on each keystroke, with the `-i` flag.
- **Display Mode**: Automatically display the text content without user interaction using the `-d` flag.
- Custom Speed: Adjust the speed of automatic content display with the `-s` flag (higher values mean faster display).
- **Version Information**: Check the current version of CodeFlow with the `-version` flag.

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/han-nwin/codeflow
    cd codeflow
    ```
2. Make installation script executable
    ```bash
    chmod +x install.sh
    ```
3. Run installation script
    ```bash
    ./install.sh
    ```
4. Verify installation
    ```bash
    codeflow -version
    ```

## Options
- -t <theme>: Specify a syntax highlighting theme (default: tokyonight-night).
- -listtheme: List all available themes.
- -d: Enable display mode (automatic content display with progress bar).
- -s <speed>: Set the speed of automatic display in display mode (default: 20, higher values are faster).
- -i: Enable interactive mode (reveal content word by word on each keystroke).
- -version: Display the program version.

## Examples usage

#### Normal Usage
By default, running `codeflow` with a file path and no flags will display the content using the default theme (`tokyonight-night`) in a non-interactive manner. The text is displayed all at once without animations or progress tracking.
```bash
codeflow myfile.txt
```
#### Specify A Theme
```bash
codeflow -t dracula myfile.txt
```
List Available Themes
```bash
codeflow -listtheme
```


#### Display Mode with Custom Speed
```bash
codeflow -d -s 50 myfile.txt
```
> You can specify a theme by using -t as well.
#### Interactive Mode
```bash
codeflow -i myfile.txt
```
> You can specify a theme by using -t as well.

