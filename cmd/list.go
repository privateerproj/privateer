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
		raids := SortAvailableAndRequested()
		if len(raids) < 1 {
			if viper.GetBool("available") {
				fmt.Fprintln(writer, "No raids present in the binaries path:", viper.GetString("binaries-path"))
			} else {
				fmt.Fprintln(writer, "No raids requested in the current configuration.")
			}
		} else {
			if viper.GetBool("available") {
				// list only the available raids
				fmt.Fprintln(writer, "| Raid \t | Available \t|")
				for raidName, raidStatus := range raids {
					if raidStatus.Available {
						fmt.Fprintf(writer, "| %s \t | %t \t|\n", raidName, raidStatus.Available)
					}
				}
			} else {
				// print all raids requested and available
				fmt.Fprintln(writer, "| Raid \t | Available \t| Requested \t|")
				for raidName, raidStatus := range raids {
					fmt.Fprintf(writer, "| %s \t | %t \t| %t \t|\n", raidName, raidStatus.Available, raidStatus.Requested)
				}
			}
		}
		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().BoolP("available", "a", false, "Review the fleet! List all raids that have been installed.")
	viper.BindPFlag("available", listCmd.PersistentFlags().Lookup("available"))
}

// GetRequestedRaids returns a list of raid names requested in the config
func GetRequestedRaids() (raidNames []string) {
	services := viper.GetStringMap("services")
	for serviceName := range services {
		raidName := viper.GetString("services." + serviceName + ".raid")
		if raidName != "" && !contains(raidNames, raidName) {
			raidNames = append(raidNames, raidName)
		}
	}
	return
}

func contains(slice []string, search string) bool {
	for _, found := range slice {
		if found == search {
			return true
		}
	}
	return false
}

// GetAvailableRaids returns a list of raids found in the binaries path
func GetAvailableRaids() (raidNames []string) {
	raidPaths, _ := hcplugin.Discover("*", viper.GetString("binaries-path"))
	for _, raidPath := range raidPaths {
		raidNames = append(raidNames, path.Base(raidPath))
	}
	return
}

type RaidStatus struct {
	Available bool
	Requested bool
}

func SortAvailableAndRequested() map[string]RaidStatus {
	raids := make(map[string]RaidStatus)
	requestedRaids := GetRequestedRaids()
	availableRaids := GetAvailableRaids()

	// loop through available raids, then requested raids to make sure available raids have requested status
	for _, availableRaid := range availableRaids {
		raids[availableRaid] = RaidStatus{Available: true}
	}
	for _, requestedRaid := range requestedRaids {
		if contains(availableRaids, requestedRaid) {
			raids[requestedRaid] = RaidStatus{Available: true, Requested: true}
		} else {
			raids[requestedRaid] = RaidStatus{Available: false, Requested: true}
		}
	}
	return raids
}
