#!/bin/bash

# Install script for Codeflow CLI

BINARY_NAME="codeflow"
DEFAULT_INSTALL_DIR="/usr/local/bin"
USER_INSTALL_DIR="$HOME/.local/bin"

echo "🚀 Welcome to the Codeflow CLI installer! 🚀"

# -------------------------------------------------
# Prompt for install directory with a default option
echo ""
echo "📁 The default installation directory is: $DEFAULT_INSTALL_DIR"
read -p "❓ Do you want to install in the default directory? [Y/n]: " RESPONSE

# Normalize the response to lowercase
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
  echo "❌ Build failed. Please ensure you have Go installed and configured correctly."
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
  echo "🔑 Permission required to install in $INSTALL_DIR. Prompting for sudo..."
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
    echo "set -U fish_user_paths $INSTALL_DIR \$fish_user_paths" >> "$SHELL_CONFIG"
  else
    echo "export PATH=\$PATH:$INSTALL_DIR" >> "$SHELL_CONFIG"
  fi
  echo ""
  echo "🔄 Reload your shell or run 'source $SHELL_CONFIG' to apply the changes."
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
  echo "🎉 Installation successful! You can now use '$BINARY_NAME' from anywhere."
  echo "ℹ️  Use 'codeflow -version' to verify installation."
else
  echo ""
  echo "⚠️ Installation completed, but $BINARY_NAME is not in your PATH. Please add $INSTALL_DIR to your PATH manually."
fi
