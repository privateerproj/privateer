package cmd

import (
	"log"
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
			log.Print("???")
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listCmd.PersistentFlags().BoolP("available", "a", false, "Consult our options! List all raids that have been installed.")
	viper.BindPFlag("available", listCmd.PersistentFlags().Lookup("available"))
}
