package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/mswiente/pkm.ai/internal/config"
	"github.com/mswiente/pkm.ai/internal/meeting"
)

func newMeetingCommand(cfg *config.Config) *cobra.Command {
	meetingCmd := &cobra.Command{
		Use:   "meeting",
		Short: "Manage meeting notes",
	}
	meetingCmd.AddCommand(newMeetingCreateCommand(cfg))
	return meetingCmd
}

func newMeetingCreateCommand(cfg *config.Config) *cobra.Command {
	var (
		title            string
		date             string
		project          string
		participantsFlag string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a meeting note",
		Long:  `Creates a meeting note in 00-inbox/ from the meeting template.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := mustVault(cfg)
			if err != nil {
				return err
			}
			r := newRenderer(cfg)

			var participants []string
			if participantsFlag != "" {
				for _, p := range strings.Split(participantsFlag, ",") {
					if p := strings.TrimSpace(p); p != "" {
						participants = append(participants, p)
					}
				}
			}

			opts := meeting.Options{
				Title:        title,
				Date:         date,
				Project:      project,
				Participants: participants,
			}
			return meeting.Run(v, r, cfg, opts)
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Meeting title (prompted if not provided in interactive mode)")
	cmd.Flags().StringVar(&date, "date", "", "Meeting date (YYYY-MM-DD); default is today")
	cmd.Flags().StringVar(&project, "project", "", "Related project or area")
	cmd.Flags().StringVar(&participantsFlag, "participants", "", "Comma-separated list of participants")

	return cmd
}
