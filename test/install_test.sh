#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
INSTALL_SCRIPT="$REPO_ROOT/install.sh"

release_json_with_checksums() {
        cat <<'EOF'
{
    "assets": [
        {
            "browser_download_url": "https://example.test/pvtr_Linux_x86_64.tar.gz"
        },
        {
            "browser_download_url": "https://example.test/checksums.txt"
        }
    ]
}
EOF
}

release_json_without_checksums() {
        cat <<'EOF'
{
    "assets": [
        {
            "browser_download_url": "https://example.test/pvtr_Linux_x86_64.tar.gz"
        }
    ]
}
EOF
}

release_json_with_checksums_compact() {
    printf '%s' '{"assets":[{"browser_download_url":"https://example.test/pvtr_Linux_x86_64.tar.gz"},{"browser_download_url":"https://example.test/checksums.txt"}]}'
}

make_mock_commands() {
    local mock_bin="$1"

    cat <<'EOF' > "$mock_bin/uname"
#!/bin/bash
case "$1" in
    -s)
        printf '%s\n' "${MOCK_UNAME_S:-Linux}"
        ;;
    -m)
        printf '%s\n' "${MOCK_UNAME_M:-x86_64}"
        ;;
    *)
        /usr/bin/uname "$@"
        ;;
esac
EOF

    cat <<'EOF' > "$mock_bin/curl"
#!/bin/bash
set -euo pipefail

output_file=""
url=""

while (($#)); do
    case "$1" in
        -o)
            output_file="$2"
            shift 2
            ;;
        -f|-S|-s|-L|-SL|-fSL)
            shift
            ;;
        -*)
            shift
            ;;
        *)
            url="$1"
            shift
            ;;
    esac
done

if [[ -z "$url" ]]; then
    exit 1
fi

if [[ "$url" == *"/releases/latest" ]]; then
    if [[ -n "$output_file" ]]; then
        printf '%s' "$MOCK_RELEASE_JSON" > "$output_file"
    else
        printf '%s' "$MOCK_RELEASE_JSON"
    fi
    exit 0
fi

if [[ "$url" == *"checksums.txt" ]]; then
    if [[ "${MOCK_FAIL_CHECKSUM_DOWNLOAD:-0}" == "1" ]]; then
        exit 22
    fi

    if [[ -n "$output_file" ]]; then
        printf '%s' "$MOCK_CHECKSUMS_CONTENT" > "$output_file"
    else
        printf '%s' "$MOCK_CHECKSUMS_CONTENT"
    fi
    exit 0
fi

if [[ -n "$output_file" ]]; then
    printf 'mock archive' > "$output_file"
else
    printf 'mock archive'
fi
EOF

    cat <<'EOF' > "$mock_bin/tar"
#!/bin/bash
set -euo pipefail

target_dir=""

while (($#)); do
    case "$1" in
        -C)
            target_dir="$2"
            shift 2
            ;;
        *)
            shift
            ;;
    esac
done

mkdir -p "$target_dir"
printf '#!/bin/sh\nexit 0\n' > "$target_dir/pvtr"
EOF

    cat <<'EOF' > "$mock_bin/chmod"
#!/bin/bash
exit 0
EOF

    cat <<'EOF' > "$mock_bin/sha256sum"
#!/bin/bash
printf '%s  %s\n' "${MOCK_ACTUAL_CHECKSUM:-expected-checksum}" "$1"
EOF

    chmod +x "$mock_bin/uname" "$mock_bin/curl" "$mock_bin/tar" "$mock_bin/chmod" "$mock_bin/sha256sum"
}

assert_contains() {
    local haystack="$1"
    local needle="$2"

    if [[ "$haystack" != *"$needle"* ]]; then
        echo "expected output to contain: $needle"
        echo "$haystack"
        exit 1
    fi
}

run_install() {
    local work_dir="$1"
    local install_dir="$2"

    local mock_bin="$work_dir/mock-bin"
    mkdir -p "$mock_bin"
    make_mock_commands "$mock_bin"

    (
        export PATH="$mock_bin:$PATH"
        export HOME="$work_dir/home"
        mkdir -p "$HOME"
        export MOCK_RELEASE_JSON MOCK_CHECKSUMS_CONTENT MOCK_ACTUAL_CHECKSUM MOCK_FAIL_CHECKSUM_DOWNLOAD
        bash -c 'source "$1"; download_latest_release "$2"' _ "$INSTALL_SCRIPT" "$install_dir"
    )
}

test_installs_when_checksum_matches() {
    local work_dir
    work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN

    local install_dir="$work_dir/install"
    MOCK_RELEASE_JSON="$(release_json_with_checksums)"
    MOCK_CHECKSUMS_CONTENT='expected-checksum  pvtr_Linux_x86_64.tar.gz'
    MOCK_ACTUAL_CHECKSUM='expected-checksum'
    MOCK_FAIL_CHECKSUM_DOWNLOAD=0

    run_install "$work_dir" "$install_dir"

    [[ -f "$install_dir/pvtr" ]]
}

test_fails_when_checksum_asset_missing() {
    local work_dir
    work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN

    local install_dir="$work_dir/install"
    MOCK_RELEASE_JSON="$(release_json_without_checksums)"
    MOCK_CHECKSUMS_CONTENT=''
    MOCK_ACTUAL_CHECKSUM='expected-checksum'
    MOCK_FAIL_CHECKSUM_DOWNLOAD=0

    local output
    if output=$(run_install "$work_dir" "$install_dir" 2>&1); then
        echo "expected install to fail when checksum asset is missing"
        exit 1
    fi

    assert_contains "$output" "No checksums file found in release assets"
}

test_fails_when_checksum_download_fails() {
    local work_dir
    work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN

    local install_dir="$work_dir/install"
    MOCK_RELEASE_JSON="$(release_json_with_checksums)"
    MOCK_CHECKSUMS_CONTENT='expected-checksum  pvtr_Linux_x86_64.tar.gz'
    MOCK_ACTUAL_CHECKSUM='expected-checksum'
    MOCK_FAIL_CHECKSUM_DOWNLOAD=1

    local output
    if output=$(run_install "$work_dir" "$install_dir" 2>&1); then
        echo "expected install to fail when checksum download fails"
        exit 1
    fi

    assert_contains "$output" "Failed to download checksums file"
}

test_fails_when_checksum_mismatches() {
    local work_dir
    work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN

    local install_dir="$work_dir/install"
    MOCK_RELEASE_JSON="$(release_json_with_checksums)"
    MOCK_CHECKSUMS_CONTENT='expected-checksum  pvtr_Linux_x86_64.tar.gz'
    MOCK_ACTUAL_CHECKSUM='different-checksum'
    MOCK_FAIL_CHECKSUM_DOWNLOAD=0

    local output
    if output=$(run_install "$work_dir" "$install_dir" 2>&1); then
        echo "expected install to fail when checksum mismatches"
        exit 1
    fi

    assert_contains "$output" "CHECKSUM MISMATCH"
}

test_installs_when_release_json_is_compact() {
    local work_dir
    work_dir="$(mktemp -d)"
    trap 'rm -rf "$work_dir"' RETURN

    local install_dir="$work_dir/install"
    MOCK_RELEASE_JSON="$(release_json_with_checksums_compact)"
    MOCK_CHECKSUMS_CONTENT='expected-checksum  pvtr_Linux_x86_64.tar.gz'
    MOCK_ACTUAL_CHECKSUM='expected-checksum'
    MOCK_FAIL_CHECKSUM_DOWNLOAD=0

    run_install "$work_dir" "$install_dir"

    [[ -f "$install_dir/pvtr" ]]
}

test_installs_when_checksum_matches
test_fails_when_checksum_asset_missing
test_fails_when_checksum_download_fails
test_fails_when_checksum_mismatches
test_installs_when_release_json_is_compact
