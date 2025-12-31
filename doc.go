// Package mysql provides a MySQL database adapter for toutago-datamapper.
//
// This adapter implements the datamapper Adapter interface for MySQL databases,
// enabling full CRUD operations, bulk inserts, transactions, and custom query execution.
//
// # Features
//
//   - Full CRUD operations (Create, Read, Update, Delete)
//   - Bulk insert support for efficient batch operations
//   - Named parameter substitution ({param_name})
//   - Auto-generated ID handling (auto-increment)
//   - Optimistic locking support
//   - Connection pooling configuration
//   - Custom SQL execution and stored procedures
//   - CQRS pattern support via source configuration
//
// # Quick Start
//
// Register the MySQL adapter with your datamapper:
//
//	import (
//	    "github.com/toutaio/toutago-datamapper/engine"
//	    "github.com/toutaio/toutago-datamapper-mysql"
//	)
//
//	mapper, err := engine.NewMapper("config.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer mapper.Close()
//
//	// Register MySQL adapter
//	mapper.RegisterAdapter("mysql", func(source config.Source) (adapter.Adapter, error) {
//	    return mysql.NewMySQLAdapter(source.Connection)
//	})
//
// # Configuration
//
// Define MySQL connection in YAML:
//
//	sources:
//	  - name: users_db
//	    type: mysql
//	    config:
//	      host: localhost
//	      port: 3306
//	      user: myapp
//	      password: ${MYSQL_PASSWORD}
//	      database: myapp_db
//	      max_connections: 20
//	      max_idle: 5
//
// # Usage
//
// Use through datamapper API:
//
//	user := &User{ID: "1"}
//	err := mapper.Fetch(context.Background(), "User", user)
//
//	user.Name = "Updated Name"
//	err = mapper.Update(context.Background(), "User", user)
//
// # Connection Pooling
//
// The adapter supports connection pooling with configurable parameters:
//   - max_connections: Maximum open connections
//   - max_idle: Maximum idle connections
//   - conn_max_age_seconds: Connection lifetime
//
// # Thread Safety
//
// The adapter uses database/sql which provides connection pooling and
// is safe for concurrent use across multiple goroutines.
//
// # Version
//
// This is version 0.1.0 - requires toutago-datamapper v0.1.0 or higher.
// Requires Go 1.22 or higher.
package mysql
