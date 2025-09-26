# Privateer

## Simplifying Validation for Infrastructure Engineers

Privateer has been meticulously crafted with infrastructure engineers in mind. If you're seeking to validate your resources against regulations, taxonomies, or standards, Privateer is your trusted companion. With a user-friendly interface and powerful features, you can now effortlessly navigate the complexities of resource validation.

## Using Privateer Plugins

There are several key benefits to Privateer Plugins:

- **Community-Driven Plugins:** Our open development model ensures that Plugins are crafted and maintained collaboratively by the community, reflecting a wealth of expertise and insights.
- **Comprehensive Resource Validation:** Privateer empowers you to validate a diverse array of resources in a single execution. No more piecemeal validation processes; instead, experience efficiency and thoroughness in one go.
- **Consistent Machine-Readable Output:** Regardless of the specific Plugin, you're guaranteed a standardized, machine-readable test output. This consistency simplifies the automation and integration of test results, enabling seamless decision-making.
- **Empowering Service Providers:** Privateer finds its calling in projects like Compliant Financial Infrastructure and Common Cloud Controls within FINOS. Service providers can leverage Privateer Plugins developed by FINOS to certify resources for use in regulated industries, such as insurance and banking.

## Install the Privateer CLI

### Option 1: Install via Script

Run the following command to install Privateer using the provided install.sh script:

```sh
/bin/bash -c "$(curl -sSL https://raw.githubusercontent.com/privateerproj/privateer/03ced90caae9f3c9203eb7f82f2c46ccf2ff15fc/install.sh)"
```

### Option 2: Download from Releases

Download the latest release from [GitHub Releases](https://github.com/privateerproj/privateer/releases).

### Build Privateer from Source

To build privateer from source, follow these steps below: 

1. Clone the Repository

    ```sh
    git clone https://github.com/privateerproj/privateer.git
    cd privateer
    ```

2. Installing Dependencies

    ```sh
    go mod tidy
    ```

3. Building Privateer

    ```sh
    make release
    ```

## Install Privateer Plugins

Plugins are built and maintained by the community. Choose the plugin(s) that you wish to run, and install them to your binaries path.

- **Default Path:** $HOME/.privateer/bin
- **Customize via CLI:** Use `--binaries-path` in your CLI command to change the path to your binaries.
- **Customize via config:** Specify a custom binaries path in your config via the top level value `binaries-path: your/bin/path`

## Configuration

1. **Create a Configuration File**: Craft a configuration file (e.g., `config.yml`) that specifies the plugins you intend to run and any necessary configuration options. Include secrets and settings required by the plugin. Refer to the specific plugin's documentation for precise details.
1. **Output Directory (Optional)**: If desired, define an output directory in your configuration. Privateer will generate log and result files for each plugin in this directory. Results files are available in both JSON and YAML formats.
1. **Advanced Config Management**: Privateer's roadmap includes plans for integrating with systems like etcd and Consul to enhance configuration and secret management.

> [!NOTE] 
> If your configuration file is stored in a non-default location, specify its file path using the -c or --config flag.

### Example Config.yml

```yaml
loglevel: Debug
write-directory: sample_output
services:
  my-cloud-service1:
    plugin: example
    test-suite:
      - tlp_red
      # - tlp_amber
      # - tlp_green
      # - tlp_clea
```

## Common Commands

Here are some common commands you can use with Privateer:

- `help` / `-h` / `--help`: Display help information about Privateer and its commands.
- `run`: Execute the specified plugin(s).
- `generate-plugin`: Generate a new plugin based on a FINOS Common Cloud Controls catalog.
- `list`: Show plugins requested by your configuration and whether they're installed.
  - `list -a`: Show all plugins that have been installed or requested in the current config.
- `version`: Display version details.
- `completion`: Generate autocompletion scripts for bash, fish, powershell, or zsh.

## Command Line Options

### Global Flags

All commands support these global flags:

- `-b, --binaries-path`: Path to the directory where plugins are installed (default: `$HOME/.privateer/bin`)
- `-c, --config`: Configuration file, JSON or YAML (default: `config.yml`)
- `-h, --help`: Display help information
- `-l, --loglevel`: Log level - trace, debug, info, warn, error, off (default: "error")
- `-s, --service`: Named service to execute from the config
- `--silent`: Only show essential log information
- `-t, --test-suites`: Named set of test sets to execute from the plugin (default: "default")
- `--write`: Keep detailed result outputs in files (default: true)
- `-w, --write-directory`: Directory to write evaluation results to (default: "evaluation_results")

### Command-Specific Options

#### generate-plugin
- `-p, --source-path`: The source file to generate the plugin from
- `-n, --service-name`: The name of the service (e.g. 'ECS', 'AKS', 'GCS')
- `-o, --output-dir`: Output directory for the generated plugin (default: "generated-plugin/")
- `--local-templates`: Path to local templates instead of downloading latest

#### list
- `-a, --all`: Show all plugins that have been installed or requested in the current config

## Output Customization

Privateer generates logs and results files for each plugin. The output location is controlled by the `-w, --write-directory` flag.

- **Log Results**: `<write-directory>/<plugin_name>/<plugin_name>.log`
- **Plugin Results**: `<write-directory>/<plugin_name>/results.yaml`
- **Default Directory**: `evaluation_results`

### Log Levels

Control logging verbosity with the `-l, --loglevel` flag:
- **trace**: Most verbose, shows all debug information
- **debug**: Debug information
- **info**: General information
- **warn**: Warning messages
- **error**: Error messages only (default)
- **off**: No logging

Use `--silent` to show only essential log information regardless of log level.
