package migrate

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

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
  elsa migration down ddl --to 00001         # Rollback DDL migrations down to ID 00001`,
		Args: cobra.ExactArgs(1),
		RunE: runDown,
	}

	downStepCount   int
	downToMigration string
)

func init() {
	downCmd.Flags().IntVarP(&downStepCount, "step", "s", 1, "Number of migrations to rollback")
	downCmd.Flags().StringVarP(&downToMigration, "to", "t", "", "Rollback migrations down to specific ID")
}

func runDown(cmd *cobra.Command, args []string) error {
	migrationType := args[0]

	// Validate migration type
	if migrationType != "ddl" && migrationType != "dml" {
		return fmt.Errorf("migration type must be 'ddl' or 'dml', got: %s", migrationType)
	}

	// Get available migrations
	migrations, err := getAvailableMigrations(migrationType)
	if err != nil {
		return fmt.Errorf("failed to get available migrations: %v", err)
	}

	if len(migrations) == 0 {
		fmt.Printf("ℹ️  No %s migrations found\n", strings.ToUpper(migrationType))
		return nil
	}

	// Get applied migrations from database
	appliedMigrations, err := getAppliedMigrations(migrationType)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	// Filter applied migrations
	appliedMigrationsList := filterAppliedMigrations(migrations, appliedMigrations)

	if len(appliedMigrationsList) == 0 {
		fmt.Printf("ℹ️  No %s migrations have been applied\n", strings.ToUpper(migrationType))
		return nil
	}

	// Sort applied migrations in reverse order (newest first)
	sort.Slice(appliedMigrationsList, func(i, j int) bool {
		if isSequentialID(appliedMigrationsList[i].ID) && isSequentialID(appliedMigrationsList[j].ID) {
			seqI, _ := strconv.Atoi(appliedMigrationsList[i].ID)
			seqJ, _ := strconv.Atoi(appliedMigrationsList[j].ID)
			return seqI > seqJ
		}
		return appliedMigrationsList[i].ID > appliedMigrationsList[j].ID
	})

	// Determine which migrations to rollback
	var migrationsToRollback []Migration
	if downToMigration != "" {
		migrationsToRollback = filterMigrationsFromID(appliedMigrationsList, downToMigration)
	} else {
		if downStepCount > len(appliedMigrationsList) {
			downStepCount = len(appliedMigrationsList)
		}
		migrationsToRollback = appliedMigrationsList[:downStepCount]
	}

	if len(migrationsToRollback) == 0 {
		fmt.Printf("ℹ️  No %s migrations to rollback\n", strings.ToUpper(migrationType))
		return nil
	}

	// Rollback migrations
	fmt.Printf("🔄 Rolling back %d %s migration(s)...\n", len(migrationsToRollback), strings.ToUpper(migrationType))

	for _, migration := range migrationsToRollback {
		if err := rollbackMigration(migration, migrationType); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %v", migration.ID, err)
		}
		fmt.Printf("✅ Rolled back: %s_%s\n", migration.ID, migration.Name)
	}

	fmt.Printf("🎉 Successfully rolled back %d %s migration(s)\n", len(migrationsToRollback), strings.ToUpper(migrationType))
	return nil
}

func filterAppliedMigrations(available []Migration, applied []string) []Migration {
	appliedMap := make(map[string]bool)
	for _, id := range applied {
		appliedMap[id] = true
	}

	var appliedList []Migration
	for _, m := range available {
		if appliedMap[m.ID] {
			appliedList = append(appliedList, m)
		}
	}

	return appliedList
}

func filterMigrationsFromID(migrations []Migration, targetID string) []Migration {
	var result []Migration
	for _, m := range migrations {
		result = append(result, m)
		if m.ID == targetID {
			break
		}
	}
	return result
}

func rollbackMigration(migration Migration, migrationType string) error {
	// Find the corresponding down migration file
	downFilePath := strings.Replace(migration.Path, ".up.sql", ".down.sql", 1)

	// Check if down migration file exists
	if _, err := os.Stat(downFilePath); os.IsNotExist(err) {
		return fmt.Errorf("down migration file not found: %s", downFilePath)
	}

	// Read down migration file
	content, err := os.ReadFile(downFilePath)
	if err != nil {
		return fmt.Errorf("failed to read down migration file: %v", err)
	}

	// TODO: Execute SQL against database
	fmt.Printf("   Executing rollback: %s\n", string(content))

	// TODO: Remove migration from applied list in database

	return nil
}
