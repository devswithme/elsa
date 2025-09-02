package generate

import (
	"fmt"
	"os"

	"github.com/risoftinc/elsa/internal/generate"
	"github.com/spf13/cobra"
)

var GenerateCmd = &cobra.Command{
	Use:   "generate [directory]",
	Short: "Generate code by finding files with elsabuild tags",
	Long: `Generate searches for files with '//go:build elsabuild' tags in the specified directory or current directory and subdirectories.
This command is useful for identifying files that need to be processed during the build process.

Usage:
  elsa generate           # Search current directory and subdirectories
  elsa generate [dir]    # Search specified directory and subdirectories
  elsa gen               # Alias for generate command
  elsa gen [dir]         # Same as generate with directory

Examples:
  elsa generate              # Find all files with elsabuild tags in current directory
  elsa generate dependencies # Find all files with elsabuild tags in dependencies directory
  elsa gen dependencies/http # Find all files with elsabuild tags in dependencies/http directory`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var targetDir string
		if len(args) > 0 {
			targetDir = args[0]
		}

		generator := generate.NewGenerator()
		if err := generator.GenerateDependencies(targetDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}
