package cli

import (
	"github.com/spf13/cobra"
	"github.com/mswiente/pkm.ai/internal/config"
	"github.com/mswiente/pkm.ai/internal/note"
)

func newNoteCommand(cfg *config.Config) *cobra.Command {
	noteCmd := &cobra.Command{
		Use:   "note",
		Short: "Manage notes in the vault",
	}
	noteCmd.AddCommand(newNoteMoveCommand(cfg))
	return noteCmd
}

func newNoteMoveCommand(cfg *config.Config) *cobra.Command {
	var (
		noteType   string
		status     string
		dryRun     bool
	)

	cmd := &cobra.Command{
		Use:   "move <filename> <folder>",
		Short: "Move a note to a target vault folder",
		Long: `Moves a note to the specified folder and updates its frontmatter.

The source file is searched in the inbox first, then the entire vault.

Target folder accepts full names, shorthand, or numbers:
  knowledge (4)   → 04-knowledge/
  projects  (2)   → 02-projects/
  resources (5)   → 05-resources/
  decisions (6)   → 06-decisions/
  areas     (3)   → 03-areas/
  archive   (9)   → 09-archive/

Frontmatter is updated automatically:
  type:    inferred from the target folder (overridable with --type)
  status:  inbox→draft, any→archived when moving to archive (overridable with --status)
  updated: set to today`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := mustVault(cfg)
			if err != nil {
				return err
			}
			opts := note.MoveOptions{
				Filename: args[0],
				Folder:   args[1],
				Type:     noteType,
				Status:   status,
				DryRun:   dryRun,
			}
			return note.Move(v, opts)
		},
	}

	cmd.Flags().StringVar(&noteType, "type", "", "Override inferred note type")
	cmd.Flags().StringVar(&status, "status", "", "Override inferred status")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would change without moving")

	return cmd
}
