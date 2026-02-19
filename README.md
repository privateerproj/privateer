# Privateer

**Privateer** is a validation framework that simplifies infrastructure testing and compliance validation. Built with infrastructure engineers in mind, Privateer helps accelerate security and compliance validation of any software asset.

## Key Features

- **Community-Driven Plugins**: Plugins are crafted and maintained collaboratively by the community or privately within your organization
- **Comprehensive Resource Validation**: Validate diverse resources in a single execution, regardless of how many resources or validations you need to queue
- **Consistent Machine-Readable Output**: Standardized output simplifies automation and integration
- **Plugin Generation**: Generate plugin scaffolding from [Gemara](https://gemara.openssf.org) Layer 2 schema catalogs with the push of a button

## Quick Start

**For detailed ecosystem documentation, visit [privateerproj.com](https://privateerproj.com)**

### Step 1: Star this Repo

Click the star at the top right of this page so that you can find it easily the next time you sign in to GitHub.

### Step 2: Choose Your Installation Method

#### Option 1: Install via Homebrew

```bash
brew install privateerproj/tap/privateer
```

#### Option 2: Install via Script

```bash
/bin/bash -c "$(curl -sSL https://raw.githubusercontent.com/privateerproj/privateer/main/install.sh)"
```

#### Option 3: Download from Releases

Download the latest release from [GitHub Releases](https://github.com/privateerproj/privateer/releases).

#### Option 4: Build from Source

```bash
git clone https://github.com/privateerproj/privateer.git
cd privateer
go mod tidy
make build
```

### Step 3: Choose Your Plugins

We do not currently maintain an authoritative list of community plugins, but a good place to start would be the OpenSSF's [plugin](https://github.com/revanite-io/pvtr-github-repo) for scanning GitHub repos against the _Open Source Project Security Baseline_.

### Step 4: Install & Verify Your Plugins

Plugin installation is currently left to the user. The default location for plugin binaries is `$HOME/.privateer/bin`. You may specify a different location at runtime via `--binaries-path` if you install your plugins elsewhere.

To review the plugins you have installed, run `privateer list -a`.

## Contributing

We welcome contributions! See our [Contributing Guidelines](https://github.com/privateerproj/privateer?tab=contributing-ov-file) for details.

All contributions are covered by the [Apache 2 License](https://github.com/privateerproj/privateer?tab=Apache-2.0-1-ov-file) at the time the pull request is opened, and all community interactions are governed by our [Code of Conduct](https://github.com/privateerproj/privateer?tab=coc-ov-file).

### Local Development Prerequisites

- **Go 1.25.1 or later** - Required for building Privateer and running tests
- **Make** - For using the Makefile build targets

### Testing

Run all tests:

```bash
make test
```

Run tests with coverage:

```bash
make testcov
```

### Available Make Targets

- `make binary` - Build the binary
- `make test` - Run tests and vet checks
- `make testcov` - Run tests with coverage report
- `make tidy` - Clean up go.mod dependencies
- `make release` - Build release binaries for all platforms
- `make build` - Alias for `tidy test binary`

### Project Structure

```bash
privateer/
├── cmd/              # CLI commands (run, list, generate-plugin, etc.)
├── test/             # Test data and fixtures
├── build/            # Build scripts and CI configurations
├── main.go           # Application entry point
└── go.mod            # Go module dependencies
```

## Security

For vulnerability reporting, please reference our [Security Policy](https://github.com/privateerproj/privateer?tab=security-ov-file). For security questions, please search our closed issues and open a new issue if your question has not yet been answered.

## Helpful Links

- **[Privateer SDK](https://github.com/privateerproj/privateer-sdk)** - SDK for developing Privateer plugins
- **[Privateer Documentation](https://privateerproj.com)** - Complete documentation site
- **[Example Plugin](https://github.com/privateerproj/raid-wireframe)** - Reference implementation
