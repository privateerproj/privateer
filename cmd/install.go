package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// installCmd represents the version command
var installCmd = &cobra.Command{
	Use:   "equip",
	Short: "Stock the Armory! Install raids that are not supported by default.",
	Long:  ``, // TODO
	Run: func(cmd *cobra.Command, args []string) {
		// This command will be a bit more complex,
		// as it will require some type of validation that
		// the specified project is compatible with Privateer
		fmt.Print("equip called")
	},
}

func init() {
	installCmd.PersistentFlags().BoolP("store", "s", false, "Github repo to source the raid from.")
	viper.BindPFlag("store", installCmd.PersistentFlags().Lookup("store"))
}
