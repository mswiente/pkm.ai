package skill

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mswiente/pkm.ai/skills"
)

// Skill represents an available skill definition.
type Skill struct {
	Name        string
	Description string
}

// List returns all skills embedded in the binary.
func List() ([]Skill, error) {
	entries, err := fs.ReadDir(skills.FS, ".")
	if err != nil {
		return nil, err
	}
	var result []Skill
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".md")
		data, _ := skills.FS.ReadFile(e.Name())
		result = append(result, Skill{
			Name:        name,
			Description: extractDescription(string(data)),
		})
	}
	return result, nil
}

// Install copies the named skill to ~/.claude/commands/<name>.md.
// If name is empty, all skills are installed.
func Install(name string) error {
	targetDir, err := commandsDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("create commands dir: %w", err)
	}

	available, err := List()
	if err != nil {
		return err
	}

	installed := 0
	for _, s := range available {
		if name != "" && s.Name != name {
			continue
		}
		data, err := skills.FS.ReadFile(s.Name + ".md")
		if err != nil {
			return fmt.Errorf("read skill %s: %w", s.Name, err)
		}
		dst := filepath.Join(targetDir, s.Name+".md")
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", dst, err)
		}
		fmt.Printf("Installed: /%s\n", s.Name)
		fmt.Printf("  %s\n", dst)
		installed++
	}

	if installed == 0 {
		if name != "" {
			return fmt.Errorf("skill %q not found", name)
		}
		fmt.Println("No skills to install.")
	}
	return nil
}

// IsInstalled reports whether the named skill is already installed.
func IsInstalled(name string) bool {
	dir, err := commandsDir()
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(dir, name+".md"))
	return err == nil
}

func commandsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, ".claude", "commands"), nil
}

func extractDescription(content string) string {
	if !strings.HasPrefix(content, "---\n") {
		return ""
	}
	rest := content[4:]
	end := strings.Index(rest, "\n---")
	if end == -1 {
		return ""
	}
	for _, line := range strings.Split(rest[:end], "\n") {
		if strings.HasPrefix(line, "description:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "description:"))
		}
	}
	return ""
}
