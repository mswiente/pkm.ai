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
		full        bool
	)

	cmd := &cobra.Command{
		Use:   "inbox",
		Short: "Analyze and report on inbox notes",
		Long: `Prints a structured report of inbox notes for AI-assisted processing.

Use this command in a Claude Code session to get a full picture of what is in
the inbox and what topic pages already exist in 04-knowledge/. Claude Code can
then propose distillation and call pkm knowledge / pkm note move to apply changes.

--full outputs the complete body of each note and appends the current
04-knowledge/index.md so Claude Code has everything it needs in one pass.`,
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
				Full:   full,
			}
			return process.Run(v, opts)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "Process a single file (basename only)")
	cmd.Flags().BoolVar(&all, "all", false, "Process all inbox files (default behavior)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show report only, no modifications")
	cmd.Flags().BoolVar(&apply, "apply", false, "Apply suggested changes (not yet implemented)")
	cmd.Flags().BoolVar(&interactive, "interactive", false, "Walk through each inbox note with a single-key action menu")
	cmd.Flags().BoolVar(&full, "full", false, "Output complete note bodies and the current knowledge index")

	return cmd
}
