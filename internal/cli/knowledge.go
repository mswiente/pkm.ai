package cli

import (
	"encoding/json"
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
		noteFile  string
		action    string
		filedTo   string
		updated   string
		created   string
		fromStdin bool
		dryRun    bool
	)

	cmd := &cobra.Command{
		Use:   "append-log",
		Short: "Append a processing entry to 04-knowledge/log.md",
		Long: `Records what happened to one or more notes during a distill session.
Appends to 04-knowledge/log.md (creates it if absent). All entries from a
single invocation share one "## <date>" heading.

Single-entry example:
  pkm knowledge append-log \
    --note 2026-04-03-readwise-article.md \
    --action "distill + file" \
    --filed-to resources \
    --updated software-engineering-philosophy \
    --created llm-patterns

Batch example (one heading for all entries):
  printf '%s\n' '[
    {"note": "a.md", "action": "distill + file", "filed_to": "resources", "updated": ["ai-agents"]},
    {"note": "b.md", "action": "file-only", "filed_to": "resources"},
    {"note": "c.md", "action": "skip", "filed_to": "archive"}
  ]' | pkm knowledge append-log --from-stdin`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var entries []process.LogEntry

			if fromStdin {
				data, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("read stdin: %w", err)
				}
				if err := json.Unmarshal(data, &entries); err != nil {
					return fmt.Errorf("parse JSON entries from stdin: %w", err)
				}
				if len(entries) == 0 {
					return fmt.Errorf("no entries provided on stdin")
				}
				for i, e := range entries {
					if e.NoteFilename == "" {
						return fmt.Errorf("entry %d: \"note\" is required", i)
					}
					if e.Action == "" {
						return fmt.Errorf("entry %d: \"action\" is required", i)
					}
				}
			} else {
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
				entries = []process.LogEntry{entry}
			}

			today := time.Now().Format("2006-01-02")

			if dryRun {
				fmt.Printf("[dry-run] Would append to 04-knowledge/log.md:\n")
				fmt.Printf("  Date: %s\n", today)
				for _, e := range entries {
					fmt.Printf("  - Note:     %s\n", e.NoteFilename)
					fmt.Printf("    Action:   %s\n", e.Action)
					if e.FiledTo != "" {
						fmt.Printf("    Filed to: %s/\n", e.FiledTo)
					}
					if len(e.Updated) > 0 {
						fmt.Printf("    Updated:  %s\n", strings.Join(e.Updated, ", "))
					}
					if len(e.Created) > 0 {
						fmt.Printf("    Created:  %s\n", strings.Join(e.Created, ", "))
					}
				}
				return nil
			}

			v, err := mustVault(cfg)
			if err != nil {
				return err
			}

			if err := process.AppendToLog(v, today, entries); err != nil {
				return err
			}
			fmt.Printf("Appended %d entr%s to 04-knowledge/log.md\n", len(entries), plural(len(entries)))
			return nil
		},
	}

	cmd.Flags().StringVar(&noteFile, "note", "", "Basename of the processed note")
	cmd.Flags().StringVar(&action, "action", "", `Action taken, e.g. "distill + file"`)
	cmd.Flags().StringVar(&filedTo, "filed-to", "", "Target folder the note was moved to")
	cmd.Flags().StringVar(&updated, "updated", "", "Comma-separated slugs of updated topic pages")
	cmd.Flags().StringVar(&created, "created", "", "Comma-separated slugs of newly created topic pages")
	cmd.Flags().BoolVar(&fromStdin, "from-stdin", false, "Read a JSON array of log entries from stdin, written under one shared date heading")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be written without making changes")
	return cmd
}

func plural(n int) string {
	if n == 1 {
		return "y"
	}
	return "ies"
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
