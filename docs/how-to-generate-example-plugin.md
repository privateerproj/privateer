# How to Generate Example Plugin

## What is this?

This guide allows you to not have to start from scratch on a new plugin.

> [!NOTE]
> Currently the generator only supports common cloud controls yaml files. This file is needed to be used to setup the plugin's test suite and test sets.  All that is needed to be done after that is writing the actual tests.

- get [`CCC.VPC_2025.01.yaml`](https://github.com/finos/common-cloud-controls/releases/download/v2025.01.VPC/CCC.VPC_2025.01.yaml) file from [Common Cloud Controls Repository Releases page](https://github.com/finos/common-cloud-controls/releases)
> [!NOTE]
> Version may change or you may need to expand the `Assets` section to find the latest yaml file.

- in the root of this repository run the following command:

    ```bash
    privateer generate-plugin -p ~/path/to/CCC.VPC_2025.01.yaml -n example
    ```

- this will generate an example plugin in the `generated-plugin` folder at the root of the repository
- Go into the newly generated directory:

    ```bash
    cd generated_plugin
    ```

- run the following command to build the plugin:

    ```bash
    cp config-example.yml config.yml
    make binary
    ```

- to run the plugin by itself in debug mode:

    ```bash
    ./example debug --service my-cloud-service1
    ```

> [!TIP]
> If you use a different service name, make sure the service name matches what is in the config.yml in the root of the repository.

> [!IMPORTANT]
> `test_output/[service_name]` folder should include a log file and a yaml file for each test suite
>
> example: `test_output/my-cloud-service1/my-cloud-service1.log` and `test_output/my-cloud-service1/tlp_red.yml`

- to run the plugin from privateer, do the following:

    ```bash
    cp example $HOME/.privateer/bin
    cp config.yml ../
    cd ..
    privateer run
    ```
