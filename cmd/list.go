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
func GetRequestedRaids() (raids []string) {
	if viper.Get("Raids") != nil {
		raidsVars := viper.Get("Raids").(map[string]interface{})
		for raidName := range raidsVars {
			raids = append(raids, raidName)
		}
	}
	return
}

// GetAvailableRaids returns a list of raids found in the binaries path
func GetAvailableRaids() (raids []string) {
	raidPaths, _ := hcplugin.Discover("*", viper.GetString("binaries-path"))
	for _, raidPath := range raidPaths {
		raids = append(raids, path.Base(raidPath))
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
		raids[availableRaid] = RaidStatus{Available: true, Requested: false}
		for _, requestedRaid := range requestedRaids {
			if requestedRaid == availableRaid {
				raids[availableRaid] = RaidStatus{Available: true, Requested: true}
			}
		}
	}
	// loop through requested raids to make sure all requested raids are represented
	for _, requestedRaid := range requestedRaids {
		if _, ok := raids[requestedRaid]; !ok {
			raids[requestedRaid] = RaidStatus{Available: false, Requested: true}
		}
	}
	return raids
}
