package migrate

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.risoftinc.com/elsa/constants"
)

var (
	infoCmd = &cobra.Command{
		Use:   "info [ddl|dml]",
		Short: "Show detailed migration information",
		Long: `Show detailed migration information for DDL or DML operations.
		
Examples:
  elsa migration info                    # Show info of all migrations
  elsa migration info ddl               # Show info of DDL migrations
  elsa migration info dml               # Show info of DML migrations
  elsa migration info ddl --path custom/migrations  # Use custom migration path`,
		Args: cobra.MaximumNArgs(1),
		RunE: runInfo,
	}

	infoCustomPath string
	infoConnection string
)

func runInfo(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// Show info for both DDL and DML
		fmt.Println(constants.InfoOverviewHeader)
		fmt.Println(strings.Repeat("=", 50))

		if err := showMigrationInfo("ddl"); err != nil {
			return err
		}

		fmt.Println()
		if err := showMigrationInfo("dml"); err != nil {
			return err
		}

		return nil
	}

	migrationType := args[0]
	if migrationType != constants.MigrationTypeDDL && migrationType != constants.MigrationTypeDML {
		return fmt.Errorf(constants.ErrInvalidMigrationType, migrationType)
	}

	return showMigrationInfo(migrationType)
}

func init() {
	infoCmd.Flags().StringVarP(&infoCustomPath, "path", "p", "", "Custom migration path")
	infoCmd.Flags().StringVarP(&infoConnection, "connection", "c", "", "Database connection string")
}

func showMigrationInfo(migrationType string) error {
	fmt.Printf(constants.InfoDDLHeader, strings.ToUpper(migrationType))
	fmt.Printf(constants.InfoDDLSeparator, strings.Repeat("-", 40))

	// Get available migrations
	migrations, err := GetAvailableMigrationsWithPath(migrationType, infoCustomPath)
	if err != nil {
		return fmt.Errorf(constants.ErrFailedShowInfo, err)
	}

	if len(migrations) == 0 {
		fmt.Printf(constants.InfoNoMigrationsFound, strings.ToUpper(migrationType))
		return nil
	}

	// Get applied migrations from database (optional, only if connection is provided)
	var appliedMigrations []string
	var err2 error

	// Try to get applied migrations if connection is provided
	appliedMigrations, err2 = getAppliedMigrationsWithConnection(migrationType, infoConnection)
	if err2 != nil {
		fmt.Printf(constants.InfoWarningDBConnect, err2)
		fmt.Printf(constants.InfoShowingFileBasedInfo)
		appliedMigrations = []string{} // Empty slice to show all as pending
	}

	// Create applied map for quick lookup
	appliedMap := make(map[string]bool)
	for _, id := range appliedMigrations {
		appliedMap[id] = true
	}

	// Display detailed migration info
	for i, migration := range migrations {
		if i > 0 {
			fmt.Println()
		}

		status := "❌ Pending"
		if appliedMap[migration.ID] {
			status = "✅ Applied"
		}

		fmt.Printf("   Migration ID: %s\n", migration.ID)
		fmt.Printf("   Name: %s\n", migration.Name)
		fmt.Printf("   Status: %s\n", status)
		fmt.Printf("   Type: %s\n", strings.ToUpper(migrationType))

		// Show file paths
		upPath := migration.Path
		downPath := strings.Replace(upPath, ".up.sql", constants.DownMigrationExtension, 1)

		fmt.Printf("   Up File: %s\n", upPath)
		fmt.Printf("   Down File: %s\n", downPath)

		// Show file sizes
		if upSize, err := getFileSize(upPath); err == nil {
			fmt.Printf("   Up File Size: %d bytes\n", upSize)
		}

		if downSize, err := getFileSize(downPath); err == nil {
			fmt.Printf("   Down File Size: %d bytes\n", downSize)
		}

		// Show preview of migration content
		if content, err := os.ReadFile(upPath); err == nil {
			preview := string(content)
			if len(preview) > 200 {
				preview = preview[:200] + "..."
			}
			fmt.Printf("   Preview:\n%s\n", indentString(preview, 6))
		}
	}

	// Summary
	total := len(migrations)
	applied := len(appliedMigrations)
	pending := total - applied

	fmt.Printf("\n📊 Summary:\n")
	fmt.Printf("   Total %s Migrations: %d\n", strings.ToUpper(migrationType), total)
	fmt.Printf("   Applied: %d\n", applied)
	fmt.Printf("   Pending: %d\n", pending)

	return nil
}

func getFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func indentString(s string, indent int) string {
	lines := strings.Split(s, "\n")
	indentedLines := make([]string, len(lines))

	for i, line := range lines {
		indentedLines[i] = strings.Repeat(" ", indent) + line
	}

	return strings.Join(indentedLines, "\n")
}
