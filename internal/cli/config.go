package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mswiente/pkm.ai/internal/config"
)

func newConfigCommand(cfg *config.Config) *cobra.Command {
	var (
		setVaultPath string
		show         bool
	)

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Show or update pkm configuration",
		Long:  fmt.Sprintf("Manages pkm configuration stored at %s.", config.ConfigFilePath()),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("set-vault-path") && !show {
				return fmt.Errorf("no flag provided — use --show or --set-vault-path\nRun: pkm config --help")
			}

			if cmd.Flags().Changed("set-vault-path") {
				if setVaultPath == "" {
					return fmt.Errorf("vault path cannot be empty")
				}
				if err := config.SetVaultPath(setVaultPath); err != nil {
					return err
				}
				fmt.Printf("vault_path set to: %s\n", setVaultPath)
				fmt.Printf("Config file: %s\n", config.ConfigFilePath())
			}

			if show {
				printConfig(cfg)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&setVaultPath, "set-vault-path", "", "Set the vault path and save to config file")
	cmd.Flags().BoolVar(&show, "show", false, "Show current configuration")

	return cmd
}

func printConfig(cfg *config.Config) {
	fmt.Printf("Config file:       %s\n", config.ConfigFilePath())
	fmt.Println()
	if cfg.VaultPath == "" {
		fmt.Println("vault_path:        (not set) — run: pkm config --set-vault-path \"/path/to/vault\"")
	} else {
		fmt.Printf("vault_path:        %s\n", cfg.VaultPath)
	}
	fmt.Printf("inbox_path:        %s\n", cfg.InboxPath)
	fmt.Printf("daily_path:        %s\n", cfg.DailyPath)
	fmt.Printf("templates_path:    %s\n", cfg.TemplatesPath)
	fmt.Printf("editor:            %s\n", cfg.Editor)
	fmt.Printf("filename_timezone: %s\n", cfg.FilenameTimezone)
	fmt.Printf("default_source:    %s\n", cfg.DefaultSource)
}
