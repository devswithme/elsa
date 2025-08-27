package cmd

import (
	"fmt"
	"runtime"

	"github.com/risoftinc/elsa/cmd/migrate"
	"github.com/spf13/cobra"
)

var (
	version   string
	goVersion = runtime.Version()
	platform  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "elsa",
		Short: "Elsa - Engineer’s Little Smart Assistant",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(getBanner(version) + `
		
Examples:
  elsa migration   		Database migration commands
  elsa --version		Show version information
  elsa --help			Show help information`)
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
}

// SetVersionInfo sets the version information for the application
func SetVersionInfo(v string) {
	version = v
	// Override the version command to use our version info
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(customVersionTemplate())
}

func customVersionTemplate() string {
	return fmt.Sprintf("ELSA v%s (CLI)\ngo version %s %s\nLearn more at: https://risoftinc.com\n",
		version, goVersion, platform,
	)
}

func getBanner(v string) string {
	return fmt.Sprintf(`Developer productivity toolkit for Go.
      ___           ___       ___           ___     
     /\  \         /\__\     /\  \         /\  \    
    /::\  \       /:/  /    /::\  \       /::\  \   
   /:/\:\  \     /:/  /    /:/\ \  \     /:/\:\  \  
  /::\~\:\  \   /:/  /    _\:\~\ \  \   /::\~\:\  \ 
 /:/\:\ \:\__\ /:/__/    /\ \:\ \ \__\ /:/\:\ \:\__\
 \:\~\:\ \/__/ \:\  \    \:\ \:\ \/__/ \/__\:\/:/  /
  \:\ \:\__\    \:\  \    \:\ \:\__\        \::/  / 
   \:\ \/__/     \:\  \    \:\/:/  /        /:/  /  
    \:\__\        \:\__\    \::/  /        /:/  /   
     \/__/         \/__/     \/__/         \/__/    V %s
(migration, scaffolding, project runner and task automation)`, v)
}
