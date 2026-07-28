#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
INSTALL_SCRIPT="$REPO_ROOT/install.sh"

make_mock_commands() {
    local mock_bin="$1"

    cat <<'EOF' > "$mock_bin/uname"
#!/bin/bash
case "$1" in
    -s) printf '%s\n' "${MOCK_UNAME_S:-Linux}" ;;
    -m) printf '%s\n' "${MOCK_UNAME_M:-x86_64}" ;;
    *) /usr/bin/uname "$@" ;;
esac
EOF

    # Serves the release archive only under MOCK_ARCHIVE_NAME, so a request for
    # any other name 404s the way GitHub would. Every requested URL is appended
    # to MOCK_URL_LOG so tests can assert which release was fetched.
    cat <<'EOF' > "$mock_bin/curl"
#!/bin/bash
set -euo pipefail

output_file=""
url=""
while (($#)); do
    case "$1" in
        -o) output_file="$2"; shift 2 ;;
        --retry | --retry-delay) shift 2 ;;
        -*) shift ;;
        *) url="$1"; shift ;;
    esac
done

[[ -n "$url" ]] || exit 1
printf '%s\n' "$url" >> "${MOCK_URL_LOG:-/dev/null}"

if [[ "$url" == *"checksums.txt" ]]; then
    [[ "${MOCK_FAIL_CHECKSUM_DOWNLOAD:-0}" == "1" ]] && exit 22
    [[ "$(basename "$url")" == "$MOCK_CHECKSUMS_NAME" ]] || exit 22
    printf '%s\n' "$MOCK_CHECKSUMS_CONTENT" > "$output_file"
    exit 0
fi

[[ "$(basename "$url")" == "$MOCK_ARCHIVE_NAME" ]] || exit 22
printf 'mock archive' > "$output_file"
EOF

    cat <<'EOF' > "$mock_bin/tar"
#!/bin/bash
set -euo pipefail

if [[ "${MOCK_TAR_FAIL:-0}" == "1" ]]; then
    echo "tar: pvtr: Not found in archive" >&2
    exit 1
fi

target_dir=""
while (($#)); do
    case "$1" in
        -C) target_dir="$2"; shift 2 ;;
        *) shift ;;
    esac
done

mkdir -p "$target_dir"
printf '#!/bin/sh\nexit 0\n' > "$target_dir/pvtr"
EOF

    cat <<'EOF' > "$mock_bin/sha256sum"
#!/bin/bash
printf '%s  %s\n' "${MOCK_ACTUAL_CHECKSUM:-expected-checksum}" "$1"
EOF

    chmod +x "$mock_bin/uname" "$mock_bin/curl" "$mock_bin/tar" "$mock_bin/sha256sum"
}

fail() {
    echo "FAIL: ${FUNCNAME[1]}: $1"
    [[ $# -gt 1 ]] && echo "$2"
    exit 1
}

assert_contains() {
    [[ "$1" == *"$2"* ]] || fail "expected output to contain: $2" "$1"
}

assert_not_contains() {
    [[ "$1" != *"$2"* ]] || fail "expected output NOT to contain: $2" "$1"
}

assert_equals() {
    [[ "$1" == "$2" ]] || fail "expected '$2', got '$1'"
}

# Evaluates a snippet with install.sh sourced, mocked commands, and an isolated
# HOME. Extra arguments are visible to the snippet as "$1", "$2", ... Stdin comes
# from SANDBOX_STDIN, which defaults to /dev/null so a stray prompt cannot hang.
# Sourcing must not trigger main; that guard is exercised by test_source_does_not_install.
run_in_sandbox() {
    local work_dir="$1" snippet="$2"
    shift 2

    local mock_bin="$work_dir/mock-bin"
    mkdir -p "$mock_bin" "$work_dir/home"
    make_mock_commands "$mock_bin"

    (
        export PATH="$mock_bin:$PATH"
        export HOME="$work_dir/home"
        export MOCK_UNAME_S MOCK_UNAME_M MOCK_ARCHIVE_NAME MOCK_CHECKSUMS_CONTENT
        export MOCK_CHECKSUMS_NAME MOCK_ACTUAL_CHECKSUM MOCK_FAIL_CHECKSUM_DOWNLOAD
        export MOCK_TAR_FAIL MOCK_URL_LOG SHELL PVTR_VERSION
        bash -c 'source "$1"; snippet="$2"; shift 2; eval "$snippet"' \
            _ "$INSTALL_SCRIPT" "$snippet" "$@" < "${SANDBOX_STDIN:-/dev/null}"
    )
}

# Assigns the caller's `work_dir` and resets mock state. Must be called
# directly, not in a command substitution, so the assignments survive.
setup() {
    MOCK_UNAME_S=Linux
    MOCK_UNAME_M=x86_64
    MOCK_ARCHIVE_NAME=pvtr_Linux_x86_64.tar.gz
    MOCK_CHECKSUMS_NAME=checksums.txt
    MOCK_CHECKSUMS_CONTENT='expected-checksum  pvtr_Linux_x86_64.tar.gz'
    MOCK_ACTUAL_CHECKSUM=expected-checksum
    MOCK_FAIL_CHECKSUM_DOWNLOAD=0
    MOCK_TAR_FAIL=0
    SHELL=/bin/bash
    PVTR_VERSION=""
    SANDBOX_STDIN=/dev/null
    work_dir="$(mktemp -d)"
    MOCK_URL_LOG="$work_dir/urls.txt"
}

# --- platform detection ------------------------------------------------------

test_detects_linux_arm64_from_aarch64() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    MOCK_UNAME_M=aarch64
    assert_equals "$(run_in_sandbox "$work_dir" 'detect_platform; echo "$PLATFORM"')" "Linux_arm64"
}

test_detects_universal_darwin_archive() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    # Darwin archives are universal, so an unrecognised arch must not matter.
    MOCK_UNAME_S=Darwin
    MOCK_UNAME_M=some_future_arch
    assert_equals "$(run_in_sandbox "$work_dir" 'detect_platform; echo "$PLATFORM"')" "Darwin_all"
}

test_rejects_unsupported_os() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    MOCK_UNAME_S=MINGW64_NT-10.0
    local output
    output="$(run_in_sandbox "$work_dir" 'detect_platform' 2>&1)" \
        && fail "expected detect_platform to reject an unsupported OS"
    assert_contains "$output" "Unsupported OS"
}

test_rejects_unsupported_arch() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    MOCK_UNAME_M=riscv64
    local output
    output="$(run_in_sandbox "$work_dir" 'detect_platform' 2>&1)" \
        && fail "expected detect_platform to reject an unsupported architecture"
    assert_contains "$output" "Unsupported architecture"
}

# detect_platform reports failure through `fail`, which exits. If it were called
# in a command substitution that exit would unwind only the subshell, and main
# would carry on and download with an empty platform.
test_unsupported_platform_aborts_main() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    MOCK_UNAME_S=MINGW64_NT-10.0
    local output
    output="$(run_in_sandbox "$work_dir" 'main -p "$1"' "$work_dir/install" 2>&1)" \
        && fail "expected main to abort on an unsupported platform"
    assert_contains "$output" "Unsupported OS"
    assert_not_contains "$output" "Failed to download"
    [[ ! -e "$MOCK_URL_LOG" ]] || fail "must not attempt a download after platform detection fails"
}

# --- prerequisites -----------------------------------------------------------

# A function named `command` shadows the builtin, which lets these tests
# simulate a missing tool without rebuilding PATH.
test_requires_a_checksum_tool() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    local output
    output="$(run_in_sandbox "$work_dir" \
        'command() { [[ "$2" == sha256sum || "$2" == shasum ]] && return 1; builtin command "$@"; }
         ensure_prerequisites' 2>&1)" \
        && fail "expected ensure_prerequisites to require a checksum tool"
    assert_contains "$output" "Neither sha256sum nor shasum is available"
}

test_requires_curl() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    local output
    output="$(run_in_sandbox "$work_dir" \
        'command() { [[ "$2" == curl ]] && return 1; builtin command "$@"; }
         ensure_prerequisites' 2>&1)" \
        && fail "expected ensure_prerequisites to require curl"
    assert_contains "$output" "curl is required"
}

test_prerequisites_pass_on_a_normal_system() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    run_in_sandbox "$work_dir" 'ensure_prerequisites' \
        || fail "expected ensure_prerequisites to pass on this machine"
}

# Missing tools must be reported before anything is fetched.
test_missing_tool_aborts_before_download() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    local output
    output="$(run_in_sandbox "$work_dir" \
        'command() { [[ "$2" == sha256sum || "$2" == shasum ]] && return 1; builtin command "$@"; }
         main -p "$1"' "$work_dir/install" 2>&1)" \
        && fail "expected main to abort when no checksum tool is available"
    assert_contains "$output" "Neither sha256sum nor shasum is available"
    assert_not_contains "$output" "Checksum mismatch"
    [[ ! -e "$MOCK_URL_LOG" ]] || fail "must not download before checking prerequisites"
}

# --- download and verification ----------------------------------------------

test_installs_when_checksum_matches() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    local output
    output="$(run_in_sandbox "$work_dir" 'download_release "$1" "$2" "$3" "$4"' \
        "$work_dir/install" Linux_x86_64 https://example.test/latest/download latest 2>&1)" \
        || fail "expected install to succeed when checksum matches" "$output"
    [[ -f "$work_dir/install/pvtr" ]] || fail "expected pvtr to be installed"
    [[ -x "$work_dir/install/pvtr" ]] || fail "expected pvtr to be executable"
}

test_falls_back_to_legacy_archive_name() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    MOCK_ARCHIVE_NAME=privateer_Linux_x86_64.tar.gz
    MOCK_CHECKSUMS_CONTENT='expected-checksum  privateer_Linux_x86_64.tar.gz'

    local output
    output="$(run_in_sandbox "$work_dir" 'download_release "$1" "$2" "$3" "$4"' \
        "$work_dir/install" Linux_x86_64 https://example.test/latest/download latest 2>&1)" \
        || fail "expected fallback to legacy privateer_* archive name" "$output"
    assert_contains "$output" "privateer_Linux_x86_64.tar.gz"
}

# v0.21.0 and earlier named the checksums file privateer_<version>_checksums.txt.
test_falls_back_to_legacy_checksums_name() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    MOCK_ARCHIVE_NAME=privateer_Linux_x86_64.tar.gz
    MOCK_CHECKSUMS_NAME=privateer_0.21.0_checksums.txt
    MOCK_CHECKSUMS_CONTENT='expected-checksum  privateer_Linux_x86_64.tar.gz'

    local output
    output="$(run_in_sandbox "$work_dir" 'download_release "$1" "$2" "$3" "$4"' \
        "$work_dir/install" Linux_x86_64 https://example.test/download/v0.21.0 v0.21.0 2>&1)" \
        || fail "expected fallback to the legacy checksums file name" "$output"
    assert_contains "$output" "Verifying against privateer_0.21.0_checksums.txt"
    [[ -f "$work_dir/install/pvtr" ]] || fail "expected pvtr to be installed"
}

# The legacy name embeds a version number, so it is unknowable for "latest" and
# must not be probed there.
test_latest_does_not_probe_legacy_checksums_name() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    run_in_sandbox "$work_dir" 'download_release "$1" "$2" "$3" "$4"' \
        "$work_dir/install" Linux_x86_64 https://example.test/latest/download latest > /dev/null
    assert_not_contains "$(cat "$MOCK_URL_LOG")" "privateer_latest_checksums.txt"
}

test_fails_when_no_archive_found() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    MOCK_ARCHIVE_NAME=pvtr_Linux_riscv64.tar.gz

    local output
    output="$(run_in_sandbox "$work_dir" 'download_release "$1" "$2" "$3" "$4"' \
        "$work_dir/install" Linux_x86_64 https://example.test/latest/download latest 2>&1)" \
        && fail "expected install to fail when no archive matches"
    assert_contains "$output" "Failed to download a release archive"
}

test_fails_when_checksum_download_fails() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    MOCK_FAIL_CHECKSUM_DOWNLOAD=1

    local output
    output="$(run_in_sandbox "$work_dir" 'download_release "$1" "$2" "$3" "$4"' \
        "$work_dir/install" Linux_x86_64 https://example.test/latest/download latest 2>&1)" \
        && fail "expected install to fail when checksums download fails"
    assert_contains "$output" "Failed to download a checksums file"
    [[ ! -f "$work_dir/install/pvtr" ]] || fail "must not install an unverified binary"
}

test_fails_when_archive_missing_from_checksums() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    MOCK_CHECKSUMS_CONTENT='expected-checksum  some_other_archive.tar.gz'

    local output
    output="$(run_in_sandbox "$work_dir" 'download_release "$1" "$2" "$3" "$4"' \
        "$work_dir/install" Linux_x86_64 https://example.test/latest/download latest 2>&1)" \
        && fail "expected install to fail when the archive has no checksum entry"
    assert_contains "$output" "No entry for"
    [[ ! -f "$work_dir/install/pvtr" ]] || fail "must not install an unverified binary"
}

test_fails_when_checksum_mismatches() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    MOCK_ACTUAL_CHECKSUM=different-checksum

    local output
    output="$(run_in_sandbox "$work_dir" 'download_release "$1" "$2" "$3" "$4"' \
        "$work_dir/install" Linux_x86_64 https://example.test/latest/download latest 2>&1)" \
        && fail "expected install to fail when checksum mismatches"
    assert_contains "$output" "Checksum mismatch"
    [[ ! -f "$work_dir/install/pvtr" ]] || fail "must not install a tampered binary"
}

# The binary is staged and renamed into place, so a failure part-way through
# must not leave anything at the install path.
test_leaves_no_binary_when_extraction_fails() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    MOCK_TAR_FAIL=1

    local output
    output="$(run_in_sandbox "$work_dir" 'download_release "$1" "$2" "$3" "$4"' \
        "$work_dir/install" Linux_x86_64 https://example.test/latest/download latest 2>&1)" \
        && fail "expected install to fail when extraction fails"
    assert_contains "$output" "does not contain a pvtr binary"
    [[ ! -e "$work_dir/install/pvtr" ]] || fail "must not leave a partial binary behind"
}

test_removes_staged_file_on_success() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    run_in_sandbox "$work_dir" 'download_release "$1" "$2" "$3" "$4"' \
        "$work_dir/install" Linux_x86_64 https://example.test/latest/download latest > /dev/null
    local leftovers
    leftovers="$(find "$work_dir/install" -name '.pvtr.download.*' | wc -l | tr -d ' ')"
    assert_equals "$leftovers" "0"
}

# --- version selection -------------------------------------------------------

test_defaults_to_latest_release() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    assert_equals "$(run_in_sandbox "$work_dir" 'release_base_url latest')" \
        "https://github.com/privateerproj/privateer/releases/latest/download"
}

test_pins_requested_version() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    assert_equals "$(run_in_sandbox "$work_dir" 'release_base_url v0.22.0')" \
        "https://github.com/privateerproj/privateer/releases/download/v0.22.0"
}

test_normalizes_version_without_v_prefix() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    assert_equals "$(run_in_sandbox "$work_dir" 'release_base_url 0.22.0')" \
        "https://github.com/privateerproj/privateer/releases/download/v0.22.0"
}

test_version_flag_reaches_the_download_url() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    run_in_sandbox "$work_dir" 'main -p "$1" -v v0.22.0 -y' "$work_dir/install" > /dev/null
    assert_contains "$(cat "$MOCK_URL_LOG")" "/releases/download/v0.22.0/pvtr_Linux_x86_64.tar.gz"
}

test_version_env_var_is_honoured() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    PVTR_VERSION=v0.21.0
    run_in_sandbox "$work_dir" 'main -p "$1" -y' "$work_dir/install" > /dev/null
    assert_contains "$(cat "$MOCK_URL_LOG")" "/releases/download/v0.21.0/"
}

test_version_flag_beats_env_var() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    PVTR_VERSION=v0.21.0
    run_in_sandbox "$work_dir" 'main -p "$1" -v v0.22.0 -y' "$work_dir/install" > /dev/null
    assert_contains "$(cat "$MOCK_URL_LOG")" "/releases/download/v0.22.0/"
    assert_not_contains "$(cat "$MOCK_URL_LOG")" "v0.21.0"
}

# --- install directory -------------------------------------------------------

test_fails_before_downloading_when_dir_not_writable() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    if [[ "$EUID" -eq 0 ]]; then
        echo "SKIP: ${FUNCNAME[0]} (running as root)"
        return
    fi

    local readonly_dir="$work_dir/readonly"
    mkdir -p "$readonly_dir"
    chmod 555 "$readonly_dir"

    local output
    output="$(run_in_sandbox "$work_dir" 'main -p "$1"' "$readonly_dir/bin" 2>&1)" \
        && fail "expected main to fail when the install dir is not writable"
    assert_contains "$output" "not writable"
    [[ ! -e "$MOCK_URL_LOG" ]] || fail "must not download before checking writability"
    chmod 755 "$readonly_dir"
}

test_relative_install_dir_is_made_absolute() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    local output
    output="$(cd "$work_dir" && run_in_sandbox "$work_dir" 'main -p "$1" -y' relative-bin 2>&1)"
    assert_contains "$output" "Installed $work_dir/relative-bin/pvtr"
    assert_contains "$(cat "$work_dir/home/.bashrc")" "export PATH=\"$work_dir/relative-bin:\$PATH\""
}

# A './' left in the path would leak into the PATH line printed to the user.
test_install_dir_is_normalized() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    assert_equals "$(run_in_sandbox "$work_dir" 'abspath ./e2e/bin')" "$PWD/e2e/bin"
    assert_equals "$(run_in_sandbox "$work_dir" 'abspath /opt//pvtr/./bin/')" "/opt/pvtr/bin"
    assert_equals "$(run_in_sandbox "$work_dir" 'abspath /')" "/"
    assert_equals "$(run_in_sandbox "$work_dir" 'abspath /already/absolute')" "/already/absolute"
}

# --- PATH handling -----------------------------------------------------------

test_path_is_not_modified_without_consent() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    local output
    output="$(run_in_sandbox "$work_dir" 'update_path "$1" false' "$work_dir/install" 2>&1)"
    assert_contains "$output" "add $work_dir/install to your PATH"
    [[ ! -f "$work_dir/home/.bashrc" ]] || fail "must not edit shell config without consent"
}

test_path_is_modified_with_consent() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    run_in_sandbox "$work_dir" 'update_path "$1" true' "$work_dir/install" > /dev/null
    assert_contains "$(cat "$work_dir/home/.bashrc")" "export PATH=\"$work_dir/install:\$PATH\""

    local output
    output="$(run_in_sandbox "$work_dir" 'update_path "$1" true' "$work_dir/install" 2>&1)"
    assert_contains "$output" "already references"
    assert_equals "$(grep -c "$work_dir/install" "$work_dir/home/.bashrc")" "1"
}

test_prompt_accepts_yes() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    mkdir -p "$work_dir/home"
    printf 'y\n' > "$work_dir/reply"
    SANDBOX_STDIN="$work_dir/reply"

    run_in_sandbox "$work_dir" 'is_interactive() { return 0; }; update_path "$1" false' \
        "$work_dir/install" > /dev/null 2>&1
    assert_contains "$(cat "$work_dir/home/.bashrc")" "export PATH=\"$work_dir/install:\$PATH\""
}

# Ctrl-D makes `read` return non-zero, which under `set -e` would otherwise
# abort the script after the binary is already installed.
test_eof_at_prompt_declines_without_aborting() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    local output
    output="$(run_in_sandbox "$work_dir" \
        'is_interactive() { return 0; }; update_path "$1" false; echo REACHED_END' \
        "$work_dir/install" 2>&1)" \
        || fail "expected EOF at the prompt to decline, not abort" "$output"
    assert_contains "$output" "REACHED_END"
    assert_contains "$output" "add $work_dir/install to your PATH"
    [[ ! -f "$work_dir/home/.bashrc" ]] || fail "must not edit shell config on EOF"
}

# macOS terminals start login shells, which read .bash_profile, not .bashrc.
test_bash_on_macos_uses_bash_profile() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    MOCK_UNAME_S=Darwin
    run_in_sandbox "$work_dir" 'update_path "$1" true' "$work_dir/install" > /dev/null
    assert_contains "$(cat "$work_dir/home/.bash_profile")" "export PATH=\"$work_dir/install:\$PATH\""
    [[ ! -f "$work_dir/home/.bashrc" ]] || fail "must not write .bashrc on macOS"
}

test_bash_on_linux_uses_bashrc() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    run_in_sandbox "$work_dir" 'update_path "$1" true' "$work_dir/install" > /dev/null
    assert_contains "$(cat "$work_dir/home/.bashrc")" "export PATH=\"$work_dir/install:\$PATH\""
    [[ ! -f "$work_dir/home/.bash_profile" ]] || fail "must not write .bash_profile on Linux"
}

test_fish_gets_native_path_syntax() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    SHELL=/usr/local/bin/fish
    run_in_sandbox "$work_dir" 'update_path "$1" true' "$work_dir/install" > /dev/null
    assert_contains "$(cat "$work_dir/home/.config/fish/config.fish")" "fish_add_path \"$work_dir/install\""
}

test_commented_out_path_line_is_not_treated_as_configured() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    mkdir -p "$work_dir/home"
    printf '# export PATH="%s:$PATH"\n' "$work_dir/install" > "$work_dir/home/.bashrc"

    run_in_sandbox "$work_dir" 'update_path "$1" true' "$work_dir/install" > /dev/null
    assert_equals "$(grep -c "^export PATH" "$work_dir/home/.bashrc")" "1"
}

# --- argument handling -------------------------------------------------------

test_rejects_unexpected_arguments() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    local output
    output="$(run_in_sandbox "$work_dir" 'main -y stray-argument' 2>&1)" \
        && fail "expected main to reject a stray positional argument"
    assert_contains "$output" "unexpected argument: stray-argument"
}

test_rejects_unknown_flag() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    local output
    output="$(run_in_sandbox "$work_dir" 'main -Z' 2>&1)" \
        && fail "expected main to reject an unknown flag"
    assert_contains "$output" "Usage: install.sh"
}

# --- invocation modes --------------------------------------------------------

# The README installs via `bash -c "$(curl ...)"`, where BASH_SOURCE is unset.
test_runs_when_piped_to_bash() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    local output
    output="$(/bin/bash -c "$(cat "$INSTALL_SCRIPT")" -- -h < /dev/null)" \
        || fail "expected piped invocation to run main"
    assert_contains "$output" "Usage: install.sh"
}

test_source_does_not_install() {
    local work_dir; setup
    trap 'rm -rf "$work_dir"' RETURN

    local output
    output="$(run_in_sandbox "$work_dir" ':' 2>&1)"
    assert_equals "$output" ""
}

test_detects_linux_arm64_from_aarch64
test_detects_universal_darwin_archive
test_rejects_unsupported_os
test_rejects_unsupported_arch
test_unsupported_platform_aborts_main
test_requires_a_checksum_tool
test_requires_curl
test_prerequisites_pass_on_a_normal_system
test_missing_tool_aborts_before_download
test_installs_when_checksum_matches
test_falls_back_to_legacy_archive_name
test_falls_back_to_legacy_checksums_name
test_latest_does_not_probe_legacy_checksums_name
test_fails_when_no_archive_found
test_fails_when_checksum_download_fails
test_fails_when_archive_missing_from_checksums
test_fails_when_checksum_mismatches
test_leaves_no_binary_when_extraction_fails
test_removes_staged_file_on_success
test_defaults_to_latest_release
test_pins_requested_version
test_normalizes_version_without_v_prefix
test_version_flag_reaches_the_download_url
test_version_env_var_is_honoured
test_version_flag_beats_env_var
test_fails_before_downloading_when_dir_not_writable
test_relative_install_dir_is_made_absolute
test_install_dir_is_normalized
test_path_is_not_modified_without_consent
test_path_is_modified_with_consent
test_prompt_accepts_yes
test_eof_at_prompt_declines_without_aborting
test_bash_on_macos_uses_bash_profile
test_bash_on_linux_uses_bashrc
test_fish_gets_native_path_syntax
test_commented_out_path_line_is_not_treated_as_configured
test_rejects_unexpected_arguments
test_rejects_unknown_flag
test_runs_when_piped_to_bash
test_source_does_not_install

echo "All install.sh tests passed."
