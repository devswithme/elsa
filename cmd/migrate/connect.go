package migrate

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.risoftinc.com/elsa/constants"
	"go.risoftinc.com/elsa/internal/database"
)

var (
	connectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to database",
		Long: `Connect to database with interactive configuration, flags, or auto-load from .env file.
		
Examples:
  elsa migration connect                                    # Interactive database connection
  elsa migration connect --connection "sqlite://elsa.db"    # Connect using connection string flag
  elsa migration connect -c "mysql://root:pass@localhost:3306/myapp"`,
		RunE: runConnect,
	}

	connectionFlag string
)

func init() {
	connectCmd.Flags().StringVarP(&connectionFlag, "connection", "c", "", "Database connection string")
}

func runConnect(cmd *cobra.Command, args []string) error {
	var config *database.DatabaseConfig

	if connectionFlag != "" {
		// Use connection string from flag
		config = database.DefaultConfig()
		config.ConnectionString = connectionFlag

		// Parse connection string to extract individual components
		if parsed := database.ParseConnectionString(connectionFlag); parsed != nil {
			config.Driver = parsed.Driver
			config.Host = parsed.Host
			config.Port = parsed.Port
			config.Username = parsed.Username
			config.Password = parsed.Password
			config.Database = parsed.Database
			config.SSLMode = parsed.SSLMode
			config.Charset = parsed.Charset
		}

		fmt.Printf(constants.InfoUsingConnectionFlag)
		fmt.Printf(constants.ConnectionInfoFormat, config.ConnectionString)
		fmt.Printf(constants.DriverInfoFormat, config.Driver)
		fmt.Printf(constants.HostInfoFormat, config.Host)
		fmt.Printf(constants.PortInfoFormat, config.Port)
		fmt.Printf(constants.DatabaseInfoFormat, config.Database)
		fmt.Printf(constants.UsernameInfoFormat, config.Username)
		if config.Password != "" {
			fmt.Printf(constants.PasswordInfoFormat, strings.Repeat("*", len(config.Password)))
		}
	} else {
		// Try to load from .env file first, then interactive if not found
		config = database.LoadFromEnv()

		// Check if we have a valid connection string from .env
		if config.ConnectionString != "" {
			fmt.Printf(constants.InfoLoadedFromEnv)
			fmt.Printf(constants.ConnectionInfoFormat, config.ConnectionString)
			fmt.Printf(constants.DriverInfoFormat, config.Driver)
			fmt.Printf(constants.HostInfoFormat, config.Host)
			fmt.Printf(constants.PortInfoFormat, config.Port)
			fmt.Printf(constants.DatabaseInfoFormat, config.Database)
			fmt.Printf(constants.UsernameInfoFormat, config.Username)
			if config.Password != "" {
				fmt.Printf(constants.PasswordInfoFormat, strings.Repeat("*", len(config.Password)))
			}
		} else {
			// No .env file or invalid config, use interactive
			fmt.Printf(constants.InfoNoEnvFile)
			config = getInteractiveConfig()
		}
	}

	// Test connection
	fmt.Printf("\n%s", constants.InfoTestingConnection)
	db, err := database.Connect(config)
	if err != nil {
		return fmt.Errorf("❌ Connection failed: %v", err, err)
	}

	fmt.Printf(constants.SuccessConnected)
	fmt.Printf(constants.ConnectionInfoFormat, config.GetConnectionString())

	// Test migration table
	fmt.Printf("\n%s", constants.InfoEnsuringTable)
	executor := database.NewMigrationExecutor(db)
	if err := executor.EnsureMigrationTable(); err != nil {
		return fmt.Errorf("❌ Failed to create migration table: %v", err)
	}

	fmt.Printf("✅ Migration table ready!\n")

	// Save configuration to .env file if not using flag (to avoid overwriting existing config)
	if connectionFlag == "" {
		if err := saveConfigToEnv(config); err != nil {
			fmt.Printf("⚠️  Warning: Could not save configuration to .env file: %v\n", err)
		} else {
			fmt.Printf("💾 Configuration saved to .env file\n")
		}
	}

	return nil
}

func getInteractiveConfig() *database.DatabaseConfig {
	config := database.DefaultConfig()
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("🔧 Database Connection Configuration\n")
	fmt.Printf("=====================================\n\n")

	fmt.Printf("Choose configuration method:\n")
	fmt.Printf("1. Single connection string (recommended)\n")
	fmt.Printf("2. Individual parameters\n")
	fmt.Printf("Select option (1-2): ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "1" {
		return getSingleConnectionString(reader, config)
	}

	// Driver selection for individual parameters
	fmt.Printf("\nAvailable drivers:\n")
	fmt.Printf("1. SQLite (file-based)\n")
	fmt.Printf("2. MySQL\n")
	fmt.Printf("3. PostgreSQL\n")
	fmt.Printf("Current: %s\n", config.Driver)

	fmt.Printf("Select driver (1-3) or press Enter for current: ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		config.Driver = "sqlite"
	case "2":
		config.Driver = "mysql"
	case "3":
		config.Driver = "postgres"
	}

	// Database-specific configuration
	switch config.Driver {
	case "sqlite":
		config = getSQLiteConfig(reader, config)
	case "mysql":
		config = getMySQLConfig(reader, config)
	case "postgres":
		config = getPostgreSQLConfig(reader, config)
	}

	// Generate connection string from individual parameters
	config.ConnectionString = config.GetConnectionString()

	return config
}

func getSingleConnectionString(reader *bufio.Reader, config *database.DatabaseConfig) *database.DatabaseConfig {
	fmt.Printf("\n🔗 Single Connection String Configuration\n")
	fmt.Printf("==========================================\n\n")

	fmt.Printf("Examples:\n")
	fmt.Printf("  SQLite:   sqlite://elsa.db\n")
	fmt.Printf("  MySQL:    mysql://root:password@localhost:3306/myapp?charset=utf8mb4\n")
	fmt.Printf("  PostgreSQL: postgresql://user:pass@localhost:5432/myapp?sslmode=disable\n\n")

	fmt.Printf("Enter connection string: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "" {
		config.ConnectionString = input
		// Parse the connection string to extract individual components
		if parsed := database.ParseConnectionString(input); parsed != nil {
			config.Driver = parsed.Driver
			config.Host = parsed.Host
			config.Port = parsed.Port
			config.Username = parsed.Username
			config.Password = parsed.Password
			config.Database = parsed.Database
			config.SSLMode = parsed.SSLMode
			config.Charset = parsed.Charset
		}
	}

	return config
}

func getSQLiteConfig(reader *bufio.Reader, config *database.DatabaseConfig) *database.DatabaseConfig {
	fmt.Printf("\n📁 SQLite Configuration\n")
	fmt.Printf("Database file path (default: %s): ", config.Database)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		config.Database = input
	}
	return config
}

func getMySQLConfig(reader *bufio.Reader, config *database.DatabaseConfig) *database.DatabaseConfig {
	fmt.Printf("\n🐬 MySQL Configuration\n")

	fmt.Printf("Host (default: %s): ", config.Host)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		config.Host = input
	}

	fmt.Printf("Port (default: %s): ", config.Port)
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		config.Port = input
	}

	fmt.Printf("Username (default: %s): ", config.Username)
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		config.Username = input
	}

	fmt.Printf("Password: ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		config.Password = input
	}

	fmt.Printf("Database name: ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		config.Database = input
	}

	fmt.Printf("Charset (default: %s): ", config.Charset)
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		config.Charset = input
	}

	return config
}

func getPostgreSQLConfig(reader *bufio.Reader, config *database.DatabaseConfig) *database.DatabaseConfig {
	fmt.Printf("\n🐘 PostgreSQL Configuration\n")

	fmt.Printf("Host (default: %s): ", config.Host)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		config.Host = input
	}

	fmt.Printf("Port (default: %s): ", config.Port)
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		config.Port = input
	}

	fmt.Printf("Username (default: %s): ", config.Username)
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		config.Username = input
	}

	fmt.Printf("Password: ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		config.Password = input
	}

	fmt.Printf("Database name: ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		config.Database = input
	}

	fmt.Printf("SSL Mode (default: %s): ", config.SSLMode)
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		config.SSLMode = input
	}

	return config
}

func saveConfigToEnv(config *database.DatabaseConfig) error {
	// Read existing .env file if it exists
	existingContent := ""
	if _, err := os.Stat(".env"); err == nil {
		content, err := os.ReadFile(".env")
		if err == nil {
			existingContent = string(content)
		}
	}

	// Parse existing content to check what keys are already present
	existingKeys := make(map[string]bool)
	if existingContent != "" {
		lines := strings.Split(existingContent, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					existingKeys[key] = true
				}
			}
		}
	}

	// Prepare new content
	var newContent strings.Builder

	// Write existing content first
	if existingContent != "" {
		newContent.WriteString(existingContent)
		// Add newline if content doesn't end with one
		if !strings.HasSuffix(existingContent, "\n") {
			newContent.WriteString("\n")
		}
		newContent.WriteString("\n")
	} else {
		// New .env file
		newContent.WriteString("# Database Configuration\n")
	}

	// Add MIGRATE_CONNECTION if not exists
	if !existingKeys["MIGRATE_CONNECTION"] {
		newContent.WriteString(fmt.Sprintf("MIGRATE_CONNECTION=%s\n", config.ConnectionString))
	}

	return os.WriteFile(".env", []byte(newContent.String()), 0644)
}

// GetDatabaseConnection is a helper function to get database connection
// that can be used by other migration commands
func GetDatabaseConnection() (*database.DatabaseConfig, error) {
	// Try to load from .env file first
	config := database.LoadFromEnv()

	// Check if we have a valid connection string from .env
	if config.ConnectionString != "" {
		return config, nil
	}

	// No .env file or invalid config
	return nil, fmt.Errorf("no database connection configured. Please run 'elsa migration connect' first")
}
