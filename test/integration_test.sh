#!/bin/sh

set -x

STATUS=0

# expect_exit asserts that the previous command's exit code matches.
# Usage: expect_exit <expected_code> <actual_code> <description>
expect_exit() {
  expected=$1
  actual=$2
  description=$3
  if [ "$actual" -ne "$expected" ]; then
    echo "ERROR: $description: expected exit $expected, got $actual"
    STATUS=1
  fi
}

# Require GITHUB_TOKEN to be set
if [ -z "$GITHUB_TOKEN" ]; then
  echo "ERROR: GITHUB_TOKEN environment variable is not set"
  echo "You can do the following to set it:"
  echo "  \`gh auth login\` and follow the prompts to authenticate with GitHub"
  echo "  export GITHUB_TOKEN=\$(gh auth token)"
  exit 1
fi

# Run basic pvtr commands to verify installation
./pvtr completion || STATUS=1
./pvtr env || STATUS=1
./pvtr help || STATUS=1
./pvtr list || STATUS=1
./pvtr version || STATUS=1

# generate-plugin: exit-code contract checks.
# Codes are defined in privateer-sdk/command/run.go:
#   TestPass=0  TestFail=1  Aborted=2  InternalError=3  BadUsage=4

# BadUsage (4): no required flags supplied
./pvtr generate-plugin
expect_exit 4 $? "generate-plugin with no flags"

# BadUsage (4): some but not all required flags supplied
./pvtr generate-plugin --source-path=test/data/OSPS_Baseline_2026_02.yaml --service-name=testsvc
expect_exit 4 $? "generate-plugin missing --organization"

# Set up a temp workspace for the remaining generate-plugin cases.
# Use an absolute catalog path derived from the script location so the test is
# robust to whatever cwd the caller invokes us from.
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CATALOG_PATH="$SCRIPT_DIR/data/OSPS_Baseline_2026_02.yaml"
GP_WORK="$(mktemp -d)"
GP_TEMPLATES_OK="$GP_WORK/templates-ok"
GP_TEMPLATES_PARTIAL="$GP_WORK/templates-partial"
GP_OUT_HAPPY="$GP_WORK/out-happy"
GP_OUT_PARTIAL="$GP_WORK/out-partial"
GP_OUT_INTERNAL="$GP_WORK/out-internal"
mkdir -p "$GP_TEMPLATES_OK" "$GP_TEMPLATES_PARTIAL"

# Renderable template for the happy path.
printf 'service: {{.ServiceName}}\norg: {{.Organization}}\n' > "$GP_TEMPLATES_OK/info.yaml.tmpl"

# Same renderable template plus one with an unregistered template function so
# tmpl.Parse fails -- exercises the partial-templating path (TestFail).
cp "$GP_TEMPLATES_OK/info.yaml.tmpl" "$GP_TEMPLATES_PARTIAL/info.yaml.tmpl"
printf '{{ thisFunctionDoesNotExist "x" }}\n' > "$GP_TEMPLATES_PARTIAL/broken.tmpl"

# InternalError (3): catalog source path that does not exist
./pvtr generate-plugin \
  --source-path="$GP_WORK/does-not-exist.yaml" \
  --service-name=testsvc \
  --organization=testorg \
  --local-templates="$GP_TEMPLATES_OK" \
  --output-dir="$GP_OUT_INTERNAL"
expect_exit 3 $? "generate-plugin with bogus source-path"

# TestPass (0): real catalog, fully renderable templates
./pvtr generate-plugin \
  --source-path="$CATALOG_PATH" \
  --service-name=testsvc \
  --organization=testorg \
  --local-templates="$GP_TEMPLATES_OK" \
  --output-dir="$GP_OUT_HAPPY"
expect_exit 0 $? "generate-plugin happy path"

# TestFail (1): one template renders, the other fails to parse
./pvtr generate-plugin \
  --source-path="$CATALOG_PATH" \
  --service-name=testsvc \
  --organization=testorg \
  --local-templates="$GP_TEMPLATES_PARTIAL" \
  --output-dir="$GP_OUT_PARTIAL"
expect_exit 1 $? "generate-plugin partial templating failure"

# The catalog file should be written even when some templates fail (regression
# guard for the Walk fix in privateer-sdk/command/generate-plugin.go).
if ! ls "$GP_OUT_PARTIAL"/data/catalogs/catalog_*.yaml >/dev/null 2>&1; then
  echo "ERROR: partial-templating run did not write the catalog file"
  STATUS=1
fi

PLUGIN_DIR="./plugins"
CONFIG_FILE="./test_config.yml"

# Ensure cleanup happens even on unexpected exits or signals
trap 'rm -rf "$PLUGIN_DIR" "$CONFIG_FILE" "$GP_WORK" evaluation_results' EXIT

# Install pvtr-github-repo-scanner plugin
./pvtr install ossf/pvtr-github-repo-scanner -b "$PLUGIN_DIR" || { echo "ERROR: Failed to install plugin"; exit 1; }

# Generate config for testing against the repo
# Tracing is disabled here to prevent GITHUB_TOKEN from appearing in logs
set +x
cat > "$CONFIG_FILE" <<EOF
loglevel: trace
write-directory: evaluation_results
write: true
output: yaml
services:
  privateer:
    plugin: ossf/pvtr-github-repo-scanner
    policy:
      catalogs:
        - osps-baseline-2026-02
      applicability:
        - Maturity Level 1
    vars:
      owner: privateerproj
      repo: privateer
      token: ${GITHUB_TOKEN}
EOF
set -x

# Run pvtr with the plugin.
# Test results themselves may pass or fail (exit 0 or 1) -- both are acceptable
# here. Anything higher (Aborted=2, InternalError=3, BadUsage=4) means the
# plugin itself tipped over and should fail the integration check.
./pvtr run -b "$PLUGIN_DIR" -c "$CONFIG_FILE"
RUN_EXIT=$?
if [ "$RUN_EXIT" -gt 1 ]; then
  echo "ERROR: pvtr run exited $RUN_EXIT (expected 0 or 1)"
  STATUS=1
fi

exit $STATUS
