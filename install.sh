#!/bin/bash

set -e

# Constants
DEFAULT_INSTALL_DIR="$HOME/.privateer/bin"
PVTR_REPO="privateerproj/privateer"
LATEST_RELEASE_URL="https://api.github.com/repos/${PVTR_REPO}/releases/latest"

# Detect OS (darwin = macOS, linux = Linux, msys or cygwin = Windows)
OS=""
case "$(uname -s)" in
    Darwin)
        OS="darwin"
        ;;
    Linux)
        OS="linux"
        ;;
    CYGWIN*|MSYS*|MINGW*)
        OS="windows"
        ;;
    *)
        echo "Unsupported Environment: $(uname -s)"
        exit 1
        ;;
esac

# Detect Architecture (x86_64 = amd64, arm64 = arm64, i386/i686 = 386)
ARCH=""
case "$(uname -m)" in
    x86_64)
        ARCH="x86_64"
        ;;
    i386)
        ARCH="i386"
        ;;
    arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported Architecture: $(uname -m)"
        exit 1
        ;;
esac

extract_download_urls() {
    local release_json="$1"

    printf '%s' "$release_json" \
        | tr '\n' ' ' \
        | grep -Eo '"browser_download_url"[[:space:]]*:[[:space:]]*"[^"]+"' \
        | sed -E 's/^"browser_download_url"[[:space:]]*:[[:space:]]*"([^"]+)"$/\1/'
}

find_release_asset_url() {
    local release_json="$1"
    local asset_pattern="$2"

    extract_download_urls "$release_json" | grep -Ei "$asset_pattern" | head -n 1
}

download_latest_release() {
    local install_dir="$1"
    local install_file="$install_dir/pvtr"

    # Ensure the directory exists
    mkdir -p "$install_dir"

    # Fetch release metadata once
    local release_json
    release_json=$(curl -s "${LATEST_RELEASE_URL}")

    # Build the grep pattern based on OS
    local pattern
    if [[ "$OS" == "darwin" ]]; then
        pattern="${OS}"
    else
        pattern="${OS}.*${ARCH}"
    fi

    # Fetch the download URL for the latest release binary
    local url
    url=$(find_release_asset_url "$release_json" "$pattern")

    if [[ -z "$url" ]]; then
        echo "Failed to fetch the download URL for the latest release."
        exit 1
    fi

    # Fetch the checksums file URL
    local checksums_url
    checksums_url=$(find_release_asset_url "$release_json" 'checksums\.txt$')

    if [[ -z "$checksums_url" ]]; then
        echo "ERROR: No checksums file found in release assets. Refusing to install an unverified binary."
        exit 1
    fi

    # Create a temporary directory for download and verification
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap "rm -rf '$tmp_dir'" EXIT

    local archive_name
    archive_name=$(basename "$url")
    local tmp_archive="$tmp_dir/$archive_name"

    echo "Downloading from: $url"

    # Download the archive to a temporary file
    curl -fSL -o "$tmp_archive" "$url"

    echo "Verifying checksum..."
    local tmp_checksums="$tmp_dir/checksums.txt"
    if ! curl -fSL -o "$tmp_checksums" "$checksums_url"; then
        echo "ERROR: Failed to download checksums file. Refusing to install an unverified binary."
        exit 1
    fi

    # Extract the expected checksum for our archive
    local expected_checksum
    expected_checksum=$(grep "  ${archive_name}$" "$tmp_checksums" | awk '{print $1}')

    if [[ -z "$expected_checksum" ]]; then
        # Try alternate format (single space separator)
        expected_checksum=$(grep " ${archive_name}$" "$tmp_checksums" | awk '{print $1}')
    fi

    if [[ -z "$expected_checksum" ]]; then
        echo "ERROR: Could not find checksum for ${archive_name} in checksums file."
        echo "Aborting installation for security. To skip verification, download manually."
        exit 1
    fi

    # Compute actual checksum
    local actual_checksum
    if command -v sha256sum &>/dev/null; then
        actual_checksum=$(sha256sum "$tmp_archive" | awk '{print $1}')
    elif command -v shasum &>/dev/null; then
        actual_checksum=$(shasum -a 256 "$tmp_archive" | awk '{print $1}')
    else
        echo "ERROR: No sha256sum or shasum found. Cannot verify checksum."
        exit 1
    fi

    if [[ "$actual_checksum" != "$expected_checksum" ]]; then
        echo "CHECKSUM MISMATCH!"
        echo "  Expected: $expected_checksum"
        echo "  Actual:   $actual_checksum"
        echo "The downloaded file may have been tampered with. Aborting."
        exit 1
    fi

    echo "Checksum verified OK."

    # Extract the verified archive
    tar xf "$tmp_archive" -C "$install_dir"

    # Ensure the binary is executable
    chmod +x "$install_file"

    echo "Downloaded binary to $install_file"
}

update_path() {
    local install_dir="$1"

    # Check if the install directory is already in PATH
    if [[ ":$PATH:" != *":$install_dir:"* ]]; then
        echo "$install_dir is not in the PATH."

        # Detect current shell
        current_shell=$(basename "$SHELL")

        case "$current_shell" in
            bash)
                config_file="$HOME/.bash_profile"
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
    download_latest_release "$install_dir"

    # Ensure the binary is accessible via PATH
    update_path "$install_dir"

    echo "pvtr installation complete!"
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
