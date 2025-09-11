package new

import (
	"github.com/spf13/cobra"
	"go.risoftinc.com/elsa/constants"
	"go.risoftinc.com/elsa/internal/new"
)

var (
	NewCmd = &cobra.Command{
		Use:   constants.NewUse,
		Short: constants.NewShort,
		Long:  constants.NewLong,
		Args:  cobra.ExactArgs(2),
		RunE:  runNew,
	}

	// Flags
	moduleFlag  string
	outputFlag  string
	forceFlag   bool
	refreshFlag bool
)

func init() {
	NewCmd.Flags().StringVarP(&moduleFlag, "module", "m", "", constants.NewFlagModuleUsage)
	NewCmd.Flags().StringVarP(&outputFlag, "output", "o", "", constants.NewFlagOutputUsage)
	NewCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, constants.NewFlagForceUsage)
	NewCmd.Flags().BoolVar(&refreshFlag, "refresh", false, constants.NewFlagRefreshUsage)

	// Module is optional now - will auto-generate from project name if not provided
}

func runNew(cmd *cobra.Command, args []string) error {
	projectName := args[1]

	// Create project options
	options := new.NewProjectOptions(
		args[0],     // template name
		projectName, // project name
		moduleFlag,  // module name (auto-generated or provided)
		outputFlag,  // output directory
		forceFlag,   // force flag
		refreshFlag, // refresh flag
	)

	// Create template manager
	templateManager := new.NewTemplateManager()

	// Create project with all logic handled internally
	return templateManager.CreateProjectWithOutput(options)
}
