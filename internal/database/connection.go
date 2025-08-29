package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/risoftinc/elsa/constants"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Driver           string `mapstructure:"driver"`
	Host             string `mapstructure:"host"`
	Port             string `mapstructure:"port"`
	Username         string `mapstructure:"username"`
	Password         string `mapstructure:"password"`
	Database         string `mapstructure:"database"`
	SSLMode          string `mapstructure:"sslmode"`
	Charset          string `mapstructure:"charset"`
	ConnectionString string `mapstructure:"connection_string"`
}

// DefaultConfig returns default database configuration
func DefaultConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Driver:           constants.DriverSQLite,
		Host:             constants.DefaultHost,
		Port:             constants.DefaultPort,
		Username:         constants.DefaultUsername,
		Password:         constants.DefaultPassword,
		Database:         constants.DefaultDatabase,
		SSLMode:          constants.DefaultSSLMode,
		Charset:          constants.DefaultCharset,
		ConnectionString: "",
	}
}

// LoadFromEnv loads database configuration from environment variables
func LoadFromEnv() *DatabaseConfig {
	config := DefaultConfig()

	// Try to load from .env file first
	if err := loadDotEnv(); err == nil {
		// If .env file loaded successfully, try to get MIGRATE_CONNECTION
		if connectionString := os.Getenv("MIGRATE_CONNECTION"); connectionString != "" {
			config.ConnectionString = connectionString
			// Parse connection string to extract individual components
			if parsed := ParseConnectionString(connectionString); parsed != nil {
				config.Driver = parsed.Driver
				config.Host = parsed.Host
				config.Port = parsed.Port
				config.Username = parsed.Username
				config.Password = parsed.Password
				config.Database = parsed.Database
				config.SSLMode = parsed.SSLMode
				config.Charset = parsed.Charset
			}
			return config
		}
	}

	// Fallback to system environment variables
	if connectionString := os.Getenv("MIGRATE_CONNECTION"); connectionString != "" {
		config.ConnectionString = connectionString
		// Parse connection string to extract individual components
		if parsed := ParseConnectionString(connectionString); parsed != nil {
			config.Driver = parsed.Driver
			config.Host = parsed.Host
			config.Port = parsed.Port
			config.Username = parsed.Username
			config.Password = parsed.Password
			config.Database = parsed.Database
			config.SSLMode = parsed.SSLMode
			config.Charset = parsed.Charset
		}
		return config
	}

	// Return default config if MIGRATE_CONNECTION not set
	return config
}

// loadDotEnv loads environment variables from .env file
func loadDotEnv() error {
	// Check if .env file exists
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return fmt.Errorf(".env file not found")
	}

	// Read .env file
	content, err := os.ReadFile(".env")
	if err != nil {
		return fmt.Errorf("failed to read .env file: %v", err)
	}

	// Parse .env file content
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value pairs
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				// Remove quotes if present
				if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
					value = strings.Trim(value, "\"")
				}

				// Set environment variable
				os.Setenv(key, value)
			}
		}
	}

	return nil
}

// Connect establishes database connection using GORM
func Connect(config *DatabaseConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(config.Driver) {
	case constants.DriverSQLite:
		dialector = sqlite.Open(config.Database)
	case constants.DriverMySQL:
		dsn := fmt.Sprintf(constants.MySQLDSNFormat,
			config.Username, config.Password, config.Host, config.Port, config.Database, config.Charset)
		dialector = mysql.Open(dsn)
	case constants.DriverPostgres, constants.DriverPostgreSQL:
		dsn := fmt.Sprintf(constants.PostgresDSNFormat,
			config.Host, config.Username, config.Password, config.Database, config.Port, config.SSLMode, constants.PostgresTimeZone)
		dialector = postgres.Open(dsn)
	default:
		return nil, fmt.Errorf(constants.ErrUnsupportedDriver, config.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable all GORM logging
	})
	if err != nil {
		return nil, fmt.Errorf(constants.ErrFailedConnect, err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrFailedGetDB, err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf(constants.ErrFailedPing, err)
	}

	return db, nil
}

// GetConnectionString returns formatted connection string for display
func (c *DatabaseConfig) GetConnectionString() string {
	if c.ConnectionString != "" {
		return c.ConnectionString
	}

	switch strings.ToLower(c.Driver) {
	case constants.DriverSQLite:
		return fmt.Sprintf(constants.SQLiteConnectionFormat, c.Database)
	case constants.DriverMySQL:
		return fmt.Sprintf(constants.MySQLConnectionFormat,
			c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset)
	case constants.DriverPostgres, constants.DriverPostgreSQL:
		return fmt.Sprintf(constants.PostgresConnectionFormat,
			c.Username, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
	default:
		return "unknown driver"
	}
}

// ParseConnectionString parses a connection string into DatabaseConfig
func ParseConnectionString(connectionString string) *DatabaseConfig {
	config := &DatabaseConfig{}

	// Handle SQLite
	if strings.HasPrefix(connectionString, constants.SQLitePrefix) {
		config.Driver = constants.DriverSQLite
		config.Database = strings.TrimPrefix(connectionString, constants.SQLitePrefix)
		return config
	}

	// Handle MySQL
	if strings.HasPrefix(connectionString, constants.MySQLPrefix) {
		config.Driver = constants.DriverMySQL
		// mysql://username:password@host:port/database?charset=utf8mb4
		connectionString = strings.TrimPrefix(connectionString, constants.MySQLPrefix)

		// Split by @ to separate credentials from host
		parts := strings.Split(connectionString, "@")
		if len(parts) == 2 {
			// Parse credentials
			credParts := strings.Split(parts[0], ":")
			if len(credParts) >= 2 {
				config.Username = credParts[0]
				config.Password = credParts[1]
			}

			// Parse host and database
			hostDB := parts[1]
			// Split by / to separate host:port from database
			dbParts := strings.Split(hostDB, "/")
			if len(dbParts) == 2 {
				config.Database = dbParts[1]
				// Split by ? to remove query parameters
				if strings.Contains(config.Database, "?") {
					config.Database = strings.Split(config.Database, "?")[0]
				}

				// Parse host:port
				hostPort := strings.Split(dbParts[0], ":")
				if len(hostPort) == 2 {
					config.Host = hostPort[0]
					config.Port = hostPort[1]
				} else {
					config.Host = dbParts[0]
					config.Port = constants.DefaultPort
				}
			}
		}

		// Parse charset from query parameters
		if strings.Contains(connectionString, "charset=") {
			queryParts := strings.Split(connectionString, "charset=")
			if len(queryParts) == 2 {
				config.Charset = strings.Split(queryParts[1], "&")[0]
			}
		} else {
			config.Charset = constants.DefaultCharset
		}

		return config
	}

	// Handle PostgreSQL
	if strings.HasPrefix(connectionString, constants.PostgresPrefix) || strings.HasPrefix(connectionString, constants.PostgresAltPrefix) {
		config.Driver = constants.DriverPostgres
		// postgresql://username:password@host:port/database?sslmode=disable
		connectionString = strings.TrimPrefix(connectionString, constants.PostgresPrefix)
		connectionString = strings.TrimPrefix(connectionString, constants.PostgresAltPrefix)

		// Split by @ to separate credentials from host
		parts := strings.Split(connectionString, "@")
		if len(parts) == 2 {
			// Parse credentials
			credParts := strings.Split(parts[0], ":")
			if len(credParts) >= 2 {
				config.Username = credParts[0]
				config.Password = credParts[1]
			}

			// Parse host and database
			hostDB := parts[1]
			// Split by / to separate host:port from database
			dbParts := strings.Split(hostDB, "/")
			if len(dbParts) == 2 {
				config.Database = dbParts[1]
				// Split by ? to remove query parameters
				if strings.Contains(config.Database, "?") {
					config.Database = strings.Split(config.Database, "?")[0]
				}

				// Parse host:port
				hostPort := strings.Split(dbParts[0], ":")
				if len(hostPort) == 2 {
					config.Host = hostPort[0]
					config.Port = hostPort[1]
				} else {
					config.Host = dbParts[0]
					config.Port = constants.DefaultPortPostgres
				}
			}
		}

		// Parse sslmode from query parameters
		if strings.Contains(connectionString, "sslmode=") {
			queryParts := strings.Split(connectionString, "sslmode=")
			if len(queryParts) == 2 {
				config.SSLMode = strings.Split(queryParts[1], "&")[0]
			}
		} else {
			config.SSLMode = constants.DefaultSSLMode
		}

		return config
	}

	return nil
}
