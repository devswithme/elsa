package generate

import (
	"github.com/spf13/cobra"
)

var GenCmd = &cobra.Command{
	Use:   "gen [directory]",
	Short: "Alias for generate command",
	Long: `Gen is an alias for the generate command.
It searches for files with '//go:build elsabuild' tags in the specified directory or current directory and subdirectories.

Usage:
  elsa gen           # Same as 'elsa generate'
  elsa gen [dir]    # Same as 'elsa generate [dir]'

Examples:
  elsa gen          # Find all files with elsabuild tags in current directory
  elsa gen database # Find all files with elsabuild tags in database directory`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Execute the same logic as generate command
		GenerateCmd.Run(cmd, args)
	},
}
