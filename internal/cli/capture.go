package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/mswiente/pkm.ai/internal/capture"
	"github.com/mswiente/pkm.ai/internal/config"
)

func newCaptureCommand(cfg *config.Config) *cobra.Command {
	var (
		title     string
		source    string
		tagsFlag  string
		typeHint  string
		editor    bool
		clipboard bool
	)

	cmd := &cobra.Command{
		Use:   "capture [text]",
		Short: "Capture a new note to the inbox",
		Long: `Capture creates a new markdown note in 00-inbox/.

Input sources (in priority order):
  --clipboard    read from system clipboard
  stdin pipe     echo "text" | pkm capture
  [text]         positional argument
  (none)         creates a note shell with empty content`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := mustVault(cfg)
			if err != nil {
				return err
			}
			r := newRenderer(cfg)

			var tags []string
			if tagsFlag != "" {
				for _, t := range strings.Split(tagsFlag, ",") {
					if t := strings.TrimSpace(t); t != "" {
						tags = append(tags, t)
					}
				}
			}

			opts := capture.Options{
				Text:       strings.Join(args, " "),
				Title:      title,
				Source:     source,
				Tags:       tags,
				TypeHint:   typeHint,
				OpenEditor: editor,
				Clipboard:  clipboard,
			}
			return capture.Run(v, r, opts)
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Note title")
	cmd.Flags().StringVar(&source, "source", "", "Source: manual|chatgpt|claude-code|copilot-cli|other")
	cmd.Flags().StringVar(&tagsFlag, "tags", "", "Comma-separated tags")
	cmd.Flags().StringVar(&typeHint, "type-hint", "", "Hint for note classification (e.g. knowledge, decision)")
	cmd.Flags().BoolVar(&editor, "editor", false, "Open in editor before saving")
	cmd.Flags().BoolVar(&clipboard, "clipboard", false, "Read content from clipboard")

	return cmd
}
