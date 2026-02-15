package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"text/tabwriter"
)

func TestEnvCmd_ContainsExpectedFields(t *testing.T) {
	var buf bytes.Buffer
	writer = tabwriter.NewWriter(&buf, 1, 1, 1, ' ', 0)

	buildVersion = "1.2.3"
	buildGitCommitHash = "a1b2c3d"
	buildTime = "2026-01-01T00:00:00Z"

	envCmd.Run(envCmd, []string{})

	output := buf.String()

	expectedFields := []string{
		"Binary:",
		"Config:",
		"Plugins Dir:",
		"Plugins:",
		"Version:",
		"Commit:",
		"Build Time:",
		"Go Version:",
		"OS/Arch:",
	}
	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("expected output to contain %q, got:\n%s", field, output)
		}
	}
}

func TestEnvCmd_DisplaysBuildInfo(t *testing.T) {
	var buf bytes.Buffer
	writer = tabwriter.NewWriter(&buf, 1, 1, 1, ' ', 0)

	buildVersion = "test-version"
	buildGitCommitHash = "e4f5a6b"
	buildTime = "2026-02-15T00:00:00Z"

	envCmd.Run(envCmd, []string{})

	output := buf.String()

	if !strings.Contains(output, "test-version") {
		t.Errorf("expected output to contain version 'test-version', got:\n%s", output)
	}
	if !strings.Contains(output, "e4f5a6b") {
		t.Errorf("expected output to contain commit 'e4f5a6b', got:\n%s", output)
	}
	if !strings.Contains(output, "2026-02-15T00:00:00Z") {
		t.Errorf("expected output to contain build time, got:\n%s", output)
	}
	if !strings.Contains(output, runtime.Version()) {
		t.Errorf("expected output to contain Go version %q, got:\n%s", runtime.Version(), output)
	}
	expectedArch := runtime.GOOS + "/" + runtime.GOARCH
	if !strings.Contains(output, expectedArch) {
		t.Errorf("expected output to contain OS/Arch %q, got:\n%s", expectedArch, output)
	}
}

func TestEnvCmd_ShowsBinaryPath(t *testing.T) {
	var buf bytes.Buffer
	writer = tabwriter.NewWriter(&buf, 1, 1, 1, ' ', 0)

	envCmd.Run(envCmd, []string{})

	output := buf.String()

	execPath, err := os.Executable()
	if err == nil && !strings.Contains(output, execPath) {
		t.Errorf("expected output to contain executable path %q, got:\n%s", execPath, output)
	}
}

func TestDiscoverPluginNames_EmptyDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "privateer-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	result := discoverPluginNames(tmpDir)
	if result != "none" {
		t.Errorf("expected 'none' for empty dir, got: %s", result)
	}
}

func TestDiscoverPluginNames_NonexistentDir(t *testing.T) {
	result := discoverPluginNames("/nonexistent/path")
	if result != "none" {
		t.Errorf("expected 'none' for nonexistent dir, got: %s", result)
	}
}

func TestDiscoverPluginNames_FiltersPrivateer(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "privateer-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	for _, name := range []string{"privateer", "privateer-foo", "my-plugin", "other-tool"} {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0755); err != nil {
			t.Fatalf("failed to create file %s: %v", name, err)
		}
	}

	result := discoverPluginNames(tmpDir)

	if strings.Contains(result, "privateer") {
		t.Errorf("expected privateer binaries to be filtered out, got: %s", result)
	}
	if !strings.Contains(result, "my-plugin") {
		t.Errorf("expected 'my-plugin' in result, got: %s", result)
	}
	if !strings.Contains(result, "other-tool") {
		t.Errorf("expected 'other-tool' in result, got: %s", result)
	}
}
