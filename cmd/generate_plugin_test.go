package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
)

func TestGeneratePluginCmd_Registration(t *testing.T) {
	var buf bytes.Buffer
	c := newTestCLI(&buf, "1.0.0", "hash", "date")
	c.addGenPluginCmd()

	cmd, _, err := c.rootCmd.Find([]string{"generate-plugin"})
	if err != nil {
		t.Fatalf("expected to find generate-plugin command: %v", err)
	}
	if cmd == nil || cmd.Name() != "generate-plugin" {
		t.Fatalf("expected generate-plugin command to be registered")
	}
	if cmd.Short != "Generate a new plugin" {
		t.Errorf("unexpected short description: %s", cmd.Short)
	}
}

func TestGeneratePluginCmd_Flags(t *testing.T) {
	var buf bytes.Buffer
	c := newTestCLI(&buf, "1.0.0", "hash", "date")
	c.addGenPluginCmd()

	cmd, _, err := c.rootCmd.Find([]string{"generate-plugin"})
	if err != nil {
		t.Fatalf("expected to find generate-plugin command: %v", err)
	}

	tests := []struct {
		name         string
		shorthand    string
		defaultValue string
	}{
		{"source-path", "p", ""},
		{"local-templates", "", ""},
		{"service-name", "n", ""},
		{"organization", "g", ""},
		{"output-dir", "o", "generated-plugin/"},
	}

	for _, tc := range tests {
		flag := cmd.Flags().Lookup(tc.name)
		if flag == nil {
			t.Errorf("expected flag %q to be registered", tc.name)
			continue
		}
		if flag.Shorthand != tc.shorthand {
			t.Errorf("flag %q: expected shorthand %q, got %q", tc.name, tc.shorthand, flag.Shorthand)
		}
		if flag.DefValue != tc.defaultValue {
			t.Errorf("flag %q: expected default value %q, got %q", tc.name, tc.defaultValue, flag.DefValue)
		}
	}
}

func TestGeneratePluginCmd_ViperBindings(t *testing.T) {
	var buf bytes.Buffer
	c := newTestCLI(&buf, "1.0.0", "hash", "date")
	c.addGenPluginCmd()

	cmd, _, err := c.rootCmd.Find([]string{"generate-plugin"})
	if err != nil {
		t.Fatalf("expected to find generate-plugin command: %v", err)
	}

	if err := cmd.Flags().Set("output-dir", "custom-dir/"); err != nil {
		t.Fatalf("unexpected error setting flag: %v", err)
	}
	if val := viper.GetString("output-dir"); val != "custom-dir/" {
		t.Errorf("expected viper to reflect flag value 'custom-dir/', got %q", val)
	}
}
