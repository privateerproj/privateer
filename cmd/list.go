package cmd

import (
	"os"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer/run"
)

var (
	cmdName = "list"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   cmdName,
	Short: "Consult the Charts! List all raids that have been installed",
	Long:  `TODO - Long description`,
	Run: func(cmd *cobra.Command, args []string) {
		writer := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

		if viper.GetBool("available") {
			fmt.Fprintln(writer, "| Available \t| Requested \t|")
			raids := getAvailableAndRequestedRaids()
			for available, requested := range raids {
				fmt.Fprintf(writer, "| %s\t| %s \t|\n", available, requested)
			}
			writer.Flush()	
		} else {
			raids := getRequestedAndAvailableRaids()
			fmt.Fprintln(writer, "| Requested\t| Available \t|")
			for requested, available := range raids {
				fmt.Fprintf(writer, "| %s\t| %s \t|\n", requested, available)
			}
			writer.Flush()	
		}
	},
}

// List all available raids and whether or not they're requested in the active config
func getAvailableAndRequestedRaids() map[string]string {
	raids := make(map[string]string)
	requestedRaids := run.GetRequestedRaids()
	availableRaids := run.GetAvailableRaids()
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
	requestedRaids := run.GetRequestedRaids()
	availableRaids := run.GetAvailableRaids()
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

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().BoolP("available", "a", false, "Inventory the Armory! List all raids that have been installed.")
	viper.BindPFlag("available", listCmd.PersistentFlags().Lookup("available"))
}
