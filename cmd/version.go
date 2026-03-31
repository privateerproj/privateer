package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (c *CLI) addVersionCmd() {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display version details.",
		Long:  `Display the version, git commit hash, and build timestamp of this pvtr build. Use the --verbose flag to see all details.`,
		Run: func(cmd *cobra.Command, args []string) {
			if viper.GetBool("verbose") {
				_, _ = fmt.Fprintf(c.writer, "Version:\t%s\n", c.buildVersion)
				_, _ = fmt.Fprintf(c.writer, "Commit:\t%s\n", c.buildGitCommitHash)
				_, _ = fmt.Fprintf(c.writer, "Build Time:\t%s\n", c.buildTime)
				err := c.writer.Flush()
				if err != nil {
					log.Printf("Error flushing writer: %v", err)
				}
			} else {
				fmt.Println(c.buildVersion)
			}
		},
	}

	versionCmd.Flags().BoolP("verbose", "v", false, "Display full version details including commit and build time")
	_ = viper.BindPFlag("verbose", versionCmd.Flags().Lookup("verbose"))

	c.rootCmd.AddCommand(versionCmd)
}
