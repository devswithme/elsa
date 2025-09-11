package migrate

import (
	"github.com/spf13/cobra"
	"go.risoftinc.com/elsa/constants"
)

var (
	// migrateCmd represents the migrate command
	migrateCmd = &cobra.Command{
		Use:   constants.MigrateUse,
		Short: constants.MigrateShort,
		Long:  constants.MigrateLong,
	}
)

// Execute adds all child commands to the migrate command
func Execute() error {
	return migrateCmd.Execute()
}

// MigrateCmd returns the migrate command
func MigrateCmd() *cobra.Command {
	return migrateCmd
}

func init() {
	// Add subcommands
	migrateCmd.AddCommand(createCmd)
	migrateCmd.AddCommand(connectCmd)
	migrateCmd.AddCommand(upCmd)
	migrateCmd.AddCommand(downCmd)
	migrateCmd.AddCommand(refreshCmd)
	migrateCmd.AddCommand(statusCmd)
	migrateCmd.AddCommand(infoCmd)
}
