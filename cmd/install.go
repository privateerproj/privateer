package cmd

import (
	"github.com/privateerproj/privateer-sdk/command/harness"
)

func (c *CLI) addInstallCmd() {
	c.rootCmd.AddCommand(harness.GetInstallCmd(func() harness.Writer { return c.writer }))
}
