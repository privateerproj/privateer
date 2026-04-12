#!/bin/sh

set -x

STATUS=0

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

PLUGIN_DIR="./plugins"
CONFIG_FILE="./test_config.yml"

# Ensure cleanup happens even on unexpected exits or signals
trap 'rm -rf "$PLUGIN_DIR" "$CONFIG_FILE" evaluation_results' EXIT

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
    plugin: pvtr-github-repo-scanner
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

# Run pvtr with the plugin
./pvtr run -b "$PLUGIN_DIR" -c "$CONFIG_FILE" || STATUS=1

exit $STATUS
