package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create [ddl|dml] [migration_name]",
		Short: "Create new database migration",
		Long: `Create new database migration files for DDL or DML operations.
By default, uses timestamp format (YYYYMMDDHHMMSSmmm) for migration IDs, including milliseconds.
		
Examples:
  elsa migration create ddl create_users_table          			# Uses timestamp with milliseconds (default)
  elsa migration create dml seed_users_data             			# Uses timestamp with milliseconds (default)
  elsa migration create ddl create_table --sequential   			# Uses sequential format
  elsa migration create ddl create_table --path custom/migrations  	# Custom folder path`,
		Args: cobra.ExactArgs(2),
		RunE: runCreate,
	}

	useTimestamp  bool
	useSequential bool
	customPath    string
)

func init() {
	createCmd.Flags().BoolVarP(&useTimestamp, "timestamp", "t", true, "Use timestamp format (YYYYMMDDHHMMSSmmm) with milliseconds - default")
	createCmd.Flags().BoolVarP(&useSequential, "sequential", "s", false, "Use sequential format (00001, 00002, etc.) instead of timestamp")
	createCmd.Flags().StringVarP(&customPath, "path", "p", "", "Custom folder path for migrations (default: database/migration/[ddl|dml])")
}

func runCreate(cmd *cobra.Command, args []string) error {
	migrationType := args[0]
	migrationName := args[1]

	// Validate migration type
	if migrationType != "ddl" && migrationType != "dml" {
		return fmt.Errorf("migration type must be 'ddl' or 'dml', got: %s", migrationType)
	}

	// Determine migration directory
	var migrationDir string
	if customPath != "" {
		// Use custom path with migration type subfolder
		migrationDir = filepath.Join(customPath, migrationType)
	} else {
		// Use default path
		migrationDir = filepath.Join("database", "migration", migrationType)
	}

	// Validate folder format consistency
	if err := validateFolderFormatConsistency(migrationDir, useSequential); err != nil {
		return err
	}

	// Generate migration ID
	var migrationID string
	if useSequential {
		// Use sequential format: 00001, 00002, etc.
		nextSeq, err := getNextSequentialNumber(migrationType)
		if err != nil {
			return fmt.Errorf("failed to get next sequential number: %v", err)
		}
		migrationID = fmt.Sprintf("%05d", nextSeq)
		fmt.Printf("🔢 Using sequential format: %s\n", migrationID)
	} else {
		// Use timestamp format: YYYYMMDDHHMMSSmmm (default) - includes milliseconds
		migrationID = time.Now().Format("20060102150405") + fmt.Sprintf("%03d", time.Now().Nanosecond()/1000000)
		fmt.Printf("📅 Using timestamp format: %s\n", migrationID)
	}

	// Create migration directory if not exists
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		return fmt.Errorf("failed to create migration directory: %v", err)
	}

	// Generate file names
	upFileName := fmt.Sprintf("%s_%s.up.sql", migrationID, migrationName)
	downFileName := fmt.Sprintf("%s_%s.down.sql", migrationID, migrationName)

	upFilePath := filepath.Join(migrationDir, upFileName)
	downFilePath := filepath.Join(migrationDir, downFileName)

	// Create up migration file
	upContent := generateUpMigrationContent(migrationType, migrationName)
	if err := os.WriteFile(upFilePath, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to create up migration file: %v", err)
	}

	// Create down migration file
	downContent := generateDownMigrationContent(migrationType, migrationName)
	if err := os.WriteFile(downFilePath, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to create down migration file: %v", err)
	}

	fmt.Printf("✅ Created migration files:\n")
	fmt.Printf("   Folder: %s\n", migrationDir)
	fmt.Printf("   Up:   %s\n", upFilePath)
	fmt.Printf("   Down: %s\n", downFilePath)

	return nil
}

func getNextSequentialNumber(migrationType string) (int, error) {
	var migrationDir string
	if customPath != "" {
		// Use custom path with migration type subfolder
		migrationDir = filepath.Join(customPath, migrationType)
	} else {
		// Use default path
		migrationDir = filepath.Join("database", "migration", migrationType)
	}

	// Check if directory exists
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		return 1, nil
	}

	// Read directory and find highest sequential number
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return 0, err
	}

	maxSeq := 0
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		// Extract sequential number from filename (e.g., "00001_create_table.up.sql")
		if len(file.Name()) >= 5 {
			seqStr := file.Name()[:5]
			if seq, err := fmt.Sscanf(seqStr, "%d", &maxSeq); err == nil && seq > 0 {
				// Keep track of max sequence
			}
		}
	}

	return maxSeq + 1, nil
}

func generateUpMigrationContent(migrationType, migrationName string) string {
	switch migrationType {
	case "ddl":
		return fmt.Sprintf(`-- Migration: %s
-- Type: DDL
-- Description: %s

-- Add your DDL statements here
-- Example:
-- CREATE TABLE users (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     email VARCHAR(255) UNIQUE NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

`, migrationName, migrationName)
	case "dml":
		return fmt.Sprintf(`-- Migration: %s
-- Type: DML
-- Description: %s

-- Add your DML statements here
-- Example:
-- INSERT INTO users (name, email) VALUES 
--     ('John Doe', 'john@example.com'),
--     ('Jane Smith', 'jane@example.com');

`, migrationName, migrationName)
	default:
		return ""
	}
}

func generateDownMigrationContent(migrationType, migrationName string) string {
	switch migrationType {
	case "ddl":
		return fmt.Sprintf(`-- Rollback Migration: %s
-- Type: DDL
-- Description: %s

-- Add your rollback DDL statements here
-- Example:
-- DROP TABLE IF EXISTS users;

`, migrationName, migrationName)
	case "dml":
		return fmt.Sprintf(`-- Rollback Migration: %s
-- Type: DML
-- Description: %s

-- Add your rollback DML statements here
-- Example:
-- DELETE FROM users WHERE email IN ('john@example.com', 'jane@example.com');

`, migrationName, migrationName)
	default:
		return ""
	}
}

// validateFolderFormatConsistency checks if the folder already contains migrations
// and ensures the new migration uses the same format
func validateFolderFormatConsistency(migrationDir string, useSequential bool) error {
	// Check if directory exists and has files
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		// Directory doesn't exist, no validation needed
		return nil
	}

	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %v", err)
	}

	// Filter only .sql files
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	if len(sqlFiles) == 0 {
		// No existing migrations, no validation needed
		return nil
	}

	// Determine existing format by checking first migration file
	firstFile := sqlFiles[0]
	existingFormat := determineMigrationFormat(firstFile)

	// Check if all existing files use the same format
	for _, fileName := range sqlFiles {
		fileFormat := determineMigrationFormat(fileName)
		if fileFormat != existingFormat {
			return fmt.Errorf("folder contains mixed migration formats. Found both %s and %s formats. Please use only one format per folder", existingFormat, fileFormat)
		}
	}

	// Check if new migration format matches existing format
	newFormat := "sequential"
	if !useSequential {
		newFormat = "timestamp"
	}

	if existingFormat != newFormat {
		return fmt.Errorf("cannot create %s migration in folder that already contains %s migrations. Please use --%s flag or create in a different folder", newFormat, existingFormat, existingFormat)
	}

	return nil
}

// determineMigrationFormat determines if a migration file uses sequential or timestamp format
func determineMigrationFormat(fileName string) string {
	// Remove .up.sql or .down.sql suffix
	baseName := strings.TrimSuffix(strings.TrimSuffix(fileName, ".up.sql"), ".down.sql")

	// Extract the ID part (before first underscore)
	parts := strings.Split(baseName, "_")
	if len(parts) == 0 {
		return "unknown"
	}

	id := parts[0]

	// Check if it's sequential (5 digits) or timestamp (17 digits)
	if len(id) == 5 {
		// Check if it's numeric
		if _, err := fmt.Sscanf(id, "%d", new(int)); err == nil {
			return "sequential"
		}
	} else if len(id) == 17 {
		// Check if it's timestamp format (YYYYMMDDHHMMSSmmm)
		if _, err := time.Parse("20060102150405", id[:14]); err == nil {
			return "timestamp"
		}
	}

	return "unknown"
}
