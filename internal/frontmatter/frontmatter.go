package frontmatter

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Note holds all supported frontmatter fields across all note types.
// Fields are ordered to match the spec's required output order.
type Note struct {
	Title        string   `yaml:"title"`
	Type         string   `yaml:"type"`
	Status       string   `yaml:"status"`
	Source       string   `yaml:"source,omitempty"`
	Created      string   `yaml:"created"`
	Updated      string   `yaml:"updated,omitempty"`
	Tags         []string `yaml:"tags,omitempty"`
	TypeHint     string   `yaml:"type_hint,omitempty"`
	// Meeting-specific
	Date         string   `yaml:"date,omitempty"`
	Participants []string `yaml:"participants,omitempty"`
	Project      string   `yaml:"project,omitempty"`
	// Decision-specific
	DecisionDate string   `yaml:"decision_date,omitempty"`
	RelatedNotes []string `yaml:"related_notes,omitempty"`
}

// Parse extracts and unmarshals the YAML frontmatter from markdown content.
// The content must start with "---\n". Returns the parsed Note, the body
// (content after the closing "---"), and any error.
func Parse(content []byte) (Note, string, error) {
	s := string(content)
	if !strings.HasPrefix(s, "---\n") {
		return Note{}, s, fmt.Errorf("no frontmatter found (content must start with ---)")
	}

	// Find closing ---
	rest := s[4:]
	end := strings.Index(rest, "\n---")
	if end == -1 {
		return Note{}, s, fmt.Errorf("frontmatter not closed (no closing ---)")
	}

	fmStr := rest[:end]
	body := rest[end+4:] // skip "\n---"
	// Skip optional newline after closing ---
	body = strings.TrimPrefix(body, "\n")

	var n Note
	if err := yaml.Unmarshal([]byte(fmStr), &n); err != nil {
		return Note{}, body, fmt.Errorf("parse frontmatter: %w", err)
	}

	return n, body, nil
}

// FormatTags formats a []string as a YAML inline sequence: [tag1, tag2] or [].
func FormatTags(tags []string) string {
	if len(tags) == 0 {
		return "[]"
	}
	return "[" + strings.Join(tags, ", ") + "]"
}

// FormatParticipants formats participants the same way as tags.
func FormatParticipants(participants []string) string {
	return FormatTags(participants)
}

// MarshalSimple serializes a subset of Note fields into a YAML frontmatter block,
// preserving the spec-required field order. It uses string formatting (not
// yaml.Marshal) so the order is guaranteed and matches the templates exactly.
func MarshalSimple(n Note) []byte {
	var b bytes.Buffer
	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("title: %s\n", n.Title))
	b.WriteString(fmt.Sprintf("type: %s\n", n.Type))
	b.WriteString(fmt.Sprintf("status: %s\n", n.Status))
	if n.Source != "" {
		b.WriteString(fmt.Sprintf("source: %s\n", n.Source))
	}
	b.WriteString(fmt.Sprintf("created: %s\n", n.Created))
	if n.Updated != "" {
		b.WriteString(fmt.Sprintf("updated: %s\n", n.Updated))
	}
	b.WriteString(fmt.Sprintf("tags: %s\n", FormatTags(n.Tags)))
	if n.TypeHint != "" {
		b.WriteString(fmt.Sprintf("type_hint: %s\n", n.TypeHint))
	}
	if n.Date != "" {
		b.WriteString(fmt.Sprintf("date: %s\n", n.Date))
	}
	if len(n.Participants) > 0 {
		b.WriteString(fmt.Sprintf("participants: %s\n", FormatParticipants(n.Participants)))
	}
	if n.Project != "" {
		b.WriteString(fmt.Sprintf("project: %s\n", n.Project))
	}
	if n.DecisionDate != "" {
		b.WriteString(fmt.Sprintf("decision_date: %s\n", n.DecisionDate))
	}
	if len(n.RelatedNotes) > 0 {
		b.WriteString(fmt.Sprintf("related_notes: %s\n", FormatTags(n.RelatedNotes)))
	}
	b.WriteString("---\n")
	return b.Bytes()
}
