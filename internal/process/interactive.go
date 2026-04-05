package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"

	"github.com/mswiente/pkm.ai/internal/note"
	"github.com/mswiente/pkm.ai/internal/vault"
)

// action maps a keypress to a vault folder alias (or special command).
var actions = []struct {
	key    byte
	label  string
	folder string // empty = special
}{
	{'k', "knowledge", "knowledge"},
	{'p', "projects", "projects"},
	{'d', "decisions", "decisions"},
	{'r', "resources", "resources"},
	{'a', "archive", "archive"},
	{'s', "skip", ""},
	{'q', "quit", ""},
}

// RunInteractive walks through inbox files one by one, showing each note and
// prompting for a single-key action.
func RunInteractive(v *vault.Vault) error {
	paths, err := v.ListInbox()
	if err != nil {
		return fmt.Errorf("list inbox: %w", err)
	}
	if len(paths) == 0 {
		fmt.Println("Inbox is empty.")
		return nil
	}

	fmt.Printf("Processing %d inbox note(s). Press key to act on each.\n\n", len(paths))

	for i, path := range paths {
		entry, err := analyzeFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: skip %s: %v\n", filepath.Base(path), err)
			continue
		}

		printEntryCompact(entry, i+1, len(paths))

		action, err := promptAction()
		if err != nil {
			return err
		}

		switch action {
		case "quit":
			fmt.Println("\nStopped.")
			return nil
		case "skip":
			fmt.Println(" skipped")
		default:
			fmt.Println()
			err := note.Move(v, note.MoveOptions{
				Filename: entry.Filename,
				Folder:   action,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
		}
		fmt.Println()
	}

	fmt.Println("Inbox processed.")
	return nil
}

func printEntryCompact(e inboxEntry, idx, total int) {
	fmt.Printf("[%d/%d] %s\n", idx, total, e.Filename)
	fmt.Printf("  Title:   %s\n", orDash(e.Note.Title))
	fmt.Printf("  Type:    %s\n", orDash(e.Note.Type))
	if len(e.Note.Tags) > 0 {
		fmt.Printf("  Tags:    %s\n", strings.Join(e.Note.Tags, ", "))
	} else {
		fmt.Printf("  Tags:    (none)\n")
	}
	fmt.Printf("  Size:    %.1f KB\n", float64(e.FileSize)/1024)

	if e.BodyPreview != "" {
		fmt.Printf("  Preview: %s\n", e.BodyPreview)
	}

	fmt.Println()

	// Build action line
	var parts []string
	for _, a := range actions {
		parts = append(parts, fmt.Sprintf("[%c]%s", a.key, a.label))
	}
	fmt.Printf("  %s\n", strings.Join(parts, "  "))
	fmt.Print("  > ")
}

// promptAction reads a single keypress and returns the action name.
func promptAction() (string, error) {
	// Put stdin in raw mode so we read one byte without Enter
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		// Fallback: read a line if raw mode fails (e.g. in tests/pipes)
		return promptActionLine()
	}
	defer term.Restore(fd, oldState)

	buf := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return "", fmt.Errorf("read input: %w", err)
		}
		key := buf[0]

		// Handle ctrl-c / ctrl-d
		if key == 3 || key == 4 {
			fmt.Println()
			return "quit", nil
		}

		for _, a := range actions {
			if key == a.key {
				fmt.Printf("%c", key) // echo the key
				if a.folder != "" {
					return a.folder, nil
				}
				return a.label, nil
			}
		}
		// Ignore unknown keys
	}
}

// promptActionLine is a fallback for non-terminal environments.
func promptActionLine() (string, error) {
	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(strings.ToLower(input))
	for _, a := range actions {
		if input == string(a.key) || input == a.label {
			if a.folder != "" {
				return a.folder, nil
			}
			return a.label, nil
		}
	}
	return "skip", nil
}
