package cmd

import (
	"bytes"
	"testing"
)

func TestLoginCmds_Registration(t *testing.T) {
	var buf bytes.Buffer
	c := newTestCLI(&buf, "1.0.0", "hash", "date")
	c.addLoginCmds()

	for _, name := range []string{"login", "logout"} {
		cmd, _, err := c.rootCmd.Find([]string{name})
		if err != nil {
			t.Fatalf("expected to find %s command: %v", name, err)
		}
		if cmd == nil || cmd.Name() != name {
			t.Errorf("expected %s command to be registered", name)
		}
	}
}
