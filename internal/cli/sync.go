package cli

import (
	"github.com/spf13/cobra"
	"github.com/mswiente/pkm.ai/internal/config"
)

func newSyncCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync content from external services into your vault",
	}
	cmd.AddCommand(newSyncReadwiseCommand(cfg))
	return cmd
}
