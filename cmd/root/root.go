package root

import (
	"fmt"
	"runtime"

	"github.com/risoftinc/elsa/cmd/elsafile"
	"github.com/risoftinc/elsa/cmd/migrate"
	"github.com/risoftinc/elsa/cmd/watch"
	"github.com/risoftinc/elsa/constants"
	"github.com/risoftinc/elsa/internal/root"
	"github.com/spf13/cobra"
)

var (
	version   string
	goVersion = runtime.Version()
	platform  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   constants.RootUse,
		Short: constants.RootShort,
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			handler := root.NewCommandHandler()
			return handler.HandleRootCommand(cmd, args, version)
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add migration command
	rootCmd.AddCommand(migrate.MigrateCmd())

	// Add watch command
	rootCmd.AddCommand(watch.WatchCmd)

	// Add Elsafile commands
	rootCmd.AddCommand(elsafile.InitCmd)
	rootCmd.AddCommand(elsafile.RunCmd)
	rootCmd.AddCommand(elsafile.ListCmd)
}

// SetVersionInfo sets the version information for the application
func SetVersionInfo(v string) {
	version = v
	// Override the version command to use our version info
	rootCmd.Version = version

	displayHelper := root.NewDisplayHelper()
	rootCmd.SetVersionTemplate(displayHelper.GetVersionTemplate(version, goVersion, platform))
}
