package cmd

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

var (
	// raids listed here should only be production-ready, and use the latest-stable version
	approvedRaids = map[string]string{
		"wireframe": "https://github.com/privateerproj/raid-hello-world/releases/download/v0.0.0/wireframe",
	}

	buildTypes = map[string]string{
		"darwin":       "Darwin_all.tar.gz",
		"linuxarm64":   "Linux_arm64.tar.gz",
		"linux386":     "Linux_i386.tar.gz",
		"windowsarm64": "Windows_arm64.zip",
		"windows386":   "Windows_i386.zip",
		"windowsamd64": "Windows_x86_64.zip",
	}
)

// StartApprovedRaid will run a single raid after ensuring it is installed
// Approved raids are listed in run/approved-raids.go
func StartApprovedRaid(raidName string) (err error) {
	err = installIfNotPresent(raidName)
	if err != nil {
		logger.Error(fmt.Sprintf(
			"Installation failed for raid '%s': %v", raidName, err))
		return
	}
	logger.Trace(fmt.Sprintf(
		"Beginning raid '%s'", raidName))
	err = Run()
	return
}

func installIfNotPresent(raidName string) (err error) {
	installed := false
	raids := GetRaids()
	for _, raid := range raids {
		if raid.Name == raidName {
			installed = true
			logger.Trace("Raid already installed.")
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
	// url, err := getDownloadPath(raidName)
	// if err != nil {
	// 	return
	// }
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

func getDownloadPath(raidName string) (string, error) {
	url := approvedRaids[raidName]
	if url == "" {
		return "", errors.New("raid not found")
	}
	return fmt.Sprintf("%s_%s", url, getExtension()), nil
}

func getExtension() string {
	if runtime.GOOS == "darwin" {
		return buildTypes["darwin"]
	}
	buildType := runtime.GOOS + runtime.GOARCH
	if buildTypes[buildType] != "" {
		return buildTypes[buildType]
	}
	if runtime.GOOS == "windows" {
		return buildTypes["windowsamd64"]
	}
	return buildTypes["linux386"]
}

// TODO: clean this up
func Unzip(src string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			logger.Error(err.Error())
		}
		// TODO: delete zip file here
	}()
	dest := viper.GetString("binaries-path")
	os.MkdirAll(dest, 0755)

	for _, file := range r.File {
		err := extractAndWrite(file, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractAndWrite(file *zip.File, dest string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			logger.Error(err.Error())
		}
	}()
	path := filepath.Join(dest, file.Name)

	if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", path)
	}

	if file.FileInfo().IsDir() {
		os.MkdirAll(path, file.Mode())
	} else {
		os.MkdirAll(filepath.Dir(path), file.Mode())
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				logger.Error(err.Error())
			}
		}()

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}
	return nil
}
