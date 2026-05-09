package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/privateerproj/privateer-sdk/command"
)

// runFn is the function used by runCmd to execute plugins. It's a package-level
// variable so tests can swap it out to avoid actually running plugins.
var runFn = func(c *CLI) int {
	c.setupCloseHandler()
	return command.Run(c.logger, command.GetPlugins)
}

func (c *CLI) addRunCmd() {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run plugins that have been specified in the config.",
		Long: `
When everything is battoned down, it is time to run forth.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			c.logger.Trace("run called")
			exitCode := runFn(c)
			if exitCode != 0 {
				os.Exit(exitCode)
			}
		},
	}
	c.rootCmd.AddCommand(runCmd)
}
