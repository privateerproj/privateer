package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "sally",
		Short: "When everything is battoned down, it is time to sally forth",
		Long:  ``, // TODO
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("sally called\n")
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
}
