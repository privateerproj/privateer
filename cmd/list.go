package cmd

import (
	"fmt"
	"log"
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
		if len(raids) == 0 {
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

var requestedRaidPackages []*RaidPkg

// GetRequestedRaids returns a list of raid names requested in the config
func GetRequestedRaids() []*RaidPkg {
	if len(requestedRaidPackages) > 0 {
		return requestedRaidPackages
	}
	services := viper.GetStringMap("services")
	for serviceName := range services {
		raidName := viper.GetString("services." + serviceName + ".raid")
		if raidName != "" && !contains(requestedRaidPackages, raidName) {
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

func contains(slice []*RaidPkg, search string) bool {
	for _, raid := range slice {
		if raid.Name == search {
			return true
		}
	}
	return false
}

func SortAvailableAndRequested() map[string]*RaidPkg {
	output := make(map[string]*RaidPkg)
	requestedRaids := GetRequestedRaids()
	log.Printf("requestedRaids: %v", requestedRaids)
	availableRaids := GetAvailableRaids()
	log.Printf("availableRaids: %v", availableRaids)

	// loop through available raids, then requested raids to make sure available raids have requested status
	for _, raid := range availableRaids {
		output[raid.Name] = raid
		for _, requested := range requestedRaids {
			if raid.Name == requested.Name {
				raid.Requested = true
			}
		}
	}
	return output
}
