package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mswiente/pkm.ai/internal/vault"
)

const (
	knowledgeSubdir = "04-knowledge"
	indexFilename   = "index.md"
	logFilename     = "log.md"
)

// KnowledgeIndex is the parsed content of 04-knowledge/index.md.
type KnowledgeIndex struct {
	Topics []TopicEntry
	raw    string
}

// TopicEntry is one entry in the index: [[slug]] — description.
type TopicEntry struct {
	Slug        string
	Description string
}

// String renders the index as a plain-text list for use in prompts.
func (idx *KnowledgeIndex) String() string {
	if idx.raw != "" {
		return idx.raw
	}
	return "(empty)"
}

// HasTopic returns true if the given slug already exists in the index.
func (idx *KnowledgeIndex) HasTopic(slug string) bool {
	for _, t := range idx.Topics {
		if t.Slug == slug {
			return true
		}
	}
	return false
}

// LoadIndex reads 04-knowledge/index.md. Returns an empty index if absent.
func LoadIndex(v *vault.Vault) (*KnowledgeIndex, error) {
	path := knowledgeFilePath(v, indexFilename)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &KnowledgeIndex{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read knowledge index: %w", err)
	}
	raw := string(data)
	idx := &KnowledgeIndex{raw: raw}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "- [[") {
			continue
		}
		end := strings.Index(line, "]]")
		if end < 0 {
			continue
		}
		slug := line[4:end]
		desc := ""
		// "—" is a multi-byte em dash; find it robustly
		if idx2 := strings.Index(line, "—"); idx2 >= 0 {
			desc = strings.TrimSpace(line[idx2+len("—"):])
		}
		idx.Topics = append(idx.Topics, TopicEntry{Slug: slug, Description: desc})
	}
	return idx, nil
}

// AddTopicToIndex appends a new [[slug]] entry to 04-knowledge/index.md.
// Creates the file with a header if it does not exist. No-ops if slug already present.
func AddTopicToIndex(v *vault.Vault, slug, description string) error {
	idx, err := LoadIndex(v)
	if err != nil {
		return err
	}
	if idx.HasTopic(slug) {
		return nil
	}

	dir := knowledgeDir(v)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create knowledge dir: %w", err)
	}

	path := knowledgeFilePath(v, indexFilename)
	data, readErr := os.ReadFile(path)
	var content string
	if os.IsNotExist(readErr) {
		content = "# Knowledge Index\n\nTopic pages in `04-knowledge/`. Updated by `pkm process inbox --distill`.\n\n"
	} else if readErr != nil {
		return fmt.Errorf("read index: %w", readErr)
	} else {
		content = string(data)
	}

	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += fmt.Sprintf("- [[%s]] — %s\n", slug, description)

	return os.WriteFile(path, []byte(content), 0o644)
}

// LoadTopicPage reads 04-knowledge/<slug>.md. Returns empty string if absent.
func LoadTopicPage(v *vault.Vault, slug string) (string, error) {
	path := knowledgeFilePath(v, slug+".md")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read topic %s: %w", slug, err)
	}
	return string(data), nil
}

// AppendToTopicPage appends contentToMerge to 04-knowledge/<slug>.md.
// Creates the page with frontmatter if it does not exist.
func AppendToTopicPage(v *vault.Vault, slug, title, contentToMerge string) error {
	dir := knowledgeDir(v)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create knowledge dir: %w", err)
	}

	existing, err := LoadTopicPage(v, slug)
	if err != nil {
		return err
	}

	var content string
	if existing == "" {
		today := time.Now().Format("2006-01-02")
		content = fmt.Sprintf(
			"---\ntitle: %s\ntype: knowledge\nstatus: evergreen\ncreated: %s\ntags: []\n---\n\n# %s\n\n",
			title, today, title,
		)
	} else {
		content = existing
	}

	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += contentToMerge
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	path := knowledgeFilePath(v, slug+".md")
	return os.WriteFile(path, []byte(content), 0o644)
}

// LogEntry records what happened to one inbox note during a distill run.
type LogEntry struct {
	NoteFilename string
	Action       string
	FiledTo      string
	Updated      []string // existing topic slugs updated
	Created      []string // new topic slugs created
}

// AppendToLog appends a dated batch of log entries to 04-knowledge/log.md.
func AppendToLog(v *vault.Vault, date string, entries []LogEntry) error {
	dir := knowledgeDir(v)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create knowledge dir: %w", err)
	}

	path := knowledgeFilePath(v, logFilename)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("open log: %w", err)
	}
	defer f.Close()

	fi, _ := f.Stat()
	if fi.Size() == 0 {
		fmt.Fprint(f, "# Knowledge Processing Log\n\nChronological audit trail of `pkm process inbox --distill` runs.\n")
	}

	fmt.Fprintf(f, "\n## %s\n\n", date)
	for _, e := range entries {
		slug := strings.TrimSuffix(e.NoteFilename, ".md")
		fmt.Fprintf(f, "### [[%s]]\n", slug)
		fmt.Fprintf(f, "- **Action:** %s\n", e.Action)
		if e.FiledTo != "" {
			fmt.Fprintf(f, "- **Filed to:** `%s/`\n", e.FiledTo)
		}
		if len(e.Updated) > 0 {
			fmt.Fprintf(f, "- **Updated:** %s\n", wikilinks(e.Updated))
		}
		if len(e.Created) > 0 {
			fmt.Fprintf(f, "- **Created:** %s\n", wikilinks(e.Created))
		}
		fmt.Fprintln(f)
	}

	return nil
}

func wikilinks(slugs []string) string {
	parts := make([]string, len(slugs))
	for i, s := range slugs {
		parts[i] = "[[" + s + "]]"
	}
	return strings.Join(parts, ", ")
}

func knowledgeDir(v *vault.Vault) string {
	return filepath.Join(filepath.Dir(v.InboxDir()), knowledgeSubdir)
}

func knowledgeFilePath(v *vault.Vault, filename string) string {
	return filepath.Join(knowledgeDir(v), filename)
}
