package readwise

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mswiente/pkm.ai/internal/config"
	"github.com/mswiente/pkm.ai/internal/frontmatter"
	"github.com/mswiente/pkm.ai/internal/slug"
	"github.com/mswiente/pkm.ai/internal/vault"
)

// Options controls sync behavior.
type Options struct {
	DryRun bool
	Since  *time.Time // overrides state file when set
	Limit  int        // max articles to sync; 0 = no limit
}

// Result summarises a sync run.
type Result struct {
	Synced  int
	Skipped int // already existed
}

// Run fetches Readwise Reader articles and writes them as inbox notes.
// On success, saves the updated sync state.
func Run(v *vault.Vault, cfg *config.Config, opts Options) (Result, error) {
	token := cfg.ReadwiseToken
	if token == "" {
		token = os.Getenv("READWISE_TOKEN")
	}
	if token == "" {
		return Result{}, fmt.Errorf("no Readwise token found\nRun: pkm sync readwise auth")
	}

	// Determine updatedAfter
	var updatedAfter *time.Time
	if opts.Since != nil {
		updatedAfter = opts.Since
	} else {
		state, err := LoadState(StateFilePath())
		if err != nil {
			return Result{}, fmt.Errorf("load sync state: %w", err)
		}
		updatedAfter = state.LastSyncedAt
	}

	syncedAt := time.Now().UTC()

	client := NewClient(token)
	docs, err := client.ListAll(context.Background(), updatedAfter)
	if err != nil {
		return Result{}, fmt.Errorf("fetch documents: %w", err)
	}

	// Separate top-level articles from highlights (highlights have a parent_id)
	highlights := make(map[string][]Document) // docID → []highlight
	var articles []Document
	for _, d := range docs {
		if d.ParentID != nil && *d.ParentID != "" {
			highlights[*d.ParentID] = append(highlights[*d.ParentID], d)
		} else {
			articles = append(articles, d)
		}
	}

	var result Result
	for _, doc := range articles {
		if opts.Limit > 0 && result.Synced >= opts.Limit {
			break
		}

		filename := buildDocFilename(doc)
		path := v.InboxPath(filename)

		if opts.DryRun {
			fmt.Printf("  [dry-run] %s\n", path)
			result.Synced++
			continue
		}

		content := buildNoteContent(doc, highlights[doc.ID])
		err := v.WriteNote(path, content)
		if err != nil {
			if strings.Contains(err.Error(), "file already exists") {
				result.Skipped++
				continue
			}
			return result, fmt.Errorf("write note %s: %w", filename, err)
		}
		result.Synced++
	}

	// Save state only on non-dry-run success
	if !opts.DryRun {
		state := &State{LastSyncedAt: &syncedAt}
		if err := SaveState(StateFilePath(), state); err != nil {
			return result, fmt.Errorf("save sync state: %w", err)
		}
	}

	return result, nil
}

// buildDocFilename creates a deterministic inbox filename using the article's saved_at date.
func buildDocFilename(doc Document) string {
	t := doc.SavedAt
	if t.IsZero() {
		t = doc.CreatedAt
	}
	if t.IsZero() {
		t = time.Now()
	}
	datePart := t.UTC().Format("2006-01-02-1504")
	slugPart := slug.FromTitle(doc.Title)
	return fmt.Sprintf("%s-readwise-%s.md", datePart, slugPart)
}

// buildNoteContent renders the full markdown note for a Readwise document.
func buildNoteContent(doc Document, docHighlights []Document) []byte {
	savedDate := doc.SavedAt.UTC().Format("2006-01-02")
	if doc.SavedAt.IsZero() {
		savedDate = doc.CreatedAt.UTC().Format("2006-01-02")
	}

	tags := []string{"readwise"}
	for name := range doc.Tags {
		tags = append(tags, name)
	}

	fm := frontmatter.Note{
		Title:   doc.Title,
		Type:    "resource",
		Status:  "inbox",
		Source:  "readwise",
		Created: savedDate,
		Tags:    tags,
	}

	var b strings.Builder
	b.Write(frontmatter.MarshalSimple(fm))
	b.WriteString("\n")

	// Metadata block
	if doc.SourceURL != "" {
		b.WriteString(fmt.Sprintf("> **Source:** %s\n", doc.SourceURL))
	}
	if doc.Author != "" {
		b.WriteString(fmt.Sprintf("> **Author:** %s\n", doc.Author))
	}
	b.WriteString(fmt.Sprintf("> **Saved:** %s\n", savedDate))
	if doc.ReadingTime > 0 {
		b.WriteString(fmt.Sprintf("> **Reading time:** %d min\n", doc.ReadingTime))
	}
	b.WriteString("\n")

	// Summary
	if doc.Summary != "" {
		b.WriteString(doc.Summary)
		b.WriteString("\n\n")
	}

	// Highlights
	if len(docHighlights) > 0 {
		b.WriteString("## Highlights\n\n")
		for _, h := range docHighlights {
			text := strings.TrimSpace(h.Content)
			if text == "" {
				text = strings.TrimSpace(h.Notes)
			}
			if text == "" {
				continue
			}
			// Multi-line highlights get blockquote formatting; single-line get list items
			if strings.Contains(text, "\n") {
				for _, line := range strings.Split(text, "\n") {
					b.WriteString(fmt.Sprintf("> %s\n", line))
				}
				b.WriteString("\n")
			} else {
				b.WriteString(fmt.Sprintf("- %s\n", text))
			}
		}
	}

	return []byte(b.String())
}
