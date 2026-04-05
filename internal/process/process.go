package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mswiente/pkm.ai/internal/frontmatter"
	"github.com/mswiente/pkm.ai/internal/vault"
)

// Options holds all parameters for the process inbox command.
type Options struct {
	File        string // single filename (basename); empty means all
	All         bool   // process all files (same as empty File)
	DryRun      bool
	Apply       bool
	Interactive bool
}

// inboxEntry holds per-file analysis data.
type inboxEntry struct {
	Filename    string
	Path        string
	Note        frontmatter.Note
	BodyPreview string
	FileSize    int64
	ModTime     time.Time
	Issues      []string
}

// Run analyzes inbox files and prints a structured report to stdout.
func Run(v *vault.Vault, opts Options) error {
	var paths []string

	if opts.File != "" {
		path := v.InboxPath(opts.File)
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("file not found in inbox: %s", opts.File)
		}
		paths = []string{path}
	} else {
		var err error
		paths, err = v.ListInbox()
		if err != nil {
			return fmt.Errorf("list inbox: %w", err)
		}
	}

	if len(paths) == 0 {
		fmt.Println("Inbox is empty.")
		return nil
	}

	var entries []inboxEntry
	for _, p := range paths {
		e, err := analyzeFile(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: skip %s: %v\n", filepath.Base(p), err)
			continue
		}
		entries = append(entries, e)
	}

	if opts.Apply {
		fmt.Fprintln(os.Stderr, "note: --apply is not yet implemented in this version")
	}

	printReport(v, entries, opts)
	return nil
}

func analyzeFile(path string) (inboxEntry, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return inboxEntry{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return inboxEntry{}, err
	}

	note, body, err := frontmatter.Parse(data)
	if err != nil {
		// File might not have frontmatter; still include it with a warning
		note = frontmatter.Note{}
		body = string(data)
	}

	preview := bodyPreview(body, 300)
	issues := detectIssues(note, body)

	return inboxEntry{
		Filename:    filepath.Base(path),
		Path:        path,
		Note:        note,
		BodyPreview: preview,
		FileSize:    fi.Size(),
		ModTime:     fi.ModTime(),
		Issues:      issues,
	}, nil
}

func bodyPreview(body string, max int) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return ""
	}

	// Find first non-empty paragraph
	var preview string
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			preview = line
			break
		}
	}

	if preview == "" {
		return ""
	}

	if utf8.RuneCountInString(preview) <= max {
		return preview
	}

	// Truncate at word boundary
	runes := []rune(preview)
	truncated := string(runes[:max])
	if idx := strings.LastIndex(truncated, " "); idx > max/2 {
		truncated = truncated[:idx]
	}
	return truncated + "..."
}

func detectIssues(note frontmatter.Note, body string) []string {
	var issues []string

	if note.Title == "" {
		issues = append(issues, "title: not set")
	}
	if len(note.Tags) == 0 {
		issues = append(issues, "tags: not set")
	}
	if note.TypeHint == "" {
		issues = append(issues, "type_hint: not set (could help with classification)")
	}
	if isBodyEmpty(body) {
		issues = append(issues, "no body content beyond template placeholders")
	}

	return issues
}

func isBodyEmpty(body string) bool {
	body = strings.TrimSpace(body)
	if body == "" {
		return true
	}
	lines := strings.Split(body, "\n")
	substantiveLines := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip headings, empty lines, and bare list markers
		if line == "" || line == "-" || line == "- [ ]" ||
			strings.HasPrefix(line, "#") ||
			(strings.HasPrefix(line, "- ") && len(strings.TrimSpace(line[2:])) == 0) {
			continue
		}
		substantiveLines++
	}
	return substantiveLines == 0
}

func printReport(v *vault.Vault, entries []inboxEntry, opts Options) {
	now := time.Now().Format("2006-01-02 15:04")

	fmt.Println("=== PKM Inbox Report ===")
	fmt.Printf("Generated: %s\n", now)
	fmt.Printf("Inbox: %s\n", v.InboxDir())
	fmt.Printf("Files: %d markdown file(s)\n", len(entries))

	withIssues := 0
	missingTags := 0
	noBody := 0
	statusInbox := 0

	for i, e := range entries {
		fmt.Println()
		fmt.Println("---")
		fmt.Printf("\n[%d/%d] %s\n", i+1, len(entries), e.Filename)
		fmt.Printf("  Created:  %s\n", orDash(e.Note.Created))
		fmt.Printf("  Source:   %s\n", orDash(e.Note.Source))
		fmt.Printf("  Title:    %s\n", orDash(e.Note.Title))
		fmt.Printf("  Type:     %s\n", orDash(e.Note.Type))
		fmt.Printf("  Status:   %s\n", orDash(e.Note.Status))
		if len(e.Note.Tags) > 0 {
			fmt.Printf("  Tags:     %s\n", strings.Join(e.Note.Tags, ", "))
		} else {
			fmt.Printf("  Tags:     (none)\n")
		}
		fmt.Printf("  Size:     %.1f KB\n", float64(e.FileSize)/1024)
		fmt.Printf("  Modified: %s\n", e.ModTime.Format("2006-01-02 15:04"))

		if e.BodyPreview != "" {
			fmt.Println()
			fmt.Println("  Preview:")
			// Indent the preview
			for _, line := range strings.Split(e.BodyPreview, "\n") {
				fmt.Printf("    %s\n", line)
			}
		}

		if len(e.Issues) > 0 {
			fmt.Println()
			fmt.Println("  Issues:")
			for _, issue := range e.Issues {
				fmt.Printf("    - %s\n", issue)
			}
			withIssues++
		}

		// Count stats
		for _, issue := range e.Issues {
			if strings.Contains(issue, "tags") {
				missingTags++
			}
			if strings.Contains(issue, "no body") {
				noBody++
			}
		}
		if e.Note.Status == "inbox" || e.Note.Status == "" {
			statusInbox++
		}
	}

	fmt.Println()
	fmt.Println("---")
	fmt.Println()
	fmt.Println("Summary:")
	fmt.Printf("  Total files:      %d\n", len(entries))
	fmt.Printf("  With issues:      %d\n", withIssues)
	fmt.Printf("  Missing tags:     %d\n", missingTags)
	fmt.Printf("  No body content:  %d\n", noBody)
	fmt.Printf("  Status = inbox:   %d (unprocessed)\n", statusInbox)

	if opts.DryRun {
		fmt.Println()
		fmt.Println("(dry-run: no files were modified)")
	}

	fmt.Println()
	fmt.Println("Suggested next steps for AI processing:")
	fmt.Println("  - Review each file above and propose: type, target folder, tags, title improvement")
	fmt.Println("  - For files with no body: decide whether to enrich or delete")
	fmt.Println("  - Suggested target folders: 04-knowledge/, 02-projects/, 06-decisions/")
}

func orDash(s string) string {
	if s == "" {
		return "(not set)"
	}
	return s
}
