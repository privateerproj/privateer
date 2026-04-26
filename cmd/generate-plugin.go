package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer-sdk/command"
)

func (c *CLI) addGenPluginCmd() {
	genPluginCmd := &cobra.Command{
		Use:   "generate-plugin",
		Short: "Generate a new plugin",
		Run: func(cmd *cobra.Command, args []string) {
			c.logger.Trace("generate-plugin called")
			c.setupCloseHandler()
			exitCode := command.GeneratePlugin(c.logger)
			os.Exit(int(exitCode))
		},
	}

	genPluginCmd.Flags().StringP("source-path", "p", "", "The source file to generate the plugin from")
	genPluginCmd.Flags().StringP("local-templates", "", "", "Path to a directory to use instead of downloading the latest templates")
	genPluginCmd.Flags().StringP("service-name", "n", "", "The name of the service (e.g. 'ECS, AKS, GCS')")
	genPluginCmd.Flags().StringP("organization", "g", "", "The GitHub organization for the plugin (e.g. 'privateerproj')")
	genPluginCmd.Flags().StringP("output-dir", "o", "generated-plugin/", "Pathname for the generated plugin")

	_ = viper.BindPFlag("source-path", genPluginCmd.Flags().Lookup("source-path"))
	_ = viper.BindPFlag("local-templates", genPluginCmd.Flags().Lookup("local-templates"))
	_ = viper.BindPFlag("service-name", genPluginCmd.Flags().Lookup("service-name"))
	_ = viper.BindPFlag("organization", genPluginCmd.Flags().Lookup("organization"))
	_ = viper.BindPFlag("output-dir", genPluginCmd.Flags().Lookup("output-dir"))

	c.rootCmd.AddCommand(genPluginCmd)
}
