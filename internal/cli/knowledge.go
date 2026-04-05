package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/mswiente/pkm.ai/internal/config"
	"github.com/mswiente/pkm.ai/internal/process"
)

func newKnowledgeCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "knowledge",
		Short: "Manage 04-knowledge/ topic pages, index, and log",
		Long: `Primitives for maintaining the 04-knowledge/ wiki.

Designed to be called by Claude Code after it has analyzed inbox notes:

  pkm knowledge append-topic <slug> --title "Title" < content.md
  pkm knowledge update-index <slug> --description "one-line desc"
  pkm knowledge append-log --note <file> --action <action> --filed-to <folder>`,
	}
	cmd.AddCommand(
		newKnowledgeAppendTopicCommand(cfg),
		newKnowledgeUpdateIndexCommand(cfg),
		newKnowledgeAppendLogCommand(cfg),
	)
	return cmd
}

// pkm knowledge append-topic <slug> --title "Title"
// Reads markdown content from stdin and appends it to 04-knowledge/<slug>.md.
// Creates the page with frontmatter if it does not exist.
func newKnowledgeAppendTopicCommand(cfg *config.Config) *cobra.Command {
	var title string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "append-topic <slug>",
		Short: "Append content (from stdin) to a 04-knowledge/ topic page",
		Long: `Reads markdown content from stdin and appends it to 04-knowledge/<slug>.md.
Creates the page with standard frontmatter if it does not exist.

Example:
  echo "## From [[2026-04-03-readwise-article]]\n\nKey insight here." \
    | pkm knowledge append-topic ai-agents --title "AI Agents"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slug := args[0]
			if strings.ContainsAny(slug, " /\\") {
				return fmt.Errorf("slug must be kebab-case with no spaces or slashes")
			}

			content, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("read stdin: %w", err)
			}
			if len(strings.TrimSpace(string(content))) == 0 {
				return fmt.Errorf("no content provided on stdin")
			}

			if title == "" {
				// Derive a readable title from the slug
				title = strings.ReplaceAll(slug, "-", " ")
				title = strings.Title(title) //nolint:staticcheck
			}

			v, err := mustVault(cfg)
			if err != nil {
				return err
			}

			if dryRun {
				existing, _ := process.LoadTopicPage(v, slug)
				action := "create"
				if existing != "" {
					action = "append to"
				}
				fmt.Printf("[dry-run] Would %s 04-knowledge/%s.md\n", action, slug)
				fmt.Printf("  Title: %s\n", title)
				fmt.Println("  Content preview:")
				preview := strings.TrimSpace(string(content))
				if len(preview) > 300 {
					preview = preview[:300] + "..."
				}
				for _, line := range strings.Split(preview, "\n") {
					fmt.Printf("    %s\n", line)
				}
				return nil
			}

			existing, _ := process.LoadTopicPage(v, slug)
			verb := "Created"
			if existing != "" {
				verb = "Updated"
			}

			if err := process.AppendToTopicPage(v, slug, title, string(content)); err != nil {
				return err
			}

			fmt.Printf("%s 04-knowledge/%s.md\n", verb, slug)
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Human-readable title for the topic page (derived from slug if omitted)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be written without making changes")
	return cmd
}

// pkm knowledge update-index <slug> --description "one-line desc"
// Adds a [[slug]] entry to 04-knowledge/index.md. No-ops if already present.
func newKnowledgeUpdateIndexCommand(cfg *config.Config) *cobra.Command {
	var description string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "update-index <slug>",
		Short: "Add a topic entry to 04-knowledge/index.md",
		Long: `Adds "- [[slug]] — description" to 04-knowledge/index.md.
No-op if the slug is already present. Creates index.md if it does not exist.

Example:
  pkm knowledge update-index ai-agents --description "agentic systems, LLM autonomy, multi-agent patterns"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slug := args[0]

			v, err := mustVault(cfg)
			if err != nil {
				return err
			}

			idx, err := process.LoadIndex(v)
			if err != nil {
				return err
			}

			if idx.HasTopic(slug) {
				fmt.Printf("index.md already contains [[%s]] — no change\n", slug)
				return nil
			}

			if dryRun {
				fmt.Printf("[dry-run] Would add [[%s]] — %s to 04-knowledge/index.md\n", slug, description)
				return nil
			}

			if err := process.AddTopicToIndex(v, slug, description); err != nil {
				return err
			}
			fmt.Printf("Added [[%s]] to 04-knowledge/index.md\n", slug)
			return nil
		},
	}

	cmd.Flags().StringVar(&description, "description", "", "One-line description for the index entry")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be written without making changes")
	return cmd
}

// pkm knowledge append-log --note <file> --action <action> --filed-to <folder>
//
//	[--updated slug1,slug2] [--created slug3]
func newKnowledgeAppendLogCommand(cfg *config.Config) *cobra.Command {
	var (
		noteFile string
		action   string
		filedTo  string
		updated  string
		created  string
		dryRun   bool
	)

	cmd := &cobra.Command{
		Use:   "append-log",
		Short: "Append a processing entry to 04-knowledge/log.md",
		Long: `Records what happened to a note during a distill session.
Appends to 04-knowledge/log.md (creates it if absent).

Example:
  pkm knowledge append-log \
    --note 2026-04-03-readwise-article.md \
    --action "distill + file" \
    --filed-to resources \
    --updated software-engineering-philosophy \
    --created llm-patterns`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if noteFile == "" {
				return fmt.Errorf("--note is required")
			}
			if action == "" {
				return fmt.Errorf("--action is required")
			}

			entry := process.LogEntry{
				NoteFilename: noteFile,
				Action:       action,
				FiledTo:      filedTo,
			}
			if updated != "" {
				entry.Updated = splitSlugs(updated)
			}
			if created != "" {
				entry.Created = splitSlugs(created)
			}

			today := time.Now().Format("2006-01-02")

			if dryRun {
				fmt.Printf("[dry-run] Would append to 04-knowledge/log.md:\n")
				fmt.Printf("  Date:     %s\n", today)
				fmt.Printf("  Note:     %s\n", noteFile)
				fmt.Printf("  Action:   %s\n", action)
				if filedTo != "" {
					fmt.Printf("  Filed to: %s/\n", filedTo)
				}
				if len(entry.Updated) > 0 {
					fmt.Printf("  Updated:  %s\n", strings.Join(entry.Updated, ", "))
				}
				if len(entry.Created) > 0 {
					fmt.Printf("  Created:  %s\n", strings.Join(entry.Created, ", "))
				}
				return nil
			}

			v, err := mustVault(cfg)
			if err != nil {
				return err
			}

			if err := process.AppendToLog(v, today, []process.LogEntry{entry}); err != nil {
				return err
			}
			fmt.Printf("Appended entry for %s to 04-knowledge/log.md\n", noteFile)
			return nil
		},
	}

	cmd.Flags().StringVar(&noteFile, "note", "", "Basename of the processed note (required)")
	cmd.Flags().StringVar(&action, "action", "", `Action taken, e.g. "distill + file" (required)`)
	cmd.Flags().StringVar(&filedTo, "filed-to", "", "Target folder the note was moved to")
	cmd.Flags().StringVar(&updated, "updated", "", "Comma-separated slugs of updated topic pages")
	cmd.Flags().StringVar(&created, "created", "", "Comma-separated slugs of newly created topic pages")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be written without making changes")
	return cmd
}

func splitSlugs(s string) []string {
	var result []string
	for _, part := range strings.Split(s, ",") {
		if slug := strings.TrimSpace(part); slug != "" {
			result = append(result, slug)
		}
	}
	return result
}
