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
			fmt.Fprintln(writer, "| Raid \t | Available \t| Requested \t|")
			for _, pluginPkg := range GetRaids() {
				fmt.Fprintf(writer, "| %s \t | %t \t| %t \t|\n", pluginPkg.Name, pluginPkg.Available, pluginPkg.Requested)
			}
		} else {
			// list only the available plugins
			fmt.Fprintln(writer, "| Raid \t | Requested \t|")
			for _, pluginPkg := range GetRaids() {
				if pluginPkg.Available {
					fmt.Fprintf(writer, "| %s \t | %t \t|\n", pluginPkg.Name, pluginPkg.Requested)
				}
			}
		}
		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().BoolP("all", "a", false, "Review the Fleet! List all plugins that have been installed or requested in the current config.")
	viper.BindPFlag("all", listCmd.PersistentFlags().Lookup("all"))
}

// GetRequestedRaids returns a list of plugin names requested in the config
func getRequestedRaids() (requestedRaidPackages []*RaidPkg) {
	services := viper.GetStringMap("services")
	for serviceName := range services {
		pluginName := viper.GetString("services." + serviceName + ".plugin")
		if pluginName != "" && !Contains(requestedRaidPackages, pluginName) {
			pluginPkg := NewRaidPkg(pluginName, serviceName)
			pluginPkg.Requested = true
			requestedRaidPackages = append(requestedRaidPackages, pluginPkg)
		}
	}
	return requestedRaidPackages
}

// GetAvailableRaids returns a list of plugins found in the binaries path
func getAvailableRaids() (availableRaidPackages []*RaidPkg) {
	pluginPaths, _ := hcplugin.Discover("*", viper.GetString("binaries-path"))
	for _, pluginPath := range pluginPaths {
		pluginPkg := NewRaidPkg(path.Base(pluginPath), "")
		pluginPkg.Available = true
		if strings.Contains(pluginPkg.Name,  "privateer"){
			continue
		}
		availableRaidPackages = append(availableRaidPackages, pluginPkg)
	}
	return availableRaidPackages
}

var allRaids []*RaidPkg

func GetRaids() []*RaidPkg {
	if allRaids != nil {
		return allRaids
	}
	output := make([]*RaidPkg, 0)
	for _, plugin := range getRequestedRaids() {
		if Contains(getAvailableRaids(), plugin.Name) {
			plugin.Available = true
		}
		output = append(output, plugin)
	}
	for _, plugin := range getAvailableRaids() {
		if !Contains(output, plugin.Name) {
			output = append(output, plugin)
		}
	}
	return output
}

func Contains(slice []*RaidPkg, search string) bool {
	for _, plugin := range slice {
		if plugin.Name == search {
			return true
		}
	}
	return false
}
