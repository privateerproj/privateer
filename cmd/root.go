package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/privateerproj/privateer-sdk/command"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "privateer",
	Short: "A brief description of your application",
	Long: `
Privateer CLI Quickstart
------------------------

This interface enables the quick execution of Privateer Raids,
with a shared input and output if multiple are executed.
Read more about the vision for Raids in our official documentation:
https://github.com/privateerproj/privateer/README.md

Several Privateer commands use unconventional terms
to encourage users to act carefully when using this CLI.
This is due to the fact that your Privateer config is likely
to contain secrets that can be destructive if misused.

The "sally" command will start all requested raids.
Raids are intended to directly interact with running services
and only should be used with caution and proper planning.
Never use a custom-built raid from an unknown source.

You may also streamline the creation of
a new Raid using the generate-raid command, or
the creation of Strikes for a Raid using generate-strike.
Review the help documentation for each command to learn more.

------------------------`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	command.SetBase(rootCmd)
}
