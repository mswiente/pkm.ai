package note

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mswiente/pkm.ai/internal/frontmatter"
	"github.com/mswiente/pkm.ai/internal/vault"
)

// folderToType maps canonical vault folder names to their inferred note type.
var folderToType = map[string]string{
	"02-projects":  "project",
	"04-knowledge": "knowledge",
	"05-resources": "resource",
	"06-decisions": "decision",
}

// MoveOptions holds all parameters for the note move command.
type MoveOptions struct {
	Filename string // basename of the source file
	Folder   string // target folder as typed by user
	Type     string // override inferred type; empty = infer
	Status   string // override inferred status; empty = infer
	DryRun   bool
}

// Move moves a note to the target folder, updating its frontmatter.
func Move(v *vault.Vault, opts MoveOptions) error {
	// Resolve source
	src, err := v.FindNote(opts.Filename)
	if err != nil {
		return err
	}

	// Resolve target directory
	targetDir, err := v.FolderPath(opts.Folder)
	if err != nil {
		return err
	}

	// Read and parse
	content, err := v.ReadNote(src)
	if err != nil {
		return fmt.Errorf("read note: %w", err)
	}
	note, body, err := frontmatter.Parse(content)
	if err != nil {
		return fmt.Errorf("parse frontmatter: %w", err)
	}

	// Resolve target folder canonical name for inference
	canonicalFolder, _ := vault.FolderName(opts.Folder)

	// Infer type
	newType := note.Type
	if opts.Type != "" {
		newType = opts.Type
	} else if inferred, ok := folderToType[canonicalFolder]; ok {
		newType = inferred
	}

	// Infer status
	newStatus := note.Status
	if opts.Status != "" {
		newStatus = opts.Status
	} else if canonicalFolder == "09-archive" {
		newStatus = "archived"
	} else if note.Status == "inbox" {
		newStatus = "draft"
	}

	// Updated date
	today := v.NowInTZ().Format("2006-01-02")

	// Build relative paths for display
	vaultPath := filepath.Dir(filepath.Dir(src)) // go up past the subfolder
	srcRel := relativeToVault(src, v)
	dst := filepath.Join(targetDir, filepath.Base(src))
	dstRel := relativeToVault(dst, v)
	_ = vaultPath

	// Print what will happen
	prefix := ""
	if opts.DryRun {
		prefix = "[dry-run] "
	}
	fmt.Printf("%sMoving: %s\n", prefix, srcRel)
	fmt.Printf("%s     → %s\n", prefix, dstRel)
	if newType != note.Type {
		fmt.Printf("%s  type:    %s → %s\n", prefix, note.Type, newType)
	}
	if newStatus != note.Status {
		fmt.Printf("%s  status:  %s → %s\n", prefix, note.Status, newStatus)
	}
	fmt.Printf("%s  updated: (set to %s)\n", prefix, today)

	if opts.DryRun {
		fmt.Println()
		fmt.Println("(no files were moved)")
		return nil
	}

	// Update frontmatter
	note.Type = newType
	note.Status = newStatus
	note.Updated = today

	// Reassemble content
	newContent := append(frontmatter.MarshalSimple(note), '\n')
	newContent = append(newContent, []byte(strings.TrimLeft(body, "\n"))...)

	// Write updated frontmatter to source, then rename
	if err := v.OverwriteNote(src, newContent); err != nil {
		return fmt.Errorf("update frontmatter: %w", err)
	}
	if err := v.MoveNote(src, dst); err != nil {
		// Attempt to restore original content on failure
		_ = v.OverwriteNote(src, content)
		return fmt.Errorf("move file: %w", err)
	}

	return nil
}

// relativeToVault returns a path relative to the vault root for display.
func relativeToVault(path string, v *vault.Vault) string {
	// Get the vault path by looking at InboxDir parent
	inboxDir := v.InboxDir()
	vaultRoot := filepath.Dir(inboxDir)
	rel, err := filepath.Rel(vaultRoot, path)
	if err != nil {
		return path
	}
	return rel
}
