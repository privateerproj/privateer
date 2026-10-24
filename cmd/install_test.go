package cmd

import (
	"bytes"
	"testing"
)

func TestInstallCmd_Registration(t *testing.T) {
	var buf bytes.Buffer
	c := newTestCLI(&buf, "1.0.0", "hash", "date")
	c.addInstallCmd()

	cmd, _, err := c.rootCmd.Find([]string{"install"})
	if err != nil {
		t.Fatalf("expected to find install command: %v", err)
	}
	if cmd == nil || cmd.Name() != "install" {
		t.Fatalf("expected install command to be registered")
	}

	flags := []string{"local", "from-config"}
	for _, f := range flags {
		if cmd.Flags().Lookup(f) == nil {
			t.Errorf("expected flag %q to be registered on install command", f)
		}
	}
}
