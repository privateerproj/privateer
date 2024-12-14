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

	"github.com/privateerproj/privateer-sdk/plugin"
)

// runCmd represents the sally command
var runCmd = &cobra.Command{
	Use:   "sally",
	Short: "Run raids that have been specified in the config.",
	Long: `
When everything is battoned down, it is time to sally forth.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Trace("sally called")
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

	raids := GetRaids()
	if len(raids) == 0 {
		logger.Error("no requested raids were found in " + viper.GetString("binaries-path"))
		return
	}

	// Run all plugins
	for serviceName := range viper.GetStringMap("services") {
		serviceRaidName := viper.GetString(fmt.Sprintf("services.%s.raid", serviceName))
		for _, raidPkg := range raids {
			if raidPkg.Name == serviceRaidName {
				if !raidPkg.Available {
					logger.Error("Requested raid that is not installed: " + raidPkg.Name)
					continue
				}
				client := newClient(raidPkg.Command)
				defer closeClient(raidPkg, client)

				// Connect via RPC
				var rpcClient hcplugin.ClientProtocol
				rpcClient, err = client.Client()
				if err != nil {
					return err
				}
				// Request the plugin
				var rawRaid interface{}
				rawRaid, err = rpcClient.Dispense(plugin.RaidPluginName)
				if err != nil {
					logger.Error(err.Error())
				}
				// Execute raid
				raid := rawRaid.(plugin.Raid)

				// Execute
				logger.Trace("Starting Raid: " + raidPkg.Name)
				response := raid.Start()
				if response != nil {
					raidPkg.Error = fmt.Errorf("Error running raid for %s: %v", serviceName, response)
				} else {
					raidPkg.Successful = true
				}
			}
		}
	}
	return
}

func closeClient(raidPkg *RaidPkg, client *hcplugin.Client) {
	if raidPkg.Successful {
		logger.Info(fmt.Sprintf("Raid %s completed successfully", raidPkg.Name))
	} else if raidPkg.Error != nil {
		logger.Error(raidPkg.Error.Error())
	} else {
		logger.Error(fmt.Sprintf("unexpected exit while attempting to run package: %v", raidPkg))
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
		plugin.RaidPluginName: &plugin.RaidPlugin{},
	}
	var handshakeConfig = plugin.GetHandshakeConfig()
	return hcplugin.NewClient(&hcplugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             cmd,
		Logger:          logger,
	})
}
