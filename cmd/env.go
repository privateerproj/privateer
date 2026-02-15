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

// envCmd represents the env command, which displays runtime environment details.
var envCmd = &cobra.Command{
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
		pluginNames := discoverPluginNames(pluginsDir)

		_, _ = fmt.Fprintf(writer, "Binary:\t%s\n", binaryPath)
		_, _ = fmt.Fprintf(writer, "Config:\t%s\n", configStatus)
		_, _ = fmt.Fprintf(writer, "Plugins Dir:\t%s\n", pluginsDir)
		_, _ = fmt.Fprintf(writer, "Plugins:\t%s\n", pluginNames)
		_, _ = fmt.Fprintf(writer, "Version:\t%s\n", buildVersion)
		_, _ = fmt.Fprintf(writer, "Commit:\t%s\n", buildGitCommitHash)
		_, _ = fmt.Fprintf(writer, "Build Time:\t%s\n", buildTime)
		_, _ = fmt.Fprintf(writer, "Go Version:\t%s\n", runtime.Version())
		_, _ = fmt.Fprintf(writer, "OS/Arch:\t%s/%s\n", runtime.GOOS, runtime.GOARCH)
		if err := writer.Flush(); err != nil {
			log.Printf("Error flushing writer: %v", err)
		}
	},
}

func discoverPluginNames(pluginsDir string) string {
	pluginPaths, _ := hcplugin.Discover("*", pluginsDir)
	var names []string
	for _, p := range pluginPaths {
		name := filepath.Base(p)
		if strings.Contains(name, "privateer") {
			continue
		}
		names = append(names, name)
	}
	if len(names) == 0 {
		return "none"
	}
	return strings.Join(names, ", ")
}

func init() {
	rootCmd.AddCommand(envCmd)
}
