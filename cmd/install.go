package cmd

import (
	"github.com/privateerproj/privateer-sdk/command"
)

func (c *CLI) addInstallCmd() {
	c.rootCmd.AddCommand(command.GetInstallCmd(func() command.Writer { return c.writer }))
}
