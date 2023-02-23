package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var genRaidName = "generate-raid"

// versionCmd represents the version command
var genRaidCmd = &cobra.Command{
	Use:   genRaidName,
	Short: "Generate a new raid",
	Long:  ``, // TODO
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s called", genRaidName)
	},
}

func init() {
	rootCmd.AddCommand(genRaidCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
