// Package cmd provides the command-line interface for pvtr.
// It defines the root command and all subcommands, handles configuration,
// and manages the execution flow of the application.
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer-sdk/command"
	"github.com/privateerproj/privateer-sdk/config"
)

// CLI holds the shared dependencies for all pvtr commands.
type CLI struct {
	buildVersion       string
	buildGitCommitHash string
	buildTime          string

	logger  hclog.Logger
	writer  *tabwriter.Writer
	rootCmd *cobra.Command
}

// NewCLI creates a CLI instance with the given build metadata
// and registers all subcommands.
func NewCLI(version, commitHash, builtAt string) *CLI {
	c := &CLI{
		buildVersion:       version,
		buildGitCommitHash: commitHash,
		buildTime:          builtAt,
		logger:             hclog.NewNullLogger(),
		writer:             tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
	}

	c.rootCmd = &cobra.Command{
		Use:              "pvtr",
		Short:            "pvtr root command",
		PersistentPreRun: c.persistentPreRun,
	}

	command.SetBase(c.rootCmd)
	c.rootCmd.PersistentFlags().StringP("binaries-path", "b", defaultBinariesPath(), "Path to the directory where plugins are installed")
	_ = viper.BindPFlag("binaries-path", c.rootCmd.PersistentFlags().Lookup("binaries-path"))

	c.addRunCmd()
	c.addEnvCmd()
	c.addVersionCmd()
	c.addListCmd()
	c.addGenPluginCmd()
	c.addInstallCmd()

	return c
}

// Execute runs the root command. This is called by main.main().
func (c *CLI) Execute() {
	err := c.rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// persistentPreRun initializes the logger and writer for use by all commands.
func (c *CLI) persistentPreRun(cmd *cobra.Command, args []string) {
	cfg := config.NewConfig(nil)
	c.logger = cfg.Logger

	if c.writer == nil {
		c.writer = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	}
	command.ReadConfig()
}

// defaultBinariesPath returns the default path where plugins are installed.
// It constructs a path in the user's home directory under .privateer/bin.
// If the home directory cannot be determined, it falls back to a relative
// path (./.privateer/bin) and logs a warning.
func defaultBinariesPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Warning: could not determine home directory: %v", err)
		return filepath.Join(".", ".privateer", "bin")
	}
	return filepath.Join(home, ".privateer", "bin")
}
