package cmd

import (
	"github.com/privateerproj/privateer-sdk/command/harness"
)

func (c *CLI) addPublishCmd() {
	c.rootCmd.AddCommand(harness.GetPublishCmd(func() harness.Writer { return c.writer }))
}
