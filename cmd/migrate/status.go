package migrate

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	statusCmd = &cobra.Command{
		Use:   "status [ddl|dml]",
		Short: "Show migration status",
		Long: `Show migration status for DDL or DML operations.
		
Examples:
  elsa migration status                   # Show status of all migrations
  elsa migration status ddl               # Show status of DDL migrations
  elsa migration status dml               # Show status of DML migrations`,
		Args: cobra.MaximumNArgs(1),
		RunE: runStatus,
	}
)

func runStatus(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// Show status for both DDL and DML
		fmt.Println("📊 Migration Status Overview")
		fmt.Println(strings.Repeat("=", 50))

		if err := showMigrationStatus("ddl"); err != nil {
			return err
		}

		fmt.Println()
		if err := showMigrationStatus("dml"); err != nil {
			return err
		}

		return nil
	}

	migrationType := args[0]
	if migrationType != "ddl" && migrationType != "dml" {
		return fmt.Errorf("migration type must be 'ddl' or 'dml', got: %s", migrationType)
	}

	return showMigrationStatus(migrationType)
}

func showMigrationStatus(migrationType string) error {
	fmt.Printf("\n🔧 %s Migrations:\n", strings.ToUpper(migrationType))
	fmt.Printf("%s\n", strings.Repeat("-", 30))

	// Get available migrations
	migrations, err := getAvailableMigrations(migrationType)
	if err != nil {
		return fmt.Errorf("failed to get available migrations: %v", err)
	}

	if len(migrations) == 0 {
		fmt.Printf("   No %s migrations found\n", strings.ToUpper(migrationType))
		return nil
	}

	// Get applied migrations from database
	appliedMigrations, err := getAppliedMigrations(migrationType)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	// Create applied map for quick lookup
	appliedMap := make(map[string]bool)
	for _, id := range appliedMigrations {
		appliedMap[id] = true
	}

	// Display migration status
	fmt.Printf("   %-15s %-25s %-10s\n", "ID", "Name", "Status")
	fmt.Printf("   %s\n", strings.Repeat("-", 55))

	for _, migration := range migrations {
		status := "❌ Pending"
		if appliedMap[migration.ID] {
			status = "✅ Applied"
		}

		// Truncate name if too long
		name := migration.Name
		if len(name) > 23 {
			name = name[:20] + "..."
		}

		fmt.Printf("   %-15s %-25s %-10s\n", migration.ID, name, status)
	}

	// Summary
	total := len(migrations)
	applied := len(appliedMigrations)
	pending := total - applied

	fmt.Printf("\n   Summary: %d total, %d applied, %d pending\n", total, applied, pending)

	return nil
}
