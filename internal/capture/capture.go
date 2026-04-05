package capture

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/mswiente/pkm.ai/internal/frontmatter"
	"github.com/mswiente/pkm.ai/internal/slug"
	"github.com/mswiente/pkm.ai/internal/templates"
	"github.com/mswiente/pkm.ai/internal/vault"
)

var validSources = map[string]bool{
	"manual": true, "chatgpt": true, "claude-code": true, "copilot-cli": true, "other": true,
}

// Options holds all parameters for the capture command.
type Options struct {
	Text       string
	Title      string
	Source     string
	Tags       []string
	TypeHint   string
	OpenEditor bool
	Clipboard  bool
}

// Run captures a new inbox note.
func Run(v *vault.Vault, r *templates.Renderer, opts Options) error {
	if opts.Source != "" && !validSources[opts.Source] {
		return fmt.Errorf("invalid source %q: must be one of manual, chatgpt, claude-code, copilot-cli, other", opts.Source)
	}

	body, err := resolveInput(opts)
	if err != nil {
		return err
	}

	title := resolveTitle(opts.Title, body)
	now := v.NowInTZ()
	date := now.Format("2006-01-02")
	source := opts.Source
	if source == "" {
		source = "manual"
	}

	data := templates.TemplateData{
		Title:    title,
		Date:     date,
		Source:   source,
		Tags:     opts.Tags,
		TagsYAML: frontmatter.FormatTags(opts.Tags),
		TypeHint: opts.TypeHint,
		Body:     body,
	}

	content, err := r.Render("inbox", data)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	if opts.OpenEditor {
		content, err = editInEditor(content)
		if err != nil {
			return fmt.Errorf("editor: %w", err)
		}
	}

	filename := buildFilename(now, source, title)
	path := v.InboxPath(filename)

	if err := v.WriteNote(path, content); err != nil {
		return err
	}

	fmt.Println(path)
	return nil
}

func resolveInput(opts Options) (string, error) {
	if opts.Clipboard {
		text, err := clipboard.ReadAll()
		if err != nil {
			return "", fmt.Errorf("read clipboard: %w", err)
		}
		return strings.TrimSpace(text), nil
	}
	if stdinIsPiped() {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("read stdin: %w", err)
		}
		return strings.TrimSpace(string(data)), nil
	}
	return strings.TrimSpace(opts.Text), nil
}

func resolveTitle(title, body string) string {
	if title != "" {
		return title
	}
	if body != "" {
		// Use first non-empty line of body, truncated to 60 chars
		for _, line := range strings.Split(body, "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				// Strip leading markdown heading markers
				line = strings.TrimLeft(line, "#")
				line = strings.TrimSpace(line)
				if len(line) > 60 {
					line = line[:60]
				}
				return line
			}
		}
	}
	return "Untitled Note"
}

func buildFilename(t time.Time, source, title string) string {
	datePart := t.Format("2006-01-02-1504")
	slugPart := slug.FromTitle(title)
	return fmt.Sprintf("%s-%s-%s.md", datePart, source, slugPart)
}

func stdinIsPiped() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) == 0
}

func editInEditor(content []byte) ([]byte, error) {
	f, err := os.CreateTemp("", "pkm-capture-*.md")
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(f.Name())

	if _, err := f.Write(content); err != nil {
		f.Close()
		return nil, err
	}
	f.Close()

	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = "vi"
	}

	args := []string{f.Name()}
	// GUI editors need --wait to block until closed
	if isGUIEditor(editor) {
		args = append([]string{"--wait"}, args...)
	}

	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("editor exited with error: %w", err)
	}

	return os.ReadFile(f.Name())
}

func isGUIEditor(editor string) bool {
	guiEditors := []string{"code", "cursor", "subl", "atom", "zed"}
	base := editor
	// Handle full paths like /usr/local/bin/code
	for i := len(editor) - 1; i >= 0; i-- {
		if editor[i] == '/' {
			base = editor[i+1:]
			break
		}
	}
	for _, g := range guiEditors {
		if base == g {
			return true
		}
	}
	return false
}
