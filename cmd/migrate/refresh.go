package migrate

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"go.risoftinc.com/elsa/constants"
)

var (
	refreshCmd = &cobra.Command{
		Use:   "refresh [ddl|dml]",
		Short: "Refresh all migrations",
		Long: `Refresh all migrations by rolling back all applied migrations and then applying them again.
		
Examples:
  elsa migration refresh ddl                    # Refresh all DDL migrations
  elsa migration refresh dml                    # Refresh all DML migrations
  elsa migration refresh ddl -c "mysql://user:pass@host:port/db"  # Use connection string`,
		Args: cobra.ExactArgs(1),
		RunE: runRefresh,
	}
	refreshCustomPath string
	refreshConnection string
)

func init() {
	refreshCmd.Flags().StringVarP(&refreshCustomPath, "path", "p", "", "Custom migration path")
	refreshCmd.Flags().StringVarP(&refreshConnection, "connection", "c", "", "Database connection string (e.g., mysql://user:pass@host:port/db)")
}

func runRefresh(cmd *cobra.Command, args []string) error {
	migrationType := args[0]

	// Validate migration type
	if migrationType != constants.MigrationTypeDDL && migrationType != constants.MigrationTypeDML {
		return fmt.Errorf(constants.ErrInvalidMigrationType, migrationType)
	}

	fmt.Printf(constants.RefreshHeader, strings.ToUpper(migrationType))
	fmt.Printf(constants.RefreshSeparator)

	// Step 1: Rollback all migrations
	fmt.Printf(constants.InfoStep1Rollback, strings.ToUpper(migrationType))

	// Call rollback directly with connection string
	if err := rollbackAllMigrations(migrationType, refreshCustomPath, refreshConnection); err != nil {
		return fmt.Errorf(constants.ErrFailedRollbackAll, err)
	}

	fmt.Printf(constants.SuccessAllRolledBack, strings.ToUpper(migrationType))

	// Step 2: Apply all migrations again
	fmt.Printf(constants.InfoStep2Apply, strings.ToUpper(migrationType))

	// Call apply directly with connection string
	if err := applyAllMigrations(migrationType, refreshCustomPath, refreshConnection); err != nil {
		return fmt.Errorf(constants.ErrFailedApplyAll, err)
	}

	fmt.Printf(constants.SuccessAllAppliedAgain, strings.ToUpper(migrationType))

	fmt.Printf(constants.SuccessRefreshed, strings.ToUpper(migrationType))
	fmt.Printf(constants.RefreshSuccessMessage)

	return nil
}

// rollbackAllMigrations rolls back all applied migrations
func rollbackAllMigrations(migrationType, customPath, connectionString string) error {
	// Get applied migrations
	appliedMigrations, err := getAppliedMigrationsWithConnection(migrationType, connectionString)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	if len(appliedMigrations) == 0 {
		fmt.Printf("ℹ️  No %s migrations have been applied\n", strings.ToUpper(migrationType))
		return nil
	}

	// Get available migration files
	availableMigrations, err := GetAvailableMigrationsWithPath(migrationType, customPath)
	if err != nil {
		return fmt.Errorf("failed to get available migrations: %v", err)
	}

	// Create ID to Migration mapping
	idToMigration := make(map[string]Migration)
	for _, m := range availableMigrations {
		idToMigration[m.ID] = m
	}

	// Rollback all migrations
	fmt.Printf("🔄 Rolling back %d %s migration(s)...\n", len(appliedMigrations), strings.ToUpper(migrationType))

	for _, migrationID := range appliedMigrations {
		migration, exists := idToMigration[migrationID]
		if !exists {
			fmt.Printf("⚠️  Warning: Migration file for ID %s not found, skipping\n", migrationID)
			continue
		}

		if err := rollbackMigrationWithConnection(migration, migrationType, connectionString); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %v", migrationID, err)
		}
		fmt.Printf("✅ Rolled back: %s_%s\n", migration.ID, migration.Name)
	}

	return nil
}

// applyAllMigrations applies all available migrations
func applyAllMigrations(migrationType, customPath, connectionString string) error {
	// Get available migrations
	migrations, err := GetAvailableMigrationsWithPath(migrationType, customPath)
	if err != nil {
		return fmt.Errorf("failed to get available migrations: %v", err)
	}

	if len(migrations) == 0 {
		fmt.Printf("ℹ️  No %s migrations found\n", strings.ToUpper(migrationType))
		return nil
	}

	// Apply all migrations
	fmt.Printf("🚀 Applying %d %s migration(s)...\n", len(migrations), strings.ToUpper(migrationType))

	for _, migration := range migrations {
		if err := applyMigrationWithConnection(migration, migrationType, connectionString); err != nil {
			return fmt.Errorf("failed to apply migration %s: %v", migration.ID, err)
		}
		fmt.Printf("✅ Applied: %s_%s\n", migration.ID, migration.Name)
	}

	return nil
}
