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

// A conflict on any target must abort the whole install without writing
// any target. Prior behavior walked targets in sorted order and wrote each
// one as it went, so a conflict on "claude" after "agents" had already been
// written left the filesystem in a partially-applied state.
func TestInstall_ConflictAbortsBeforeAnyWrite(t *testing.T) {
	home := t.TempDir()
	claudePath := filepath.Join(home, ".claude/skills/ordinal-cli/SKILL.md")
	if err := os.MkdirAll(filepath.Dir(claudePath), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(claudePath, []byte("pre-existing"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	agentsPath := filepath.Join(home, ".agents/skills/ordinal-cli/SKILL.md")

	results, err := Install(InstallOptions{HomeDir: home, Targets: []string{"agents", "claude"}})
	if err == nil {
		t.Fatalf("expected conflict error, got results=%+v", results)
	}
	if _, statErr := os.Stat(agentsPath); statErr == nil {
		t.Errorf("agents target was written despite claude conflict: %s exists", agentsPath)
	}
	if data, _ := os.ReadFile(claudePath); string(data) != "pre-existing" {
		t.Errorf("claude target was modified despite conflict; got %q", data)
	}
}
