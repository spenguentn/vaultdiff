package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmd_NoArgs(t *testing.T) {
	// Running the root command with --help should not error.
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("expected no error on --help, got: %v", err)
	}

	out := buf.String()
	if out == "" {
		t.Error("expected help output, got empty string")
	}
}

func TestRootCmd_HasOutputFlag(t *testing.T) {
	f := rootCmd.PersistentFlags().Lookup("output")
	if f == nil {
		t.Fatal("expected --output flag to be registered")
	}
	if f.DefValue != "text" {
		t.Errorf("expected default output 'text', got %q", f.DefValue)
	}
}

func TestRootCmd_HasMaskSecretsFlag(t *testing.T) {
	f := rootCmd.PersistentFlags().Lookup("mask-secrets")
	if f == nil {
		t.Fatal("expected --mask-secrets flag to be registered")
	}
	if f.DefValue != "true" {
		t.Errorf("expected default mask-secrets 'true', got %q", f.DefValue)
	}
}

func TestRootCmd_HasConfigFlag(t *testing.T) {
	f := rootCmd.PersistentFlags().Lookup("config")
	if f == nil {
		t.Fatal("expected --config flag to be registered")
	}
	if f.DefValue != "" {
		t.Errorf("expected empty default for --config, got %q", f.DefValue)
	}
}
