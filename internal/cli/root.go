package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/mswiente/pkm.ai/internal/config"
	"github.com/mswiente/pkm.ai/internal/templates"
	"github.com/mswiente/pkm.ai/internal/vault"
)

// NewRootCommand builds and returns the root cobra command with all subcommands wired.
func NewRootCommand(cfg *config.Config) *cobra.Command {
	root := &cobra.Command{
		Use:   "pkm",
		Short: "Personal knowledge management CLI",
		Long:  "pkm is a CLI for a markdown-first personal knowledge system built around Obsidian.",
		// Silence default error printing; we handle it in main.
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.AddCommand(
		newCaptureCommand(cfg),
		newProcessCommand(cfg),
		newDailyCommand(cfg),
		newMeetingCommand(cfg),
		newDecisionCommand(cfg),
		newNoteCommand(cfg),
		newConfigCommand(cfg),
		newSkillCommand(),
		newSyncCommand(cfg),
	)

	return root
}

// mustVault creates a Vault or returns a descriptive error.
func mustVault(cfg *config.Config) (*vault.Vault, error) {
	if err := config.RequireVaultPath(cfg); err != nil {
		return nil, err
	}
	if _, err := os.Stat(cfg.VaultPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("vault not found at %s\nCreate it with: bash setup/setup_pkm_vault.sh %q", cfg.VaultPath, cfg.VaultPath)
	}
	return vault.New(cfg)
}

// newRenderer creates a Renderer pointed at the vault's templates directory.
func newRenderer(cfg *config.Config) *templates.Renderer {
	return templates.NewRenderer(cfg.VaultPath + "/" + cfg.TemplatesPath)
}
