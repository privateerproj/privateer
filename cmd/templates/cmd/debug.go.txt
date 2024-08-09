package cmd

import (
	"log"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/spf13/cobra"
)

var (
	// debugCmd represents the base command when called without any subcommands
	debugCmd = &cobra.Command{
		Use:   "debug",
		Short: "Run the Raid in debug mode",
		Run: func(cmd *cobra.Command, args []string) {
			err := raidengine.Run(RaidName, Armory)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	runCmd.AddCommand(debugCmd) // This enables the debug command for use while working on your Raid
}
