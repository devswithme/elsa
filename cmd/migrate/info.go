package migrate

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	infoCmd = &cobra.Command{
		Use:   "info [ddl|dml]",
		Short: "Show detailed migration information",
		Long: `Show detailed migration information for DDL or DML operations.
		
Examples:
  elsa migration info                    # Show info of all migrations
  elsa migration info ddl               # Show info of DDL migrations
  elsa migration info dml               # Show info of DML migrations`,
		Args: cobra.MaximumNArgs(1),
		RunE: runInfo,
	}
)

func runInfo(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// Show info for both DDL and DML
		fmt.Println("📋 Migration Information Overview")
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
	if migrationType != "ddl" && migrationType != "dml" {
		return fmt.Errorf("migration type must be 'ddl' or 'dml', got: %s", migrationType)
	}

	return showMigrationInfo(migrationType)
}

func showMigrationInfo(migrationType string) error {
	fmt.Printf("\n🔧 %s Migrations Information:\n", strings.ToUpper(migrationType))
	fmt.Printf("%s\n", strings.Repeat("-", 40))

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
		downPath := strings.Replace(upPath, ".up.sql", ".down.sql", 1)

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
