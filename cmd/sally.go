package cmd

import (
	"fmt"
	"os"
	"io"
	"net/http"
	"strings"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/privateerproj/privateer/run"
)

var runName = "sally"

// This is just a single download link for now. Later, stub the URL at the version & infer the filename based on OS type
var approvedRaids = map[string]string{
	"wireframe": "https://github.com/privateerproj/raid-hello-world/releases/download/v0.0.0/wireframe",
}

// runCmd represents the sally command
var runCmd = &cobra.Command{
	Use:   runName,
	Short: "When everything is battoned down, it is time to sally forth.",
	Long:  `TODO - Long description`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s called\n", runName)
		if len(args) > 1 {
			fmt.Printf("Sally only accepts a single argument. Unknown args: %v\n", args[1:])
			os.Exit(1)
		} else if len(args) == 1 {
			raidName := args[0]
			configVar := make(map[string]interface{}, 1)
			configVar[raidName] = make(map[string]interface{}, 1)
			viper.Set("Raids", configVar)

			err := installIfNotPResent(raidName)
			if err != nil {
				fmt.Printf("Installation failed for raid '%s': %v\n", raidName, err)
			}

			fmt.Printf("Calling sally for raid '%s'\n", raidName) // TODO
			run.CLIContext()
		} else {
			fmt.Printf("Calling sally for all raids in config\n") // TODO
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runApprovedRaid(raidName string) {
	err := installIfNotPResent(raidName)
	if err != nil {
		fmt.Print("Error with installer logic.") // these are all just temporary logs
		os.Exit(1)
	}
	err = executeRaid(raidName)
	if err != nil {
		fmt.Print("Error with execute logic.") // don't forget to fix all these logs
		os.Exit(1)
	}
}

func installIfNotPResent(raidName string) (err error) {
	fmt.Printf("install called for %s\n", raidName)

	installed := false
	for _, raid := range run.GetAvailableRaids() {
		if raid == raidName {
			installed = true
		}
	}
	if !installed {
		fmt.Printf("installing raid: %s\n", raidName)
		err = downloadRaid(raidName)
	}
	return err
}

func executeRaid(raidName string) error {
	fmt.Printf("sally %s called\n", raidName)
	return nil
}

// DownloadFile will download a url to a local file.
// It will write as it downloads and will not load the whole file into memory.
func downloadRaid(raidName string) (err error) {
	url := approvedRaids[raidName]
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	u := strings.Split(url, "/")
	filename := u[len(u)-1]
	localpath := filepath.Join(viper.GetString("binaries-path"), filename)
	fmt.Printf("filepath: %s\n", localpath)

	out, err := os.Create(localpath)
	if err != nil {
		return
	}
	defer out.Close()

	err = out.Chmod(0755)
	if err != nil {
		return
	}

	_, err = io.Copy(out, resp.Body)
	return
}