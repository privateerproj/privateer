// Package main provides the entry point for the pvtr command-line tool.
// pvtr is a validation framework that executes plugins to perform
// compliance and security evaluations.
package main

import (
	"runtime/debug"

	"github.com/privateerproj/privateer/cmd"
)

var (
	// Version is the version string that will be replaced at build time by the associated tag.
	Version string
	// VersionPostfix is a marker for the version such as "dev", "beta", "rc", etc.
	// This is appended to the version string if it is not empty.
	VersionPostfix string
	// GitCommitHash is the git commit hash at build time.
	GitCommitHash string
	// BuiltAt is the actual build datetime string.
	BuiltAt string
)

func setVCSInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	for _, setting := range info.Settings {
	   switch setting.Key {
	   case "vcs.revision":
		  GitCommitHash = setting.Value
	   case "vcs.time":
		  BuiltAt = setting.Value
	   case "vcs.modified":
		  VersionPostfix = "-dev"
	   }
	}
}

// main is the entry point for the pvtr application.
// It formats the version string with any postfix and executes the root command.
func main() {
	if Version == "" {
		Version = "edge"
		setVCSInfo()
	}
	if VersionPostfix != "" {
		Version += VersionPostfix
	}
	cmd.NewCLI(Version, GitCommitHash, BuiltAt).Execute()
}
