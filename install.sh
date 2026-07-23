#!/usr/bin/env bash

set -euo pipefail

REPO="privateerproj/privateer"
DEFAULT_INSTALL_DIR="$HOME/.local/bin"

# Where the old installer put pvtr. That directory is really the plugin
# directory (see cmd/cli.go), so new installs go to DEFAULT_INSTALL_DIR. An
# existing install is upgraded where it stands rather than left behind on PATH
# to shadow the new one.
LEGACY_INSTALL_DIR="$HOME/.privateer/bin"

# GoReleaser names archives <project>_<Os>_<Arch>.tar.gz. The project prefix was
# renamed privateer -> pvtr, so releases published before the rename are still
# reachable, which matters because -v can pin any historical version.
ASSET_PREFIXES=(pvtr privateer)

# Global so the EXIT trap can still see it once main has returned.
TMP_DIR=""

cleanup() {
    if [[ -n "$TMP_DIR" ]]; then
        rm -rf "$TMP_DIR"
    fi
}
trap cleanup EXIT

usage() {
    cat <<EOF
Install the pvtr CLI.

Usage: install.sh [-p <dir>] [-v <version>]

  -p <dir>      Install directory (default: $DEFAULT_INSTALL_DIR)
  -v <version>  Release tag to install, e.g. v0.22.0 (default: latest)
  -h            Show this help

Environment:
  PVTR_VERSION  Same as -v; the flag wins if both are set.

When piping, pass arguments after a placeholder:
  /bin/bash -c "\$(curl -sSL <url>/install.sh)" -- -v v0.22.0
EOF
}

# Print the archive name for this machine, minus the project prefix.
detect_asset_suffix() {
    local arch
    case "$(uname -s)" in
        Darwin)
            # Darwin ships a universal binary, so architecture is irrelevant.
            printf 'Darwin_all.tar.gz'
            return 0
            ;;
        Linux) ;;
        CYGWIN*|MSYS*|MINGW*)
            echo "This script installs the Unix build of pvtr and cannot install on Windows." >&2
            echo "Download pvtr_Windows_x86_64.zip from https://github.com/${REPO}/releases," >&2
            echo "extract it, and place pvtr.exe on your PATH." >&2
            return 1
            ;;
        *)
            echo "Unsupported operating system: $(uname -s)" >&2
            return 1
            ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64) arch="x86_64" ;;
        aarch64|arm64) arch="arm64" ;;
        i386|i686) arch="i386" ;;
        *)
            echo "Unsupported architecture: $(uname -m)" >&2
            echo "pvtr publishes Linux builds for x86_64, arm64, and i386." >&2
            return 1
            ;;
    esac

    printf 'Linux_%s.tar.gz' "$arch"
}

# Release asset URLs are deterministic, so no API call (and no rate limit) is
# needed to find them. An empty version resolves to the latest release.
asset_base_url() {
    local version="$1"

    if [[ -n "$version" ]]; then
        printf 'https://github.com/%s/releases/download/%s' "$REPO" "$version"
    else
        printf 'https://github.com/%s/releases/latest/download' "$REPO"
    fi
}

# Download the release archive, trying each project prefix in turn. Echoes the
# archive filename that resolved so the caller can match it in checksums.txt.
download_archive() {
    local base_url="$1" suffix="$2" dest_dir="$3"
    local prefix archive_name status=0

    # Probes are silent: a 404 on the first prefix is the expected path for
    # releases published before the rename, not something to report.
    for prefix in "${ASSET_PREFIXES[@]}"; do
        archive_name="${prefix}_${suffix}"
        if curl -fsL -o "${dest_dir}/${archive_name}" "${base_url}/${archive_name}"; then
            printf '%s' "$archive_name"
            return 0
        fi
        status=$?
    done

    echo "Could not download a release archive for this platform from ${base_url} (curl exit ${status})." >&2
    echo "Tried: ${ASSET_PREFIXES[*]/%/_${suffix}}" >&2
    return 1
}

verify_checksum() {
    local base_url="$1" dest_dir="$2" archive_name="$3"
    local checksums="${dest_dir}/checksums.txt"

    if ! curl -fsSL -o "$checksums" "${base_url}/checksums.txt"; then
        echo "ERROR: Failed to download checksums file. Refusing to install an unverified binary." >&2
        return 1
    fi

    local expected
    expected=$(awk -v name="$archive_name" '$2 == name { print $1 }' "$checksums")
    if [[ -z "$expected" ]]; then
        echo "ERROR: Could not find a checksum for ${archive_name} in checksums.txt." >&2
        return 1
    fi

    local actual
    if command -v sha256sum >/dev/null 2>&1; then
        actual=$(sha256sum "${dest_dir}/${archive_name}" | awk '{print $1}')
    elif command -v shasum >/dev/null 2>&1; then
        actual=$(shasum -a 256 "${dest_dir}/${archive_name}" | awk '{print $1}')
    else
        echo "ERROR: Neither sha256sum nor shasum is available. Cannot verify the download." >&2
        return 1
    fi

    if [[ "$actual" != "$expected" ]]; then
        echo "CHECKSUM MISMATCH for ${archive_name}" >&2
        echo "  Expected: ${expected}" >&2
        echo "  Actual:   ${actual}" >&2
        echo "The downloaded file may have been tampered with. Aborting." >&2
        return 1
    fi
}

installed_version() {
    local binary="$1"

    [[ -x "$binary" ]] || return 1
    "$binary" version 2>/dev/null | head -n 1 | tr -d '[:space:]'
}

# Report whether the install directory is reachable, and if not, print the line
# the user needs. The script never edits shell configuration: it cannot affect
# the calling shell anyway, and writing bash syntax into another shell's config
# breaks that shell.
report_path() {
    local install_dir="$1"

    if [[ ":${PATH}:" == *":${install_dir}:"* ]]; then
        return 0
    fi

    echo
    echo "${install_dir} is not on your PATH. Add it with:"
    case "$(basename "${SHELL:-}")" in
        fish)
            echo "  fish_add_path ${install_dir}"
            ;;
        zsh)
            echo "  echo 'export PATH=\"${install_dir}:\$PATH\"' >> ~/.zshrc && exec zsh"
            ;;
        *)
            echo "  echo 'export PATH=\"${install_dir}:\$PATH\"' >> ~/.bash_profile && exec bash -l"
            ;;
    esac
}

main() {
    local install_dir="" version="${PVTR_VERSION:-}" opt

    while getopts "p:v:h" opt; do
        case "$opt" in
            p) install_dir="$OPTARG" ;;
            v) version="$OPTARG" ;;
            h) usage; return 0 ;;
            *) usage >&2; return 1 ;;
        esac
    done

    if [[ -z "$install_dir" ]]; then
        if [[ -f "${LEGACY_INSTALL_DIR}/pvtr" ]]; then
            install_dir="$LEGACY_INSTALL_DIR"
            echo "Upgrading the existing install in ${install_dir}."
            echo "New installs use ${DEFAULT_INSTALL_DIR}; pass -p to move there."
        else
            install_dir="$DEFAULT_INSTALL_DIR"
        fi
    fi

    local suffix base_url
    suffix=$(detect_asset_suffix)
    base_url=$(asset_base_url "$version")

    TMP_DIR=$(mktemp -d)

    local previous
    previous=$(installed_version "${install_dir}/pvtr" || true)

    echo "Downloading pvtr from ${base_url}"
    local archive_name
    archive_name=$(download_archive "$base_url" "$suffix" "$TMP_DIR")

    echo "Verifying checksum..."
    verify_checksum "$base_url" "$TMP_DIR" "$archive_name"
    echo "Checksum verified."

    mkdir -p "${TMP_DIR}/extract" "$install_dir"
    tar xzf "${TMP_DIR}/${archive_name}" -C "${TMP_DIR}/extract"

    # Install only the binary. The archive also carries LICENSE and README,
    # which do not belong in a bin directory.
    install -m 0755 "${TMP_DIR}/extract/pvtr" "${install_dir}/pvtr"

    local current
    current=$(installed_version "${install_dir}/pvtr" || true)
    if [[ -n "$previous" && "$previous" != "$current" ]]; then
        echo "Upgraded pvtr ${previous} -> ${current} in ${install_dir}"
    else
        echo "Installed pvtr ${current:-(version unknown)} in ${install_dir}"
    fi

    report_path "$install_dir"
}

# Run unconditionally so the script works whether it is executed directly, piped
# to bash (curl ... | bash), or run via bash -c "$(curl ...)".
main "$@"
