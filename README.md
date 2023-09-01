[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/privateerproj/privateer/badge)](https://securityscorecards.dev/viewer/?uri=github.com/privateerproj/privateer)

# Privateer

*Simplify. Validate. Elevate.*

## Introduction

Welcome to Privateer, your all-inclusive test harness designed to streamline the validation process for a wide range of resources through its innovative approach to harmonizing inputs, execution, and outputs. Aptly dubbed "Raids," these test packs serve as your go-to solution for ensuring security hardening, regulatory compliance, and adherence to taxonomy standards across your infrastructure.

### Simplifying Validation for Infrastructure Engineers

Privateer has been meticulously crafted with infrastructure engineers in mind. If you're seeking to validate your resources against regulations, taxonomies, or standards, Privateer is your trusted companion. With a user-friendly interface and powerful features, you can now effortlessly navigate the complexities of resource validation.

### Unlocking the Power of Raids

There are several key benefits to Privateer Raids:

- **Community-Driven Raids:** Our open development model ensures that Raids are crafted and maintained collaboratively by the community, reflecting a wealth of expertise and insights.
- **Comprehensive Resource Validation:** Privateer empowers you to validate a diverse array of resources in a single execution. No more piecemeal validation processes; instead, experience efficiency and thoroughness in one go.
- **Consistent Machine-Readable Output:** Regardless of the specific Raid, you're guaranteed a standardized, machine-readable test output. This consistency simplifies the automation and integration of test results, enabling seamless decision-making.
- **Empowering Service Providers:** Privateer finds its calling in projects like Compliant Financial Infrastructure and Common Cloud Controls within FINOS. Service providers can leverage Privateer Raids developed by FINOS to certify resources for use in regulated industries, such as the insurance and banking.

### Your Journey Begins Here

As we incubate this project, you have the opportunity to be part of its growth and transformation. Privateer is primed to become an essential component of automated solutions, fitting seamlessly into CI/CD pipelines and API tooling.

**Take Action: Embrace and Build**: The time to act is now. Whether you're looking to deploy existing Raids or embark on crafting new ones, Privateer stands ready to elevate your resource validation journey.

Ready to embark on your validation journey? Start using or contributing to Privateer Raids today!

For more information and support, visit Privateer on GitHub or connect with our community on [Slack](https://finos-lf.slack.com/messages/cfi).


## Usage & Quickstart Guide

Privateer empowers you to ensure the security, compliance, and integrity of your resources with ease. Here's how to dive in and make the most of this versatile tool:

### Installation

1. **Download Privateer**: Obtain the latest release of Privateer from the [GitHub repository](https://github.com/privateerproj/privateer/releases).
1. **Install Raids**: Choose the raid(s) you wish to use from the same release on GitHub. Install them to your preferred `binaries-path`. By default, this is `$HOME/privateer/bin`, but you can customize it in your configuration.

### Configuration

1. **Create a Configuration File**: Craft a configuration file (e.g., `config.yml`) that specifies the raids you intend to run and any necessary configuration options. Include secrets and settings required by the raid. Refer to the specific raid's documentation for precise details.
1. **Output Directory (Optional)**: If desired, define an output directory in your configuration. Privateer will generate log and result files for each raid in this directory. Results files are available in both JSON and YAML formats.
1. **Advanced Config Management**: Privateer's roadmap includes plans for integrating with systems like etcd and Consul to enhance configuration and secret management.

#### Example Config.yml

```yaml
loglevel: trace
WriteDirectory: test_output
Raids:
  Wireframe:
    JokeName: Jimmy
```

### Interacting with Raids

1. **List Requested Raids**: Use `privateer list` to view a list of raids you have requested in your configuration. This provides an overview of the tests you've planned.
1. **View Installed Raids**: Explore the raids available for installation using `privateer list --available`. This allows you to discover new tests you might find valuable.
1. **Run Raids**: Execute your selected raids with the `privateer sally` command. By default, Privateer will execute all raids specified in your configuration. You can also override this behavior by including the specific raid you want to run as a third value in the command, e.g., `privateer sally wireframe`.

### Customizing Output and Logs

1. **Output Logs**: Privateer generates logs and results files for each raid. Logs are saved as `<output_dir>/<raid_name>/<raid_name>.log`, and results are saved as `<output_dir>/<raid_name>/results.yaml`. Specify your output directory using the config value `WriteDirectory`. The default value is `$HOME/privateer/output`

### Tailoring Verbosity

1. **Log Verbosity**: Increase the verbosity of logs using the `-v` or `--verbose` flag. This is particularly useful for gaining deeper insights into the execution process.
1. **Silence Logs**: For a streamlined experience, silence non-essential log information by utilizing the `-s` or `--silent` flag.

### Configuration Flexibility

1. **Alternate Config Path**: If your configuration file is stored in a different location, specify its filepath using the `-c` or `--config` flag.

Now you're equipped with the essential knowledge to unleash the power of Privateer. Whether you're certifying resources or maintaining configuration integrity, Privateer stands ready to help. Embrace the journey of resource validation with confidence!

For more detailed information, refer to the [Privateer GitHub Repository](https://github.com/privateerproj/privateer) and explore the possibilities that await you.
