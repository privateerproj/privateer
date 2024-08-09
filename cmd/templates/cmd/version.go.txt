package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version details.",
	Run: func(cmd *cobra.Command, args []string) {
		writer := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		if viper.GetBool("verbose") {
			fmt.Fprintf(writer, "Version:\t%s\n", buildVersion)
			fmt.Fprintf(writer, "Commit:\t%s\n", buildGitCommitHash)
			fmt.Fprintf(writer, "Build Time:\t%s\n", buildTime)
			writer.Flush()
		} else {
			fmt.Println(buildVersion)
		}
	},
}

func init() {
	runCmd.AddCommand(versionCmd)
}
