/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runName = "sally"

// runCmd represents the sally command
var runCmd = &cobra.Command{
	Use:   runName,
	Short: "When everything is battoned down, it is time to sally forth.",
	Long:  `TODO - Long description`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s called", runName)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	runCmd.Flags().BoolP("help", "h", false, fmt.Sprintf("Give me a heading! Help for the %s command.", runName))
	runCmd.Flags().BoolP("verbose", "v", false, "Louder now! Increase log verbosity to maximum.")
	runCmd.Flags().BoolP("quiet", "q", false, "Quiet! Only show essential log information.")
}
