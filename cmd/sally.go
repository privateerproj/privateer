package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"

	"github.com/privateerproj/privateer-sdk/plugin"
)

// runCmd represents the sally command
var runCmd = &cobra.Command{
	Use:   "sally",
	Short: "Run raids that have been specified within the command or configuration.",
	Long: `
When everything is battoned down, it is time to sally forth.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Trace("sally called")
		if len(args) > 1 {
			logger.Error(fmt.Sprintf(
				"Sally only accepts a single argument; all other elements should be flags. Unknown args: %v", args[1:]))
		} else if len(args) == 1 {
			StartApprovedRaid(args[0])
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
	// logger.Debug(fmt.Sprintf(
	// 	"Using bin: %s", viper.GetString("binaries-path")))

	// var errString string
	// raids := GetRequestedRaids()
	// for _, raidName := range raids {
	// 	raid, err := NewRaidPkg(raidName, serviceName)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	raidPackages = append(raidPackages, raid)
	// }

	// if err == nil && len(raidPackages) == 0 {
	// 	// If there are no errors but also no commands run, it's probably unexpected
	// 	available := GetAvailableRaids()
	// 	errString = fmt.Sprintf("No valid raids specified. Requested: %v, Available: %v", raids, available)
	// }

	// if errString != "" {
	// 	logger.Error(fmt.Sprintf(
	// 		"Error loading plugins from config: %s", err))
	// 	return
	// }

	// logger.Trace(fmt.Sprintf("cmdSet: %s", cmdSet))

	// // Run all plugins

	// raidErrors := make([]RaidError, 0)

	// for _, cmd := range cmdSet {

	// 	client := newClient(cmd)
	// 	defer client.Kill()

	// 	// Connect via RPC
	// 	var rpcClient hcplugin.ClientProtocol
	// 	rpcClient, err = client.Client()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Request the plugin
	// 	var rawRaid interface{}
	// 	rawRaid, err = rpcClient.Dispense(plugin.RaidPluginName)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Execute raid, expecting a silent response
	// 	raid := rawRaid.(plugin.Raid)
	// 	response := raid.Start()
	// 	if response != nil {
	// 		raidErr := RaidError{
	// 			Raid: cmd.String(), // TODO: retrieve raid name from interface function
	// 			Err:  response,
	// 		}
	// 		raidErrors = append(raidErrors, raidErr)
	// 		logger.Error(fmt.Sprintf("%v", raidErrors))
	// 	} else {
	// 		logger.Info(fmt.Sprintf("Victory! Raid %s completed with successful results.", plugin.RaidPluginName))
	// 	}

	// 	if err != nil {
	// 		return
	// 	}
	// }

	// if len(raidErrors) > 0 {
	// 	// Return all raid errors to main
	// 	err = &RaidErrors{
	// 		Errors: raidErrors,
	// 	}
	// }

	// if err != nil {
	// 	// TODO: Log the config values to a file, with sensitive values redacted
	// 	switch e := err.(type) {
	// 	case *RaidErrors:
	// 		logger.Error(fmt.Sprintf(
	// 			"%d out of %d raids failed.", len(e.Errors), len(cmdSet)))
	// 		return
	// 	default:
	// 		logger.Error(err.Error())
	// 	}
	// }

	// logger.Trace(fmt.Sprintf(
	// 	"No errors encountered during plugin execution. Output directory: %s",
	// 	viper.GetString("WriteDirectory")))
	return
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
