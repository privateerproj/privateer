package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer-sdk/shared"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run plugins that have been specified in the config.",
	Long: `
When everything is battoned down, it is time to run forth.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Trace("run called")
		if len(args) > 0 {
			logger.Error(fmt.Sprintf(
				"Unknown args: %v", args))
		} else {
			Run()
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// Run executes all plugins with handling for the command line
func Run() (err error) {

	// Setup for handling SIGTERM (Ctrl+C)
	setupCloseHandler()
	logger.Trace(fmt.Sprintf(
		"Using bin: %s", viper.GetString("binaries-path")))

	plugins := GetPlugins()
	if len(plugins) == 0 {
		logger.Error("no requested plugins were found in " + viper.GetString("binaries-path"))
		return
	}

	// Run all plugins
	for serviceName := range viper.GetStringMap("services") {
		servicePluginName := viper.GetString(fmt.Sprintf("services.%s.plugin", serviceName))
		for _, pluginPkg := range plugins {
			if pluginPkg.Name == servicePluginName {
				if !pluginPkg.Available {
					logger.Error("Requested plugin that is not installed: " + pluginPkg.Name)
					continue
				}
				client := newClient(pluginPkg.Command)
				defer closeClient(pluginPkg, client)

				// Connect via RPC
				var rpcClient hcplugin.ClientProtocol
				rpcClient, err = client.Client()
				if err != nil {
					return err
				}
				// Request the plugin
				var rawPlugin interface{}
				rawPlugin, err = rpcClient.Dispense(shared.PluginName)
				if err != nil {
					logger.Error(err.Error())
				}
				// Execute plugin
				plugin := rawPlugin.(shared.Pluginer)

				// Execute
				logger.Trace("Starting Plugin: " + pluginPkg.Name)
				response := plugin.Start()
				if response != nil {
					pluginPkg.Error = fmt.Errorf("Error running plugin for %s: %v", serviceName, response)
				} else {
					pluginPkg.Successful = true
				}
			}
		}
	}
	return
}

func closeClient(pluginPkg *PluginPkg, client *hcplugin.Client) {
	if pluginPkg.Successful {
		logger.Info(fmt.Sprintf("Plugin %s completed successfully", pluginPkg.Name))
	} else if pluginPkg.Error != nil {
		logger.Error(pluginPkg.Error.Error())
	} else {
		logger.Error(fmt.Sprintf("unexpected exit while attempting to run package: %v", pluginPkg))
	}
	client.Kill()
}

// setupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
// Ref: https://golangcode.com/handle-ctrl-c-exit-in-terminal/
func setupCloseHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Error("Execution aborted - SIGTERM")
		os.Exit(0)
	}()
}

// newClient client handles the lifecycle of a plugin application
// Plugin hosts should use one Client for each plugin executable
// (this is different from the client that manages gRPC)
func newClient(cmd *exec.Cmd) *hcplugin.Client {
	var pluginMap = map[string]hcplugin.Plugin{
		shared.PluginName: &shared.Plugin{},
	}
	var handshakeConfig = shared.GetHandshakeConfig()
	return hcplugin.NewClient(&hcplugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             cmd,
		Logger:          logger,
	})
}