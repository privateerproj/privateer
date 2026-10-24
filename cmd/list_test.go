package cmd

import (
	"bytes"
	"testing"
)

func TestListCmd_Registration(t *testing.T) {
	var buf bytes.Buffer
	c := newTestCLI(&buf, "1.0.0", "hash", "date")
	c.addListCmd()

	cmd, _, err := c.rootCmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("expected to find list command: %v", err)
	}
	if cmd == nil || cmd.Name() != "list" {
		t.Fatalf("expected list command to be registered")
	}

	if cmd.Flags().Lookup("installable") == nil {
		t.Errorf("expected flag 'installable' to be registered on list command")
	}
}

func TestListCmd_WriterFn(t *testing.T) {
	var buf bytes.Buffer
	c := newTestCLI(&buf, "1.0.0", "hash", "date")
	c.addListCmd()

	cmd, _, _ := c.rootCmd.Find([]string{"list"})
	// RunE executes writerFn() to obtain the writer
	_ = cmd.RunE(cmd, []string{})
}
