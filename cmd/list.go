package cmd

import (
	"fmt"
	"path"
	"strings"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cmdName = "list"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   cmdName,
	Short: "Consult the Charts! List all plugins that have been installed.",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("all") {
			_, _ = fmt.Fprintln(writer, "| Plugin \t | Available \t| Requested \t|")
			for _, pluginPkg := range GetPlugins() {
				_, _ = fmt.Fprintf(writer, "| %s \t | %t \t| %t \t|\n", pluginPkg.Name, pluginPkg.Available, pluginPkg.Requested)
			}
		} else {
			// list only the available plugins
			_, _ = fmt.Fprintln(writer, "| Plugin \t | Requested \t|")
			for _, pluginPkg := range GetPlugins() {
				if pluginPkg.Available {
					_, _ = fmt.Fprintf(writer, "| %s \t | %t \t|\n", pluginPkg.Name, pluginPkg.Requested)
				}
			}
		}
		_ = writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().BoolP("all", "a", false, "Review the Fleet! List all plugins that have been installed or requested in the current config.")
	_ = viper.BindPFlag("all", listCmd.PersistentFlags().Lookup("all"))
}

// GetRequestedPlugins returns a list of plugin names requested in the config
func getRequestedPlugins() (requestedPluginPackages []*PluginPkg) {
	services := viper.GetStringMap("services")
	for serviceName := range services {
		pluginName := viper.GetString("services." + serviceName + ".plugin")
		pluginPkg := NewPluginPkg(pluginName, serviceName)
		pluginPkg.Requested = true
		requestedPluginPackages = append(requestedPluginPackages, pluginPkg)
	}
	return requestedPluginPackages
}

// GetAvailablePlugins returns a list of plugins found in the binaries path
func getAvailablePlugins() (availablePluginPackages []*PluginPkg) {
	pluginPaths, _ := hcplugin.Discover("*", viper.GetString("binaries-path"))
	for _, pluginPath := range pluginPaths {
		pluginPkg := NewPluginPkg(path.Base(pluginPath), "")
		pluginPkg.Available = true
		if strings.Contains(pluginPkg.Name, "privateer") {
			continue
		}
		availablePluginPackages = append(availablePluginPackages, pluginPkg)
	}
	return availablePluginPackages
}

var allPlugins []*PluginPkg

func GetPlugins() []*PluginPkg {
	if allPlugins != nil {
		return allPlugins
	}
	output := make([]*PluginPkg, 0)
	for _, plugin := range getRequestedPlugins() {
		if Contains(getAvailablePlugins(), plugin.Name) {
			plugin.Available = true
		}
		output = append(output, plugin)
	}
	for _, plugin := range getAvailablePlugins() {
		if !Contains(output, plugin.Name) {
			output = append(output, plugin)
		}
	}
	return output
}

func Contains(slice []*PluginPkg, search string) bool {
	for _, plugin := range slice {
		if plugin.Name == search {
			return true
		}
	}
	return false
}
