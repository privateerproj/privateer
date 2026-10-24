package cmd

import (
	"bytes"
	"testing"
)

func TestBenchmarkCmd_Registration(t *testing.T) {
	var buf bytes.Buffer
	c := newTestCLI(&buf, "1.0.0", "hash", "date")
	c.addBenchmarkCmd()

	cmd, _, err := c.rootCmd.Find([]string{"benchmark"})
	if err != nil {
		t.Fatalf("expected to find benchmark command: %v", err)
	}
	if cmd == nil || cmd.Name() != "benchmark" {
		t.Fatalf("expected benchmark command to be registered")
	}

	flags := []string{"service", "target", "json", "write-directory", "payload-only"}
	for _, f := range flags {
		if cmd.Flags().Lookup(f) == nil {
			t.Errorf("expected flag %q to be registered on benchmark command", f)
		}
	}
}
