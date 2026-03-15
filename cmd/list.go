package cmd

import (
	"github.com/spf13/cobra"

	"github.com/privateerproj/privateer-sdk/command"
)

// listCmd represents the list command, which displays all installed plugins.
var listCmd *cobra.Command

func init() {
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Consult the Charts! List all plugins that have been installed.",
		Run: func(cmd *cobra.Command, args []string) {
			sdkListCmd := command.GetListCmd(writer)
			sdkListCmd.Run(cmd, args)
		},
	}
	command.SetListCmdFlags(listCmd)
	rootCmd.AddCommand(listCmd)
}
