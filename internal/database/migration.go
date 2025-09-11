package database

import (
	"fmt"
	"strings"
	"time"

	"go.risoftinc.com/elsa/constants"
	"gorm.io/gorm"
)

// MigrationRecord represents a migration record in the database
type MigrationRecord struct {
	ID            uint      `gorm:"primaryKey"`
	MigrationID   string    `gorm:"uniqueIndex;not null;"` // varchar(255) for MySQL compatibility
	Name          string    `gorm:"not null;"`             // varchar(255) for MySQL compatibility
	Type          string    `gorm:"not null;"`             // varchar(10) for ddl/dml
	AppliedAt     time.Time `gorm:"not null"`
	Checksum      string    `gorm:"not null;"` // varchar(64) for checksum
	ExecutionTime int64     `gorm:"not null"`  // in milliseconds
}

// TableName specifies the table name for MigrationRecord
func (MigrationRecord) TableName() string {
	return constants.MigrationsTableName
}

// MigrationExecutor handles migration execution
type MigrationExecutor struct {
	db *gorm.DB
}

// NewMigrationExecutor creates a new migration executor
func NewMigrationExecutor(db *gorm.DB) *MigrationExecutor {
	return &MigrationExecutor{db: db}
}

// EnsureMigrationTable ensures the migrations table exists
func (me *MigrationExecutor) EnsureMigrationTable() error {
	// Check if table already exists by trying to query it
	var count int64
	err := me.db.Table(constants.MigrationsTableName).Count(&count).Error

	if err == nil {
		// Table already exists, no need to create
		return nil
	}

	// Table doesn't exist, create it with appropriate schema
	// Get database driver name to create appropriate table schema
	dbType := me.db.Dialector.Name()

	var createTableSQL string

	switch dbType {
	case constants.DriverMySQL:
		createTableSQL = constants.MySQLCreateTableSQL
	case constants.DriverPostgres:
		createTableSQL = constants.PostgresCreateTableSQL
	case constants.DriverSQLite:
		createTableSQL = constants.SQLiteCreateTableSQL
	default:
		// Fallback to generic SQL that should work on most databases
		createTableSQL = constants.GenericCreateTableSQL
	}

	return me.db.Exec(createTableSQL).Error
}

// GetAppliedMigrations retrieves all applied migrations
func (me *MigrationExecutor) GetAppliedMigrations(migrationType string) ([]string, error) {
	var records []MigrationRecord
	if err := me.db.Where(constants.TypeField+" = ?", migrationType).Find(&records).Error; err != nil {
		return nil, fmt.Errorf(constants.ErrFailedGetMigrations, err)
	}

	var migrationIDs []string
	for _, record := range records {
		migrationIDs = append(migrationIDs, record.MigrationID)
	}

	return migrationIDs, nil
}

// RecordMigration records a migration as applied
func (me *MigrationExecutor) RecordMigration(migrationID, name, migrationType, checksum string, executionTime int64) error {
	record := MigrationRecord{
		MigrationID:   migrationID,
		Name:          name,
		Type:          migrationType,
		AppliedAt:     time.Now(),
		Checksum:      checksum,
		ExecutionTime: executionTime,
	}

	if err := me.db.Create(&record).Error; err != nil {
		return fmt.Errorf(constants.ErrFailedRecord, err)
	}

	return nil
}

// RemoveMigration removes a migration record (for rollback)
func (me *MigrationExecutor) RemoveMigration(migrationID string) error {
	if err := me.db.Where(constants.MigrationIDField+" = ?", migrationID).Delete(&MigrationRecord{}).Error; err != nil {
		return fmt.Errorf(constants.ErrFailedRemove, err)
	}

	return nil
}

// ExecuteMigration executes a migration SQL file
func (me *MigrationExecutor) ExecuteMigration(sqlContent string, migrationType string) error {
	// Split SQL content by semicolon and execute each statement
	statements := splitSQLStatements(sqlContent)

	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		// Execute the SQL statement
		if err := me.db.Exec(statement).Error; err != nil {
			return fmt.Errorf(constants.ErrFailedExecute, i+1, err, statement)
		}
	}

	return nil
}

// splitSQLStatements splits SQL content into individual statements
func splitSQLStatements(sqlContent string) []string {
	// Simple split by semicolon - in production you might want more sophisticated parsing
	statements := strings.Split(sqlContent, ";")

	var result []string
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt)
		}
	}

	return result
}

// GetMigrationChecksum calculates a simple checksum for migration content
func GetMigrationChecksum(content string) string {
	// Simple hash for now - in production you might want to use crypto/sha256
	var checksum uint32
	for _, char := range content {
		checksum = ((checksum << 5) + checksum) + uint32(char)
	}
	return fmt.Sprintf("%x", checksum)
}
