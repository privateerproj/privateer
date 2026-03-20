# Contributing to pvtr

Thank you for your interest in contributing to pvtr! This document covers everything that you will need to set up and run the project locally, understand the project layout, run  the tests, and submit changes.

All contributions are covered by the [Apache 2 License](https://github.com/privateerproj/privateer?tab=Apache-2.0-1-ov-file) at the time the pull request is opened. All community interactions are governed by our [Code of Conduct](https://github.com/privateerproj/privateer?tab=coc-ov-file).

---

## Table of Contents

- [Local Development](#local-development)
  - [Prerequisites](#prerequisites)
  - [Clone and Set Up](#clone-and-set-up)
- [Make Targets](#make-targets)
  - [Core Development](#core-development)
  - [Release Builds](#release-builds)
  - [Utilities](#utilities)
- [Running Tests](#running-tests)
- [Project Structure](#project-structure)
- [Submitting Changes](#submitting-changes)
- [Privateer Ecosystem](#privateer-ecosystem)

---

## Local Development

### Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.25.1 or later | Building pvtr and running tests |
| [Make](https://www.gnu.org/software/make/) | Any recent version | Running Makefile targets |
| [Git](https://git-scm.com/) | Any recent version | Version control |

### Clone and Set Up

```bash
git clone https://github.com/privateerproj/privateer.git
cd privateer
go mod tidy
```

Build the binary:

```bash
make build
```

This runs `tidy → test → binary` in sequence and produces a `pvtr` binary in the project root.

---

## Make Targets

### Core Development

| Target | Description |
|---|---|
| `make build` | Full build: runs `tidy`, then `test`, then compiles the binary. Alias for `tidy test binary`. |
| `make binary` | Compile the `pvtr` binary only (no tests, no tidy). |
| `make test` | Run `go vet` checks and all unit tests (`go test ./...`). |
| `make testcov` | Run tests, vet checks, and generate a coverage report with total coverage percentage. |
| `make tidy` | Run `go mod tidy` to clean up module dependencies. |

### Release Builds

| Target | Description |
|---|---|
| `make release` | Full release: runs `tidy`, `test`, then builds binaries for all platforms (Linux, Windows, macOS). |
| `make release-candidate` | Build RC binaries for all platforms (tagged with `-rc` version postfix). |
| `make release-nix` | Build Linux binary only. |
| `make release-win` | Build Windows binary only. |
| `make release-mac` | Build macOS (Darwin) binary only. |

> **Note:** Production releases are handled by [GoReleaser](https://goreleaser.com/) via the `.goreleaser.yaml` configuration and the GitHub Actions release workflow. The `make release-*` targets are intended for local verification before tagging.

### Utilities

| Target | Description |
|---|---|
| `make todo` | Interactively append a to-do item to `TODO.md`. Prompts for input. |

---

## Running Tests

Run all tests and vet checks:

```bash
make test
```

Run tests with a coverage report (outputs total coverage percentage to stdout):

```bash
make testcov
```

---

## Project Structure

```
privateer/
├── .github/
│   └── workflows/      # CI configurations
├── cmd/                # CLI commands (run, list, generate-plugin, env, version)
├── test/               # Test data and fixtures
├── main.go             # Application entry point
├── main_test.go        # Top-level tests
├── go.mod              # Go module definition
├── go.sum              # Dependency checksums
├── Makefile            # Build and development targets
├── .goreleaser.yaml    # Multi-platform release configuration
├── install.sh          # Installation script (used by install-via-script option)
└── Dockerfile          # Container image definition

```     

## Submitting Changes

### Step 1: Fork and clone the repo

If you haven't already, click **Fork** at the top right of the GitHub page to create your own copy, then clone it locally:

```bash
git clone https://github.com/<your-username>/privateer.git
cd privateer
```

### Step 2: Create a branch for your changes

Never work directly on `main`. Create a branch with a short name that describes what you're doing:

```bash
git checkout -b my-feature-branch
```

> **Tip:** Good branch names are short and descriptive, e.g. `fix-list-command`, `add-env-docs`, `update-readme`.

### Step 3: Make your changes, then check what changed

Once you have made your edits, run this to see which files were modified:

```bash
git status
```

You will see a list of changed files in red (not yet staged) and green (staged and ready to commit).

### Step 4: Stage the files you want to include

Add the specific files you changed. Replace the filenames with your actual changed files:

```bash
git add your-file
```

To stage everything at once (use carefully, double-check with `git status` first):

```bash
git add .
```

### Step 5 — Commit with a sign-off

The `-s` flag automatically appends a `Signed-off-by` line to your commit, which is required to certify that your contribution is made under the project license. The `-m` flag lets you write your commit message inline:

```bash
git commit -s -m "description of the change"
```

> **Example:** `git commit -s -m "fix: correct install script in README"`

Keep your message short, describing *what* the commit does.

### Step 6 — Push your branch to GitHub

```bash
git push origin my-feature-branch
```

### Step 7 — Open a Pull Request

Go to your fork on GitHub. You will see a prompt to open a pull request for your recently pushed branch, click it, fill in a clear title and description, and submit.

> **Note:** PR titles are validated by CI. Keep them short and descriptive (the same style as your commit message). Merges require approval from a maintainer per, .

### Step 8 — Verify your build locally before pushing

Before opening a PR, it's good to run:

```bash
make build    # compiles the binary and runs all tests
make testcov  # checks your test coverage hasn't dropped
```

If either of these fails, fix the issues before pushing, CI will catch the same problems.

---

For significant features or breaking changes, please open an issue first to discuss the approach before writing code.

---

## Privateer Ecosystem

| Project | Description |
|---|---|
| [privateer](https://github.com/privateerproj/privateer) | Core CLI (this repo) |
| [privateer-sdk](https://github.com/privateerproj/privateer-sdk) | SDK for developing pvtr plugins |


**[Browse all projects →](https://github.com/privateerproj)**

For SDK usage, CLI command reference, and the developers guide, visit **[privateerproj.com/docs/developers/](https://privateerproj.com/docs/developers/)**.