package cmd

import (
	"github.com/privateerproj/privateer-sdk/command/harness"
)

func (c *CLI) addListCmd() {
	listCmd := harness.GetListCmd(func() harness.Writer { return c.writer })
	c.rootCmd.AddCommand(listCmd)
}
