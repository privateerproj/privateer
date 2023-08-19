package cmd

import (
	"fmt"
	"os"
	"path"
	"text/tabwriter"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer-sdk/command"
	"github.com/privateerproj/privateer-sdk/logging"
)

var (
	buildVersion       string
	buildGitCommitHash string
	buildTime          string

	logger hclog.Logger      // enables formatted logging (logger.Trace, etc)
	writer *tabwriter.Writer // enables bare line writing (for use in list & version)

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:              "privateer",
		Short:            "TODO: some kind of description",
		PersistentPreRun: persistentPreRun,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	cmd.PersistentFlags().StringP("binaries-path", "b", defaultBinariesPath(), "Path to the directory where raids are installed")
	viper.BindPFlag("binaries-path", cmd.PersistentFlags().Lookup("binaries-path"))

	command.InitializeConfig()
	logger = logging.GetLogger("core", viper.GetString("loglevel"), false)
	logger.Trace(fmt.Sprintf(
		"Initialized core logger: %s", viper.GetString("loglevel")))

	// writer is used for output in the list & version commands
	writer = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
}

func defaultBinariesPath() string {
	home, _ := os.UserHomeDir() // sue me
	return path.Join(home, "privateer", "bin")
}
