// Package cmd provides the command-line interface for Privateer.
// It defines the root command and all subcommands, handles configuration,
// and manages the execution flow of the application.
package cmd

import (
	"os"
	"path"
	"text/tabwriter"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer-sdk/command"
)

var (
	// buildVersion holds the version string set at build time.
	buildVersion string
	// buildGitCommitHash holds the git commit hash set at build time.
	buildGitCommitHash string
	// buildTime holds the build timestamp set at build time.
	buildTime string

	// logger enables formatted logging with methods like logger.Trace, logger.Info, etc.
	logger hclog.Logger
	// writer enables formatted tabular output for use in list and version commands.
	writer *tabwriter.Writer

	// rootCmd represents the base command when called without any subcommands.
	// It is the entry point for all Privateer commands and handles global configuration.
	rootCmd = &cobra.Command{
		Use:              "privateer",
		Short:            "privateer root command",
		PersistentPreRun: persistentPreRun,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
//
// Parameters:
//   - version: The version string to display in version command
//   - commitHash: The git commit hash to display in version command
//   - builtAt: The build timestamp to display in version command
func Execute(version, commitHash, builtAt string) {
	buildVersion = version
	buildGitCommitHash = commitHash
	buildTime = builtAt

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	command.SetBase(rootCmd)
	rootCmd.PersistentFlags().StringP("binaries-path", "b", defaultBinariesPath(), "Path to the directory where plugins are installed")
	_ = viper.BindPFlag("binaries-path", rootCmd.PersistentFlags().Lookup("binaries-path"))
}

// persistentPreRun initializes the logger and writer for use by all commands.
// It is called before every command execution and sets up lightweight logging
// that does not create files or directories on disk.
func persistentPreRun(cmd *cobra.Command, args []string) {
	loglevel := viper.GetString("loglevel")
	if loglevel == "" {
		loglevel = "error"
	}
	logger = hclog.New(&hclog.LoggerOptions{
		Level:  hclog.LevelFromString(loglevel),
		Output: os.Stderr,
	})

	// writer is used for output in the list & version commands
	writer = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	command.ReadConfig()
}

// defaultBinariesPath returns the default path where plugins are installed.
// It constructs a path in the user's home directory under .privateer/bin.
func defaultBinariesPath() string {
	home, _ := os.UserHomeDir()
	return path.Join(home, ".privateer", "bin")
}
