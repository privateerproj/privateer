package cmd

import (
	"bytes"
	"strings"
	"testing"
	"text/tabwriter"

	"github.com/spf13/viper"
)

func TestVersionCmd_DefaultOutputShowsVersionOnly(t *testing.T) {
	var buf bytes.Buffer
	writer = tabwriter.NewWriter(&buf, 1, 1, 1, ' ', 0)

	buildVersion = "1.2.3"
	buildGitCommitHash = "a1b2c3d"
	buildTime = "2026-01-01T00:00:00Z"
	viper.Set("verbose", false)
	defer viper.Reset()

	versionCmd.Run(versionCmd, []string{})

	output := buf.String()

	// Non-verbose should not write to the tabwriter at all
	if output != "" {
		t.Errorf("expected no tabwriter output in non-verbose mode, got:\n%s", output)
	}
}

func TestVersionCmd_VerboseOutputShowsAllFields(t *testing.T) {
	var buf bytes.Buffer
	writer = tabwriter.NewWriter(&buf, 1, 1, 1, ' ', 0)

	buildVersion = "2.0.0"
	buildGitCommitHash = "f1e2d3c"
	buildTime = "2026-06-15T12:00:00Z"
	viper.Set("verbose", true)
	defer viper.Reset()

	versionCmd.Run(versionCmd, []string{})

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
