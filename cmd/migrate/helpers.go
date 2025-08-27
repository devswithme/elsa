package migrate

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/risoftinc/elsa/internal/database"
)

// getAppliedMigrationsWithConnection retrieves applied migrations using connection string
func getAppliedMigrationsWithConnection(migrationType, connectionString string) ([]string, error) {
	var config *database.DatabaseConfig
	var err error

	// If connection string is provided, use it directly
	if connectionString != "" {
		config = database.ParseConnectionString(connectionString)
		if config == nil {
			return nil, fmt.Errorf("invalid connection string: %s", connectionString)
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

// applyMigrationWithConnection applies a migration using connection string
func applyMigrationWithConnection(migration Migration, migrationType, connectionString string) error {
	// Read migration file
	content, err := os.ReadFile(migration.Path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %v", err)
	}

	var config *database.DatabaseConfig

	// If connection string is provided, use it directly
	if connectionString != "" {
		config = database.ParseConnectionString(connectionString)
		if config == nil {
			return fmt.Errorf("invalid connection string: %s", connectionString)
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

// rollbackMigrationWithConnection rolls back a migration using connection string
func rollbackMigrationWithConnection(migration Migration, migrationType, connectionString string) error {
	// Read down migration file
	downFilePath := strings.Replace(migration.Path, ".up.sql", ".down.sql", 1)

	content, err := os.ReadFile(downFilePath)
	if err != nil {
		return fmt.Errorf("failed to read down migration file: %v", err)
	}

	var config *database.DatabaseConfig

	// If connection string is provided, use it directly
	if connectionString != "" {
		config = database.ParseConnectionString(connectionString)
		if config == nil {
			return fmt.Errorf("invalid connection string: %s", connectionString)
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
