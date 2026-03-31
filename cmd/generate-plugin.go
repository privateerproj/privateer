package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer-sdk/command"
)

func (c *CLI) addGenPluginCmd() {
	genPluginCmd := &cobra.Command{
		Use:   "generate-plugin",
		Short: "Generate a new plugin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.generatePlugin()
		},
		SilenceUsage: true,
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

func (c *CLI) generatePlugin() error {
	cfg, err := command.SetupTemplatingEnvironment(c.logger)
	if err != nil {
		return err
	}
	return command.GeneratePlugin(c.logger, cfg)
}
