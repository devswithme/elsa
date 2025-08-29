package constants

// Database driver constants
const (
	DriverSQLite     = "sqlite"
	DriverMySQL      = "mysql"
	DriverPostgres   = "postgres"
	DriverPostgreSQL = "postgresql"
)

// Database default values
const (
	DefaultHost         = "localhost"
	DefaultPort         = "3306"
	DefaultUsername     = "root"
	DefaultPassword     = ""
	DefaultDatabase     = "elsa.db"
	DefaultSSLMode      = "disable"
	DefaultCharset      = "utf8mb4"
	DefaultPortPostgres = "5432"
)

// Database connection string prefixes
const (
	SQLitePrefix      = "sqlite://"
	MySQLPrefix       = "mysql://"
	PostgresPrefix    = "postgresql://"
	PostgresAltPrefix = "postgres://"
)

// Database table and field constants
const (
	MigrationsTableName = "migrations"
	MigrationIDField    = "migration_id"
	NameField           = "name"
	TypeField           = "type"
	AppliedAtField      = "applied_at"
	ChecksumField       = "checksum"
	ExecutionTimeField  = "execution_time"
)

// Migration type constants
const (
	MigrationTypeDDL = "ddl"
	MigrationTypeDML = "dml"
)

// Database schema constants
const (
	// Field lengths for MySQL compatibility
	MaxMigrationIDLength = 255
	MaxNameLength        = 255
	MaxTypeLength        = 10
	MaxChecksumLength    = 64

	// Time zone for PostgreSQL
	PostgresTimeZone = "Asia/Jakarta"
)

// Error messages
const (
	ErrUnsupportedDriver   = "unsupported database driver: %s"
	ErrFailedConnect       = "failed to connect to database: %v"
	ErrFailedGetDB         = "failed to get underlying sql.DB: %v"
	ErrFailedPing          = "failed to ping database: %v"
	ErrFailedReadEnv       = "failed to read .env file: %v"
	ErrFailedRecord        = "failed to record migration: %v"
	ErrFailedRemove        = "failed to remove migration record: %v"
	ErrFailedExecute       = "failed to execute statement %d: %v\nSQL: %s"
	ErrFailedGetMigrations = "failed to get applied migrations: %v"
)

// SQL statements
const (
	// MySQL create table
	MySQLCreateTableSQL = `CREATE TABLE IF NOT EXISTS migrations (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		migration_id VARCHAR(255) NOT NULL UNIQUE,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(10) NOT NULL,
		applied_at DATETIME(3) NOT NULL,
		checksum VARCHAR(64) NOT NULL,
		execution_time BIGINT NOT NULL
	)`

	// PostgreSQL create table
	PostgresCreateTableSQL = `CREATE TABLE IF NOT EXISTS migrations (
		id BIGSERIAL PRIMARY KEY,
		migration_id VARCHAR(255) NOT NULL UNIQUE,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(10) NOT NULL,
		applied_at TIMESTAMP NOT NULL,
		checksum VARCHAR(64) NOT NULL,
		execution_time BIGINT NOT NULL
	)`

	// SQLite create table
	SQLiteCreateTableSQL = `CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		migration_id TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		applied_at DATETIME NOT NULL,
		checksum TEXT NOT NULL,
		execution_time INTEGER NOT NULL
	)`

	// Generic create table (fallback)
	GenericCreateTableSQL = `CREATE TABLE IF NOT EXISTS migrations (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		migration_id VARCHAR(255) NOT NULL UNIQUE,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(10) NOT NULL,
		applied_at TIMESTAMP NOT NULL,
		checksum VARCHAR(64) NOT NULL,
		execution_time BIGINT NOT NULL
	)`
)

// Connection string formats
const (
	SQLiteConnectionFormat   = "sqlite://%s"
	MySQLConnectionFormat    = "mysql://%s:%s@%s:%s/%s?charset=%s"
	PostgresConnectionFormat = "postgresql://%s:%s@%s:%s/%s?sslmode=%s"
)

// DSN formats
const (
	MySQLDSNFormat    = "%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local"
	PostgresDSNFormat = "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s"
)
