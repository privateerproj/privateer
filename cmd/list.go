package cmd

import (
	"github.com/privateerproj/privateer-sdk/command"
)

var (
	cmdName = "list"
)

// listCmd represents the list command
var listCmd = command.GetListCmd(writer)

func init() {
	rootCmd.AddCommand(listCmd)
}

// Re-export functions from SDK for backwards compatibility
var GetPlugins = command.GetPlugins
var Contains = command.Contains
