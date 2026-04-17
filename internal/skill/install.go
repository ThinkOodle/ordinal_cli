// Package skill provides access to the bundled Ordinal CLI skill and helpers
// for installing it into common agent skill directories.
package skill

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	skillName = "ordinal-cli"
	skillFile = "SKILL.md"
)

var targetRoots = map[string]string{
	"agents": ".agents/skills",
	"claude": ".claude/skills",
	"codex":  ".codex/skills",
}

//go:embed assets/ordinal-cli/SKILL.md
var assets embed.FS

// InstallOptions controls how the bundled skill is installed.
type InstallOptions struct {
	Targets []string
	Force   bool
	HomeDir string
}

// InstallResult describes the outcome for a single install target.
type InstallResult struct {
	Target string
	Path   string
	Status string
}

// Install writes the bundled skill into the selected target directories.
func Install(opts InstallOptions) ([]InstallResult, error) {
	homeDir := opts.HomeDir
	if homeDir == "" {
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("determining home directory: %w", err)
		}
	}

	targets, err := normalizeTargets(opts.Targets)
	if err != nil {
		return nil, err
	}

	content, err := assets.ReadFile("assets/ordinal-cli/SKILL.md")
	if err != nil {
		return nil, fmt.Errorf("reading bundled skill: %w", err)
	}

	results := make([]InstallResult, 0, len(targets))
	for _, target := range targets {
		root := filepath.Join(homeDir, targetRoots[target], skillName)
		if err := os.MkdirAll(root, 0755); err != nil {
			return results, fmt.Errorf("creating %s skill directory: %w", target, err)
		}

		path := filepath.Join(root, skillFile)
		result := InstallResult{
			Target: target,
			Path:   path,
			Status: "installed",
		}

		existing, err := os.ReadFile(path)
		switch {
		case err == nil && bytes.Equal(existing, content):
			result.Status = "unchanged"
			results = append(results, result)
			continue
		case err == nil && !opts.Force:
			return results, fmt.Errorf("skill already exists at %s; rerun with --force to overwrite", path)
		case err == nil:
			result.Status = "updated"
		case !os.IsNotExist(err):
			return results, fmt.Errorf("reading existing skill at %s: %w", path, err)
		}

		if err := os.WriteFile(path, content, 0644); err != nil {
			return results, fmt.Errorf("writing skill to %s: %w", path, err)
		}

		results = append(results, result)
	}

	return results, nil
}

func normalizeTargets(raw []string) ([]string, error) {
	if len(raw) == 0 {
		return []string{"agents", "claude", "codex"}, nil
	}

	seen := map[string]bool{}
	targets := make([]string, 0, len(raw))

	for _, entry := range raw {
		target := strings.TrimSpace(strings.ToLower(entry))
		if target == "" {
			continue
		}
		if target == "all" {
			return []string{"agents", "claude", "codex"}, nil
		}
		if _, ok := targetRoots[target]; !ok {
			return nil, fmt.Errorf("invalid skill target %q; valid targets are: agents, claude, codex, all", entry)
		}
		if !seen[target] {
			seen[target] = true
			targets = append(targets, target)
		}
	}

	if len(targets) == 0 {
		return []string{"agents", "claude", "codex"}, nil
	}

	slices.Sort(targets)
	return targets, nil
}
