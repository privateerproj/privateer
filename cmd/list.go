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
	Short: "Consult the Charts! List all raids that have been requested",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("all") {
			fmt.Fprintln(writer, "| Raid \t | Available \t| Requested \t|")
			for _, raidPkg := range GetRaidsAvailableOrRequested() {
				fmt.Fprintf(writer, "| %s \t | %t \t| %t \t|\n", raidPkg.Name, raidPkg.Available, raidPkg.Requested)
			}
		} else {
			// list only the available raids
			fmt.Fprintln(writer, "| Raid \t | Requested \t|")
			for _, raidPkg := range GetAvailableRaids() {
				fmt.Fprintf(writer, "| %s \t | %t \t|\n", raidPkg.Name, raidPkg.Requested)
			}
		}
		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().BoolP("all", "a", false, "Review the fleet! List all raids that have been installed or requested.")
	viper.BindPFlag("all", listCmd.PersistentFlags().Lookup("all"))
}

var requestedRaidPackages []*RaidPkg

// GetRequestedRaids returns a list of raid names requested in the config
func GetRequestedRaids() []*RaidPkg {
	if len(requestedRaidPackages) > 0 {
		return requestedRaidPackages
	}
	services := viper.GetStringMap("services")
	for serviceName := range services {
		raidName := viper.GetString("services." + serviceName + ".raid")
		if raidName != "" && !Contains(requestedRaidPackages, raidName) {
			raidPkg := NewRaidPkg(raidName, serviceName)
			requestedRaidPackages = append(requestedRaidPackages, raidPkg)
		}
	}
	return requestedRaidPackages
}

var availableRaidPackages []*RaidPkg

// GetAvailableRaids returns a list of raids found in the binaries path
func GetAvailableRaids() []*RaidPkg {
	if len(availableRaidPackages) != 0 {
		return availableRaidPackages
	}
	raidPaths, _ := hcplugin.Discover("*", viper.GetString("binaries-path"))
	for _, raidPath := range raidPaths {
		raidPkg := NewRaidPkg(path.Base(raidPath), "")
		raidPkg.Available = true
		availableRaidPackages = append(availableRaidPackages, raidPkg)
	}
	return availableRaidPackages
}

func Contains(slice []*RaidPkg, search string) bool {
	for _, raid := range slice {
		if raid.Name == search {
			return true
		}
	}
	return false
}

func GetRaidsAvailableOrRequested() []*RaidPkg {
	// Combine the available and requested raids
	// Mark the values for Requested and Available accordingly.

	availableRaidPackages := GetAvailableRaids()
	requestedRaidPackages := GetRequestedRaids()
	output := make([]*RaidPkg, 0)

	for _, raid := range availableRaidPackages {
		raid.Available = true
		output = append(output, raid)
	}

	for _, raid := range requestedRaidPackages {
		if !Contains(output, raid.Name) {
			raid.Requested = true
			output = append(output, raid)
		} else {
			for _, r := range output {
				if r.Name == raid.Name {
					r.Requested = true
				}
			}
		}
	}
	return output
}
