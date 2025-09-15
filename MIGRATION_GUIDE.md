# Elsa Database Migration Management Guide

**Elsa Migration** is a powerful database migration management system that provides structured, version-controlled database schema and data management for Go applications.

> ⚠️ **IMPORTANT DISCLAIMER**: Database migrations involve direct changes to your database structure and data. While we strive to make Elsa Migration as safe and reliable as possible, **we cannot be held responsible for any data loss, corruption, or damage** that may occur during migration processes. Always backup your database before running migrations, test migrations in development environments first, and use at your own risk. We remain committed to continuously improving the application for better safety and reliability.

## 🎯 Overview

Elsa Migration separates database changes into two distinct categories:

- **DDL (Data Definition Language)**: Schema changes, table creation, modifications, indexes
- **DML (Data Manipulation Language)**: Data seeding, updates, transformations, and data migrations

This separation provides better organization, cleaner production deployments, and more granular control over database changes.

## 🏗️ Why DDL and DML Separation?

### Production Benefits

1. **Deployment Safety**: DDL changes (schema) can be deployed separately from DML changes (data)
2. **Rollback Control**: Schema changes and data changes can be rolled back independently
3. **Environment Consistency**: Schema changes are applied consistently across all environments, while data changes can be environment-specific
4. **Performance**: DDL operations are typically faster and can be applied during maintenance windows
5. **Team Collaboration**: Different team members can work on schema vs. data changes without conflicts

### Development Benefits

1. **Clear Organization**: Easy to distinguish between structural and data changes
2. **Version Control**: Better tracking of what changed and when
3. **Testing**: Schema and data can be tested independently
4. **Documentation**: Clear separation makes it easier to document changes

## 📁 Project Structure

```
your-project/
├── database/
│   └── migration/
│       ├── ddl/                    # DDL migrations (schema changes)
│       │   ├── 20240101120000_create_users.up.sql
│       │   ├── 20240101120000_create_users.down.sql
│       │   ├── 20240101120001_create_roles.up.sql
│       │   └── 20240101120001_create_roles.down.sql
│       └── dml/                    # DML migrations (data changes)
│           ├── 20240101120002_seed_users.up.sql
│           ├── 20240101120002_seed_users.down.sql
│           ├── 20240101120003_seed_roles.up.sql
│           └── 20240101120003_seed_roles.down.sql
├── .env                           # Database configuration
└── your-go-files...
```

## 🚀 Quick Start

### 1. Connect to Database

Elsa Migration supports multiple ways to configure database connection:

#### Interactive Setup (Recommended for first time)
```bash
# Interactive connection setup - creates/updates .env file
elsa migration connect
```

#### Direct Connection String
```bash
# Use -c flag for direct connection
elsa migration connect -c "mysql://user:password@localhost:3306/database"
elsa migration connect -c "postgres://user:password@localhost:5432/database"
elsa migration connect -c "sqlite://database.db"
```

#### Environment Variable
```bash
# Set MIGRATE_CONNECTION environment variable
export MIGRATE_CONNECTION="mysql://user:password@localhost:3306/database"
elsa migration up ddl
```

#### .env File Configuration
```bash
# Add to .env file
echo "MIGRATE_CONNECTION=mysql://user:password@localhost:3306/database" >> .env
elsa migration up ddl
```

**Connection Priority Order:**
1. `-c` flag (highest priority)
2. `MIGRATE_CONNECTION` environment variable
3. `.env` file with `MIGRATE_CONNECTION` key
4. Individual database environment variables (`DB_DRIVER`, `DB_HOST`, etc.)

### 2. Create Migrations

```bash
# Create DDL migration (schema changes)
elsa migration create ddl create_users_table

# Create DML migration (data changes)
elsa migration create dml seed_users_data

# Create with sequential numbering
elsa migration create ddl create_products_table --sequential

# Create in custom path
elsa migration create ddl create_orders_table --path custom/migrations
```

### 3. Apply Migrations

```bash
# Apply all DDL migrations
elsa migration up ddl

# Apply all DML migrations
elsa migration up dml

# Apply specific number of migrations
elsa migration up ddl --step 2

# Apply up to specific migration
elsa migration up ddl --to 00002
```

### 4. Check Status

```bash
# Show all migration status
elsa migration status

# Show only DDL migrations
elsa migration status --ddl

# Show only DML migrations
elsa migration status --dml
```

## 📋 Migration Commands Reference

### Connection Commands

| Command | Description |
|---------|-------------|
| `elsa migration connect` | Interactive setup - creates/updates .env file |
| `elsa migration connect -c <string>` | Connect using direct connection string |

### Creation Commands

| Command | Description |
|---------|-------------|
| `elsa migration create ddl <name>` | Create DDL migration (timestamp format) |
| `elsa migration create dml <name>` | Create DML migration (timestamp format) |
| `elsa migration create ddl <name> --sequential` | Create with sequential numbering |
| `elsa migration create ddl <name> --path <path>` | Create in custom directory |

### Execution Commands

| Command | Description |
|---------|-------------|
| `elsa migration up ddl` | Apply all DDL migrations |
| `elsa migration up dml` | Apply all DML migrations |
| `elsa migration up ddl --step <n>` | Apply n DDL migrations |
| `elsa migration up ddl --to <id>` | Apply up to specific migration ID |
| `elsa migration down ddl` | Rollback last DDL migration |
| `elsa migration down dml` | Rollback last DML migration |
| `elsa migration refresh ddl` | Refresh all DDL migrations (rollback + apply) |
| `elsa migration refresh dml` | Refresh all DML migrations (rollback + apply) |

### Status Commands

| Command | Description |
|---------|-------------|
| `elsa migration status` | Show all migration status |
| `elsa migration status --ddl` | Show only DDL migrations |
| `elsa migration status --dml` | Show only DML migrations |
| `elsa migration info` | Show detailed migration information |

## 🔧 Migration File Formats

### Timestamp Format (Default)

Elsa uses timestamp format with milliseconds for unique migration IDs:

```
20240101120000123_create_users_table.up.sql
20240101120000123_create_users_table.down.sql
```

**Format**: `YYYYMMDDHHMMSSmmm_<name>.up.sql`

### Sequential Format

For projects preferring sequential numbering:

```
00001_create_users_table.up.sql
00001_create_users_table.down.sql
```

**Format**: `%05d_<name>.up.sql`

### File Naming Rules

1. **Consistency**: All migrations in a folder must use the same format
2. **Uniqueness**: Migration IDs must be unique within the same type (DDL/DML)
3. **Descriptive**: Use clear, descriptive names for migration purposes
4. **Underscores**: Use underscores to separate words in migration names

## 📝 Migration File Templates

### DDL Migration Template

**Up Migration** (`create_users_table.up.sql`):
```sql
-- Migration: create_users_table
-- Type: DDL
-- Description: Create users table with basic fields

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Add indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
```

**Down Migration** (`create_users_table.down.sql`):
```sql
-- Migration: create_users_table (Rollback)
-- Type: DDL
-- Description: Create users table with basic fields

-- Drop indexes first
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_email;

-- Drop table
DROP TABLE IF EXISTS users;
```

### DML Migration Template

**Up Migration** (`seed_users_data.up.sql`):
```sql
-- Migration: seed_users_data
-- Type: DML
-- Description: Seed initial user data

INSERT INTO users (name, email, created_at) VALUES 
    ('John Doe', 'john@example.com', NOW()),
    ('Jane Smith', 'jane@example.com', NOW()),
    ('Admin User', 'admin@example.com', NOW());
```

**Down Migration** (`seed_users_data.down.sql`):
```sql
-- Migration: seed_users_data (Rollback)
-- Type: DML
-- Description: Seed initial user data

DELETE FROM users WHERE email IN (
    'john@example.com',
    'jane@example.com',
    'admin@example.com'
);
```

## 🗄️ Database Support

Elsa Migration supports multiple database systems:

### MySQL
```bash
elsa migration connect -c "mysql://user:password@localhost:3306/database"
```

### PostgreSQL
```bash
elsa migration connect -c "postgres://user:password@localhost:5432/database"
```

### SQLite
```bash
elsa migration connect -c "sqlite://database.db"
```

## ⚙️ Configuration

### Environment Variables

#### Option 1: Single Connection String (Recommended)
```env
# Single connection string (highest priority)
MIGRATE_CONNECTION=mysql://user:password@localhost:3306/database
```

#### Option 2: Individual Database Variables (Input Only)
When using `elsa migration connect` interactively, you can input individual database details:
- Database Driver (mysql, postgres, sqlite)
- Host
- Port
- Username
- Password
- Database Name

**Note**: Individual variables are only used for input. The system automatically converts them to `MIGRATE_CONNECTION` format and saves it to `.env` file.

#### Connection String vs Individual Variables
- **`MIGRATE_CONNECTION`**: Single string, easier to manage, higher priority (saved to .env)
- **Individual variables**: Input method only - converted to MIGRATE_CONNECTION format

### Connection String Format

```
<driver>://<username>:<password>@<host>:<port>/<database>?<parameters>
```

Examples:
- `mysql://user:pass@localhost:3306/mydb`
- `postgres://user:pass@localhost:5432/mydb?sslmode=disable`
- `sqlite://./database.db`

### Automatic .env File Management

When you run `elsa migration connect` interactively:

1. **If .env doesn't exist**: Creates new .env file with `MIGRATE_CONNECTION`
2. **If .env exists**: Appends `MIGRATE_CONNECTION` to existing file
3. **If MIGRATE_CONNECTION exists**: Updates the existing value

**Important**: Even when you input individual database details (driver, host, port, username, password, database), the system automatically converts them to a single `MIGRATE_CONNECTION` string and saves it to the `.env` file.

Example conversion:
```
Input: Driver=mysql, Host=localhost, Port=3306, Username=user, Password=pass, Database=mydb
Saved to .env: MIGRATE_CONNECTION=mysql://user:pass@localhost:3306/mydb
```

This makes it easy to set up database connections for your team members!

## 🔄 Migration Lifecycle

### 1. Development Phase

```bash
# Create new DDL migration
elsa migration create ddl add_user_phone_column

# Edit the generated files
# database/migration/ddl/20240101120000_add_user_phone_column.up.sql
# database/migration/ddl/20240101120000_add_user_phone_column.down.sql

# Test the migration
elsa migration up ddl
elsa migration status
```

### 2. Testing Phase

```bash
# Apply all migrations
elsa migration up ddl
elsa migration up dml

# Verify status
elsa migration status

# Test rollback
elsa migration down ddl
elsa migration down dml
```

### 3. Production Deployment

```bash
# Apply DDL migrations first (schema changes)
elsa migration up ddl

# Verify schema changes
elsa migration status --ddl

# Apply DML migrations (data changes)
elsa migration up dml

# Verify all migrations
elsa migration status
```

## 🛠️ Best Practices

### DDL Migration Best Practices

1. **Atomic Changes**: Each migration should make one logical change
2. **Backward Compatibility**: Ensure migrations can be rolled back safely
3. **Index Management**: Create indexes in separate migrations for better performance
4. **Column Modifications**: Use separate migrations for adding, modifying, and dropping columns
5. **Foreign Keys**: Add foreign key constraints after creating all referenced tables

### DML Migration Best Practices

1. **Idempotent**: Ensure migrations can be run multiple times safely
2. **Data Validation**: Include data validation and cleanup
3. **Batch Operations**: Use batch operations for large data sets
4. **Environment Specific**: Use environment-specific data when needed
5. **Rollback Safety**: Ensure rollback operations don't affect unrelated data

### General Best Practices

1. **Naming Convention**: Use descriptive, consistent naming
2. **Version Control**: Always commit migration files to version control
3. **Testing**: Test migrations in development before production
4. **Backup**: Always backup database before applying migrations
5. **Documentation**: Document complex migrations with comments
6. **Risk Management**: Understand that migrations can cause data loss - use with caution
7. **Rollback Planning**: Always have a rollback plan before applying migrations
8. **Staging Environment**: Test migrations in staging environment that mirrors production

## 🚨 Troubleshooting

### Common Issues

#### Migration Already Applied
```bash
# Check status to see applied migrations
elsa migration status

# If migration is partially applied, you may need to manually fix the database state
```

#### Connection Issues
```bash
# Test connection
elsa migration connect -c "your_connection_string"

# Check .env file configuration
cat .env
```

#### File Format Conflicts
```bash
# Error: folder contains mixed migration formats
# Solution: Use consistent format for all migrations in the same folder
elsa migration create ddl new_migration --sequential  # or --timestamp
```

#### Migration Execution Errors
```bash
# Check migration file syntax
# Ensure SQL is valid for your database type
# Check for missing semicolons or syntax errors
```

### Recovery Procedures

> ⚠️ **WARNING**: Recovery procedures can be risky and may cause data loss. Always backup your database before attempting any recovery operations.

#### Rollback Failed Migration
```bash
# Rollback last migration
elsa migration down ddl

# Check status
elsa migration status
```

#### Refresh All Migrations
```bash
# Rollback and reapply all migrations
elsa migration refresh ddl
elsa migration refresh dml
```

#### Manual Database Fix
If migrations are in an inconsistent state:

1. **BACKUP FIRST**: Create a full database backup
2. Check migration table: `SELECT * FROM elsa_migrations;`
3. Manually remove problematic records
4. Fix database schema manually
5. Update migration status accordingly

**Note**: Manual fixes should only be performed by experienced database administrators who understand the risks involved.

## 📊 Migration Tracking

Elsa automatically tracks migrations in the `elsa_migrations` table:

| Column | Type | Description |
|--------|------|-------------|
| `id` | Primary Key | Auto-incrementing ID |
| `migration_id` | String | Migration identifier (timestamp or sequential) |
| `name` | String | Migration name |
| `type` | String | Migration type (ddl/dml) |
| `applied_at` | Timestamp | When migration was applied |
| `checksum` | String | Content checksum for integrity |
| `execution_time` | Integer | Execution time in milliseconds |

## 📚 Advanced Usage

### Custom Migration Paths

```bash
# Use custom migration directory
elsa migration create ddl new_table --path custom/migrations
elsa migration up ddl --path custom/migrations
elsa migration status --path custom/migrations
```

### Migration Information

```bash
# Get detailed migration information
elsa migration info

# Get specific migration type info
elsa migration info --ddl
elsa migration info --dml
```

### Batch Operations

```bash
# Apply multiple migrations with specific steps
elsa migration up ddl --step 5

# Apply up to specific migration
elsa migration up ddl --to 20240101120000
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Elsa Migration** - Making database management simple, organized, and production-ready! 🚀
