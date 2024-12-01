#!/bin/bash

# Install script for Codeflow CLI

BINARY_NAME="codeflow"
DEFAULT_INSTALL_DIR="/usr/local/bin"
USER_INSTALL_DIR="$HOME/.local/bin"
MIN_GO_VERSION="1.23.3"
echo "🚀 Welcome to the Codeflow CLI installer! 🚀"

# -------------------------------------------------
# Function to compare version numbers
version_ge() {
    [[ "$1" == $(echo -e "$1\n$2" | sort -V | head -n 1) ]]
}

# -------------------------------------------------
# Check if "go" is installed and determine version
if command -v go >/dev/null 2>&1; then
    INSTALLED_GO_VERSION=$(go version | awk '{print $3}' | cut -d' ' -f3 | sed 's/go//')
    echo "🟢 Go is installed. Version: $INSTALLED_GO_VERSION"
else
    INSTALLED_GO_VERSION=""
    echo "🔴 Go is not installed."
fi

# -------------------------------------------------
# Compare versions and install Go if necessary
if [ -z "$INSTALLED_GO_VERSION" ] || ! version_ge "$INSTALLED_GO_VERSION"
"$MIN_GO_VERSION"; then
    echo "🔄 Installing/Updating Go to version $MIN_GO_VERSION..."
    OS=$(uname | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
    elif [ "$ARCH" == "aarch64" ]; then
    ARCH="arm64"
    fi

    GO_TAR="go$MIN_GO_VERSION.$OS-$ARCH.tar.gz"
    curl -OL "https://golang.org/dl/$GO_TAR"

    if [ -d "/usr/local/go" ]; then
    sudo rm -rf /usr/local/go
    fi

    sudo tar -C /usr/local -xzf "$GO_TAR"
    rm "$GO_TAR"

    export PATH=$PATH:/usr/local/go/bin
    echo "🟢 Go has been installed. Please restart your terminal to ensure Go is
in your PATH."
else
    echo "🟢 Your Go version is up to date."
fi

echo ""
echo "📁 The default installation directory is: $DEFAULT_INSTALL_DIR"
read -p "❓ Do you want to install in the default directory? [Y/n]: " RESPONSE

RESPONSE=${RESPONSE,,}  # Convert to lowercase for consistency

if [[ "$RESPONSE" == "n" || "$RESPONSE" == "no" ]]; then
    echo ""
    read -p "✏️  Enter custom installation directory: " CUSTOM_DIR
    INSTALL_DIR=${CUSTOM_DIR:-$DEFAULT_INSTALL_DIR}
else
    INSTALL_DIR=$DEFAULT_INSTALL_DIR
fi

echo ""
echo "✅ Installation directory set to: $INSTALL_DIR"

# -------------------------------------------------
# Step 1: Build the binary
echo ""
echo "🔧 Building the binary..."
go build -o $BINARY_NAME cmd/main/main.go
if [ $? -ne 0 ]; then
    echo ""
    echo "❌ Build failed. Please ensure you have Go installed and configured
correctly."
    exit 1
fi
echo ""
echo "✔️ Build successful."

# -------------------------------------------------
# Step 2: Install the binary
echo ""
echo "📥 Installing the binary to $INSTALL_DIR..."
mkdir -p "$INSTALL_DIR"
if [ -w "$INSTALL_DIR" ]; then
    mv $BINARY_NAME "$INSTALL_DIR/"
else
    echo ""
    echo "🔑 Permission required to install in $INSTALL_DIR. Prompting for sudo...
"
    sudo mv $BINARY_NAME "$INSTALL_DIR/"
fi

# -------------------------------------------------
# Step 3: Add to PATH if necessary
echo ""
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "🔗 Adding $INSTALL_DIR to your PATH..."
    SHELL_CONFIG="$HOME/.bashrc" # Default shell config file
    if [[ "$SHELL" == *"zsh"* ]]; then
    SHELL_CONFIG="$HOME/.zshrc"
    elif [[ "$SHELL" == *"fish"* ]]; then
    SHELL_CONFIG="$HOME/.config/fish/config.fish"
    fi

    if [[ "$SHELL" == *"fish"* ]]; then
    echo "set -U fish_user_paths $INSTALL_DIR \$fish_user_paths" >>
"$SHELL_CONFIG"
    else
    echo "export PATH=\$PATH:$INSTALL_DIR" >> "$SHELL_CONFIG"
    fi
    echo ""
    echo "🔄 Reload your shell or run 'source $SHELL_CONFIG' to apply the
changes."
else
    echo ""
    echo "🔍 $INSTALL_DIR is already in your PATH."
fi

# -------------------------------------------------
# Step 4: Verify installation
echo ""
echo "✅ Verifying installation..."
if command -v $BINARY_NAME >/dev/null 2>&1; then
    echo ""
    echo "🎉 Installation successful! You can now use '$BINARY_NAME' from
anywhere."
    echo "ℹ️  Use 'codeflow -version' to verify installation."
else
    echo ""
    echo "⚠️ Installation completed, but $BINARY_NAME is not in your PATH.
Please add $INSTALL_DIR to your PATH manually."
fi
