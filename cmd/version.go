package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionName = "version"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   versionName,
	Short: "Display the version details for this privateer executable.",
	Long:  `TODO - Long description`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s called", versionName)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	versionCmd.Flags().BoolP("help", "h", false, fmt.Sprintf("Give me a heading! Help for the %s command.", versionName))
}
