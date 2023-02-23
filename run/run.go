package run

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer-sdk/logging"
	"github.com/privateerproj/privateer-sdk/plugin"
	"github.com/privateerproj/privateer-sdk/probeengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// CLIContext executes all plugins with handling for the command line
func CLIContext() {
	// Setup for handling SIGTERM (Ctrl+C)
	setupCloseHandler()

	cmdSet, err := getCommands()
	if err != nil {
		log.Printf("Error loading plugins from config: %s", err)
		os.Exit(2)
	}

	// Run all plugins
	if err := AllPlugins(cmdSet); err != nil {
		log.Printf("[INFO] Output directory: %s", viper.GetString("WriteDirectory"))
		switch e := err.(type) {
		case *RaidErrors:
			log.Printf("[ERROR] %d out of %d raids failed.", len(e.Errors), len(cmdSet))
			os.Exit(1) // At least one raid failed
		default:
			log.Print(utils.ReformatError(err.Error()))
			os.Exit(2) // Internal error
		}
	}
	log.Printf("[INFO] No errors encountered during plugin execution. Output directory: %s", viper.GetString("WriteDirectory"))
	os.Exit(0)
}

// AllPlugins executes specified plugins in a loop
func AllPlugins(cmdSet []*exec.Cmd) (err error) {
	spErrors := make([]RaidError, 0) // This will store any plugin errors received during execution

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
		log.Printf("[INFO] Probes all completed with successful results") // TODO: use hclogger in this file
	}
	return spErrors, nil
}

// GetRaidBinary finds provided raid in installation folder and return binary name
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
		log.Print(<-c)
		log.Printf("Execution aborted - %v", "SIGTERM")
		probeengine.CleanupTmp()
		os.Exit(0)
	}()
}

// GetRequestedRaids returns a list of raid names requested in the config
func GetRequestedRaids() (raids []string) {
	if viper.Get("Raids") != nil {
		raidsVars := viper.Get("Raids").(map[string]interface{})
		for raidName := range raidsVars {
			raids = append(raids, raidName)
		}
	}
	return
}

// GetAvailableRaides returns a list of raids found in the binaries path
func GetAvailableRaids() (raids []string) {
	raidPaths, _ := hcplugin.Discover("*", viper.GetString("binaries-path"))
	for _, raidPath := range raidPaths {
		raids = append(raids, path.Base(raidPath))
	}
	return
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
	log.Printf("[INFO] Using bin: %s", viper.GetString("binaries-path"))
	if err == nil && len(cmdSet) == 0 {
		// If there are no errors but also no commands run, it's probably unexpected
		var available []string
		GetAvailableRaids()
		err = utils.ReformatError("No valid raids specified. Requested: %v, Available: %v", raids, available)
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
	cmd.Args = append(cmd.Args, fmt.Sprintf("--varsfile=%s", viper.GetString("config")))
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
		Logger:          logging.GetLogger("core"),
	})
}
