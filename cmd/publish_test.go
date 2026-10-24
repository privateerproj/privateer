package cmd

import (
	"bytes"
	"testing"
)

func TestPublishCmd_Registration(t *testing.T) {
	var buf bytes.Buffer
	c := newTestCLI(&buf, "1.0.0", "hash", "date")
	c.addPublishCmd()

	cmd, _, err := c.rootCmd.Find([]string{"publish"})
	if err != nil {
		t.Fatalf("expected to find publish command: %v", err)
	}
	if cmd == nil || cmd.Name() != "publish" {
		t.Fatalf("expected publish command to be registered")
	}

	flags := []string{"dist", "registry", "no-sync"}
	for _, f := range flags {
		if cmd.Flags().Lookup(f) == nil {
			t.Errorf("expected flag %q to be registered on publish command", f)
		}
	}
}

func TestPublishCmd_WriterFn(t *testing.T) {
	var buf bytes.Buffer
	c := newTestCLI(&buf, "1.0.0", "hash", "date")
	c.addPublishCmd()

	cmd, _, _ := c.rootCmd.Find([]string{"publish"})
	// RunE executes writerFn() to obtain the writer
	_ = cmd.RunE(cmd, []string{})
}
