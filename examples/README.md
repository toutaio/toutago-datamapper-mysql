# MySQL Adapter Examples

This directory contains examples demonstrating various features of the MySQL adapter.

## Prerequisites

1. **Docker** (recommended) or MySQL server installed
2. **Go 1.21.5+**

## Setup

### Option 1: Using Docker (Recommended)

```bash
# Start MySQL server
docker-compose up -d

# Wait for MySQL to be ready
docker-compose logs -f mysql
# Press Ctrl+C when you see "ready for connections"

# Set environment variable
export MYSQL_PASSWORD=rootpassword
```

### Option 2: Using Local MySQL

```bash
# Create database
mysql -u root -p -e "CREATE DATABASE example_db;"

# Run init script
mysql -u root -p example_db < ../init.sql

# Set environment variable
export MYSQL_PASSWORD=your_mysql_password
```

## Running Examples

### Basic CRUD Operations

```bash
cd basic
go run main.go
```

Demonstrates:
- Inserting a user
- Fetching by ID
- Updating a user
- Deleting a user

### Bulk Operations

```bash
cd bulk
go run main.go
```

Demonstrates:
- Bulk insert of multiple users
- Fetching all users

### Advanced Features

```bash
cd advanced
go run main.go
```

Demonstrates:
- Optimistic locking
- Custom actions (stored procedures)
- Complex queries
- Error handling

## Cleanup

```bash
# Stop and remove containers
docker-compose down

# Remove volumes (optional)
docker-compose down -v
```

## Configuration Files

Each example contains:
- `main.go` - Example application code
- `sources.yaml` - Database connection configuration
- `mappings/*.yaml` - Object-to-table mapping definitions

## Environment Variables

- `MYSQL_PASSWORD` - MySQL root password (required)
- `MYSQL_HOST` - MySQL host (default: localhost)
- `MYSQL_PORT` - MySQL port (default: 3306)
- `MYSQL_DATABASE` - Database name (default: example_db)
