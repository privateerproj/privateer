package main

import (
	"os"
	"testing"
)

func TestSetVCSInfo(t *testing.T) {
	origGitCommitHash := GitCommitHash
	origBuiltAt := BuiltAt
	origVersionPostfix := VersionPostfix
	defer func() {
		GitCommitHash = origGitCommitHash
		BuiltAt = origBuiltAt
		VersionPostfix = origVersionPostfix
	}()

	setVCSInfo()
}

func TestMain_DefaultEdgeVersion(t *testing.T) {
	origVersion := Version
	origVersionPostfix := VersionPostfix
	origArgs := os.Args
	defer func() {
		Version = origVersion
		VersionPostfix = origVersionPostfix
		os.Args = origArgs
	}()

	Version = ""
	VersionPostfix = "-dev"
	os.Args = []string{"pvtr", "version"}

	main()

	if Version != "edge-dev" {
		t.Errorf("expected Version to be 'edge-dev', got %q", Version)
	}
}

func TestMain_CustomVersionWithoutPostfix(t *testing.T) {
	origVersion := Version
	origVersionPostfix := VersionPostfix
	origArgs := os.Args
	defer func() {
		Version = origVersion
		VersionPostfix = origVersionPostfix
		os.Args = origArgs
	}()

	Version = "1.0.0"
	VersionPostfix = ""
	os.Args = []string{"pvtr", "version"}

	main()

	if Version != "1.0.0" {
		t.Errorf("expected Version to be '1.0.0', got %q", Version)
	}
}
