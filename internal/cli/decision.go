package cli

import (
	"github.com/spf13/cobra"
	"github.com/mswiente/pkm.ai/internal/config"
	"github.com/mswiente/pkm.ai/internal/decision"
)

func newDecisionCommand(cfg *config.Config) *cobra.Command {
	decisionCmd := &cobra.Command{
		Use:   "decision",
		Short: "Manage decision notes",
	}
	decisionCmd.AddCommand(newDecisionCreateCommand(cfg))
	return decisionCmd
}

func newDecisionCreateCommand(cfg *config.Config) *cobra.Command {
	var (
		title     string
		project   string
		status    string
		fromStdin bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a decision note",
		Long:  `Creates a decision note in 00-inbox/ from the decision template.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := mustVault(cfg)
			if err != nil {
				return err
			}
			r := newRenderer(cfg)
			opts := decision.Options{
				Title:     title,
				Project:   project,
				Status:    status,
				FromStdin: fromStdin,
			}
			return decision.Run(v, r, cfg, opts)
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Decision title (required)")
	cmd.Flags().StringVar(&project, "project", "", "Related project or area")
	cmd.Flags().StringVar(&status, "status", "draft", "Decision status: draft|accepted|superseded")
	cmd.Flags().BoolVar(&fromStdin, "from-stdin", false, "Read decision body from stdin")

	return cmd
}
