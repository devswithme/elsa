package cmd

import (
	"fmt"
	"runtime"

	"github.com/risoftinc/elsa/cmd/elsafile"
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
		Use:                "elsa",
		Short:              "Elsa - Engineer’s Little Smart Assistant",
		DisableFlagParsing: true,
		Args:               cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				// Try to handle as Elsafile command
				handler := elsafile.NewSimpleHandler()
				if err := handler.HandleUnknownCommand(args[0]); err != nil {
					// If it's not an Elsafile command, show suggestions
					suggestions := handler.SuggestCommands(args[0])
					if len(suggestions) > 0 {
						fmt.Printf("💡 Did you mean one of these commands?\n")
						for _, suggestion := range suggestions {
							fmt.Printf("  elsa %s\n", suggestion)
						}
						fmt.Println()
					}
					return err
				}
				return nil
			}

			customRootTemplate(cmd)
			return nil
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
	rootCmd.SetVersionTemplate(customVersionTemplate())
}

func customRootTemplate(cmd *cobra.Command) {
	fmt.Println(getBanner(version) + `
		
Usage:
  elsa [flags]
  elsa [command]

Available Commands:`)

	printAvaibleCommands(cmd)

	fmt.Println(`
Flags:
  -h, --help      help for elsa
  -v, --version   version for elsa`)
}

func printAvaibleCommands(cmd *cobra.Command) {
	maxLen := 0
	for _, c := range cmd.Commands() {
		if (!c.IsAvailableCommand() || c.Hidden) && c.Name() != "help" {
			continue
		}
		if l := len(c.Name()); l > maxLen {
			maxLen = l
		}
	}

	for _, c := range cmd.Commands() {
		if (!c.IsAvailableCommand() || c.Hidden) && c.Name() != "help" {
			continue
		}
		fmt.Printf("  %-*s %s\n", maxLen+1, c.Name(), c.Short)
	}
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
