package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"

	hcplugin "github.com/hashicorp/go-plugin"

	"github.com/privateerproj/privateer/internal/config"
	"github.com/privateerproj/privateer/internal/flags"
	"github.com/privateerproj/privateer/run"
)

var (
	// See Makefile for more on how this package is built

	// Version is the main version number that is being run at the moment
	Version = "0.1.0"

	// VersionPostfix is a marker for the version. If this is "" (empty string)
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "dev" (in development), "beta", "rc", etc.
	VersionPostfix = "dev"

	// GitCommitHash references the commit id at build time
	GitCommitHash = ""

	// BuiltAt is the build date
	BuiltAt = ""
)

func main() {
	var subCommand string
	if len(os.Args) > 1 {
		subCommand = os.Args[1]
	}
	switch subCommand {
	// Ref: https://gobyexample.com/command-line-subcommands
	case "list":
		flags.List.Parse(os.Args[2:])
		listRaids()

	case "version":
		flags.Version.Parse(os.Args[2:])
		printVersion()

	default:
		flags.Run.Parse(os.Args[1:])
		run.CLIContext()
	}
}

func printVersion() {
	if VersionPostfix != "" {
		Version = fmt.Sprintf("%s-%s", Version, VersionPostfix)
	}

	fmt.Fprintf(os.Stdout, "Privateer Version: %s", Version)
	if config.Vars.Verbose != nil && *config.Vars.Verbose {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "Commit       : %s", GitCommitHash)
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "Built at     : %s", BuiltAt)
	}
}

// listRaids reads all raids declared in config and checks whether they are installed
func listRaids() {
	raidNames, err := getRaidNames()
	if err != nil {
		log.Fatalf("An error occurred while retrieving raids from config: %v", err)
	}

	raids := make(map[string]string)
	for _, raid := range raidNames {
		raidName, binErr := run.GetRaidBinary(raid)
		if binErr != nil {
			raids[raid] = fmt.Sprintf("ERROR: %v", binErr)
		} else {
			raids[filepath.Base(raidName)] = "OK"
		}
	}

	// Print output
	writer := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(writer, "| Raid\t | Installed ")
	for k, v := range raids {
		fmt.Fprintf(writer, "| %s\t | %s\n", k, v)
	}
	writer.Flush()
}

// getRaidNames returns all raids declared in config file
func getRaidNames() (raidNames []string, err error) {
	if err != nil || (config.Vars.AllRaids != nil && *config.Vars.AllRaids) {
		return hcplugin.Discover("*", config.Vars.BinariesPath)
	}
	return config.Vars.Run, nil
}
