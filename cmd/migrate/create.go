package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/risoftinc/elsa/constants"
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
	createCmd.Flags().BoolVarP(&useTimestamp, constants.TimestampFormatName, "t", true, "Use timestamp format (YYYYMMDDHHMMSSmmm) with milliseconds - default")
	createCmd.Flags().BoolVarP(&useSequential, constants.SequentialFormatName, "s", false, "Use sequential format (00001, 00002, etc.) instead of timestamp")
	createCmd.Flags().StringVarP(&customPath, "path", "p", "", "Custom folder path for migrations (default: database/migration/[ddl|dml])")
}

func runCreate(cmd *cobra.Command, args []string) error {
	migrationType := args[0]
	migrationName := args[1]

	// Validate migration type
	if migrationType != constants.MigrationTypeDDL && migrationType != constants.MigrationTypeDML {
		return fmt.Errorf(constants.ErrInvalidMigrationType, migrationType)
	}

	// Determine migration directory
	var migrationDir string
	if customPath != "" {
		// Use custom path with migration type subfolder
		migrationDir = filepath.Join(customPath, migrationType)
	} else {
		// Use default path
		migrationDir = filepath.Join(constants.DefaultMigrationBaseDir, constants.DefaultMigrationDir, migrationType)
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
			return fmt.Errorf(constants.ErrFailedGetSequential, err)
		}
		migrationID = fmt.Sprintf(constants.SequentialFormat, nextSeq)
		fmt.Printf(constants.InfoUsingSequential, migrationID)
	} else {
		// Use timestamp format: YYYYMMDDHHMMSSmmm (default) - includes milliseconds
		migrationID = time.Now().Format(constants.TimestampFormat) + fmt.Sprintf("%03d", time.Now().Nanosecond()/1000000)
		fmt.Printf(constants.InfoUsingTimestamp, migrationID)
	}

	// Create migration directory if not exists
	if err := os.MkdirAll(migrationDir, constants.MigrationDirPerm); err != nil {
		return fmt.Errorf(constants.ErrFailedCreateDir, err)
	}

	// Generate file names
	upFileName := fmt.Sprintf(constants.MigrationNameFormat+constants.UpMigrationExtension, migrationID, migrationName)
	downFileName := fmt.Sprintf(constants.MigrationNameFormat+constants.DownMigrationExtension, migrationID, migrationName)

	upFilePath := filepath.Join(migrationDir, upFileName)
	downFilePath := filepath.Join(migrationDir, downFileName)

	// Create up migration file
	upContent := generateUpMigrationContent(migrationType, migrationName)
	if err := os.WriteFile(upFilePath, []byte(upContent), constants.MigrationFilePerm); err != nil {
		return fmt.Errorf(constants.ErrFailedCreateFile, err)
	}

	// Create down migration file
	downContent := generateDownMigrationContent(migrationType, migrationName)
	if err := os.WriteFile(downFilePath, []byte(downContent), constants.MigrationFilePerm); err != nil {
		return fmt.Errorf(constants.ErrFailedCreateFile, err)
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
	case constants.MigrationTypeDDL:
		return fmt.Sprintf(constants.DDLUpTemplate, migrationName, migrationName)
	case constants.MigrationTypeDML:
		return fmt.Sprintf(constants.DMLUpTemplate, migrationName, migrationName)
	default:
		return ""
	}
}

func generateDownMigrationContent(migrationType, migrationName string) string {
	switch migrationType {
	case constants.MigrationTypeDDL:
		return fmt.Sprintf(constants.DDLDownTemplate, migrationName, migrationName)
	case constants.MigrationTypeDML:
		return fmt.Sprintf(constants.DMLDownTemplate, migrationName, migrationName)
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
	newFormat := constants.SequentialFormatName
	if !useSequential {
		newFormat = constants.TimestampFormatName
	}

	if existingFormat != newFormat {
		return fmt.Errorf("cannot create %s migration in folder that already contains %s migrations. Please use --%s flag or create in a different folder", newFormat, existingFormat, existingFormat)
	}

	return nil
}

// determineMigrationFormat determines if a migration file uses sequential or timestamp format
func determineMigrationFormat(fileName string) string {
	// Remove .up.sql or .down.sql suffix
	baseName := strings.TrimSuffix(strings.TrimSuffix(fileName, ".up.sql"), constants.DownMigrationExtension)

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
			return constants.SequentialFormatName
		}
	} else if len(id) == 17 {
		// Check if it's timestamp format (YYYYMMDDHHMMSSmmm)
		if _, err := time.Parse(constants.TimestampFormat, id[:14]); err == nil {
			return constants.TimestampFormatName
		}
	}

	return "unknown"
}

// GetMigrationPath returns the migration directory path for a given type
// This function can be used by other migration commands to get the correct path
func GetMigrationPath(migrationType string, customPath string) string {
	if customPath != "" {
		return filepath.Join(customPath, migrationType)
	}
	return filepath.Join(constants.DefaultMigrationBaseDir, constants.DefaultMigrationDir, migrationType)
}

// GetAvailableMigrationsWithPath returns available migrations from a specific path
// This function can be used by other migration commands to get migrations with custom path
func GetAvailableMigrationsWithPath(migrationType string, customPath string) ([]Migration, error) {
	migrationDir := GetMigrationPath(migrationType, customPath)

	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		return []Migration{}, nil
	}

	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), constants.UpMigrationExtension) {
			continue
		}

		// Parse filename: 00001_create_table.up.sql
		parts := strings.Split(strings.TrimSuffix(file.Name(), constants.UpMigrationExtension), "_")
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
