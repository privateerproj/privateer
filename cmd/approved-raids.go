package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

// Approved Raids are raids that users may quickly install via the CLI
// These are to be thoroughly vetted and must be version locked here
// If a user installs a different version locally, it will not be overriden here

var wireframeCmd = &cobra.Command{
	Use:   "wireframe",
	Short: "Run the example raid. Useful for playing around with Privateer without touching any services.",
	Long:  ``, // TODO
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("wireframe called")
	},
}

func init() {
	runCmd.AddCommand(wireframeCmd)
}
