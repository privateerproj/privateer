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
	Short: "Consult the Charts! List all raids that have been installed",
	Long:  ``, // TODO
	Run: func(cmd *cobra.Command, args []string) {
		var raids map[string]string
		if viper.GetBool("available") {
			raids = getAvailableAndRequestedRaids()
		} else {
			raids = getRequestedAndAvailableRaids()
		}
		if len(raids) < 1 {
			if viper.GetBool("available") {
				fmt.Fprintln(writer, "No raids present in the binaries path:", viper.GetString("binaries-path"))
			} else {
				fmt.Fprintln(writer, "No raids requested in the current configuration.")
			}
		} else {
			fmt.Fprintln(writer, "| Available \t| Requested \t|")
			for available, requested := range raids {
				fmt.Fprintf(writer, "| %s\t| %s \t|\n", available, requested)
			}
		}
		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().BoolP("available", "a", false, "Inventory the Armory! List all raids that have been installed.")
	viper.BindPFlag("available", listCmd.PersistentFlags().Lookup("available"))
}

// List all available raids and whether or not they're requested in the active config
// Currently lists raids including the file extension. TODO: should we change that?
func getAvailableAndRequestedRaids() map[string]string {
	raids := make(map[string]string)
	requestedRaids := GetRequestedRaids()
	availableRaids := GetAvailableRaids()
	for _, availableRaid := range availableRaids {
		raids[availableRaid] = "Not Requested"
		for _, requestedRaid := range requestedRaids {
			if requestedRaid == availableRaid {
				raids[availableRaid] = "Requested"
			}
		}
	}
	return raids
}

// List only raids requested in the active config and whether or not they're available
func getRequestedAndAvailableRaids() map[string]string {
	raids := make(map[string]string)
	requestedRaids := GetRequestedRaids()
	availableRaids := GetAvailableRaids()
	for _, requestedRaid := range requestedRaids {
		raids[requestedRaid] = "Not Found"
		for _, availableRaid := range availableRaids {
			if requestedRaid == availableRaid {
				raids[requestedRaid] = "Yes"
			}
		}
	}
	return raids
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
