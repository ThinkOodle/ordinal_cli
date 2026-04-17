package cmd

import (
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/skill"
	"github.com/spf13/cobra"
)

var (
	skillInstallTargets []string
	skillInstallForce   bool
)

func init() {
	rootCmd.AddCommand(skillCmd)
	skillCmd.AddCommand(skillInstallCmd)

	skillInstallCmd.Flags().StringSliceVar(&skillInstallTargets, "target", []string{"all"}, "Install target(s): agents, claude, codex, all")
	skillInstallCmd.Flags().BoolVar(&skillInstallForce, "force", false, "Overwrite existing skill files when content differs")
}

var skillCmd = &cobra.Command{
	Use:         "skill",
	Short:       "Manage bundled AI agent skills",
	Long:        "Install the bundled Ordinal CLI usage skill for agents that operate the installed `ordinal` binary.",
	Annotations: map[string]string{skipOutputValidationAnnotation: "true"},
}

var skillInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the bundled Ordinal CLI skill",
	Long:  "Install the bundled Ordinal CLI usage SKILL.md for agents that need to operate the installed `ordinal` program, writing it into standard skill directories such as ~/.agents/skills, ~/.claude/skills, and ~/.codex/skills.",
	RunE: func(cmd *cobra.Command, args []string) error {
		results, err := skill.Install(skill.InstallOptions{
			Targets: skillInstallTargets,
			Force:   skillInstallForce,
		})
		if err != nil {
			return err
		}

		for _, result := range results {
			fmt.Printf("%s: %s (%s)\n", result.Target, result.Path, result.Status)
		}

		return nil
	},
}
