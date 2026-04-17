package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/config"
)

// TestAuth_RejectsEmptyOrWhitespaceKey locks in that `ordinal auth ""` and
// `ordinal auth "   "` fail up front. cobra.ExactArgs(1) only checks argument
// count, so without an explicit TrimSpace guard a blank key would be saved to
// config and then sent in the Authorization header on every subsequent
// request, surfacing as an opaque 401 long after the fact.
func TestAuth_RejectsEmptyOrWhitespaceKey(t *testing.T) {
	tests := []struct {
		name string
		arg  string
	}{
		{"empty quoted arg", ""},
		{"single space", " "},
		{"tabs and spaces", "  \t "},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			t.Setenv("ORDINAL_API_KEY", "")
			t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
			t.Setenv("ORDINAL_VERBOSE", "")

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs([]string{"auth", tc.arg})

			err := rootCmd.Execute()
			if err == nil {
				t.Fatalf("expected error for blank api key arg %q", tc.arg)
			}
			if !strings.Contains(err.Error(), "empty") {
				t.Errorf("expected 'empty' error, got: %v", err)
			}

			// Confirm the config file was not written with a blank key.
			// SaveAPIKey would create it under $HOME/.config/ordinal.
			if _, err := os.Stat(filepath.Join(home, ".config", "ordinal", "config.yaml")); err == nil {
				cfg, loadErr := config.Load()
				if loadErr == nil && cfg.APIKey != "" {
					t.Errorf("expected config unchanged, got saved api key %q", cfg.APIKey)
				}
			}
		})
	}
}

// TestAuth_TrimsSurroundingWhitespace locks in that a key padded with
// surrounding whitespace (easy to introduce via shell history or copy-paste)
// is saved trimmed, so it matches the canonical form expected by the server.
func TestAuth_TrimsSurroundingWhitespace(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ORDINAL_API_KEY", "")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"auth", "  real-key  "})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("auth with padded key: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}
	if cfg.APIKey != "real-key" {
		t.Errorf("expected saved key 'real-key', got %q", cfg.APIKey)
	}
}
