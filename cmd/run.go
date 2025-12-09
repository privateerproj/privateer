package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/privateerproj/privateer-sdk/command"
)

// runCmd represents the run command, which executes all plugins specified in the configuration.
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

func init() {
	rootCmd.AddCommand(runCmd)
}

// Run executes all plugins with handling for the command line.
// It sets up signal handlers for graceful shutdown and returns the exit code
// from the plugin execution.
//
// Returns:
//   - exitCode: The exit code from the plugin execution (0 for success, non-zero for failure)
func Run() (exitCode int) {
	// Setup for handling SIGTERM (Ctrl+C)
	setupCloseHandler()
	
	return command.Run(logger, GetPlugins)
}

// setupCloseHandler creates a signal listener on a new goroutine which will notify
// the program if it receives an interrupt from the OS (SIGINT or SIGTERM).
// When an interrupt is received, it logs an error message and exits with the
// Aborted exit code.
//
// Reference: https://golangcode.com/handle-ctrl-c-exit-in-terminal/
func setupCloseHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Error("Test execution was aborted by user")
		os.Exit(int(command.Aborted))
	}()
}
