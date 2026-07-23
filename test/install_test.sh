#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
INSTALL_SCRIPT="$REPO_ROOT/install.sh"

FAILURES=0

# The mock curl serves a fixed set of asset names. Anything else 404s the way
# the real release endpoint would, which is what exercises the prefix fallback
# and the unsupported-platform paths.
make_mock_commands() {
    local mock_bin="$1"

    cat <<'EOF' > "$mock_bin/uname"
#!/usr/bin/env bash
case "$1" in
    -s) printf '%s\n' "${MOCK_UNAME_S:-Linux}" ;;
    -m) printf '%s\n' "${MOCK_UNAME_M:-x86_64}" ;;
    *) /usr/bin/uname "$@" ;;
esac
EOF

    cat <<'EOF' > "$mock_bin/curl"
#!/usr/bin/env bash
set -euo pipefail

output_file=""
url=""
while (($#)); do
    case "$1" in
        -o) output_file="$2"; shift 2 ;;
        -*) shift ;;
        *) url="$1"; shift ;;
    esac
done

asset="${url##*/}"

emit() {
    if [[ -n "$output_file" ]]; then
        printf '%s' "$1" > "$output_file"
    else
        printf '%s' "$1"
    fi
}

if [[ "$asset" == "checksums.txt" ]]; then
    if [[ "${MOCK_FAIL_CHECKSUM_DOWNLOAD:-0}" == "1" ]]; then
        exit 22
    fi
    emit "$MOCK_CHECKSUMS_CONTENT"
    exit 0
fi

# 404 for any asset the release does not contain.
case " ${MOCK_AVAILABLE_ASSETS} " in
    *" ${asset} "*) ;;
    *) exit 22 ;;
esac

emit "archive:${asset}"
EOF

    # Stands in for tar: writes the payload the real archive would contain.
    cat <<'EOF' > "$mock_bin/tar"
#!/usr/bin/env bash
set -euo pipefail

target_dir=""
while (($#)); do
    case "$1" in
        -C) target_dir="$2"; shift 2 ;;
        *) shift ;;
    esac
done

mkdir -p "$target_dir"
{
    printf '#!/bin/sh\n'
    printf 'printf "%s\\n"\n' "${MOCK_BINARY_VERSION:-0.0.0}"
} > "$target_dir/pvtr"
chmod +x "$target_dir/pvtr"
printf 'license\n' > "$target_dir/LICENSE"
EOF

    cat <<'EOF' > "$mock_bin/sha256sum"
#!/usr/bin/env bash
printf '%s  %s\n' "${MOCK_ACTUAL_CHECKSUM:-expected-checksum}" "$1"
EOF

    chmod +x "$mock_bin/uname" "$mock_bin/curl" "$mock_bin/tar" "$mock_bin/sha256sum"
}

# Run install.sh with mocks in front of PATH. Echoes combined output; returns the
# script's exit status.
run_install() {
    local work_dir="$1"
    shift

    local mock_bin="$work_dir/mock-bin"
    mkdir -p "$mock_bin" "$work_dir/home"
    make_mock_commands "$mock_bin"

    (
        export PATH="$mock_bin:$PATH"
        export HOME="$work_dir/home"
        export SHELL="/bin/bash"
        bash "$INSTALL_SCRIPT" "$@" 2>&1
    )
}

fail() {
    echo "FAIL: $1"
    [[ $# -lt 2 ]] || echo "$2"
    FAILURES=$((FAILURES + 1))
}

assert_contains() {
    local haystack="$1" needle="$2" context="$3"

    if [[ "$haystack" != *"$needle"* ]]; then
        fail "$context: expected output to contain '$needle'" "$haystack"
    fi
}

setup_defaults() {
    MOCK_UNAME_S="Linux"
    MOCK_UNAME_M="x86_64"
    MOCK_AVAILABLE_ASSETS="pvtr_Linux_x86_64.tar.gz"
    MOCK_CHECKSUMS_CONTENT='expected-checksum  pvtr_Linux_x86_64.tar.gz'
    MOCK_ACTUAL_CHECKSUM="expected-checksum"
    MOCK_FAIL_CHECKSUM_DOWNLOAD=0
    MOCK_BINARY_VERSION="0.22.0"
    export MOCK_UNAME_S MOCK_UNAME_M MOCK_AVAILABLE_ASSETS MOCK_CHECKSUMS_CONTENT \
        MOCK_ACTUAL_CHECKSUM MOCK_FAIL_CHECKSUM_DOWNLOAD MOCK_BINARY_VERSION
}

test_installs_to_default_dir() {
    local work_dir; work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN
    setup_defaults

    local output
    if ! output=$(run_install "$work_dir"); then
        fail "default install should succeed" "$output"
        return
    fi

    [[ -x "$work_dir/home/.local/bin/pvtr" ]] || fail "expected binary in ~/.local/bin"
    # The archive's LICENSE must not be dropped into the bin directory.
    [[ ! -e "$work_dir/home/.local/bin/LICENSE" ]] || fail "LICENSE leaked into the install dir"
    assert_contains "$output" "Installed pvtr 0.22.0" "default install"
}

# Regression: uname -m reports aarch64 on Linux arm64, which the previous script
# rejected as unsupported even though the release contains that archive.
test_installs_on_linux_aarch64() {
    local work_dir; work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN
    setup_defaults
    MOCK_UNAME_M="aarch64"
    MOCK_AVAILABLE_ASSETS="pvtr_Linux_arm64.tar.gz"
    MOCK_CHECKSUMS_CONTENT='expected-checksum  pvtr_Linux_arm64.tar.gz'

    local output
    if ! output=$(run_install "$work_dir"); then
        fail "aarch64 install should succeed" "$output"
        return
    fi

    [[ -x "$work_dir/home/.local/bin/pvtr" ]] || fail "expected binary for aarch64"
}

test_darwin_uses_universal_archive() {
    local work_dir; work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN
    setup_defaults
    MOCK_UNAME_S="Darwin"
    MOCK_UNAME_M="arm64"
    MOCK_AVAILABLE_ASSETS="pvtr_Darwin_all.tar.gz"
    MOCK_CHECKSUMS_CONTENT='expected-checksum  pvtr_Darwin_all.tar.gz'

    local output
    if ! output=$(run_install "$work_dir"); then
        fail "darwin install should succeed" "$output"
    fi
}

# Releases published before the project rename still carry privateer_* archives.
test_falls_back_to_legacy_asset_prefix() {
    local work_dir; work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN
    setup_defaults
    MOCK_AVAILABLE_ASSETS="privateer_Linux_x86_64.tar.gz"
    MOCK_CHECKSUMS_CONTENT='expected-checksum  privateer_Linux_x86_64.tar.gz'

    local output
    if ! output=$(run_install "$work_dir" -v v0.21.2); then
        fail "legacy prefix install should succeed" "$output"
        return
    fi

    [[ -x "$work_dir/home/.local/bin/pvtr" ]] || fail "expected binary from legacy-prefix archive"
}

test_fails_when_no_asset_matches() {
    local work_dir; work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN
    setup_defaults
    MOCK_AVAILABLE_ASSETS=""

    local output
    if output=$(run_install "$work_dir"); then
        fail "install should fail when no archive resolves" "$output"
        return
    fi

    assert_contains "$output" "Could not download a release archive" "missing asset"
}

test_fails_on_unsupported_architecture() {
    local work_dir; work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN
    setup_defaults
    MOCK_UNAME_M="mips"

    local output
    if output=$(run_install "$work_dir"); then
        fail "install should fail on an unsupported architecture" "$output"
        return
    fi

    assert_contains "$output" "Unsupported architecture: mips" "unsupported arch"
}

test_windows_directs_to_manual_install() {
    local work_dir; work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN
    setup_defaults
    MOCK_UNAME_S="MINGW64_NT-10.0"

    local output
    if output=$(run_install "$work_dir"); then
        fail "install should fail on Windows rather than extracting a zip with tar" "$output"
        return
    fi

    assert_contains "$output" "cannot install on Windows" "windows"
}

test_fails_when_checksum_mismatches() {
    local work_dir; work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN
    setup_defaults
    MOCK_ACTUAL_CHECKSUM="different-checksum"

    local output
    if output=$(run_install "$work_dir"); then
        fail "install should fail on checksum mismatch" "$output"
        return
    fi

    assert_contains "$output" "CHECKSUM MISMATCH" "checksum mismatch"
    [[ ! -e "$work_dir/home/.local/bin/pvtr" ]] || fail "binary installed despite checksum mismatch"
}

test_fails_when_checksum_entry_missing() {
    local work_dir; work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN
    setup_defaults
    MOCK_CHECKSUMS_CONTENT='expected-checksum  some_other_archive.tar.gz'

    local output
    if output=$(run_install "$work_dir"); then
        fail "install should fail when the archive has no checksum entry" "$output"
        return
    fi

    assert_contains "$output" "Could not find a checksum" "missing checksum entry"
}

test_fails_when_checksum_download_fails() {
    local work_dir; work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN
    setup_defaults
    MOCK_FAIL_CHECKSUM_DOWNLOAD=1

    local output
    if output=$(run_install "$work_dir"); then
        fail "install should fail when checksums cannot be downloaded" "$output"
        return
    fi

    assert_contains "$output" "Failed to download checksums file" "checksum download failure"
}

test_respects_install_path_flag() {
    local work_dir; work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN
    setup_defaults

    local target="$work_dir/custom/bin"
    local output
    if ! output=$(run_install "$work_dir" -p "$target"); then
        fail "-p install should succeed" "$output"
        return
    fi

    [[ -x "$target/pvtr" ]] || fail "expected binary at the -p path"
}

test_pinned_version_uses_tagged_url() {
    local work_dir; work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN
    setup_defaults

    local output
    if ! output=$(run_install "$work_dir" -v v0.22.0); then
        fail "pinned install should succeed" "$output"
        return
    fi

    assert_contains "$output" "releases/download/v0.22.0" "pinned version"
}

# Existing installs live in the plugin directory. Upgrading in place keeps a
# stale binary from shadowing the new one on PATH.
test_upgrades_existing_legacy_install_in_place() {
    local work_dir; work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN
    setup_defaults

    mkdir -p "$work_dir/home/.privateer/bin"
    printf '#!/bin/sh\nprintf "0.21.2\\n"\n' > "$work_dir/home/.privateer/bin/pvtr"
    chmod +x "$work_dir/home/.privateer/bin/pvtr"

    local output
    if ! output=$(run_install "$work_dir"); then
        fail "legacy upgrade should succeed" "$output"
        return
    fi

    assert_contains "$output" "Upgraded pvtr 0.21.2 -> 0.22.0" "legacy upgrade"
    [[ ! -e "$work_dir/home/.local/bin/pvtr" ]] || fail "legacy install should not be duplicated in ~/.local/bin"
}

test_reports_path_without_editing_shell_config() {
    local work_dir; work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN
    setup_defaults

    mkdir -p "$work_dir/home"
    printf '# untouched\n' > "$work_dir/home/.bash_profile"

    local output
    output=$(run_install "$work_dir")

    assert_contains "$output" "is not on your PATH" "path reporting"
    [[ "$(cat "$work_dir/home/.bash_profile")" == "# untouched" ]] \
        || fail "install.sh must not write to shell configuration"
}

test_installs_to_default_dir
test_installs_on_linux_aarch64
test_darwin_uses_universal_archive
test_falls_back_to_legacy_asset_prefix
test_fails_when_no_asset_matches
test_fails_on_unsupported_architecture
test_windows_directs_to_manual_install
test_fails_when_checksum_mismatches
test_fails_when_checksum_entry_missing
test_fails_when_checksum_download_fails
test_respects_install_path_flag
test_pinned_version_uses_tagged_url
test_upgrades_existing_legacy_install_in_place
test_reports_path_without_editing_shell_config

if ((FAILURES > 0)); then
    echo "${FAILURES} test(s) failed"
    exit 1
fi

echo "All install.sh tests passed"
