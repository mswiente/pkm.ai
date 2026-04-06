package claudemd

import (
	_ "embed"
	"bytes"
	"text/template"
)

//go:embed global.md
var globalTemplate string

// Render returns the global CLAUDE.md content with the given vault path substituted in.
func Render(vaultPath string) (string, error) {
	tmpl, err := template.New("global").Parse(globalTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct{ VaultPath string }{VaultPath: vaultPath}); err != nil {
		return "", err
	}
	return buf.String(), nil
}
