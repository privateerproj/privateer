package cmd

import (
	"fmt"
	"path"

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
	Short: "Consult the Blueprints! List all raids that have been installed.",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("all") {
			fmt.Fprintln(writer, "| Raid \t | Available \t| Requested \t|")
			for _, raidPkg := range GetRaids() {
				fmt.Fprintf(writer, "| %s \t | %t \t| %t \t|\n", raidPkg.Name, raidPkg.Available, raidPkg.Requested)
			}
		} else {
			// list only the available raids
			fmt.Fprintln(writer, "| Raid \t | Requested \t|")
			for _, raidPkg := range GetRaids() {
				if raidPkg.Available {
					fmt.Fprintf(writer, "| %s \t | %t \t|\n", raidPkg.Name, raidPkg.Requested)
				}
			}
		}
		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().BoolP("all", "a", false, "Review the Fleet! List all raids that have been installed or requested in the current config.")
	viper.BindPFlag("all", listCmd.PersistentFlags().Lookup("all"))
}

// GetRequestedRaids returns a list of raid names requested in the config
func getRequestedRaids() (requestedRaidPackages []*RaidPkg) {
	services := viper.GetStringMap("services")
	for serviceName := range services {
		raidName := viper.GetString("services." + serviceName + ".raid")
		if raidName != "" && !Contains(requestedRaidPackages, raidName) {
			raidPkg := NewRaidPkg(raidName, serviceName)
			raidPkg.Requested = true
			requestedRaidPackages = append(requestedRaidPackages, raidPkg)
		}
	}
	return requestedRaidPackages
}

// GetAvailableRaids returns a list of raids found in the binaries path
func getAvailableRaids() (availableRaidPackages []*RaidPkg) {
	raidPaths, _ := hcplugin.Discover("*", viper.GetString("binaries-path"))
	for _, raidPath := range raidPaths {
		raidPkg := NewRaidPkg(path.Base(raidPath), "")
		raidPkg.Available = true
		availableRaidPackages = append(availableRaidPackages, raidPkg)
	}
	return availableRaidPackages
}

var allRaids []*RaidPkg

func GetRaids() []*RaidPkg {
	if allRaids != nil {
		return allRaids
	}
	output := make([]*RaidPkg, 0)
	for _, raid := range getRequestedRaids() {
		if Contains(getAvailableRaids(), raid.Name) {
			raid.Available = true
		}
		output = append(output, raid)
	}
	for _, raid := range getAvailableRaids() {
		if !Contains(output, raid.Name) {
			output = append(output, raid)
		}
	}
	return output
}

func Contains(slice []*RaidPkg, search string) bool {
	for _, raid := range slice {
		if raid.Name == search {
			return true
		}
	}
	return false
}
