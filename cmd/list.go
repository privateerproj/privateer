/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdName = "list"

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   cmdName,
	Short: "Consult the Carts! List all raids that have been installed",
	Long:  `TODO - Long description`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s called", cmdName)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listCmd.Flags().BoolP("help", "h", false, fmt.Sprintf("Give me a heading! Help for the %s command.", cmdName))
	listCmd.Flags().BoolP("quiet", "q", false, "Quiet! Only show the raids that are planned by the config file.")
}
