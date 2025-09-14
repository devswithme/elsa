package make

import (
	"github.com/spf13/cobra"
	"go.risoftinc.com/elsa/internal/make"
)

// MakeCmd represents the make command
var MakeCmd = &cobra.Command{
	Use:   "make <template-type> <name>",
	Short: "Generate files from templates",
	Long: `Generate files from templates using the configured template types.

Examples:
  elsa make repository user_repository
  elsa make service user_service
  elsa make repository health/health_repository
  elsa make repository user_repository --refresh`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		refresh, _ := cmd.Flags().GetBool("refresh")
		command := make.NewMakeCommand()
		command.SetRefresh(refresh)
		return command.Execute(args)
	},
}

func init() {
	// Add flags
	MakeCmd.Flags().Bool("refresh", false, "Force refresh templates from remote repository")

	// Add subcommands
	MakeCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List available template types",
		RunE: func(cmd *cobra.Command, args []string) error {
			command := make.NewMakeCommand()
			return command.ListAvailableTypes()
		},
	})
}
