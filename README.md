# MySQL Adapter for Toutā DataMapper

[![CI](https://github.com/toutaio/toutago-datamapper-mysql/actions/workflows/ci.yml/badge.svg)](https://github.com/toutaio/toutago-datamapper-mysql/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/toutaio/toutago-datamapper-mysql.svg)](https://pkg.go.dev/github.com/toutaio/toutago-datamapper-mysql)
[![Go Report Card](https://goreportcard.com/badge/github.com/toutaio/toutago-datamapper-mysql)](https://goreportcard.com/report/github.com/toutaio/toutago-datamapper-mysql)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> A MySQL database adapter implementation for [toutago-datamapper](https://github.com/toutaio/toutago-datamapper), providing full CRUD operations, bulk inserts, transactions, and custom query execution. Part of the **Toutā Framework** ecosystem.

## Features

- ✅ Full CRUD operations (Create, Read, Update, Delete)
- ✅ Bulk insert support for efficient batch operations
- ✅ Named parameter substitution (`{param_name}`)
- ✅ Auto-generated ID handling (auto-increment)
- ✅ Optimistic locking support
- ✅ Connection pooling configuration
- ✅ Custom SQL execution and stored procedures
- ✅ CQRS pattern support via source configuration

## Installation

```bash
go get github.com/toutaio/toutago-datamapper-mysql
```

## Quick Start

### 1. Define Your Configuration

Create a `sources.yaml` file with MySQL connection details:

```yaml
sources:
  - name: users_db
    type: mysql
    config:
      host: localhost
      port: 3306
      user: myapp
      password: ${MYSQL_PASSWORD}  # From environment
      database: myapp_db
      ssl: "false"
      max_connections: 20
      max_idle: 5
      conn_max_age_seconds: 3600
```

### 2. Configure Mappings

Create a `mappings/users.yaml` file:

```yaml
object: User
source: users_db

mappings:
  - name: fetch_by_id
    type: fetch
    statement: "SELECT id, name, email, created_at FROM users WHERE id = {id}"
    properties:
      - object: ID
        data: id
      - object: Name
        data: name
      - object: Email
        data: email
      - object: CreatedAt
        data: created_at

  - name: insert
    type: insert
    statement: users
    properties:
      - object: Name
        data: name
      - object: Email
        data: email
    generated:
      - object: ID
        data: id
```

### 3. Use in Your Application

```go
package main

import (
    "context"
    "log"
    
    "github.com/toutaio/toutago-datamapper/engine"
    mysql "github.com/toutaio/toutago-datamapper-mysql"
)

func main() {
    // Create engine
    mapper := engine.NewMapper()
    
    // Register MySQL adapter
    mapper.RegisterAdapter("mysql", mysql.NewMySQLAdapter())
    
    // Load configuration
    if err := mapper.LoadConfig("sources.yaml", "mappings/"); err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Fetch a user
    user := &User{}
    err := mapper.Fetch(ctx, "User", "fetch_by_id", 
        map[string]interface{}{"id": 123}, user)
    if err != nil {
        log.Fatal(err)
    }
    
    // Insert a new user
    newUser := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    err = mapper.Insert(ctx, "User", "insert", newUser)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Created user with ID: %d", newUser.ID)
}

type User struct {
    ID        int64
    Name      string
    Email     string
    CreatedAt string
}
```

## Configuration Options

### Source Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `host` | string | MySQL server hostname | `localhost` |
| `port` | int | MySQL server port | `3306` |
| `user` | string | Database user | `root` |
| `password` | string | Database password | `""` |
| `database` | string | Database name | Required |
| `ssl` | string | SSL mode (`true`, `false`, `skip-verify`) | `false` |
| `max_connections` | int | Maximum open connections | `10` |
| `max_idle` | int | Maximum idle connections | `5` |
| `conn_max_age_seconds` | int | Connection max age in seconds | `3600` |

### Parameter Substitution

The adapter supports named parameter placeholders in queries:

```yaml
statement: "SELECT * FROM users WHERE name = {name} AND status = {status}"
```

Parameters are replaced with MySQL's `?` placeholders at runtime.

## Advanced Features

### Bulk Insert

```yaml
- name: bulk_insert
  type: insert
  statement: users
  bulk: true
  properties:
    - object: Name
      data: name
    - object: Email
      data: email
```

```go
users := []interface{}{
    &User{Name: "Alice", Email: "alice@example.com"},
    &User{Name: "Bob", Email: "bob@example.com"},
}

err := mapper.InsertBulk(ctx, "User", "bulk_insert", users)
```

### Custom Actions (Stored Procedures)

```yaml
actions:
  - name: archive_old_users
    statement: "CALL archive_users_older_than({days})"
    parameters:
      - object: Days
        data: days
    result:
      multi: true
      properties:
        - object: ArchivedCount
          data: archived_count
```

### Optimistic Locking

```yaml
- name: update
  type: update
  statement: users
  identifier:
    - object: ID
      data: id
  condition:
    - object: Version
      data: version
  properties:
    - object: Name
      data: name
    - object: Version
      data: version
```

## Error Handling

The adapter returns standard errors from `github.com/toutaio/toutago-datamapper/adapter`:

- `adapter.ErrNotFound` - Record not found
- `adapter.ErrConnection` - Connection failure
- `adapter.ErrValidation` - Constraint violation
- `adapter.ErrConflict` - Optimistic locking conflict

## Testing

```bash
# Unit tests
go test -v

# Integration tests (requires MySQL)
docker-compose up -d mysql
go test -v -tags=integration
docker-compose down
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Related Projects

- [toutago-datamapper](https://github.com/toutaio/toutago-datamapper) - Core data mapper library
- [toutago-datamapper-postgres](https://github.com/toutaio/toutago-datamapper-postgres) - PostgreSQL adapter

---

Part of the **Toutā Framework** - A modular Go framework emphasizing interface-first design and zero framework lock-in.

