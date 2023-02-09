package flags

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/privateerproj/privateer/internal/config"
)

// Run flags relate to the primary privateer execution
var Run *flag.FlagSet

// List flags manage the view of installed binaries
var List *flag.FlagSet

// Version flags relate to the version information for this Privateer installation
var Version *flag.FlagSet

func init() {
	Run = flag.NewFlagSet("privateer", flag.ExitOnError)
	List = flag.NewFlagSet("privateer list", flag.ExitOnError)
	Version = flag.NewFlagSet("privateer version", flag.ExitOnError)

	config.Vars.Init()

	addAllFlag(Run)

	addAllFlag(List)

	addVerboseFlag(Version)
}

func addVarsFileFlag(flagSet *flag.FlagSet) {
	config.Vars.VarsFile = flagSet.String("config-file", defaultConfigPath(), "Location for config vars file.")
}

func addVerboseFlag(flagSet *flag.FlagSet) {
	config.Vars.Verbose = flagSet.Bool("v", *config.Vars.Verbose, "Display extended version information")
}

func addAllFlag(flagSet *flag.FlagSet) {
	config.Vars.AllRaids = flagSet.Bool("all", *config.Vars.AllRaids, "Include all installed packs, not just those specified within the provided config")
}

func defaultConfigPath() string {
	workDir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Join(workDir, "config.yml")
}
