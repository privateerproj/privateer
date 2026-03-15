package install

import (
	"fmt"
	"runtime"
	"strings"
)

// InferGitHubReleaseBase returns the GitHub releases download base URL for the given source and version.
// source should be a GitHub repo URL (e.g. https://github.com/owner/repo). latest is the version (e.g. "0.19.2");
// the download tag will be "v" + latest unless latest already starts with "v". Returns ("", false) if source
// does not look like a GitHub URL or latest is empty.
func InferGitHubReleaseBase(source, latest string) (base string, ok bool) {
	source = strings.TrimSpace(source)
	latest = strings.TrimSpace(latest)
	if source == "" || latest == "" {
		return "", false
	}
	source = strings.TrimSuffix(source, "/")
	source = strings.TrimSuffix(source, ".git")
	if !strings.Contains(source, "github.com") {
		return "", false
	}
	tag := latest
	if !strings.HasPrefix(tag, "v") {
		tag = "v" + tag
	}
	return source + "/releases/download/" + tag, true
}

// InferArtifactFilename returns the release artifact filename for the current platform.
// Convention: {binaryName}_{OS}_{arch}.{ext} with OS in Darwin/Linux/Windows,
// arch in all/arm64/i386/x86_64, ext .tar.gz for Darwin/Linux and .zip for Windows.
func InferArtifactFilename(binaryName string) (string, error) {
	binaryName = strings.TrimSpace(binaryName)
	if binaryName == "" {
		return "", fmt.Errorf("binary name is required")
	}

	var osLabel string
	switch runtime.GOOS {
	case "darwin":
		osLabel = "Darwin"
	case "linux":
		osLabel = "Linux"
	case "windows":
		osLabel = "Windows"
	default:
		return "", fmt.Errorf("unsupported GOOS %q", runtime.GOOS)
	}

	var archLabel string
	var ext string
	if runtime.GOOS == "darwin" {
		archLabel = "all"
		ext = ".tar.gz"
	} else {
		switch runtime.GOARCH {
		case "amd64":
			archLabel = "x86_64"
		case "386":
			archLabel = "i386"
		case "arm64":
			archLabel = "arm64"
		default:
			return "", fmt.Errorf("unsupported GOARCH %q", runtime.GOARCH)
		}
		if runtime.GOOS == "windows" {
			ext = ".zip"
		} else {
			ext = ".tar.gz"
		}
	}

	return binaryName + "_" + osLabel + "_" + archLabel + ext, nil
}
