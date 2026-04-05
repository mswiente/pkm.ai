package readwise

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// State tracks the last successful sync timestamp.
type State struct {
	LastSyncedAt *time.Time `json:"last_synced_at"`
}

// StateFilePath returns the path to the sync state file.
// Respects XDG_CONFIG_HOME if set, otherwise ~/.config/pkm/readwise_sync_state.json.
func StateFilePath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "pkm", "readwise_sync_state.json")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "pkm", "readwise_sync_state.json")
}

// LoadState reads the sync state from path. Returns an empty state if the file doesn't exist.
func LoadState(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &State{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read state file: %w", err)
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse state file: %w", err)
	}
	return &s, nil
}

// SaveState writes the sync state to path, creating parent directories as needed.
func SaveState(path string, s *State) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create state dir: %w", err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write state file: %w", err)
	}
	return nil
}
