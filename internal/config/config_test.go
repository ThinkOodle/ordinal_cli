package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_DefaultsWhenNoConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ORDINAL_API_KEY", "")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != DefaultOutputFormat {
		t.Errorf("expected default output format %q, got %q", DefaultOutputFormat, cfg.OutputFormat)
	}
}

func TestSaveAPIKey_Roundtrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := SaveAPIKey("abc123"); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	path, err := ConfigFilePath()
	if err != nil {
		t.Fatalf("config path: %v", err)
	}
	if filepath.Dir(path) == "" {
		t.Fatalf("unexpected empty dir for config path %q", path)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected config file at %s: %v", path, err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.APIKey != "abc123" {
		t.Errorf("expected abc123, got %q", cfg.APIKey)
	}
}

func TestLoad_BrokenConfigPreservesEnvAPIKey(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ORDINAL_API_KEY", "env-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	dir := filepath.Join(home, ".config", configDirName)
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(dir, configFileName+"."+configFileType)
	// Deliberately malformed YAML so Viper's parser fails.
	if err := os.WriteFile(path, []byte("api_key: [not: valid"), 0600); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg, err := Load()
	if err == nil {
		t.Fatalf("expected a read error for malformed config")
	}
	if cfg == nil {
		t.Fatalf("expected a non-nil config even on broken file")
	}
	if cfg.APIKey != "env-key" {
		t.Errorf("expected env ORDINAL_API_KEY to survive broken config, got %q", cfg.APIKey)
	}
	if cfg.OutputFormat != DefaultOutputFormat {
		t.Errorf("expected default output format to survive, got %q", cfg.OutputFormat)
	}
}

// TestSaveAPIKey_RecoversFromMalformedConfig locks in the recovery path the
// root command's PersistentPreRunE promises: if the saved config file can't
// be parsed, `ordinal auth` must still succeed so the user isn't stranded.
// The unparseable bytes are preserved in a timestamped backup so any other
// settings the user cared about remain recoverable by hand.
func TestSaveAPIKey_RecoversFromMalformedConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dir := filepath.Join(home, ".config", configDirName)
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(dir, configFileName+"."+configFileType)
	original := []byte("api_key: [not: valid\noutput_format: csv\n")
	if err := os.WriteFile(path, original, 0600); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err := SaveAPIKey("new-key"); err != nil {
		t.Fatalf("expected recovery to succeed, got: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("post-recovery Load: %v", err)
	}
	if cfg.APIKey != "new-key" {
		t.Errorf("expected api key new-key, got %q", cfg.APIKey)
	}

	// The unparseable original must survive somewhere so the user can
	// hand-recover settings they may have cared about (e.g. output_format).
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("readdir: %v", err)
	}
	var backup []byte
	for _, e := range entries {
		if e.Name() == filepath.Base(path) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("read backup %s: %v", e.Name(), err)
		}
		if string(data) == string(original) {
			backup = data
			break
		}
	}
	if backup == nil {
		t.Errorf("expected a backup of the original unparseable config in %s, found entries: %v", dir, entries)
	}
}

func TestGetAPIKey_Priority(t *testing.T) {
	cfg := &Config{APIKey: "from-config"}
	if got, err := GetAPIKey("flag-key", cfg); err != nil || got != "flag-key" {
		t.Errorf("flag should win: got=%q err=%v", got, err)
	}
	if got, err := GetAPIKey("", cfg); err != nil || got != "from-config" {
		t.Errorf("config should be used when flag empty: got=%q err=%v", got, err)
	}
	if _, err := GetAPIKey("", &Config{}); err == nil {
		t.Errorf("expected error when both are empty")
	}
}
