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
		logger.Trace("%s called", runName)
		if len(args) > 1 {
			logger.Error(fmt.Sprintf(
				"Sally only accepts a single argument. Unknown args: %v",args[1:]))
		} else if len(args) == 1 {
			raidName := args[0]
			configVar := make(map[string]interface{}, 1)
			configVar[raidName] = make(map[string]interface{}, 1)
			viper.Set("Raids", configVar)

			runApprovedRaid(raidName)
		} else {
			logger.Trace("Sequentially executing all raids in config") // TODO
		}
	},
}

func init() {

	rootCmd.AddCommand(runCmd)
}

func runApprovedRaid(raidName string) (err error) {
	err = installIfNotPResent(raidName)
	if err != nil {
		logger.Error(fmt.Sprintf(
			"Installation failed for raid '%s': %v", raidName, err))
		return
	}
	logger.Trace(fmt.Sprintf(
		"Beginning raid '%s'", raidName))
	// TODO get the logger set up in the run commands. 
	// ...Do those benefit from being outside of cmd?
	err = run.CLIContext()
	if err != nil {
		logger.Error("Error with execute logic.") // don't forget to fix all these logs
	}
	return
}

func installIfNotPResent(raidName string) (err error) {

	installed := false
	for _, raid := range run.GetAvailableRaids() {
		if raid == raidName {
			installed = true
		}
	}
	if !installed {
		logger.Trace(fmt.Sprintf(
			"Installing raid: %s", raidName))
		err = downloadRaid(raidName)
	}
	return err
}

func downloadRaid(raidName string) (err error) {
	url := approvedRaids[raidName]
	logger.Trace(fmt.Sprintf(
		"Attempting download from: %s", url))
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	u := strings.Split(url, "/")
	f := u[len(u)-1]
	localpath := filepath.Join(viper.GetString("binaries-path"), f)
	logger.Trace(fmt.Sprintf(
		"Creating file: %s", localpath))
	out, err := os.Create(localpath)
	if err != nil {
		return
	}
	defer out.Close()

	logger.Trace("Setting file permissions to 0755")
	err = out.Chmod(0755)
	if err != nil {
		return
	}

	_, err = io.Copy(out, resp.Body)
	return
}