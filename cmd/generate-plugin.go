package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer-sdk/command"
)

var (
	// genPluginCmd represents the generate-plugin command.
	// It generates a new pvtr plugin from a source file using templates.
	genPluginCmd = &cobra.Command{
		Use:   "generate-plugin",
		Short: "Generate a new plugin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generatePlugin()
		},
		SilenceUsage: true,
	}
)

func init() {
	genPluginCmd.PersistentFlags().StringP("source-path", "p", "", "The source file to generate the plugin from")
	genPluginCmd.PersistentFlags().StringP("local-templates", "", "", "Path to a directory to use instead of downloading the latest templates")
	genPluginCmd.PersistentFlags().StringP("service-name", "n", "", "The name of the service (e.g. 'ECS, AKS, GCS')")
	genPluginCmd.PersistentFlags().StringP("organization", "g", "", "The GitHub organization for the plugin (e.g. 'privateerproj')")
	genPluginCmd.PersistentFlags().StringP("output-dir", "o", "generated-plugin/", "Pathname for the generated plugin")

	_ = viper.BindPFlag("source-path", genPluginCmd.PersistentFlags().Lookup("source-path"))
	_ = viper.BindPFlag("local-templates", genPluginCmd.PersistentFlags().Lookup("local-templates"))
	_ = viper.BindPFlag("service-name", genPluginCmd.PersistentFlags().Lookup("service-name"))
	_ = viper.BindPFlag("organization", genPluginCmd.PersistentFlags().Lookup("organization"))
	_ = viper.BindPFlag("output-dir", genPluginCmd.PersistentFlags().Lookup("output-dir"))

	rootCmd.AddCommand(genPluginCmd)
}

// generatePlugin validates the plugin configuration and generates a new plugin
// from templates. It returns any errors encountered to the caller.
func generatePlugin() error {
	cfg, err := command.SetupTemplatingEnvironment(logger)
	if err != nil {
		return err
	}

	return command.GeneratePlugin(logger, cfg)
}
