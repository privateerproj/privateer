package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer-sdk/command"
)

func (c *CLI) addListCmd() {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Consult the Charts! List all plugins that have been installed.",
		Run: func(cmd *cobra.Command, args []string) {
			sdkListCmd := command.GetListCmd(c.writer)
			sdkListCmd.Run(cmd, args)
		},
	}

	listCmd.PersistentFlags().BoolP("all", "a", false, "Review the Fleet! List all plugins that have been installed or requested in the current config")
	_ = viper.BindPFlag("all", listCmd.PersistentFlags().Lookup("all"))

	listCmd.PersistentFlags().Bool("installed", false, "List only plugins that are installed locally")
	_ = viper.BindPFlag("installed", listCmd.PersistentFlags().Lookup("installed"))

	listCmd.PersistentFlags().Bool("installable", false, "List vetted plugins from the registry that are available to install")
	_ = viper.BindPFlag("installable", listCmd.PersistentFlags().Lookup("installable"))

	c.rootCmd.AddCommand(listCmd)
}
