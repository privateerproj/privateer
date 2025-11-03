package cmd

import (
	"github.com/privateerproj/privateer-sdk/command"
)

// Re-export types from SDK for backwards compatibility
type PluginError = command.PluginError
type PluginErrors = command.PluginErrors
type PluginPkg = command.PluginPkg

var NewPluginPkg = command.NewPluginPkg
