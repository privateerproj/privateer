package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestVersionCmd_DefaultOutputShowsVersionOnly(t *testing.T) {
	var buf bytes.Buffer
	c := newTestCLI(&buf, "1.2.3", "a1b2c3d", "2026-01-01T00:00:00Z")
	c.addVersionCmd()
	defer viper.Reset()

	c.rootCmd.SetArgs([]string{"version"})
	if err := c.rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error executing version command: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "1.2.3") {
		t.Errorf("expected output to contain version '1.2.3', got:\n%s", output)
	}
	// Non-verbose should only show version, not commit or build time
	if strings.Contains(output, "Commit:") {
		t.Errorf("expected no commit info in non-verbose mode, got:\n%s", output)
	}
}

func TestVersionCmd_VerboseFlagOutputShowsAllFields(t *testing.T) {
	var buf bytes.Buffer
	c := newTestCLI(&buf, "2.0.0", "f1e2d3c", "2026-06-15T12:00:00Z")
	c.addVersionCmd()
	defer viper.Reset()

	c.rootCmd.SetArgs([]string{"version", "--verbose"})
	if err := c.rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error executing version --verbose: %v", err)
	}

	output := buf.String()

	expectedFields := []string{
		"Version:", "2.0.0",
		"Commit:", "f1e2d3c",
		"Build Time:", "2026-06-15T12:00:00Z",
	}
	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("expected verbose output to contain %q, got:\n%s", field, output)
		}
	}
}
