package cmd

import (
	"fmt"
	"runtime"

	"github.com/risoftinc/elsa/cmd/migrate"
	"github.com/risoftinc/elsa/cmd/watch"
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
		
Usage:
  elsa [flags]
  elsa [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  migration   Database migration commands
  watch       Watch Go files and auto-restart on changes

Flags:
  -h, --help      help for elsa
  -v, --version   version for elsa`)
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
