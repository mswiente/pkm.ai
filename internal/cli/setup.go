package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/mswiente/pkm.ai/internal/claudemd"
	"github.com/mswiente/pkm.ai/internal/config"
	"github.com/mswiente/pkm.ai/internal/skill"
)

func newSetupCommand(cfg *config.Config) *cobra.Command {
	var (
		dryRun      bool
		skillsOnly  bool
		claudeMDOnly bool
	)

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Bootstrap Claude Code integration for the PKM system",
		Long: `Installs PKM skills and writes a global ~/.claude/CLAUDE.md so that
Claude Code has full PKM context in every working directory.

  pkm setup              # install skills + write ~/.claude/CLAUDE.md
  pkm setup --dry-run    # show what would be written without making changes
  pkm setup --skills-only
  pkm setup --claude-md-only`,
		RunE: func(cmd *cobra.Command, args []string) error {
			doSkills := !claudeMDOnly
			doClaudeMD := !skillsOnly

			if doSkills {
				if err := runSkillsInstall(dryRun); err != nil {
					return err
				}
			}

			if doClaudeMD {
				if err := runClaudeMDInstall(cfg, dryRun); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be written without making changes")
	cmd.Flags().BoolVar(&skillsOnly, "skills-only", false, "Only install skills, skip CLAUDE.md")
	cmd.Flags().BoolVar(&claudeMDOnly, "claude-md-only", false, "Only write CLAUDE.md, skip skills")

	return cmd
}

func runSkillsInstall(dryRun bool) error {
	if dryRun {
		skills, err := skill.List()
		if err != nil {
			return err
		}
		home, _ := os.UserHomeDir()
		dir := filepath.Join(home, ".claude", "commands")
		for _, s := range skills {
			fmt.Printf("[dry-run] Would install /%s → %s/%s.md\n", s.Name, dir, s.Name)
		}
		return nil
	}
	return skill.Install("")
}

func runClaudeMDInstall(cfg *config.Config, dryRun bool) error {
	if err := config.RequireVaultPath(cfg); err != nil {
		return fmt.Errorf("vault path required to generate CLAUDE.md — run: pkm config --set-vault-path <path>\n%w", err)
	}

	content, err := claudemd.Render(cfg.VaultPath)
	if err != nil {
		return fmt.Errorf("render CLAUDE.md: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}
	claudeDir := filepath.Join(home, ".claude")
	dst := filepath.Join(claudeDir, "CLAUDE.md")

	if dryRun {
		fmt.Printf("[dry-run] Would write %s (%d bytes)\n", dst, len(content))
		fmt.Println("--- preview ---")
		fmt.Print(content)
		fmt.Println("--- end preview ---")
		return nil
	}

	// Warn if file already exists.
	if _, err := os.Stat(dst); err == nil {
		fmt.Printf("Overwriting existing %s\n", dst)
	}

	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		return fmt.Errorf("create .claude dir: %w", err)
	}
	if err := os.WriteFile(dst, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", dst, err)
	}

	fmt.Printf("Written: %s\n", dst)
	return nil
}
