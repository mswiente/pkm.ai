package meeting

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mswiente/pkm.ai/internal/config"
	"github.com/mswiente/pkm.ai/internal/frontmatter"
	"github.com/mswiente/pkm.ai/internal/slug"
	"github.com/mswiente/pkm.ai/internal/templates"
	"github.com/mswiente/pkm.ai/internal/vault"
)

// Options holds all parameters for the meeting create command.
type Options struct {
	Title        string
	Date         string   // YYYY-MM-DD; empty means today
	Project      string
	Participants []string
}

// Run creates a meeting note in the inbox.
func Run(v *vault.Vault, r *templates.Renderer, cfg *config.Config, opts Options) error {
	title, err := resolveTitle(opts.Title)
	if err != nil {
		return err
	}

	var t time.Time
	if opts.Date != "" {
		t, err = time.Parse("2006-01-02", opts.Date)
		if err != nil {
			return fmt.Errorf("invalid date %q: use YYYY-MM-DD format", opts.Date)
		}
	} else {
		t = v.NowInTZ()
	}

	date := t.Format("2006-01-02")

	data := templates.TemplateData{
		Title:            title,
		Date:             date,
		Source:           cfg.DefaultSource,
		ParticipantsYAML: frontmatter.FormatParticipants(opts.Participants),
		Project:          opts.Project,
	}

	content, err := r.Render("meeting", data)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	now := v.NowInTZ()
	filename := fmt.Sprintf("%s-%s-%s.md",
		now.Format("2006-01-02-1504"),
		cfg.DefaultSource,
		slug.FromTitle(title),
	)
	path := v.InboxPath(filename)

	if err := v.WriteNote(path, content); err != nil {
		return err
	}

	fmt.Println(path)
	return nil
}

func resolveTitle(title string) (string, error) {
	if title != "" {
		return strings.TrimSpace(title), nil
	}
	// If stdin is a terminal, prompt for title
	if isTerminal() {
		fmt.Print("Meeting title: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			t := strings.TrimSpace(scanner.Text())
			if t != "" {
				return t, nil
			}
		}
		return "", fmt.Errorf("title is required")
	}
	return "", fmt.Errorf("--title is required in non-interactive mode")
}

func isTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
