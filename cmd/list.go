package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer-sdk/command"
)

var listCmd *cobra.Command

func init() {
	// Create a wrapper command that gets the writer at runtime
	// This is necessary because writer is nil at package init time
	// and only gets initialized in persistentPreRun()
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Consult the Charts! List all plugins that have been installed.",
		Run: func(cmd *cobra.Command, args []string) {
			// Get the list command from SDK with the initialized writer
			// At this point, writer has been initialized by persistentPreRun()
			sdkListCmd := command.GetListCmd(writer)
			// Execute the SDK command's Run function
			sdkListCmd.Run(cmd, args)
		},
	}

	// Add flags (matching the SDK command)
	listCmd.PersistentFlags().BoolP("all", "a", false, "Review the Fleet! List all plugins that have been installed or requested in the current config")
	_ = viper.BindPFlag("all", listCmd.PersistentFlags().Lookup("all"))

	rootCmd.AddCommand(listCmd)
}

// Re-export functions from SDK for backwards compatibility
var GetPlugins = command.GetPlugins
var Contains = command.Contains
