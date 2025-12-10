package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// versionCmd represents the version command, which displays version information
// about the Privateer build including version, commit hash, and build time.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version details.",
	Long:  `Display the version, git commit hash, and build timestamp of this Privateer build. Use the --verbose flag to see all details.`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("verbose") {
			_, _ = fmt.Fprintf(writer, "Version:\t%s\n", buildVersion)
			_, _ = fmt.Fprintf(writer, "Commit:\t%s\n", buildGitCommitHash)
			_, _ = fmt.Fprintf(writer, "Build Time:\t%s\n", buildTime)
			err := writer.Flush()
			if err != nil {
				log.Printf("Error flushing writer: %v", err)
			}
		} else {
			fmt.Println(buildVersion)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
