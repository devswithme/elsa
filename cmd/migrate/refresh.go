package migrate

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	refreshCmd = &cobra.Command{
		Use:   "refresh [ddl|dml]",
		Short: "Refresh database migrations",
		Long: `Refresh database migrations for DDL or DML operations.
This will rollback all applied migrations and then apply them again.
		
Examples:
  elsa migration refresh ddl                    # Refresh all DDL migrations
  elsa migration refresh dml                    # Refresh all DML migrations`,
		Args: cobra.ExactArgs(1),
		RunE: runRefresh,
	}
)

func runRefresh(cmd *cobra.Command, args []string) error {
	migrationType := args[0]

	// Validate migration type
	if migrationType != "ddl" && migrationType != "dml" {
		return fmt.Errorf("migration type must be 'ddl' or 'dml', got: %s", migrationType)
	}

	fmt.Printf("🔄 Refreshing %s migrations...\n", strings.ToUpper(migrationType))

	// Step 1: Rollback all applied migrations
	fmt.Printf("📤 Step 1: Rolling back all applied %s migrations...\n", strings.ToUpper(migrationType))

	// Create a temporary command to execute down
	tempDownCmd := &cobra.Command{}
	tempDownCmd.Flags().Int("step", 999, "Number of migrations to rollback")
	tempDownCmd.Flags().String("to", "", "Rollback migrations down to specific ID")

	// Execute down command with maximum step
	downStepCount = 999
	downToMigration = ""

	if err := runDown(tempDownCmd, args); err != nil {
		return fmt.Errorf("failed to rollback migrations: %v", err)
	}

	fmt.Printf("✅ Successfully rolled back all %s migrations\n", strings.ToUpper(migrationType))

	// Step 2: Apply all migrations again
	fmt.Printf("📥 Step 2: Applying all %s migrations...\n", strings.ToUpper(migrationType))

	// Reset up command flags
	stepCount = 0
	toMigration = ""

	if err := runUp(tempDownCmd, args); err != nil {
		return fmt.Errorf("failed to apply migrations: %v", err)
	}

	fmt.Printf("✅ Successfully applied all %s migrations\n", strings.ToUpper(migrationType))
	fmt.Printf("🎉 %s migrations refreshed successfully!\n", strings.ToUpper(migrationType))

	return nil
}
