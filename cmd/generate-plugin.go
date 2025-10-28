package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer-sdk/command"
)

var (
	// genPluginCmd represents the generate-plugin command
	genPluginCmd = &cobra.Command{
		Use:   "generate-plugin",
		Short: "Generate a new plugin",
		Run: func(cmd *cobra.Command, args []string) {
			generatePlugin()
		},
	}
)

func init() {
	genPluginCmd.PersistentFlags().StringP("source-path", "p", "", "The source file to generate the plugin from")
	genPluginCmd.PersistentFlags().StringP("local-templates", "", "", "Path to a directory to use instead of downloading the latest templates")
	genPluginCmd.PersistentFlags().StringP("service-name", "n", "", "The name of the service (e.g. 'ECS, AKS, GCS')")
	genPluginCmd.PersistentFlags().StringP("output-dir", "o", "generated-plugin/", "Pathname for the generated plugin")

	_ = viper.BindPFlag("source-path", genPluginCmd.PersistentFlags().Lookup("source-path"))
	_ = viper.BindPFlag("local-templates", genPluginCmd.PersistentFlags().Lookup("local-templates"))
	_ = viper.BindPFlag("service-name", genPluginCmd.PersistentFlags().Lookup("service-name"))
	_ = viper.BindPFlag("output-dir", genPluginCmd.PersistentFlags().Lookup("output-dir"))

	rootCmd.AddCommand(genPluginCmd)
}

func generatePlugin() {
	templatesDir, sourcePath, outputDir, serviceName, err := command.SetupTemplatingEnvironment(logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = command.GeneratePlugin(logger, templatesDir, sourcePath, outputDir, serviceName)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
