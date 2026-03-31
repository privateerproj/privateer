package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/privateerproj/privateer-sdk/command"
)

func (c *CLI) addRunCmd() {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run plugins that have been specified in the config.",
		Long: `
When everything is battoned down, it is time to run forth.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			c.logger.Trace("run called")
			exitCode := c.run()
			os.Exit(int(exitCode))
		},
	}
	c.rootCmd.AddCommand(runCmd)
}

// run executes all plugins with handling for the command line.
func (c *CLI) run() (exitCode int) {
	c.setupCloseHandler()
	return command.Run(c.logger, command.GetPlugins)
}

// setupCloseHandler creates a signal listener on a new goroutine which will notify
// the program if it receives an interrupt from the OS (SIGINT or SIGTERM).
func (c *CLI) setupCloseHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		c.logger.Error("Test execution was aborted by user")
		os.Exit(int(command.Aborted))
	}()
}
