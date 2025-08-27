package migrate

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show migration status",
		Long: `Show the status of all migrations (DDL and DML).
		
Examples:
  elsa migration status                   # Show status of all migrations
  elsa migration status --ddl            # Show only DDL migrations
  elsa migration status --dml            # Show only DML migrations
  elsa migration status --path custom/migrations  # Use custom migration path`,
		RunE: runStatus,
	}

	showDDL          bool
	showDML          bool
	statusCustomPath string
)

func init() {
	statusCmd.Flags().BoolVarP(&showDDL, "ddl", "d", false, "Show only DDL migrations")
	statusCmd.Flags().BoolVarP(&showDML, "dml", "m", false, "Show only DML migrations")
	statusCmd.Flags().StringVarP(&statusCustomPath, "path", "p", "", "Custom migration path")
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Determine which migration types to show
	showTypes := []string{}
	if showDDL {
		showTypes = append(showTypes, "ddl")
	} else if showDML {
		showTypes = append(showTypes, "dml")
	} else {
		// Show both if no specific type specified
		showTypes = []string{"ddl", "dml"}
	}

	fmt.Printf("📊 Migration Status Overview\n")
	fmt.Printf("==================================================\n\n")

	for _, migrationType := range showTypes {
		if err := showMigrationStatus(migrationType); err != nil {
			return fmt.Errorf("failed to show %s migration status: %v", migrationType, err)
		}
		fmt.Println()
	}

	return nil
}

func showMigrationStatus(migrationType string) error {
	// Get available migrations
	availableMigrations, err := GetAvailableMigrationsWithPath(migrationType, statusCustomPath)
	if err != nil {
		return fmt.Errorf("failed to get available migrations: %v", err)
	}

	// Get applied migrations from database
	appliedMigrations, err := getAppliedMigrations(migrationType)
	if err != nil {
		// If database connection fails, show only file-based status
		fmt.Printf("⚠️  Warning: Could not connect to database: %v\n", err)
		fmt.Printf("   Showing file-based status only (no applied/pending info)\n\n")

		// Display available migrations without status
		fmt.Printf("🔧 %s Migrations:\n", strings.ToUpper(migrationType))
		fmt.Printf("   ID                 Name\n")
		fmt.Printf("   ----------------------------------------\n")

		for _, migration := range availableMigrations {
			displayID := migration.ID
			if len(displayID) < 17 {
				displayID = displayID + strings.Repeat(" ", 17-len(displayID))
			}
			fmt.Printf("   %s  %s\n", displayID, migration.Name)
		}

		fmt.Printf("\n   Summary: %d total (database status unavailable)\n", len(availableMigrations))
		return nil
	}

	// Create applied migrations map for quick lookup
	appliedMap := make(map[string]bool)
	for _, id := range appliedMigrations {
		appliedMap[id] = true
	}

	// Sort migrations by ID
	sort.Slice(availableMigrations, func(i, j int) bool {
		// Handle both sequential and timestamp formats
		if isSequentialID(availableMigrations[i].ID) && isSequentialID(availableMigrations[j].ID) {
			seqI, _ := strconv.Atoi(availableMigrations[i].ID)
			seqJ, _ := strconv.Atoi(availableMigrations[j].ID)
			return seqI < seqJ
		}
		return availableMigrations[i].ID < availableMigrations[j].ID
	})

	// Display header
	fmt.Printf("🔧 %s Migrations:\n", strings.ToUpper(migrationType))
	fmt.Printf("   ID                 Name                    Status\n")
	fmt.Printf("   -------------------------------------------------------\n")

	// Display each migration
	for _, migration := range availableMigrations {
		status := "❌ Pending"
		if appliedMap[migration.ID] {
			status = "✅ Applied"
		}

		// Format ID for display (ensure consistent width)
		displayID := migration.ID
		if len(displayID) < 17 {
			displayID = displayID + strings.Repeat(" ", 17-len(displayID))
		}

		fmt.Printf("   %s  %-20s  %s\n", displayID, migration.Name, status)
	}

	// Show summary
	appliedCount := len(appliedMigrations)
	pendingCount := len(availableMigrations) - appliedCount

	fmt.Printf("\n   Summary: %d total, %d applied, %d pending\n",
		len(availableMigrations), appliedCount, pendingCount)

	return nil
}
