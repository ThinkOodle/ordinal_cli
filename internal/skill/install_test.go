package skill

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstall_DefaultTargets(t *testing.T) {
	home := t.TempDir()
	results, err := Install(InstallOptions{HomeDir: home})
	if err != nil {
		t.Fatalf("install failed: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if _, err := os.Stat(r.Path); err != nil {
			t.Errorf("expected file at %s: %v", r.Path, err)
		}
		if r.Status != "installed" {
			t.Errorf("expected installed, got %q", r.Status)
		}
	}
}

func TestInstall_UnchangedOnRerun(t *testing.T) {
	home := t.TempDir()
	if _, err := Install(InstallOptions{HomeDir: home, Targets: []string{"claude"}}); err != nil {
		t.Fatalf("first install: %v", err)
	}
	results, err := Install(InstallOptions{HomeDir: home, Targets: []string{"claude"}})
	if err != nil {
		t.Fatalf("second install: %v", err)
	}
	if len(results) != 1 || results[0].Status != "unchanged" {
		t.Errorf("expected unchanged, got %+v", results)
	}
}

func TestInstall_ForceOverwrite(t *testing.T) {
	home := t.TempDir()
	path := filepath.Join(home, ".claude/skills/ordinal-cli/SKILL.md")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte("old content"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	// Without --force, should error
	if _, err := Install(InstallOptions{HomeDir: home, Targets: []string{"claude"}}); err == nil {
		t.Errorf("expected error when file exists without force")
	}

	results, err := Install(InstallOptions{HomeDir: home, Targets: []string{"claude"}, Force: true})
	if err != nil {
		t.Fatalf("force install: %v", err)
	}
	if len(results) != 1 || results[0].Status != "updated" {
		t.Errorf("expected updated, got %+v", results)
	}
}

func TestInstall_InvalidTarget(t *testing.T) {
	home := t.TempDir()
	if _, err := Install(InstallOptions{HomeDir: home, Targets: []string{"nonsense"}}); err == nil {
		t.Errorf("expected error for invalid target")
	}
}
