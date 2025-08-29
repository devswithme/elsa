package constants

// Migration command constants
const (
	// Command usage
	MigrateUse = "migration"

	// Command descriptions
	MigrateShort = "Database migration commands"

	// Command long description
	MigrateLong = `Database migration commands for managing DDL and DML changes.
		
Examples:
  elsa migration connect                                            Connect to database interactively
  elsa migration connect -c "sqlite://elsa.db"                     Connect using connection string flag
  elsa migration create ddl create_users_table                      Create new DDL migration (timestamp with ms)
  elsa migration create dml seed_users_data                         Create new DML migration (timestamp with ms)
  elsa migration create ddl create_table --sequential               Create with sequential format
  elsa migration create ddl create_table --path custom/migrations   Custom folder path
  elsa migration up ddl                                             Apply all DDL migrations
  elsa migration down dml                                           Rollback last DML migration
  elsa migration status                                             Show migration status
  elsa migration refresh ddl                                        Refresh all DDL migrations`
)

// Migration type constants (reusing from database constants)
const (
// MigrationTypeDDL and MigrationTypeDML are defined in constants/database.go
)

// Migration file extensions
const (
	UpMigrationExtension   = ".up.sql"
	DownMigrationExtension = ".down.sql"
)

// Migration directory structure
const (
	DefaultMigrationBaseDir = "database"
	DefaultMigrationDir     = "migration"
)

// Migration file format constants
const (
	TimestampFormatName  = "timestamp"
	SequentialFormatName = "sequential"
	TimestampFormat      = "20060102150405"
	SequentialFormat     = "%05d"
	MigrationNameFormat  = "%s_%s"
)

// Migration file permissions
const (
	MigrationDirPerm  = 0755
	MigrationFilePerm = 0644
)

// Error messages
const (
	ErrInvalidMigrationType       = "migration type must be 'ddl' or 'dml', got: %s"
	ErrInvalidConnectionString    = "invalid connection string: %s"
	ErrFailedConnectDB            = "failed to connect to database: %v"
	ErrFailedEnsureTable          = "failed to ensure migration table: %v"
	ErrFailedReadFile             = "failed to read migration file: %v"
	ErrFailedExecuteMigration     = "failed to execute migration: %v"
	ErrFailedRecordMigration      = "failed to record migration: %v"
	ErrFailedRollbackMigration    = "failed to execute rollback migration: %v"
	ErrFailedRemoveRecord         = "failed to remove migration record: %v"
	ErrFailedCreateDir            = "failed to create migration directory: %v"
	ErrFailedCreateFile           = "failed to create migration file: %v"
	ErrFailedGetSequential        = "failed to get next sequential number: %v"
	ErrFailedValidateFolder       = "failed to validate folder format consistency: %v"
	ErrFailedGetAppliedMigrations = "failed to get applied migrations: %v"
	ErrFailedApplyMigration       = "failed to apply migration %s: %v"
	ErrFailedRollbackAll          = "failed to rollback migrations: %v"
	ErrFailedApplyAll             = "failed to apply migrations: %v"
	ErrFailedShowStatus           = "failed to show %s migration status: %v"
	ErrFailedShowInfo             = "failed to get available migrations: %v"
	ErrFailedConnectDBStatus      = "Could not connect to database: %v"
	ErrFailedConnectDBInfo        = "Could not connect to database: %v"
)

// Success messages
const (
	SuccessExecuted        = "   " + SuccessEmoji + " Executed in %dms\n"
	SuccessRolledBack      = "   " + SuccessEmoji + " Rolled back in %dms\n"
	SuccessConnected       = SuccessEmoji + " Successfully connected to database!\n"
	SuccessTableExists     = SuccessEmoji + " Migration table exists and is ready!\n"
	SuccessAllApplied      = SuccessEmoji + " All %s migrations are already applied\n"
	SuccessAllRolledBack   = SuccessEmoji + " All %s migrations rolled back successfully\n"
	SuccessAllAppliedAgain = SuccessEmoji + " All %s migrations applied successfully\n"
	SuccessRefreshed       = PartyEmoji + " Successfully refreshed all %s migrations!\n"
)

// Info messages
const (
	InfoUsingSequential           = NumberEmoji + " Using sequential format: %s\n"
	InfoUsingTimestamp            = CalendarEmoji + " Using timestamp format: %s\n"
	InfoUsingConnectionFlag       = LinkEmoji + " Using connection string from flag:\n"
	InfoLoadedFromEnv             = FolderEmoji + " Loaded configuration from .env file:\n"
	InfoNoEnvFile                 = InfoEmoji + " No .env file found or invalid configuration\n"
	InfoTestingConnection         = PlugEmoji + " Testing database connection...\n"
	InfoEnsuringTable             = ClipboardEmoji + " Ensuring migration table exists...\n"
	InfoNoMigrations              = InfoEmoji + " No %s migrations have been applied\n"
	InfoNoMigrationsToRollback    = InfoEmoji + " No %s migrations to rollback\n"
	InfoNoMigrationsFoundStatus   = InfoEmoji + " No %s migrations found\n"
	InfoNoMigrationsToApply       = InfoEmoji + " No %s migrations to apply\n"
	InfoNoMigrationsToRollbackAll = InfoEmoji + " No %s migrations have been applied\n"
	InfoApplyingMigrations        = RocketEmoji + " Applying %d %s migration(s)...\n"
	InfoRollingBackMigrations     = RestartEmoji + " Rolling back %d %s migration(s)...\n"
	InfoRefreshingMigrations      = RestartEmoji + " Refreshing all %s migrations...\n"
	InfoStep1Rollback             = OutboxEmoji + " Step 1: Rolling back all applied %s migrations...\n"
	InfoStep2Apply                = InboxEmoji + " Step 2: Applying all %s migrations...\n"
	InfoMigrationTableReady       = SuccessEmoji + " Migration table ready!\n"
	InfoConfigSaved               = FloppyDiskEmoji + " Configuration saved to .env file\n"
	InfoWarningSaveConfig         = WarningEmoji + " Warning: Could not save configuration to .env file: %v\n"
	InfoWarningFileNotFound       = WarningEmoji + " Warning: Migration file for ID %s not found, skipping\n"
	InfoWarningDBConnect          = WarningEmoji + " Warning: Could not connect to database: %v\n"
	InfoShowingFileBased          = "   Showing file-based status only (no applied/pending info)\n"
	InfoShowingFileBasedInfo      = "   Showing file-based information only\n"
	InfoDatabaseStatusUnavailable = "   Summary: %d total (database status unavailable)\n"
)

// Connection info format
const (
	ConnectionInfoFormat = "   Connection: %s\n"
	DriverInfoFormat     = "   Driver:     %s\n"
	HostInfoFormat       = "   Host:       %s\n"
	PortInfoFormat       = "   Port:       %s\n"
	DatabaseInfoFormat   = "   Database:   %s\n"
	UsernameInfoFormat   = "   Username:   %s\n"
	PasswordInfoFormat   = "   Password:  %s\n"
)

// Migration file content templates
const (
	// DDL migration templates
	DDLUpTemplate = `-- Migration: %s
-- Type: DDL
-- Description: %s

-- Add your DDL statements here
-- Example:
-- CREATE TABLE table_name (
--     id INT PRIMARY KEY AUTO_INCREMENT,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

`

	DDLDownTemplate = `-- Migration: %s (Rollback)
-- Type: DDL
-- Description: %s

-- Add your rollback DDL statements here
-- Example:
-- DROP TABLE IF EXISTS table_name;

`

	// DML migration templates
	DMLUpTemplate = `-- Migration: %s
-- Type: DML
-- Description: %s

-- Add your DML statements here
-- Example:
-- INSERT INTO table_name (name, created_at) VALUES 
--     ('Sample Data 1', NOW()),
--     ('Sample Data 2', NOW());

`

	DMLDownTemplate = `-- Migration: %s (Rollback)
-- Type: DML
-- Description: %s

-- Add your rollback DML statements here
-- Example:
-- DELETE FROM table_name WHERE name IN ('Sample Data 1', 'Sample Data 2');

`
)

// Status and info display constants
const (
	StatusOverviewHeader       = ChartEmoji + " Migration Status Overview\n"
	StatusOverviewSeparator    = "==================================================\n\n"
	StatusDDLHeader            = WrenchEmoji + " %s Migrations:\n"
	StatusTableHeader          = "   ID                 Name\n"
	StatusTableSeparator       = "   ----------------------------------------\n"
	StatusPending              = ErrorEmoji + " Pending"
	StatusApplied              = SuccessEmoji + " Applied"
	InfoOverviewHeader         = ClipboardEmoji + " Migration Information Overview\n"
	InfoOverviewSeparator      = "==================================================\n"
	InfoDDLHeader              = WrenchEmoji + " %s Migrations Information:\n"
	InfoDDLSeparator           = "%s\n"
	InfoNoMigrationsFound      = "   No %s migrations found\n"
	InfoMigrationDetails       = "   Migration: %s\n"
	InfoMigrationStatus        = "   Status: %s\n"
	InfoMigrationPath          = "   Path: %s\n"
	InfoMigrationType          = "   Type: %s\n"
	InfoMigrationAppliedAt     = "   Applied At: %s\n"
	InfoMigrationExecutionTime = "   Execution Time: %dms\n"
	InfoMigrationChecksum      = "   Checksum: %s\n"
	RefreshHeader              = RestartEmoji + " Refreshing all %s migrations...\n"
	RefreshSeparator           = "==================================================\n\n"
	RefreshSuccessMessage      = "   All migrations have been rolled back and reapplied.\n"
)
