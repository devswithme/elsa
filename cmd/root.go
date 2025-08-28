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
			customRootTemplate(cmd)
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

func customRootTemplate(cmd *cobra.Command) {
	fmt.Println(getBanner(version) + `
		
Usage:
  elsa [flags]
  elsa [command]

Available Commands:`)

	for _, c := range cmd.Commands() {
		if (!c.IsAvailableCommand() || c.Hidden) && c.Name() != "help" {
			continue
		}
		fmt.Printf("  %-12s %s\n", c.Name(), c.Short)
	}

	fmt.Println(`
Flags:
  -h, --help      help for elsa
  -v, --version   version for elsa`)
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
