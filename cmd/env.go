package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (c *CLI) addEnvCmd() {
	envCmd := &cobra.Command{
		Use:     "env",
		Aliases: []string{"info"},
		Short:   "Display runtime environment details.",
		Long:    `Display the binary path, config location, plugins directory, installed plugins, version, and build information.`,
		Run: func(cmd *cobra.Command, args []string) {
			binaryPath, err := os.Executable()
			if err != nil {
				binaryPath = "unknown"
			}

			configFile := viper.ConfigFileUsed()
			var configStatus string
			if configFile == "" {
				configStatus = "none"
			} else if _, err := os.Stat(configFile); err == nil {
				configStatus = fmt.Sprintf("%s (found)", configFile)
			} else {
				configStatus = fmt.Sprintf("%s (not found)", configFile)
			}

			pluginsDir := viper.GetString("binaries-path")
			pluginNames := discoverPluginNames(c.logger, pluginsDir)

			_, _ = fmt.Fprintf(c.writer, "Binary:\t%s\n", binaryPath)
			_, _ = fmt.Fprintf(c.writer, "Config:\t%s\n", configStatus)
			_, _ = fmt.Fprintf(c.writer, "Plugins Dir:\t%s\n", pluginsDir)
			_, _ = fmt.Fprintf(c.writer, "Plugins:\t%s\n", pluginNames)
			_, _ = fmt.Fprintf(c.writer, "Version:\t%s\n", c.buildVersion)
			_, _ = fmt.Fprintf(c.writer, "Commit:\t%s\n", c.buildGitCommitHash)
			_, _ = fmt.Fprintf(c.writer, "Build Time:\t%s\n", c.buildTime)
			_, _ = fmt.Fprintf(c.writer, "Go Version:\t%s\n", runtime.Version())
			_, _ = fmt.Fprintf(c.writer, "OS/Arch:\t%s/%s\n", runtime.GOOS, runtime.GOARCH)
			if err := c.writer.Flush(); err != nil {
				log.Printf("Error flushing writer: %v", err)
			}
		},
	}
	c.rootCmd.AddCommand(envCmd)
}

func discoverPluginNames(logger interface{ Error(string, ...interface{}) }, pluginsDir string) string {
	pluginPaths, err := hcplugin.Discover("*", pluginsDir)
	if err != nil {
		logger.Error("error discovering plugins", "dir", pluginsDir, "error", err)
	}
	var names []string
	for _, p := range pluginPaths {
		name := filepath.Base(p)
		nameWithoutExt := strings.TrimSuffix(name, filepath.Ext(name))
		if nameWithoutExt == "pvtr" || nameWithoutExt == "privateer" {
			continue
		}
		names = append(names, name)
	}
	if len(names) == 0 {
		return "none"
	}
	return strings.Join(names, ", ")
}
