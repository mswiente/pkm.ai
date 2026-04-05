package decision

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mswiente/pkm.ai/internal/config"
	"github.com/mswiente/pkm.ai/internal/slug"
	"github.com/mswiente/pkm.ai/internal/templates"
	"github.com/mswiente/pkm.ai/internal/vault"
)

var validStatuses = map[string]string{
	"draft":      "Draft",
	"accepted":   "Accepted",
	"superseded": "Superseded",
}

// Options holds all parameters for the decision create command.
type Options struct {
	Title     string
	Project   string
	Status    string // draft|accepted|superseded; default "draft"
	FromStdin bool   // read body from stdin
}

// Run creates a decision note in the inbox.
func Run(v *vault.Vault, r *templates.Renderer, cfg *config.Config, opts Options) error {
	if opts.Title == "" {
		return fmt.Errorf("--title is required")
	}

	status := opts.Status
	if status == "" {
		status = "draft"
	}
	statusLabel, ok := validStatuses[status]
	if !ok {
		return fmt.Errorf("invalid status %q: must be draft, accepted, or superseded", status)
	}

	var body string
	if opts.FromStdin {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("read stdin: %w", err)
		}
		body = strings.TrimSpace(string(data))
	}

	now := v.NowInTZ()
	date := now.Format("2006-01-02")

	data := templates.TemplateData{
		Title:       strings.TrimSpace(opts.Title),
		Date:        date,
		Project:     opts.Project,
		Status:      status,
		StatusLabel: statusLabel,
		Body:        body,
	}

	content, err := r.Render("decision", data)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	filename := fmt.Sprintf("%s-%s-%s.md",
		now.Format("2006-01-02-1504"),
		cfg.DefaultSource,
		slug.FromTitle(opts.Title),
	)
	path := v.InboxPath(filename)

	if err := v.WriteNote(path, content); err != nil {
		return err
	}

	fmt.Println(path)
	return nil
}
