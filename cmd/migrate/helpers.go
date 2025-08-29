package migrate

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/risoftinc/elsa/constants"
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
			return nil, fmt.Errorf(constants.ErrInvalidConnectionString, connectionString)
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
		return nil, fmt.Errorf(constants.ErrFailedConnectDB, err)
	}

	// Get migration executor
	executor := database.NewMigrationExecutor(db)

	// Ensure migration table exists
	if err := executor.EnsureMigrationTable(); err != nil {
		return nil, fmt.Errorf(constants.ErrFailedEnsureTable, err)
	}

	// Get applied migrations
	return executor.GetAppliedMigrations(migrationType)
}

// applyMigrationWithConnection applies a migration using connection string
func applyMigrationWithConnection(migration Migration, migrationType, connectionString string) error {
	// Read migration file
	content, err := os.ReadFile(migration.Path)
	if err != nil {
		return fmt.Errorf(constants.ErrFailedReadFile, err)
	}

	var config *database.DatabaseConfig

	// If connection string is provided, use it directly
	if connectionString != "" {
		config = database.ParseConnectionString(connectionString)
		if config == nil {
			return fmt.Errorf(constants.ErrInvalidConnectionString, connectionString)
		}
	} else {
		// Load database configuration from .env
		config = database.LoadFromEnv()
	}

	// Connect to database
	db, err := database.Connect(config)
	if err != nil {
		return fmt.Errorf(constants.ErrFailedConnectDB, err)
	}

	// Get migration executor
	executor := database.NewMigrationExecutor(db)

	// Ensure migration table exists
	if err := executor.EnsureMigrationTable(); err != nil {
		return fmt.Errorf(constants.ErrFailedEnsureTable, err)
	}

	// Execute migration
	startTime := time.Now()
	if err := executor.ExecuteMigration(string(content), migrationType); err != nil {
		return fmt.Errorf(constants.ErrFailedExecuteMigration, err)
	}
	executionTime := time.Since(startTime).Milliseconds()

	// Record migration as applied
	checksum := database.GetMigrationChecksum(string(content))
	if err := executor.RecordMigration(migration.ID, migration.Name, migrationType, checksum, executionTime); err != nil {
		return fmt.Errorf(constants.ErrFailedRecordMigration, err)
	}

	fmt.Printf(constants.SuccessExecuted, executionTime)

	return nil
}

// rollbackMigrationWithConnection rolls back a migration using connection string
func rollbackMigrationWithConnection(migration Migration, migrationType, connectionString string) error {
	// Read down migration file
	downFilePath := strings.Replace(migration.Path, constants.UpMigrationExtension, constants.DownMigrationExtension, 1)

	content, err := os.ReadFile(downFilePath)
	if err != nil {
		return fmt.Errorf(constants.ErrFailedReadFile, err)
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
		return fmt.Errorf(constants.ErrFailedConnectDB, err)
	}

	// Get migration executor
	executor := database.NewMigrationExecutor(db)

	// Ensure migration table exists
	if err := executor.EnsureMigrationTable(); err != nil {
		return fmt.Errorf(constants.ErrFailedEnsureTable, err)
	}

	// Execute rollback migration
	startTime := time.Now()
	if err := executor.ExecuteMigration(string(content), migrationType); err != nil {
		return fmt.Errorf(constants.ErrFailedRollbackMigration, err)
	}
	executionTime := time.Since(startTime).Milliseconds()

	// Remove migration record
	if err := executor.RemoveMigration(migration.ID); err != nil {
		return fmt.Errorf(constants.ErrFailedRemoveRecord, err)
	}

	fmt.Printf(constants.SuccessRolledBack, executionTime)

	return nil
}
