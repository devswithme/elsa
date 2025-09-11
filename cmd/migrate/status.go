package migrate

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"go.risoftinc.com/elsa/constants"
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
		showTypes = append(showTypes, constants.MigrationTypeDDL)
	} else if showDML {
		showTypes = append(showTypes, constants.MigrationTypeDML)
	} else {
		// Show both if no specific type specified
		showTypes = []string{constants.MigrationTypeDDL, constants.MigrationTypeDML}
	}

	fmt.Printf(constants.StatusOverviewHeader)
	fmt.Printf(constants.StatusOverviewSeparator)

	for _, migrationType := range showTypes {
		if err := showMigrationStatus(migrationType); err != nil {
			return fmt.Errorf(constants.ErrFailedShowStatus, migrationType, err)
		}
		fmt.Println()
	}

	return nil
}

func showMigrationStatus(migrationType string) error {
	// Get available migrations
	availableMigrations, err := GetAvailableMigrationsWithPath(migrationType, statusCustomPath)
	if err != nil {
		return fmt.Errorf(constants.ErrFailedShowStatus, migrationType, err)
	}

	// Get applied migrations from database
	appliedMigrations, err := getAppliedMigrations(migrationType)
	if err != nil {
		// If database connection fails, show only file-based status
		fmt.Printf(constants.InfoWarningDBConnect, err)
		fmt.Printf(constants.InfoShowingFileBased)

		// Display available migrations without status
		fmt.Printf(constants.StatusDDLHeader, strings.ToUpper(migrationType))
		fmt.Printf(constants.StatusTableHeader)
		fmt.Printf(constants.StatusTableSeparator)

		for _, migration := range availableMigrations {
			displayID := migration.ID
			if len(displayID) < 17 {
				displayID = displayID + strings.Repeat(" ", 17-len(displayID))
			}
			fmt.Printf("   %s  %s\n", displayID, migration.Name)
		}

		fmt.Printf(constants.InfoDatabaseStatusUnavailable, len(availableMigrations))
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
		status := constants.StatusPending
		if appliedMap[migration.ID] {
			status = constants.StatusApplied
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
