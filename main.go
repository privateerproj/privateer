// Package main provides the entry point for the Privateer command-line tool.
// Privateer is a security testing framework that executes plugins to perform
// security evaluations and assessments.
package main

import (
	"fmt"

	"github.com/privateerproj/privateer/cmd"
)

var (
	// Version is the version string that will be replaced at build time by the associated tag.
	Version = "0.0.0"
	// VersionPostfix is a marker for the version such as "dev", "beta", "rc", etc.
	// This is appended to the version string if it is not empty.
	VersionPostfix = "dev"
	// GitCommitHash is the git commit hash at build time.
	GitCommitHash = ""
	// BuiltAt is the actual build datetime string.
	BuiltAt = ""
)

// main is the entry point for the Privateer application.
// It formats the version string with any postfix and executes the root command.
func main() {
	if VersionPostfix != "" {
		Version = fmt.Sprintf("%s-%s", Version, VersionPostfix)
	}
	cmd.Execute(Version, GitCommitHash, BuiltAt)
}
