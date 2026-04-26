package cmd

import (
	"os"

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
