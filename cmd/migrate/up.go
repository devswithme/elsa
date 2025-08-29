package migrate

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/risoftinc/elsa/constants"
	"github.com/risoftinc/elsa/internal/database"
	"github.com/spf13/cobra"
)

var (
	upCmd = &cobra.Command{
		Use:   "up [ddl|dml]",
		Short: "Apply database migrations",
		Long: `Apply database migrations for DDL or DML operations.
		
Examples:
  elsa migration up ddl                    # Apply all DDL migrations
  elsa migration up dml                    # Apply all DML migrations
  elsa migration up ddl --step 2           # Apply 2 DDL migrations
  elsa migration up ddl --to 00002         # Apply DDL migrations up to ID 00002
  elsa migration up ddl --path custom/migrations  # Use custom migration path`,
		Args: cobra.ExactArgs(1),
		RunE: runUp,
	}

	upStepCount   int
	upToMigration string
	upCustomPath  string
	upConnection  string
)

func init() {
	upCmd.Flags().IntVarP(&upStepCount, "step", "s", 0, "Number of migrations to apply")
	upCmd.Flags().StringVarP(&upToMigration, "to", "t", "", "Apply migrations up to specific ID")
	upCmd.Flags().StringVarP(&upCustomPath, "path", "p", "", "Custom migration path")
	upCmd.Flags().StringVarP(&upConnection, "connection", "c", "", "Direct database connection string (e.g., mysql://user:pass@host:port/db)")
}

func runUp(cmd *cobra.Command, args []string) error {
	migrationType := args[0]

	// Validate migration type
	if migrationType != constants.MigrationTypeDDL && migrationType != constants.MigrationTypeDML {
		return fmt.Errorf(constants.ErrInvalidMigrationType, migrationType)
	}

	// Get available migrations
	migrations, err := getAvailableMigrations(migrationType)
	if err != nil {
		return fmt.Errorf(constants.ErrFailedGetAppliedMigrations, err)
	}

	if len(migrations) == 0 {
		fmt.Printf(constants.InfoNoMigrationsFoundStatus, strings.ToUpper(migrationType))
		return nil
	}

	// Get applied migrations from database
	appliedMigrations, err := getAppliedMigrations(migrationType)
	if err != nil {
		return fmt.Errorf(constants.ErrFailedGetAppliedMigrations, err)
	}

	// Filter pending migrations
	pendingMigrations := filterPendingMigrations(migrations, appliedMigrations)

	if len(pendingMigrations) == 0 {
		fmt.Printf(constants.SuccessAllApplied, strings.ToUpper(migrationType))
		return nil
	}

	// Determine which migrations to apply
	var migrationsToApply []Migration
	if upStepCount > 0 {
		if upStepCount > len(pendingMigrations) {
			upStepCount = len(pendingMigrations)
		}
		migrationsToApply = pendingMigrations[:upStepCount]
	} else if upToMigration != "" {
		migrationsToApply = filterMigrationsToID(pendingMigrations, upToMigration)
	} else {
		migrationsToApply = pendingMigrations
	}

	if len(migrationsToApply) == 0 {
		fmt.Printf(constants.InfoNoMigrationsToApply, strings.ToUpper(migrationType))
		return nil
	}

	// Apply migrations
	fmt.Printf(constants.InfoApplyingMigrations, len(migrationsToApply), strings.ToUpper(migrationType))

	for _, migration := range migrationsToApply {
		if err := applyMigration(migration, migrationType); err != nil {
			return fmt.Errorf(constants.ErrFailedApplyMigration, migration.ID, err)
		}
		fmt.Printf("✅ Applied: %s_%s\n", migration.ID, migration.Name)
	}

	fmt.Printf("🎉 Successfully applied %d %s migration(s)\n", len(migrationsToApply), strings.ToUpper(migrationType))
	return nil
}

type Migration struct {
	ID   string
	Name string
	Path string
}

func getAvailableMigrations(migrationType string) ([]Migration, error) {
	return GetAvailableMigrationsWithPath(migrationType, upCustomPath)
}

func isSequentialID(id string) bool {
	_, err := strconv.Atoi(id)
	return err == nil
}

func getAppliedMigrations(migrationType string) ([]string, error) {
	var config *database.DatabaseConfig
	var err error

	// If connection flag is provided, use it directly
	if upConnection != "" {
		config = database.ParseConnectionString(upConnection)
		if config == nil {
			return nil, fmt.Errorf("invalid connection string: %s", upConnection)
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

func filterPendingMigrations(available []Migration, applied []string) []Migration {
	appliedMap := make(map[string]bool)
	for _, id := range applied {
		appliedMap[id] = true
	}

	var pending []Migration
	for _, m := range available {
		if !appliedMap[m.ID] {
			pending = append(pending, m)
		}
	}

	return pending
}

func filterMigrationsToID(migrations []Migration, targetID string) []Migration {
	var result []Migration
	for _, m := range migrations {
		result = append(result, m)
		if m.ID == targetID {
			break
		}
	}
	return result
}

func applyMigration(migration Migration, migrationType string) error {
	// Read migration file
	content, err := os.ReadFile(migration.Path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %v", err)
	}

	var config *database.DatabaseConfig

	// If connection flag is provided, use it directly
	if upConnection != "" {
		config = database.ParseConnectionString(upConnection)
		if config == nil {
			return fmt.Errorf("invalid connection string: %s", upConnection)
		}
	} else {
		// Load database configuration from .env
		config = database.LoadFromEnv()
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

	// Execute migration
	startTime := time.Now()
	if err := executor.ExecuteMigration(string(content), migrationType); err != nil {
		return fmt.Errorf("failed to execute migration: %v", err)
	}
	executionTime := time.Since(startTime).Milliseconds()

	// Record migration as applied
	checksum := database.GetMigrationChecksum(string(content))
	if err := executor.RecordMigration(migration.ID, migration.Name, migrationType, checksum, executionTime); err != nil {
		return fmt.Errorf("failed to record migration: %v", err)
	}

	fmt.Printf("   ✅ Executed in %dms\n", executionTime)

	return nil
}
