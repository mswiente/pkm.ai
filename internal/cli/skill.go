package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mswiente/pkm.ai/internal/skill"
)

func newSkillCommand() *cobra.Command {
	skillCmd := &cobra.Command{
		Use:   "skill",
		Short: "Manage Claude Code skills",
	}
	skillCmd.AddCommand(newSkillInstallCommand())
	skillCmd.AddCommand(newSkillListCommand())
	return skillCmd
}

func newSkillInstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "install [name]",
		Short: "Install skills to ~/.claude/commands/",
		Long: `Installs pkm skill definitions as Claude Code slash commands.

Without a name, all available skills are installed.
With a name, only that skill is installed.

After installation, invoke the skill in any Claude Code session with /<name>.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) == 1 {
				name = args[0]
			}
			return skill.Install(name)
		},
	}
}

func newSkillListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available skills",
		RunE: func(cmd *cobra.Command, args []string) error {
			skills, err := skill.List()
			if err != nil {
				return err
			}
			if len(skills) == 0 {
				fmt.Println("No skills available.")
				return nil
			}
			for _, s := range skills {
				installed := " "
				if skill.IsInstalled(s.Name) {
					installed = "✓"
				}
				fmt.Printf("[%s] %-20s  %s\n", installed, s.Name, s.Description)
			}
			return nil
		},
	}
}
