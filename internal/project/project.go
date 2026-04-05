package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mswiente/pkm.ai/internal/vault"
)

const projectsSubdir = "02-projects"

// UpdateOptions holds all parameters for creating or updating a project note.
type UpdateOptions struct {
	Slug          string
	Title         string
	ProjectStatus string // active | on-hold | archived (frontmatter status field)
	Intent        string // replaces the ## Intent section when non-empty
	CurrentStatus string // replaces the ## Current Status section when non-empty
	NextSteps     string // replaces the ## Next Steps section when non-empty
	PlanContent   string // appended to ## Plan History with a date-stamped heading
	PlanHeading   string // heading text for the plan history entry (default: today's date)
	DryRun        bool
}

// UpdateResult describes what happened.
type UpdateResult struct {
	Path    string
	IsNew   bool
	Patched []string // names of sections that were updated
}

// Update creates or patches a project note at 02-projects/<slug>.md.
// Returns whether the note was newly created and a description of changes.
func Update(v *vault.Vault, opts UpdateOptions) (*UpdateResult, error) {
	dir := projectsDir(v)
	path := filepath.Join(dir, opts.Slug+".md")

	existing, readErr := os.ReadFile(path)
	isNew := os.IsNotExist(readErr)

	var content string
	var patched []string

	today := time.Now().Format("2006-01-02")

	if isNew {
		content = buildNote(opts, today)
		patched = []string{"(created)"}
	} else {
		content, patched = patchNote(string(existing), opts, today)
	}

	if opts.DryRun {
		return &UpdateResult{Path: path, IsNew: isNew, Patched: patched}, nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create projects dir: %w", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return nil, fmt.Errorf("write project note: %w", err)
	}
	return &UpdateResult{Path: path, IsNew: isNew, Patched: patched}, nil
}

// Load reads the raw content of a project note. Returns "" if not found.
func Load(v *vault.Vault, slug string) (string, error) {
	path := filepath.Join(projectsDir(v), slug+".md")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read project %s: %w", slug, err)
	}
	return string(data), nil
}

// List returns all project slugs in 02-projects/.
func List(v *vault.Vault) ([]string, error) {
	entries, err := os.ReadDir(projectsDir(v))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read projects dir: %w", err)
	}
	var slugs []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			slugs = append(slugs, strings.TrimSuffix(e.Name(), ".md"))
		}
	}
	return slugs, nil
}

// Exists reports whether a project note exists for the given slug.
func Exists(v *vault.Vault, slug string) bool {
	_, err := os.Stat(filepath.Join(projectsDir(v), slug+".md"))
	return err == nil
}

// buildNote creates a complete new project note from opts.
func buildNote(opts UpdateOptions, today string) string {
	title := opts.Title
	if title == "" {
		title = titleFromSlug(opts.Slug)
	}
	status := opts.ProjectStatus
	if status == "" {
		status = "active"
	}
	intent := blankIfEmpty(opts.Intent, "(to be defined)")
	curStatus := blankIfEmpty(opts.CurrentStatus, "(to be described)")
	nextSteps := blankIfEmpty(opts.NextSteps, "(to be defined)")

	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("title: %s\n", title))
	b.WriteString("type: project\n")
	b.WriteString(fmt.Sprintf("status: %s\n", status))
	b.WriteString("source: claude-code\n")
	b.WriteString(fmt.Sprintf("created: %s\n", today))
	b.WriteString(fmt.Sprintf("updated: %s\n", today))
	b.WriteString("tags: [project]\n")
	b.WriteString("---\n\n")
	b.WriteString("## Intent\n\n")
	b.WriteString(strings.TrimSpace(intent) + "\n\n")
	b.WriteString("## Current Status\n\n")
	b.WriteString(strings.TrimSpace(curStatus) + "\n\n")
	b.WriteString("## Next Steps\n\n")
	b.WriteString(strings.TrimSpace(nextSteps) + "\n\n")

	if opts.PlanContent != "" {
		b.WriteString("## Plan History\n\n")
		b.WriteString(planEntry(opts.PlanHeading, opts.PlanContent, today))
	}

	return b.String()
}

// patchNote applies opts to an existing note, returning the new content and
// a list of section names that were changed.
func patchNote(existing string, opts UpdateOptions, today string) (string, []string) {
	var patched []string

	existing = setFrontmatterField(existing, "updated", today)

	if opts.ProjectStatus != "" {
		existing = setFrontmatterField(existing, "status", opts.ProjectStatus)
		patched = append(patched, "status")
	}
	if opts.Intent != "" {
		existing = replaceSection(existing, "## Intent", opts.Intent)
		patched = append(patched, "Intent")
	}
	if opts.CurrentStatus != "" {
		existing = replaceSection(existing, "## Current Status", opts.CurrentStatus)
		patched = append(patched, "Current Status")
	}
	if opts.NextSteps != "" {
		existing = replaceSection(existing, "## Next Steps", opts.NextSteps)
		patched = append(patched, "Next Steps")
	}
	if opts.PlanContent != "" {
		existing = appendToPlanHistory(existing, opts.PlanHeading, opts.PlanContent, today)
		patched = append(patched, "Plan History")
	}

	return existing, patched
}

// replaceSection replaces the body of a ## heading section with newBody.
// If the section doesn't exist it is appended at the end.
func replaceSection(content, heading, newBody string) string {
	idx := strings.Index(content, "\n"+heading+"\n")
	if idx < 0 {
		// Try at very start of content (no leading newline)
		if strings.HasPrefix(content, heading+"\n") {
			idx = -1 // handled below
		} else {
			// Section not found — append
			if !strings.HasSuffix(content, "\n") {
				content += "\n"
			}
			return content + "\n" + heading + "\n\n" + strings.TrimSpace(newBody) + "\n"
		}
	}

	// bodyStart: first char after the heading line's newline
	var bodyStart int
	if idx < 0 {
		bodyStart = len(heading) + 1 // after "## Heading\n"
	} else {
		bodyStart = idx + 1 + len(heading) + 1 // after "\n## Heading\n"
	}

	// Find start of next ## section or end of content
	nextIdx := strings.Index(content[bodyStart:], "\n## ")
	var bodyEnd int
	if nextIdx < 0 {
		bodyEnd = len(content)
	} else {
		bodyEnd = bodyStart + nextIdx + 1 // keep the \n that precedes ##
	}

	return content[:bodyStart] + "\n" + strings.TrimSpace(newBody) + "\n\n" + content[bodyEnd:]
}

// appendToPlanHistory appends a new dated entry to the ## Plan History section.
// Creates the section if absent.
func appendToPlanHistory(content, heading, planContent, today string) string {
	entry := planEntry(heading, planContent, today)
	histHeading := "\n## Plan History\n"

	if !strings.Contains(content, histHeading) {
		// Create section at end
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		return content + "\n## Plan History\n\n" + entry
	}

	// Append at end of the Plan History section (i.e. end of file)
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return content + entry
}

func planEntry(heading, content, today string) string {
	h := heading
	if h == "" {
		h = today
	}
	return fmt.Sprintf("### %s\n\n%s\n\n", h, strings.TrimSpace(content))
}

// setFrontmatterField replaces the value of a YAML frontmatter field in-place.
func setFrontmatterField(content, field, value string) string {
	prefix := field + ": "
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, prefix) {
			lines[i] = prefix + value
			return strings.Join(lines, "\n")
		}
	}
	return content
}

func projectsDir(v *vault.Vault) string {
	return filepath.Join(filepath.Dir(v.InboxDir()), projectsSubdir)
}

func titleFromSlug(slug string) string {
	words := strings.Split(slug, "-")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func blankIfEmpty(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}
