package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func (c *CLI) addVersionCmd() {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display version details.",
		Long:  `Display the version, git commit hash, and build timestamp of this pvtr build. Use the --verbose flag to see all details.`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				log.Printf("Error reading verbose flag: %v", err)
			}
			if verbose {
				_, _ = fmt.Fprintf(c.writer, "Version:\t%s\n", c.buildVersion)
				_, _ = fmt.Fprintf(c.writer, "Commit:\t%s\n", c.buildGitCommitHash)
				_, _ = fmt.Fprintf(c.writer, "Build Time:\t%s\n", c.buildTime)
				if err := c.writer.Flush(); err != nil {
					log.Printf("Error flushing writer: %v", err)
				}
			} else {
				_, _ = fmt.Fprintf(c.writer, "%s\n", c.buildVersion)
				if err := c.writer.Flush(); err != nil {
					log.Printf("Error flushing writer: %v", err)
				}
			}
		},
	}

	versionCmd.Flags().BoolP("verbose", "v", false, "Display full version details including commit and build time")

	c.rootCmd.AddCommand(versionCmd)
}
