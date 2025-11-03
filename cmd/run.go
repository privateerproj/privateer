package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/privateerproj/privateer-sdk/command"
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

func init() {
	rootCmd.AddCommand(runCmd)
}

// Run executes all plugins with handling for the command line
func Run() (exitCode int) {
	// Setup for handling SIGTERM (Ctrl+C)
	setupCloseHandler()
	
	return command.Run(logger, GetPlugins)
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
		os.Exit(int(command.Aborted))
	}()
}
