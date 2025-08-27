package migrate

import (
	"github.com/spf13/cobra"
)

var (
	// migrateCmd represents the migrate command
	migrateCmd = &cobra.Command{
		Use:   "migration",
		Short: "Database migration commands",
		Long: `Database migration commands for managing DDL and DML changes.
		
Examples:
  elsa migration connect                                            Connect to database interactively
  elsa migration connect -c "sqlite://elsa.db"                     Connect using connection string flag
  elsa migration create ddl create_users_table                      Create new DDL migration (timestamp with ms)
  elsa migration create dml seed_users_data                         Create new DML migration (timestamp with ms)
  elsa migration create ddl create_table --sequential               Create with sequential format
  elsa migration create ddl create_table --path custom/migrations   Custom folder path
  elsa migration up ddl                                             Apply all DDL migrations
  elsa migration down dml                                           Rollback last DML migration
  elsa migration status                                             Show migration status
  elsa migration refresh ddl                                        Refresh all DDL migrations`,
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
