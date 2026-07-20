package cmd

import (
	"github.com/privateerproj/privateer-sdk/command/harness"
)

func (c *CLI) addBenchmarkCmd() {
	c.rootCmd.AddCommand(harness.GetBenchmarkCmd(func() harness.Writer { return c.writer }))
}
