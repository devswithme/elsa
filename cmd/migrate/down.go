package migrate

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/risoftinc/elsa/constants"
	"github.com/risoftinc/elsa/internal/database"
	"github.com/spf13/cobra"
)

var (
	downCmd = &cobra.Command{
		Use:   "down [ddl|dml]",
		Short: "Rollback database migrations",
		Long: `Rollback database migrations for DDL or DML operations.
		
Examples:
  elsa migration down ddl                    # Rollback last DDL migration
  elsa migration down dml                    # Rollback last DML migration
  elsa migration down ddl --step 2           # Rollback 2 DDL migrations
  elsa migration down ddl --to 00002         # Rollback DDL migrations down to ID 00002
  elsa migration down ddl --from 00005       # Rollback DDL migrations from ID 00005
  elsa migration down ddl --path custom/migrations  # Use custom migration path`,
		Args: cobra.ExactArgs(1),
		RunE: runDown,
	}

	downAll           bool
	downStepCount     int
	downToMigration   string
	downFromMigration string
	downCustomPath    string
	downConnection    string
)

func init() {
	downCmd.Flags().IntVarP(&downStepCount, "step", "s", 0, "Number of migrations to rollback")
	downCmd.Flags().StringVarP(&downToMigration, "to", "t", "", "Rollback migrations down to specific ID")
	downCmd.Flags().StringVarP(&downFromMigration, "from", "f", "", "Rollback migrations from specific ID")
	downCmd.Flags().StringVarP(&downCustomPath, "path", "p", "", "Custom migration path")
	downCmd.Flags().StringVarP(&downConnection, "connection", "c", "", "Database connection string (e.g., mysql://user:pass@host:port/db)")
	downCmd.Flags().BoolVarP(&downAll, "all", "a", false, "Rollback all migrations")
}

func runDown(cmd *cobra.Command, args []string) error {
	migrationType := args[0]

	// Validate migration type
	if migrationType != constants.MigrationTypeDDL && migrationType != constants.MigrationTypeDML {
		return fmt.Errorf(constants.ErrInvalidMigrationType, migrationType)
	}

	// Get applied migrations from database
	appliedMigrations, err := getAppliedMigrationsForDown(migrationType)
	if err != nil {
		return fmt.Errorf(constants.ErrFailedGetAppliedMigrations, err)
	}

	if len(appliedMigrations) == 0 {
		fmt.Printf(constants.InfoNoMigrations, strings.ToUpper(migrationType))
		return nil
	}

	// Get available migration files to map IDs to names
	availableMigrations, err := GetAvailableMigrationsWithPath(migrationType, downCustomPath)
	if err != nil {
		return fmt.Errorf(constants.ErrFailedGetAppliedMigrations, err)
	}

	// Create ID to Migration mapping
	idToMigration := make(map[string]Migration)
	for _, m := range availableMigrations {
		idToMigration[m.ID] = m
	}

	// Determine which migrations to rollback
	var migrationsToRollback []string
	if downAll {
		migrationsToRollback = appliedMigrations
	} else if downStepCount > 0 {
		if downStepCount > len(appliedMigrations) {
			downStepCount = len(appliedMigrations)
		}
		migrationsToRollback = appliedMigrations[len(appliedMigrations)-downStepCount:]
	} else if downToMigration != "" {
		migrationsToRollback = filterMigrationsFromID(appliedMigrations, downToMigration)
	} else if downFromMigration != "" {
		migrationsToRollback = filterMigrationsFromID(appliedMigrations, downFromMigration)
	} else {
		// Default: rollback last migration
		migrationsToRollback = appliedMigrations[len(appliedMigrations)-1:]
	}

	if len(migrationsToRollback) == 0 {
		fmt.Printf(constants.InfoNoMigrationsToRollback, strings.ToUpper(migrationType))
		return nil
	}

	// Sort migrations in reverse order (newest first) for rollback
	sort.Slice(migrationsToRollback, func(i, j int) bool {
		// Handle both sequential and timestamp formats
		if isSequentialID(migrationsToRollback[i]) && isSequentialID(migrationsToRollback[j]) {
			seqI, _ := strconv.Atoi(migrationsToRollback[i])
			seqJ, _ := strconv.Atoi(migrationsToRollback[j])
			return seqI > seqJ
		}
		return migrationsToRollback[i] > migrationsToRollback[j]
	})

	// Rollback migrations
	fmt.Printf("🔄 Rolling back %d %s migration(s)...\n", len(migrationsToRollback), strings.ToUpper(migrationType))

	for _, migrationID := range migrationsToRollback {
		migration, exists := idToMigration[migrationID]
		if !exists {
			fmt.Printf("⚠️  Warning: Migration file for ID %s not found, skipping\n", migrationID)
			continue
		}

		if err := rollbackMigration(migration, migrationType); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %v", migrationID, err)
		}
		fmt.Printf("✅ Rolled back: %s_%s\n", migration.ID, migration.Name)
	}

	fmt.Printf("🎉 Successfully rolled back %d %s migration(s)\n", len(migrationsToRollback), strings.ToUpper(migrationType))
	return nil
}

func filterMigrationsFromID(migrations []string, targetID string) []string {
	var result []string
	found := false

	for _, id := range migrations {
		if id == targetID {
			found = true
		}
		if found {
			result = append(result, id)
		}
	}

	return result
}

// getAppliedMigrationsForDown retrieves applied migrations from database for down command
func getAppliedMigrationsForDown(migrationType string) ([]string, error) {
	var config *database.DatabaseConfig
	var err error

	// If connection flag is provided, use it directly
	if downConnection != "" {
		config = database.ParseConnectionString(downConnection)
		if config == nil {
			return nil, fmt.Errorf("invalid connection string: %s", downConnection)
		}
	} else {
		// Get database configuration using helper function
		config, err = GetDatabaseConnection()
		if err != nil {
			return nil, err
		}
	}

	// Connect to database
	db, err := database.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Get migration executor
	executor := database.NewMigrationExecutor(db)

	// Ensure migration table exists
	if err := executor.EnsureMigrationTable(); err != nil {
		return nil, fmt.Errorf("failed to ensure migration table: %v", err)
	}

	// Get applied migrations
	return executor.GetAppliedMigrations(migrationType)
}

func rollbackMigration(migration Migration, migrationType string) error {
	// Read down migration file
	downFilePath := strings.Replace(migration.Path, ".up.sql", constants.DownMigrationExtension, 1)

	content, err := os.ReadFile(downFilePath)
	if err != nil {
		return fmt.Errorf("failed to read down migration file: %v", err)
	}

	var config *database.DatabaseConfig

	// If connection flag is provided, use it directly
	if downConnection != "" {
		config = database.ParseConnectionString(downConnection)
		if config == nil {
			return fmt.Errorf("invalid connection string: %s", downConnection)
		}
	} else {
		// Get database configuration using helper function
		var err error
		config, err = GetDatabaseConnection()
		if err != nil {
			return err
		}
	}

	// Connect to database
	db, err := database.Connect(config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Get migration executor
	executor := database.NewMigrationExecutor(db)

	// Ensure migration table exists
	if err := executor.EnsureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %v", err)
	}

	// Execute rollback migration
	startTime := time.Now()
	if err := executor.ExecuteMigration(string(content), migrationType); err != nil {
		return fmt.Errorf("failed to execute rollback migration: %v", err)
	}
	executionTime := time.Since(startTime).Milliseconds()

	// Remove migration record
	if err := executor.RemoveMigration(migration.ID); err != nil {
		return fmt.Errorf("failed to remove migration record: %v", err)
	}

	fmt.Printf("   ✅ Rolled back in %dms\n", executionTime)

	return nil
}
