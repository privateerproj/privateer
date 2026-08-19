package cmd

import (
	"github.com/privateerproj/privateer-sdk/command/harness"
)

func (c *CLI) addLoginCmds() {
	c.rootCmd.AddCommand(harness.GetLoginCmd(func() harness.Writer { return c.writer }))
	c.rootCmd.AddCommand(harness.GetLogoutCmd(func() harness.Writer { return c.writer }))
}
