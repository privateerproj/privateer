package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/privateerproj/privateer/run"
)

var runName = "sally"

// runCmd represents the sally command
var runCmd = &cobra.Command{
	Use:   runName,
	Short: "When everything is battoned down, it is time to sally forth.",
	Long:  `TODO - Long description`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Trace("%s called", runName)
		if len(args) > 1 {
			logger.Error(fmt.Sprintf(
				"Sally only accepts a single argument. Unknown args: %v",args[1:]))
		} else if len(args) == 1 {
			run.StartApprovedRaid(args[0])
		} else {
			logger.Trace("Sequentially executing all raids in config") // TODO
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}