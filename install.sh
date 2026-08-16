#!/usr/bin/env bash

set -euo pipefail

REPO="privateerproj/privateer"
DEFAULT_INSTALL_DIR="$HOME/.privateer/bin"
DEFAULT_VERSION="latest"

# Set by detect_platform. A global rather than a stdout value because `fail`
# calls `exit`, which only unwinds the subshell inside a `$(...)` and would
# leave the caller running with an empty platform.
PLATFORM=""

usage() {
    cat <<EOF
Usage: install.sh [-p install_dir] [-v version] [-y] [-h]

  -p  Install directory (default: $DEFAULT_INSTALL_DIR)
  -v  Release to install, e.g. v0.22.0 (default: $DEFAULT_VERSION)
  -y  Add the install directory to PATH in your shell config without prompting
  -h  Show this help

Environment:
  PVTR_VERSION  Same as -v. The flag wins if both are set.
EOF
}

fail() {
    echo "ERROR: $*" >&2
    exit 1
}

# Overridden by the test suite to exercise the prompt without a pty.
is_interactive() {
    [[ -t 0 && -t 1 ]]
}

# Absolute path with redundant '.' and '/' segments collapsed, so the PATH line
# printed to the user is clean. `realpath` is absent from a stock macOS, so this
# stays in shell. '..' is deliberately left in place: resolving it lexically
# would give the wrong answer under a symlinked parent, and a PATH entry
# containing '..' still works.
abspath() {
    local path="$1" part result=""
    local -a parts=()

    [[ "$path" == /* ]] || path="$PWD/$path"

    local IFS='/'
    read -r -a parts <<< "$path"
    for part in "${parts[@]}"; do
        [[ -z "$part" || "$part" == "." ]] && continue
        result="$result/$part"
    done

    printf '%s\n' "${result:-/}"
}

detect_platform() {
    local os arch
    os="$(uname -s)"
    arch="$(uname -m)"

    case "$os" in
        Darwin)
            # Darwin archives are universal binaries, so the arch is irrelevant.
            PLATFORM="Darwin_all"
            return
            ;;
        Linux) ;;
        *)
            fail "Unsupported OS: $os. Download a release directly from https://github.com/${REPO}/releases"
            ;;
    esac

    case "$arch" in
        x86_64 | amd64) arch="x86_64" ;;
        aarch64 | arm64) arch="arm64" ;;
        i386 | i686) arch="i386" ;;
        *)
            fail "Unsupported architecture: $arch. Download a release directly from https://github.com/${REPO}/releases"
            ;;
    esac

    PLATFORM="Linux_${arch}"
}

release_base_url() {
    local version="$1"

    if [[ "$version" == "latest" ]]; then
        printf '%s\n' "https://github.com/${REPO}/releases/latest/download"
        return
    fi

    [[ "$version" == v* ]] || version="v$version"
    printf '%s\n' "https://github.com/${REPO}/releases/download/${version}"
}

# Checked by ensure_prerequisites, so this never has to report a missing tool.
# Reporting one here would be swallowed anyway: the caller runs it in a command
# substitution, where `fail` would exit only the subshell and leave an empty
# digest that reads as a checksum mismatch.
sha256() {
    if command -v sha256sum > /dev/null 2>&1; then
        sha256sum "$1" | awk '{print $1}'
    else
        shasum -a 256 "$1" | awk '{print $1}'
    fi
}

# Fail before downloading rather than part-way through.
ensure_prerequisites() {
    local tool
    for tool in curl tar awk; do
        command -v "$tool" > /dev/null 2>&1 || fail "$tool is required but was not found."
    done

    command -v sha256sum > /dev/null 2>&1 || command -v shasum > /dev/null 2>&1 \
        || fail "Neither sha256sum nor shasum is available; cannot verify the download."
}

# Fail before downloading rather than after, so `-p /usr/local/bin` (which needs
# root on most systems) reports something actionable.
ensure_writable_dir() {
    local dir="$1" ancestor="$1"

    if [[ -d "$dir" ]]; then
        [[ -w "$dir" ]] || fail "$dir is not writable. Re-run with sudo, or choose another directory with -p."
        return
    fi

    while [[ ! -d "$ancestor" ]]; do
        ancestor="$(dirname "$ancestor")"
    done

    [[ -w "$ancestor" ]] || fail "Cannot create $dir: $ancestor is not writable. Re-run with sudo, or choose another directory with -p."
}

cleanup() {
    [[ -n "${tmp_dir:-}" ]] && rm -rf "$tmp_dir"
    [[ -n "${staged_binary:-}" ]] && rm -f "$staged_binary"
    return 0
}

# Downloads the first candidate asset name that exists and echoes that name.
# Asset names have changed across releases, so a pinned older version needs the
# name that release actually used. The probe is deliberately quiet (-s, no -S):
# a miss is expected and the caller reports the aggregate failure.
fetch_first_asset() {
    local base_url="$1" dest="$2"
    shift 2

    local name
    for name in "$@"; do
        if curl -fsL --retry 3 --retry-delay 1 -o "$dest" "$base_url/$name"; then
            printf '%s\n' "$name"
            return 0
        fi
    done
    return 1
}

download_release() {
    local install_dir="$1" platform="$2" base_url="$3" version="$4"
    local archive checksums expected actual
    local -a checksum_names=("checksums.txt")

    # Not local: the EXIT trap runs after this function returns.
    tmp_dir="$(mktemp -d)"
    staged_binary=""
    trap cleanup EXIT

    # v0.22.0 and earlier published archives as privateer_* rather than pvtr_*.
    archive="$(fetch_first_asset "$base_url" "$tmp_dir/archive.tar.gz" \
        "pvtr_${platform}.tar.gz" "privateer_${platform}.tar.gz")" \
        || fail "Failed to download a release archive for ${platform} from $base_url"
    echo "Downloaded $archive"

    # v0.21.0 and earlier named the checksums file privateer_<version>_checksums.txt.
    # Only worth probing for a pinned version; "latest" is always v0.22.0 or newer.
    if [[ "$version" != "latest" ]]; then
        checksum_names+=("privateer_${version#v}_checksums.txt")
    fi

    checksums="$(fetch_first_asset "$base_url" "$tmp_dir/checksums.txt" "${checksum_names[@]}")" \
        || fail "Failed to download a checksums file from $base_url; refusing to install an unverified binary."
    echo "Verifying against $checksums"

    expected="$(awk -v name="$archive" '$2 == name { print $1 }' "$tmp_dir/checksums.txt")"
    [[ -n "$expected" ]] || fail "No entry for $archive in $checksums; refusing to install an unverified binary."

    actual="$(sha256 "$tmp_dir/archive.tar.gz")"
    [[ "$actual" == "$expected" ]] || fail "Checksum mismatch for $archive.
  Expected: $expected
  Actual:   $actual
The download is corrupt or has been altered in transit."
    echo "Checksum verified."

    tar -xzf "$tmp_dir/archive.tar.gz" -C "$tmp_dir" pvtr \
        || fail "Archive $archive does not contain a pvtr binary."

    # Stage inside the install dir, then rename. A rename within one filesystem
    # is atomic, so an interrupted install cannot leave a truncated pvtr behind,
    # and an in-use or code-signed binary is replaced rather than overwritten.
    mkdir -p "$install_dir" || fail "Failed to create $install_dir"
    staged_binary="$(mktemp "$install_dir/.pvtr.download.XXXXXX")"
    cp "$tmp_dir/pvtr" "$staged_binary"
    chmod 755 "$staged_binary"
    mv -f "$staged_binary" "$install_dir/pvtr"
    staged_binary=""

    echo "Installed $install_dir/pvtr"
}

update_path() {
    local install_dir="$1" assume_yes="$2"
    local config_file path_line consent reply

    case ":$PATH:" in
        *":$install_dir:"*) return ;;
    esac

    case "$(basename "${SHELL:-}")" in
        bash)
            # macOS terminals start login shells, which read .bash_profile and
            # never source .bashrc.
            if [[ "$(uname -s)" == "Darwin" ]]; then
                config_file="$HOME/.bash_profile"
            else
                config_file="$HOME/.bashrc"
            fi
            path_line="export PATH=\"$install_dir:\$PATH\""
            ;;
        zsh)
            config_file="$HOME/.zshrc"
            path_line="export PATH=\"$install_dir:\$PATH\""
            ;;
        fish)
            config_file="$HOME/.config/fish/config.fish"
            path_line="fish_add_path \"$install_dir\""
            ;;
        *)
            echo "To use pvtr, add $install_dir to your PATH."
            return
            ;;
    esac

    # Ignore commented-out lines so a leftover comment does not read as configured.
    if grep -vE '^[[:space:]]*#' "$config_file" 2> /dev/null | grep -qF -- "$install_dir"; then
        echo "$config_file already references $install_dir. Restart your shell to use pvtr."
        return
    fi

    consent="$assume_yes"
    if [[ "$consent" != true ]] && is_interactive; then
        reply=""
        # `read` returns non-zero at EOF (Ctrl-D). Treat that as declining rather
        # than letting `set -e` abort an otherwise complete install.
        read -r -p "Add $install_dir to PATH in $config_file? [y/N] " reply || reply=""
        if [[ "$reply" == [yY]* ]]; then
            consent=true
        fi
    fi

    if [[ "$consent" == true ]]; then
        mkdir -p "$(dirname "$config_file")"
        printf '\n%s\n' "$path_line" >> "$config_file"
        echo "Added $install_dir to PATH in $config_file. Restart your shell to use pvtr."
    else
        echo "To use pvtr, add $install_dir to your PATH:"
        echo "  $path_line"
    fi
}

main() {
    local install_dir="$DEFAULT_INSTALL_DIR" assume_yes=false
    local version="${PVTR_VERSION:-$DEFAULT_VERSION}"

    while getopts "p:v:yh" opt; do
        case "$opt" in
            p) install_dir="$OPTARG" ;;
            v) version="$OPTARG" ;;
            y) assume_yes=true ;;
            h)
                usage
                exit 0
                ;;
            *)
                usage >&2
                exit 1
                ;;
        esac
    done
    shift $((OPTIND - 1))

    if [[ $# -gt 0 ]]; then
        echo "ERROR: unexpected argument: $1" >&2
        usage >&2
        exit 1
    fi

    install_dir="$(abspath "$install_dir")"
    ensure_prerequisites
    detect_platform
    ensure_writable_dir "$install_dir"
    download_release "$install_dir" "$PLATFORM" "$(release_base_url "$version")" "$version"
    update_path "$install_dir" "$assume_yes"
    echo "pvtr installation complete."
}

# BASH_SOURCE is unset under `curl | bash` and `bash -c`, where $0 is the
# fallback; when sourced (see test/install_test.sh) the two differ.
if [[ "${BASH_SOURCE[0]:-$0}" == "$0" ]]; then
    main "$@"
fi
