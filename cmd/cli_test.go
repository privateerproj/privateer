package cmd

import (
	"testing"
)

func TestNewCLI_RegistersAllSubcommands(t *testing.T) {
	c := NewCLI("1.0.0", "abc123", "2026-01-01T00:00:00Z")

	expectedCommands := []string{
		"run",
		"env",
		"version",
		"list",
		"generate-plugin",
		"install",
		"publish",
		"login",
		"logout",
		"benchmark",
	}

	registered := make(map[string]bool)
	for _, cmd := range c.rootCmd.Commands() {
		registered[cmd.Name()] = true
	}

	for _, name := range expectedCommands {
		if !registered[name] {
			t.Errorf("expected subcommand %q to be registered, but it was not", name)
		}
	}
}

func TestNewCLI_LoggerIsNonNil(t *testing.T) {
	c := NewCLI("1.0.0", "abc123", "2026-01-01T00:00:00Z")
	if c.logger == nil {
		t.Fatal("expected logger to be non-nil after NewCLI, got nil")
	}
	c.logger.Trace("regression: logger must be safe to call before persistentPreRun")
}

func TestNewCLI_WriterIsNonNil(t *testing.T) {
	c := NewCLI("1.0.0", "abc123", "2026-01-01T00:00:00Z")
	if c.writer == nil {
		t.Fatal("expected writer to be non-nil after NewCLI, got nil")
	}
	_, err := c.writer.Write([]byte("regression: writer must be safe to call before persistentPreRun"))
	if err != nil {
		t.Fatalf("unexpected error writing to default writer: %v", err)
	}
}
func TestRunCmd_DoesNotPanicWithoutPersistentPreRun(t *testing.T) {
	c := NewCLI("1.0.0", "abc123", "2026-01-01T00:00:00Z")

	originalRunFn := runFn
	runFn = func(_ *CLI) int { return 0 }
	defer func() { runFn = originalRunFn }()

	runCmd, _, err := c.rootCmd.Find([]string{"run"})
	if err != nil {
		t.Fatalf("expected to find 'run' subcommand: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("runCmd.Run panicked without persistentPreRun: %v", r)
		}
	}()
	runCmd.Run(runCmd, []string{})
}

// The harness flags (hub-url, autoinstall) are registered by the SDK, not by
// this repo, so nothing here fails to compile if the call is dropped — hence
// this check.
func TestNewCLI_RegistersHarnessFlags(t *testing.T) {
	c := NewCLI("1.0.0", "abc123", "2026-01-01T00:00:00Z")

	for _, name := range []string{"hub-url", "autoinstall"} {
		if c.rootCmd.PersistentFlags().Lookup(name) == nil {
			t.Errorf("expected persistent flag %q to be registered, but it was not", name)
		}
	}
}
