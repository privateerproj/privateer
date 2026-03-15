package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FromImage installs a binary from an OCI image into destPath. binaryPathInImage is the path
// inside the image (empty for default "/plugin"); binaryName is the local file name to write.
// Uses docker or podman from PATH: create container, cp file out, rm container.
func FromImage(imageRef, binaryPathInImage, destPath, binaryName string) error {
	_ = binaryName
	pathInImage := strings.TrimSpace(binaryPathInImage)
	if pathInImage == "" {
		pathInImage = "/plugin"
	}

	cli, err := findContainerCLI()
	if err != nil {
		return err
	}

	containerName := "pvtr-install-" + filepath.Base(destPath)
	defer func() {
		_ = exec.Command(cli, "rm", "-f", containerName).Run()
	}()

	// create (pull if needed, no run)
	create := exec.Command(cli, "create", "--name", containerName, imageRef)
	create.Stdout = nil
	create.Stderr = nil
	if err := create.Run(); err != nil {
		return fmt.Errorf("%s create: %w (is the image public and pullable?)", cli, err)
	}

	// cp container:path -> destPath
	cp := exec.Command(cli, "cp", containerName+":"+pathInImage, destPath)
	cp.Stdout = nil
	cp.Stderr = nil
	if err := cp.Run(); err != nil {
		return fmt.Errorf("%s cp: %w (path %q in image may not exist)", cli, err, pathInImage)
	}

	_ = os.Chmod(destPath, 0755)
	return nil
}

func findContainerCLI() (string, error) {
	for _, name := range []string{"docker", "podman"} {
		path, err := exec.LookPath(name)
		if err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("install from image requires docker or podman in PATH")
}
