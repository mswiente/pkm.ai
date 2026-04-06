package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/mswiente/pkm.ai/internal/config"
	"github.com/mswiente/pkm.ai/internal/project"
)

func newProjectCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage project notes in 02-projects/",
		Long: `Primitives for maintaining 02-projects/ project notes.

Each project has one note with five sections:
  Intent        — what the project is trying to achieve (stable)
  Current Status — what was last worked on (updated each session)
  Next Steps    — what to do next (updated each session)
  Timeline      — one-line-per-session quick log (append-only)
  Plan History  — dated log of captured plans (append-only)

Designed to be called by Claude Code when processing captured plan notes:

  pkm project update pkm-ai --title "pkm.ai CLI" \
    --current-status "Readwise sync done. Distill workflow implemented." \
    --next-steps "- [ ] Add project management\n- [ ] Raycast integration" \
    --plan-heading "2026-04-06 — Distill workflow" < plan.md

  pkm project list`,
	}
	cmd.AddCommand(
		newProjectUpdateCommand(cfg),
		newProjectListCommand(cfg),
	)
	return cmd
}

// pkm project update <slug>
func newProjectUpdateCommand(cfg *config.Config) *cobra.Command {
	var (
		title         string
		projectStatus string
		intent        string
		currentStatus string
		nextSteps     string
		planHeading   string
		timelineEntry string
		dryRun        bool
	)

	cmd := &cobra.Command{
		Use:   "update <slug>",
		Short: "Create or update a project note in 02-projects/",
		Long: `Creates 02-projects/<slug>.md if it doesn't exist, or patches specific sections.

Only sections with provided flags are updated; others are left unchanged.
If stdin has content it is appended to the ## Plan History section under
--plan-heading (defaults to today's date).

Example — create a new project:
  pkm project update pkm-ai \
    --title "pkm.ai CLI" \
    --intent "A lightweight CLI for a markdown-first PKM system." \
    --current-status "Readwise sync and distill workflow implemented." \
    --next-steps "- [ ] Add project management commands"

Example — update after a work session (with plan content from stdin):
  cat plan.md | pkm project update pkm-ai \
    --current-status "Implemented pkm project commands." \
    --next-steps "- [ ] Update distill skill\n- [ ] Write tests" \
    --plan-heading "2026-04-06 — Project management"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slug := args[0]
			if strings.ContainsAny(slug, " /\\") {
				return fmt.Errorf("slug must be kebab-case with no spaces or slashes")
			}

			// Read plan content from stdin if available
			var planContent string
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				data, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("read stdin: %w", err)
				}
				planContent = strings.TrimSpace(string(data))
			}

			v, err := mustVault(cfg)
			if err != nil {
				return err
			}

			opts := project.UpdateOptions{
				Slug:          slug,
				Title:         title,
				ProjectStatus: projectStatus,
				Intent:        intent,
				CurrentStatus: currentStatus,
				NextSteps:     nextSteps,
				PlanContent:   planContent,
				PlanHeading:   planHeading,
				TimelineEntry: timelineEntry,
				DryRun:        dryRun,
			}

			result, err := project.Update(v, opts)
			if err != nil {
				return err
			}

			verb := "Updated"
			if result.IsNew {
				verb = "Created"
			}
			prefix := ""
			if dryRun {
				prefix = "[dry-run] Would have "
				verb = strings.ToLower(verb)
			}
			fmt.Printf("%s%s 02-projects/%s.md\n", prefix, verb, slug)
			if len(result.Patched) > 0 {
				fmt.Printf("  Sections: %s\n", strings.Join(result.Patched, ", "))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Project title (used when creating)")
	cmd.Flags().StringVar(&projectStatus, "status", "", "Project status: active | on-hold | archived")
	cmd.Flags().StringVar(&intent, "intent", "", "Replace the ## Intent section")
	cmd.Flags().StringVar(&currentStatus, "current-status", "", "Replace the ## Current Status section")
	cmd.Flags().StringVar(&nextSteps, "next-steps", "", "Replace the ## Next Steps section")
	cmd.Flags().StringVar(&planHeading, "plan-heading", "", "Heading for the plan history entry (default: today's date)")
	cmd.Flags().StringVar(&timelineEntry, "timeline-entry", "", "One-line summary appended to ## Timeline (defaults to --plan-heading value)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be written without making changes")

	return cmd
}

// pkm project list
func newProjectListCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all project notes in 02-projects/",
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := mustVault(cfg)
			if err != nil {
				return err
			}

			slugs, err := project.List(v)
			if err != nil {
				return err
			}
			if len(slugs) == 0 {
				fmt.Println("No project notes found in 02-projects/.")
				return nil
			}

			fmt.Printf("Projects (%d):\n", len(slugs))
			for _, slug := range slugs {
				content, _ := project.Load(v, slug)
				status := extractFrontmatterValue(content, "status")
				title := extractFrontmatterValue(content, "title")
				if title == "" {
					title = slug
				}
				marker := " "
				if status == "active" {
					marker = "●"
				} else if status == "on-hold" {
					marker = "○"
				} else if status == "archived" {
					marker = "–"
				}
				fmt.Printf("  %s %-30s  %s\n", marker, slug, title)
			}
			return nil
		},
	}
}

// extractFrontmatterValue pulls a scalar value from YAML frontmatter by field name.
func extractFrontmatterValue(content, field string) string {
	prefix := field + ": "
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}
