package cmd

import (
	"github.com/spf13/cobra"

	"github.com/privateerproj/privateer-sdk/command"
)

func (c *CLI) addInstallCmd() {
	installCmd := &cobra.Command{
		Use:   "install [plugin-name]",
		Short: "Install a vetted plugin from the registry.",
		Long:  "Resolve the plugin name to registry metadata, then download the plugin binary from the release URL into the binaries path.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sdkInstall := command.GetInstallCmd(c.writer)
			sdkInstall.SetArgs(args)
			return sdkInstall.Execute()
		},
	}
	c.rootCmd.AddCommand(installCmd)
}
