package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version details for this privateer executable",
	Long:  ``, // TODO
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("version called")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolP("help", "h", false, "Give me a heading! Help for the version command.")
}
