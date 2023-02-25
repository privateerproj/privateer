package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer-sdk/plugin"
	"github.com/privateerproj/privateer-sdk/probeengine"
	"github.com/privateerproj/privateer-sdk/utils"
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
				"Sally only accepts a single argument. Unknown args: %v", args[1:]))
		} else if len(args) == 1 {
			StartApprovedRaid(args[0])
		} else {
			logger.Trace("Sequentially executing all raids in config") // TODO
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

	cmdSet, err := getCommands()
	if err != nil {
		logger.Error(fmt.Sprintf(
			"Error loading plugins from config: %s", err))
		return
	}

	// Run all plugins
	err = AllPlugins(cmdSet)
	if err != nil {
		logger.Info(fmt.Sprintf(
			"Output directory: %s", viper.GetString("WriteDirectory")))
		switch e := err.(type) {
		case *RaidErrors:
			logger.Error(fmt.Sprintf(
				"%d out of %d raids failed.", len(e.Errors), len(cmdSet)))
			return
		default:
			logger.Error(err.Error())
		}
	}
	logger.Info(fmt.Sprintf(
		"[INFO] No errors encountered during plugin execution. Output directory: %s",
		viper.GetString("WriteDirectory")))
	return
}

// AllPlugins executes specified plugins in a loop
func AllPlugins(cmdSet []*exec.Cmd) (err error) {
	// Capture any plugin errors received during execution
	spErrors := make([]RaidError, 0)

	for _, cmd := range cmdSet {
		spErrors, err = Plugin(cmd, spErrors)
		if err != nil {
			return
		}
	}

	if len(spErrors) > 0 {
		// Return all raid errors to main
		err = &RaidErrors{
			Errors: spErrors,
		}
	}
	return
}

// Plugin executes single plugin based on the provided command
func Plugin(cmd *exec.Cmd, spErrors []RaidError) ([]RaidError, error) {
	// Launch the plugin process
	client := newClient(cmd)
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return spErrors, err
	}

	// Request the plugin
	rawSP, err := rpcClient.Dispense(plugin.RaidPluginName)
	if err != nil {
		return spErrors, err
	}

	// Execute raid, expecting a silent response
	raid := rawSP.(plugin.Raid)
	response := raid.Start()
	if response != nil {
		spErr := RaidError{
			Raid: cmd.String(), // TODO: retrieve raid name from interface function
			Err:  response,
		}
		spErrors = append(spErrors, spErr)
	} else {
		logger.Info("Probes all completed with successful results")
	}
	return spErrors, nil
}

// GetRaidBinary returns the path to the executable for the specified raid
func GetRaidBinary(name string) (binaryName string, err error) {
	name = filepath.Base(strings.ToLower(name)) // in some cases a filepath may arrive here instead of the base name
	if runtime.GOOS == "windows" && !strings.HasSuffix(name, ".exe") {
		name = fmt.Sprintf("%s.exe", name)
	}
	plugins, _ := hcplugin.Discover(name, viper.GetString("binaries-path"))
	if len(plugins) != 1 {
		err = fmt.Errorf("failed to locate requested plugin '%s' at path '%s'", name, viper.GetString("binaries-path"))
		return
	}
	binaryName = plugins[0]

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
		probeengine.CleanupTmp()
		os.Exit(0)
	}()
}

func getCommands() (cmdSet []*exec.Cmd, err error) {
	// TODO: give any exec errors a familiar format
	raids := GetRequestedRaids()
	for _, raidName := range raids {
		cmd, err := getCommand(raidName)
		if err != nil {
			break
		}
		cmdSet = append(cmdSet, cmd)
	}
	logger.Debug(fmt.Sprintf(
		"Using bin: %s", viper.GetString("binaries-path")))
	if err == nil && len(cmdSet) == 0 {
		// If there are no errors but also no commands run, it's probably unexpected
		var available []string
		GetAvailableRaids()
		err = utils.ReformatError(
			"No valid raids specified. Requested: %v, Available: %v", raids, available)
	}
	return
}

func getCommand(raid string) (cmd *exec.Cmd, err error) {
	binaryName, binErr := GetRaidBinary(raid)
	if binErr != nil {
		err = binErr
		return
	}
	cmd = exec.Command(binaryName)
	flags := fmt.Sprintf("--varsfile=%s", viper.GetString("config")) // TODO this is wonky- can't change it yet without causing an error somewhere else
	cmd.Args = append(cmd.Args, flags)
	return
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
