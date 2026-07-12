package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/privateerproj/privateer-sdk/command"
	"github.com/privateerproj/privateer-sdk/command/harness"
)

// runFn is the function used by runCmd to execute plugins. It's a package-level
// variable so tests can swap it out to avoid actually running plugins.
var runFn = func(c *CLI) int {
	c.setupCloseHandler()
	return harness.Run(context.Background(), c.writer, c.logger, harness.GetPlugins)
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
	// Execution-specific flags (output, write-directory, etc.) live on the run
	// command rather than the shared root, so their shorthands (e.g. -o) do not
	// collide with sibling subcommands like generate-plugin's --output-dir.
	command.SetRunFlags(runCmd)
	c.rootCmd.AddCommand(runCmd)
}
