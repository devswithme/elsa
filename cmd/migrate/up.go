package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

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
  elsa migration up ddl --to 00002         # Apply DDL migrations up to ID 00002`,
		Args: cobra.ExactArgs(1),
		RunE: runUp,
	}

	stepCount   int
	toMigration string
)

func init() {
	upCmd.Flags().IntVarP(&stepCount, "step", "s", 0, "Number of migrations to apply")
	upCmd.Flags().StringVarP(&toMigration, "to", "t", "", "Apply migrations up to specific ID")
}

func runUp(cmd *cobra.Command, args []string) error {
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

	// Filter pending migrations
	pendingMigrations := filterPendingMigrations(migrations, appliedMigrations)

	if len(pendingMigrations) == 0 {
		fmt.Printf("✅ All %s migrations are already applied\n", strings.ToUpper(migrationType))
		return nil
	}

	// Determine which migrations to apply
	var migrationsToApply []Migration
	if stepCount > 0 {
		if stepCount > len(pendingMigrations) {
			stepCount = len(pendingMigrations)
		}
		migrationsToApply = pendingMigrations[:stepCount]
	} else if toMigration != "" {
		migrationsToApply = filterMigrationsToID(pendingMigrations, toMigration)
	} else {
		migrationsToApply = pendingMigrations
	}

	if len(migrationsToApply) == 0 {
		fmt.Printf("ℹ️  No %s migrations to apply\n", strings.ToUpper(migrationType))
		return nil
	}

	// Apply migrations
	fmt.Printf("🚀 Applying %d %s migration(s)...\n", len(migrationsToApply), strings.ToUpper(migrationType))

	for _, migration := range migrationsToApply {
		if err := applyMigration(migration, migrationType); err != nil {
			return fmt.Errorf("failed to apply migration %s: %v", migration.ID, err)
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
	migrationDir := filepath.Join("database", "migration", migrationType)

	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		return []Migration{}, nil
	}

	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".up.sql") {
			continue
		}

		// Parse filename: 00001_create_table.up.sql
		parts := strings.Split(strings.TrimSuffix(file.Name(), ".up.sql"), "_")
		if len(parts) < 2 {
			continue
		}

		migrationID := parts[0]
		migrationName := strings.Join(parts[1:], "_")

		migrations = append(migrations, Migration{
			ID:   migrationID,
			Name: migrationName,
			Path: filepath.Join(migrationDir, file.Name()),
		})
	}

	// Sort migrations by ID
	sort.Slice(migrations, func(i, j int) bool {
		// Handle both sequential and timestamp formats
		if isSequentialID(migrations[i].ID) && isSequentialID(migrations[j].ID) {
			seqI, _ := strconv.Atoi(migrations[i].ID)
			seqJ, _ := strconv.Atoi(migrations[j].ID)
			return seqI < seqJ
		}
		return migrations[i].ID < migrations[j].ID
	})

	return migrations, nil
}

func isSequentialID(id string) bool {
	_, err := strconv.Atoi(id)
	return err == nil
}

func getAppliedMigrations(migrationType string) ([]string, error) {
	// TODO: Implement database connection and migration table query
	// For now, return empty slice (assume no migrations applied)
	return []string{}, nil
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
	// TODO: Implement actual database execution
	// For now, just simulate the process

	// Read migration file
	content, err := os.ReadFile(migration.Path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %v", err)
	}

	// TODO: Execute SQL against database
	fmt.Printf("   Executing: %s\n", string(content))

	// TODO: Record migration as applied in database

	return nil
}
