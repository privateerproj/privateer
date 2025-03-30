package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version details.",
	Long:  ``, // TODO
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("verbose") {
			_, _ = fmt.Fprintf(writer, "Version:\t%s\n", buildVersion)
			_, _ = fmt.Fprintf(writer, "Commit:\t%s\n", buildGitCommitHash)
			_, _ = fmt.Fprintf(writer, "Build Time:\t%s\n", buildTime)
			_ = writer.Flush()
		} else {
			fmt.Println(buildVersion)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
