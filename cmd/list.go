package cmd

import (
	"github.com/privateerproj/privateer-sdk/command"
)

func (c *CLI) addListCmd() {
	listCmd := command.GetListCmd(func() command.Writer { return c.writer })
	c.rootCmd.AddCommand(listCmd)
}
