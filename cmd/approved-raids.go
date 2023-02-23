package cmd

import (
	"fmt"
	"os"
	// "github.com/spf13/viper"
)

// Approved Raids are raids that users may quickly install via the CLI
// These are to be thoroughly vetted and must be version locked here
// If a user installs a different version locally, it will not be overriden here

// var wireframeCmd = &cobra.Command{
// 	Use:   "wireframe",
// 	Short: "Run the example raid. Useful for playing around with Privateer without touching any services.",
// 	Long:  ``, // TODO
// 	Run: func(cmd *cobra.Command, args []string) {
// 		runApprovedRaid("wireframe")
// 	},
// }

// func init() {
// 	runCmd.AddCommand(wireframeCmd)
// }

func runApprovedRaid(raidName string) {
	err := installIfNotPResent(raidName)
	if err != nil {
		fmt.Print("Error with installer logic.") // these are all just temporary logs
		os.Exit(1)
	}
	err = executeRaid(raidName)
	if err != nil {
		fmt.Print("Error with execute logic.")
		os.Exit(1)
	}
}

func installIfNotPResent(raidName string) error {
	fmt.Printf("install called for %s", raidName)
	return nil
}

func executeRaid(raidName string) error {
	fmt.Printf("sally %s called", raidName)
	return nil
}
