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
			exitCode := Run()
			os.Exit(int(exitCode))
		}
	},
}

const (
	TestPass = iota
	TestFail
	Aborted
	InternalError
	BadUsage
	NoTests
)

func init() {
	rootCmd.AddCommand(runCmd)
}

// Run executes all plugins with handling for the command line
func Run() (exitCode int) {

	// Setup for handling SIGTERM (Ctrl+C)
	setupCloseHandler()
	logger.Trace(fmt.Sprintf(
		"Using bin: %s", viper.GetString("binaries-path")))

	plugins := GetPlugins()
	if len(plugins) == 0 {
		logger.Error(fmt.Sprintf("no plugins were requested in config: %s", viper.GetString("binaries-path")))
		return NoTests
	}

	// Run all plugins
	var runCount int
	for serviceName := range viper.GetStringMap("services") {
		servicePluginName := viper.GetString(fmt.Sprintf("services.%s.plugin", serviceName))
		for _, pluginPkg := range plugins {
			if pluginPkg.Name == servicePluginName {
				if !pluginPkg.Available {
					logger.Error(fmt.Sprintf("requested plugin that is not installed: %s", pluginPkg.Name))
					return BadUsage
				}
				runCount++
				client := newClient(pluginPkg.Command)
				// Connect via RPC
				var rpcClient hcplugin.ClientProtocol
				rpcClient, err := client.Client()
				if err != nil {
					logger.Error(fmt.Sprintf("internal error while initializing %s RPC client: %s", serviceName, err))
					closeClient(pluginPkg, serviceName, client)
					return InternalError
				}
				// Request the plugin
				var rawPlugin interface{}
				rawPlugin, err = rpcClient.Dispense(shared.PluginName)
				if err != nil {
					logger.Error(fmt.Sprintf("internal error while dispensing RPC client: %s", err.Error()))
					closeClient(pluginPkg, serviceName, client)
					return InternalError
				}
				// Execute plugin
				plugin := rawPlugin.(shared.Pluginer)
				logger.Trace(fmt.Sprintf("Starting Plugin %v: %s", runCount, pluginPkg.Name))
				response := plugin.Start()
				if response != nil {
					pluginPkg.Error = fmt.Errorf("tests failed in plugin %s: %v", serviceName, response)
					exitCode = TestFail
				} else {
					pluginPkg.Successful = true
				}
				closeClient(pluginPkg, serviceName, client)
			}
		}
	}
	return exitCode
}

func closeClient(pluginPkg *PluginPkg, serviceName string, client *hcplugin.Client) {
	// Close the client: this doesn't work via defer because it leaves the plugin running while the next begins
	if pluginPkg.Successful {
		logger.Info(fmt.Sprintf("Plugin for %s completed successfully", serviceName))
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
		logger.Error("Test execution was aborted by user")
		os.Exit(int(Aborted))
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
