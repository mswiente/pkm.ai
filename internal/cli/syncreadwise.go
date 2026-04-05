package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"
	"github.com/mswiente/pkm.ai/internal/config"
	"github.com/mswiente/pkm.ai/internal/readwise"
)

func newSyncReadwiseCommand(cfg *config.Config) *cobra.Command {
	var dryRun bool
	var since string
	var limit int

	cmd := &cobra.Command{
		Use:   "readwise",
		Short: "Sync saved articles from Readwise Reader to your vault inbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := mustVault(cfg)
			if err != nil {
				return err
			}

			opts := readwise.Options{
				DryRun: dryRun,
				Limit:  limit,
			}

			if since != "" {
				t, err := time.Parse("2006-01-02", since)
				if err != nil {
					return fmt.Errorf("invalid --since date %q: use YYYY-MM-DD format", since)
				}
				opts.Since = &t
			}

			if dryRun {
				fmt.Println("Dry run — no files will be written.")
			}

			result, err := readwise.Run(v, cfg, opts)
			if err != nil {
				return err
			}

			if dryRun {
				fmt.Printf("Would sync %d article(s).\n", result.Synced)
			} else {
				fmt.Printf("Synced %d article(s) to inbox.", result.Synced)
				if result.Skipped > 0 {
					fmt.Printf(" Skipped %d already-saved article(s).", result.Skipped)
				}
				fmt.Println()
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be synced without writing files")
	cmd.Flags().StringVar(&since, "since", "", "Sync articles updated after this date (YYYY-MM-DD); overrides state file")
	cmd.Flags().IntVar(&limit, "limit", 0, "Maximum number of articles to sync (0 = no limit)")

	cmd.AddCommand(newSyncReadwiseAuthCommand(cfg))
	return cmd
}

func newSyncReadwiseAuthCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "Save your Readwise API token to the pkm config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Get your Readwise API token from:")
			fmt.Println("  https://readwise.io/access_token")
			fmt.Println()

			token, err := readToken()
			if err != nil {
				return fmt.Errorf("read token: %w", err)
			}
			token = strings.TrimSpace(token)
			if token == "" {
				return fmt.Errorf("no token provided")
			}

			fmt.Print("Validating token... ")
			client := readwise.NewClient(token)
			if err := client.ValidateToken(context.Background()); err != nil {
				fmt.Println("failed.")
				return err
			}
			fmt.Println("OK.")

			if err := config.SetReadwiseToken(token); err != nil {
				return fmt.Errorf("save token: %w", err)
			}

			fmt.Printf("Token saved to %s\n", config.ConfigFilePath())
			return nil
		},
	}
}

// readToken reads a token from stdin. Uses masked input when running in a terminal.
func readToken() (string, error) {
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		fmt.Print("Paste your token: ")
		b, err := term.ReadPassword(fd)
		fmt.Println() // newline after masked input
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	// Non-interactive: read a line from stdin
	var token string
	_, err := fmt.Fscan(os.Stdin, &token)
	return token, err
}
