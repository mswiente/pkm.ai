package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/mswiente/pkm.ai/internal/config"
)

// Vault wraps all filesystem interactions with the Obsidian vault.
type Vault struct {
	cfg *config.Config
	loc *time.Location
}

// New creates a Vault. Returns an error if FilenameTimezone is not a valid timezone.
func New(cfg *config.Config) (*Vault, error) {
	loc, err := time.LoadLocation(cfg.FilenameTimezone)
	if err != nil {
		return nil, fmt.Errorf("invalid filename_timezone %q: %w", cfg.FilenameTimezone, err)
	}
	return &Vault{cfg: cfg, loc: loc}, nil
}

// InboxDir returns the absolute path to the inbox folder.
func (v *Vault) InboxDir() string {
	return filepath.Join(v.cfg.VaultPath, v.cfg.InboxPath)
}

// DailyDir returns the absolute path to the daily notes folder.
func (v *Vault) DailyDir() string {
	return filepath.Join(v.cfg.VaultPath, v.cfg.DailyPath)
}

// TemplatesDir returns the absolute path to the vault's templates folder.
func (v *Vault) TemplatesDir() string {
	return filepath.Join(v.cfg.VaultPath, v.cfg.TemplatesPath)
}

// NowInTZ returns time.Now() in the configured filename timezone.
func (v *Vault) NowInTZ() time.Time {
	return time.Now().In(v.loc)
}

// InboxPath returns the full path for a file in the inbox directory.
func (v *Vault) InboxPath(filename string) string {
	return filepath.Join(v.InboxDir(), filename)
}

// DailyPath returns the full path for a file in the daily directory.
func (v *Vault) DailyPath(filename string) string {
	return filepath.Join(v.DailyDir(), filename)
}

// EnsureDir creates the directory (and parents) if it does not exist.
func (v *Vault) EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

// WriteNote writes content to path, creating parent dirs as needed.
// Returns an error if the file already exists.
func (v *Vault) WriteNote(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("file already exists: %s", path)
		}
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	_, err = f.Write(content)
	return err
}

// ReadNote reads and returns the content of a note at path.
func (v *Vault) ReadNote(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// folderAliases maps user-supplied folder identifiers to canonical vault subfolder names.
var folderAliases = map[string]string{
	"1": "01-daily", "01-daily": "01-daily",
	"2": "02-projects", "02-projects": "02-projects", "projects": "02-projects",
	"3": "03-areas", "03-areas": "03-areas", "areas": "03-areas",
	"4": "04-knowledge", "04-knowledge": "04-knowledge", "knowledge": "04-knowledge",
	"5": "05-resources", "05-resources": "05-resources", "resources": "05-resources",
	"6": "06-decisions", "06-decisions": "06-decisions", "decisions": "06-decisions",
	"9": "09-archive", "09-archive": "09-archive", "archive": "09-archive",
}

// FolderPath resolves a user-supplied folder identifier to its absolute vault path.
// Accepts full names ("04-knowledge"), shorthand ("knowledge"), or numbers ("4").
func (v *Vault) FolderPath(folder string) (string, error) {
	canonical, ok := folderAliases[folder]
	if !ok {
		return "", fmt.Errorf("unknown folder %q\nValid options: projects(2), areas(3), knowledge(4), resources(5), decisions(6), archive(9)", folder)
	}
	return filepath.Join(v.cfg.VaultPath, canonical), nil
}

// FolderName returns the canonical subfolder name for a folder identifier,
// without the vault path prefix.
func FolderName(folder string) (string, bool) {
	canonical, ok := folderAliases[folder]
	return canonical, ok
}

// MoveNote moves src to dst atomically using os.Rename.
// Creates the destination directory if needed.
// Returns an error if dst already exists.
func (v *Vault) MoveNote(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("destination already exists: %s", dst)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("create destination directory: %w", err)
	}
	return os.Rename(src, dst)
}

// OverwriteNote writes content to an existing file, replacing its contents.
func (v *Vault) OverwriteNote(path string, content []byte) error {
	return os.WriteFile(path, content, 0o644)
}

// FindNote searches for a file by basename, checking inbox first, then all
// vault subdirectories. Returns the full path or an error if not found.
func (v *Vault) FindNote(filename string) (string, error) {
	// Check inbox first
	inboxPath := v.InboxPath(filename)
	if _, err := os.Stat(inboxPath); err == nil {
		return inboxPath, nil
	}

	// Walk the vault
	var found string
	err := filepath.WalkDir(v.cfg.VaultPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}
		if !d.IsDir() && d.Name() == filename {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("search vault: %w", err)
	}
	if found == "" {
		return "", fmt.Errorf("file not found: %s", filename)
	}
	return found, nil
}

// ListInbox returns all .md files in the inbox directory, sorted alphabetically.
func (v *Vault) ListInbox() ([]string, error) {
	entries, err := os.ReadDir(v.InboxDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read inbox: %w", err)
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".md" {
			files = append(files, filepath.Join(v.InboxDir(), e.Name()))
		}
	}
	sort.Strings(files)
	return files, nil
}
