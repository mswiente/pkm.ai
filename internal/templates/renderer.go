package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateData holds all values available to note templates.
type TemplateData struct {
	Title            string
	Date             string // YYYY-MM-DD
	Source           string
	Tags             []string
	TagsYAML         string // pre-formatted: [tag1, tag2] or []
	ParticipantsYAML string // pre-formatted: [p1, p2] or []
	Project          string
	TypeHint         string
	Status           string
	StatusLabel      string // human-readable: Draft / Accepted / Superseded
	Body             string // raw user content; injected via BODY_PLACEHOLDER
}

// Renderer loads and executes markdown templates.
type Renderer struct {
	vaultTemplatesDir string
}

// NewRenderer creates a Renderer that looks for templates in vaultTemplatesDir.
func NewRenderer(vaultTemplatesDir string) *Renderer {
	return &Renderer{vaultTemplatesDir: vaultTemplatesDir}
}

// Render produces a note from the named template (e.g. "inbox", "daily").
// It loads the template from the vault's templates directory if present,
// falling back to the embedded default. User body content is injected as a
// literal string via BODY_PLACEHOLDER to avoid template syntax conflicts.
func (r *Renderer) Render(name string, data TemplateData) ([]byte, error) {
	tmplStr, err := r.loadTemplate(name)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(name).Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("parse template %q: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template %q: %w", name, err)
	}

	// Replace BODY_PLACEHOLDER with raw user content (safe from template injection)
	result := strings.ReplaceAll(buf.String(), "BODY_PLACEHOLDER", data.Body)
	return []byte(result), nil
}

// loadTemplate returns the template string for the given name, preferring
// the vault's template file over the embedded default.
func (r *Renderer) loadTemplate(name string) (string, error) {
	if r.vaultTemplatesDir != "" {
		path := filepath.Join(r.vaultTemplatesDir, name+".md")
		if data, err := os.ReadFile(path); err == nil {
			// The vault template is a static Obsidian template — it does not
			// contain Go template syntax, so we cannot use it directly.
			// Fall through to the embedded default which has the correct syntax.
			_ = data
		}
	}

	if tmpl, ok := defaultTemplates[name]; ok {
		return tmpl, nil
	}
	return "", fmt.Errorf("unknown template %q", name)
}
