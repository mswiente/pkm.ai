package cli

import (
	"github.com/spf13/cobra"
	"github.com/mswiente/pkm.ai/internal/config"
	"github.com/mswiente/pkm.ai/internal/process"
)

func newProcessCommand(cfg *config.Config) *cobra.Command {
	processCmd := &cobra.Command{
		Use:   "process",
		Short: "Process notes",
	}
	processCmd.AddCommand(newProcessInboxCommand(cfg))
	return processCmd
}

func newProcessInboxCommand(cfg *config.Config) *cobra.Command {
	var (
		file        string
		all         bool
		dryRun      bool
		apply       bool
		interactive bool
	)

	cmd := &cobra.Command{
		Use:   "inbox",
		Short: "Analyze and report on inbox notes",
		Long: `Prints a structured report of inbox notes for AI-assisted processing.

Use this command within a Claude Code session to get a summary of all inbox
notes that can then be analyzed and organized by the AI.

No files are modified unless --apply is passed (not yet implemented).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := mustVault(cfg)
			if err != nil {
				return err
			}
			if interactive {
				return process.RunInteractive(v)
			}
			opts := process.Options{
				File:   file,
				All:    all,
				DryRun: dryRun,
				Apply:  apply,
			}
			return process.Run(v, opts)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "Process a single file (basename only)")
	cmd.Flags().BoolVar(&all, "all", false, "Process all inbox files (default behavior)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show report only, no modifications")
	cmd.Flags().BoolVar(&apply, "apply", false, "Apply suggested changes (not yet implemented)")
	cmd.Flags().BoolVar(&interactive, "interactive", false, "Walk through each inbox note with a single-key action menu")

	return cmd
}
