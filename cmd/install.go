package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var installCmd = &cobra.Command{
	Use:   "equip",
	Short: "Stock the Armory! Install an official raid from the Privateer Project",
	Long:  `TODO - Long description (mention how to bin your own raid?)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("equip called")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
