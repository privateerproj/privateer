#!/bin/bash

set -e

# Constants
DEFAULT_INSTALL_DIR="$HOME/.privateer/bin"
PRIVATEER_REPO="privateerproj/privateer"
LATEST_RELEASE_URL="https://api.github.com/repos/${PRIVATEER_REPO}/releases/latest"

# Detect OS (darwin = macOS, linux = Linux, msys or cygwin = Windows)
OS=""
case "$(uname -s)" in
    Darwin)
        OS="darwin"
        ;;
    Linux)
        OS="linux"
        ;;
    CYGWIN*|MINGW32*|MSYS*|MINGW*)
        OS="windows"
        ;;
    *)
        echo "Unsupported Environment: $(uname -s)"
        exit 1
        ;;
esac

# Function to download the latest release
download_latest_release() {
    local url
    url=$(curl -s ${LATEST_RELEASE_URL} | grep "browser_download_url.*${OS}.*" | cut -d '"' -f 4)

    if [[ -z "$url" ]]; then
        echo "Failed to fetch the download URL for the latest release."
        exit 1
    fi

    echo "Downloading from: $url"
    curl -LO "$url"
}

# Function to install the binary
install_binary() {
    local install_dir="$1"
    local file_name=$(basename $(curl -s ${LATEST_RELEASE_URL} | grep "browser_download_url.*${OS}.*" | cut -d '"' -f 4))
    local install_file="$install_dir/privateer"

    echo "Installing Privateer to $install_file"

    chmod +x "$file_name"
    sudo mv "$file_name" "$install_file"
}

# Function to update the PATH if needed
update_path() {
    local install_dir="$1"

    # Check if the install directory is already in PATH
    if [[ ":$PATH:" != *":$install_dir:"* ]]; then
        echo "$install_dir is not in the PATH."

        # macOS/Linux path update
        if [[ "$OS" == "linux" || "$OS" == "darwin" ]]; then
            # Detect current shell
            current_shell=$(basename "$SHELL")

            case "$current_shell" in
                bash)
                    config_file="$HOME/.bashrc"
                    ;;
                zsh)
                    config_file="$HOME/.zshrc"
                    ;;
                fish)
                    config_file="$HOME/.config/fish/config.fish"
                    ;;
                *)
                    echo "Unsupported shell: $current_shell. You may need to manually add $install_dir to your PATH."
                    return
                    ;;
            esac

            # Check if the path is already added to the config file
            if ! grep -q "$install_dir" "$config_file"; then
                echo "export PATH=\"$install_dir:\$PATH\"" >> "$config_file"
                echo "$install_dir added to $config_file"
                source $config_file
            else
                echo "$install_dir is already in $config_file."
            fi
        fi
    else
        echo "$install_dir is already in the PATH."
    fi
}

# Main logic
main() {
    local install_dir="$DEFAULT_INSTALL_DIR"

    # Handle CLI arguments for installation path override
    while getopts "p:" opt; do
        case $opt in
            p)
                install_dir="$OPTARG"
                ;;
            *)
                echo "Usage: $0 [-p install_path]"
                exit 1
                ;;
        esac
    done

    mkdir -p "$install_dir"

    # Download the latest release
    download_latest_release

    # Install the binary
    install_binary "$install_dir"

    # Ensure the binary is accessible via PATH
    update_path "$install_dir"

    echo "Privateer installation complete!"
}

main "$@"
