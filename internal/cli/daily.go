package cli

import (
	"github.com/spf13/cobra"
	"github.com/mswiente/pkm.ai/internal/config"
	"github.com/mswiente/pkm.ai/internal/daily"
)

func newDailyCommand(cfg *config.Config) *cobra.Command {
	dailyCmd := &cobra.Command{
		Use:   "daily",
		Short: "Manage daily notes",
	}
	dailyCmd.AddCommand(newDailyCreateCommand(cfg))
	return dailyCmd
}

func newDailyCreateCommand(cfg *config.Config) *cobra.Command {
	var (
		date string
		open bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create today's daily note",
		Long:  `Creates a daily note for today (or a specified date) in 01-daily/. Idempotent: does nothing if the note already exists.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := mustVault(cfg)
			if err != nil {
				return err
			}
			r := newRenderer(cfg)
			opts := daily.Options{
				Date: date,
				Open: open,
			}
			return daily.Run(v, r, cfg, opts)
		},
	}

	cmd.Flags().StringVar(&date, "date", "", "Date for the note (YYYY-MM-DD); default is today")
	cmd.Flags().BoolVar(&open, "open", false, "Open the note in the configured editor after creation")

	return cmd
}
