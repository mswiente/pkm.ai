package daily

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/mswiente/pkm.ai/internal/config"
	"github.com/mswiente/pkm.ai/internal/templates"
	"github.com/mswiente/pkm.ai/internal/vault"
)

// Options holds all parameters for the daily create command.
type Options struct {
	Date string // YYYY-MM-DD; empty means today
	Open bool   // open with configured editor after creation
}

// Run creates a daily note. If the note for the given date already exists,
// it prints a message and exits without error.
func Run(v *vault.Vault, r *templates.Renderer, cfg *config.Config, opts Options) error {
	var t time.Time
	if opts.Date != "" {
		var err error
		t, err = time.Parse("2006-01-02", opts.Date)
		if err != nil {
			return fmt.Errorf("invalid date %q: use YYYY-MM-DD format", opts.Date)
		}
	} else {
		t = v.NowInTZ()
	}

	date := t.Format("2006-01-02")
	filename := date + ".md"
	path := v.DailyPath(filename)

	// Idempotent: skip if already exists
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Daily note for %s already exists: %s\n", date, path)
		if opts.Open {
			openEditor(cfg.Editor, path)
		}
		return nil
	}

	data := templates.TemplateData{
		Title: date,
		Date:  date,
	}

	content, err := r.Render("daily", data)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	if err := v.WriteNote(path, content); err != nil {
		return err
	}

	fmt.Println(path)

	if opts.Open {
		openEditor(cfg.Editor, path)
	}

	return nil
}

func openEditor(editor, path string) {
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Start() // fire and forget for --open
}
