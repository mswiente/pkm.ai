package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds all pkm configuration values.
type Config struct {
	VaultPath        string `yaml:"vault_path,omitempty"`
	InboxPath        string `yaml:"inbox_path,omitempty"`
	DailyPath        string `yaml:"daily_path,omitempty"`
	TemplatesPath    string `yaml:"templates_path,omitempty"`
	Editor           string `yaml:"editor,omitempty"`
	FilenameTimezone string `yaml:"filename_timezone,omitempty"`
	DefaultSource    string `yaml:"default_source,omitempty"`
	ReadwiseToken    string `yaml:"readwise_token,omitempty"`
}

// Load reads ~/.config/pkm/config.yaml and applies defaults for non-path fields.
// vault_path has no default — it must be explicitly configured.
// Returns an error only for parse failures or unexpected I/O errors.
func Load() (*Config, error) {
	cfg := &Config{}

	path := configPath()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		applyNonPathDefaults(cfg)
		return cfg, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	applyNonPathDefaults(cfg)
	cfg.VaultPath = expand(cfg.VaultPath)

	return cfg, nil
}

// SetVaultPath writes vault_path to the config file, creating it if needed.
func SetVaultPath(path string) error {
	cfgPath := configPath()
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	// Load existing config to preserve other fields
	cfg := &Config{}
	if data, err := os.ReadFile(cfgPath); err == nil {
		_ = yaml.Unmarshal(data, cfg)
	}

	cfg.VaultPath = path

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(cfgPath, data, 0o644); err != nil {
		return fmt.Errorf("write config %s: %w", cfgPath, err)
	}

	return nil
}

// SetReadwiseToken writes readwise_token to the config file, creating it if needed.
func SetReadwiseToken(token string) error {
	cfgPath := configPath()
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	cfg := &Config{}
	if data, err := os.ReadFile(cfgPath); err == nil {
		_ = yaml.Unmarshal(data, cfg)
	}

	cfg.ReadwiseToken = token

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(cfgPath, data, 0o644); err != nil {
		return fmt.Errorf("write config %s: %w", cfgPath, err)
	}

	return nil
}

// RequireVaultPath returns an error with a helpful message if vault_path is not set.
func RequireVaultPath(cfg *Config) error {
	if cfg.VaultPath == "" {
		return fmt.Errorf("vault path not configured\nRun: pkm config --set-vault-path \"/path/to/vault\"")
	}
	return nil
}

// ConfigFilePath returns the path to the config file (for display purposes).
func ConfigFilePath() string {
	return configPath()
}

func applyNonPathDefaults(cfg *Config) {
	if cfg.InboxPath == "" {
		cfg.InboxPath = "00-inbox"
	}
	if cfg.DailyPath == "" {
		cfg.DailyPath = "01-daily"
	}
	if cfg.TemplatesPath == "" {
		cfg.TemplatesPath = "07-templates"
	}
	if cfg.Editor == "" {
		cfg.Editor = defaultEditor()
	}
	if cfg.FilenameTimezone == "" {
		cfg.FilenameTimezone = "UTC"
	}
	if cfg.DefaultSource == "" {
		cfg.DefaultSource = "manual"
	}
}

func defaultEditor() string {
	if v := os.Getenv("VISUAL"); v != "" {
		return v
	}
	if v := os.Getenv("EDITOR"); v != "" {
		return v
	}
	return "vi"
}

func configPath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "pkm", "config.yaml")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "pkm", "config.yaml")
}

func expand(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return filepath.Join(home, path[1:])
}
