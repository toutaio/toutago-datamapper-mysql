// Package mysql provides a MySQL adapter implementation for toutago-datamapper.
// This adapter enables mapping domain objects to MySQL database tables with full CRUD support.
package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/toutaio/toutago-datamapper/adapter"
)

// MySQLAdapter implements the adapter.Adapter interface for MySQL databases.
type MySQLAdapter struct {
	db         *sql.DB
	dsn        string
	maxConn    int
	maxIdle    int
	connMaxAge int
}

// Config keys for MySQL adapter configuration
const (
	ConfigHost     = "host"
	ConfigPort     = "port"
	ConfigUser     = "user"
	ConfigPassword = "password"
	ConfigDatabase = "database"
	ConfigSSL      = "ssl"
	ConfigMaxConn  = "max_connections"
	ConfigMaxIdle  = "max_idle"
	ConfigConnAge  = "conn_max_age_seconds"
)

// NewMySQLAdapter creates a new MySQL adapter instance.
func NewMySQLAdapter() *MySQLAdapter {
	return &MySQLAdapter{
		maxConn:    10,
		maxIdle:    5,
		connMaxAge: 3600,
	}
}

// Name returns the adapter type identifier.
func (a *MySQLAdapter) Name() string {
	return "mysql"
}

// Connect establishes connection to MySQL database.
func (a *MySQLAdapter) Connect(ctx context.Context, config map[string]interface{}) error {
	// Extract connection parameters
	host := getStringConfig(config, ConfigHost, "localhost")
	port := getIntConfig(config, ConfigPort, 3306)
	user := getStringConfig(config, ConfigUser, "root")
	password := getStringConfig(config, ConfigPassword, "")
	database := getStringConfig(config, ConfigDatabase, "")
	ssl := getStringConfig(config, ConfigSSL, "false")

	// Optional connection pooling parameters
	if maxConn, ok := config[ConfigMaxConn].(int); ok {
		a.maxConn = maxConn
	}
	if maxIdle, ok := config[ConfigMaxIdle].(int); ok {
		a.maxIdle = maxIdle
	}
	if connAge, ok := config[ConfigConnAge].(int); ok {
		a.connMaxAge = connAge
	}

	// Build DSN
	a.dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&tls=%s",
		user, password, host, port, database, ssl)

	// Open database connection
	db, err := sql.Open("mysql", a.dsn)
	if err != nil {
		return fmt.Errorf("mysql: failed to open connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(a.maxConn)
	db.SetMaxIdleConns(a.maxIdle)

	// Verify connection
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("mysql: failed to ping database: %w", err)
	}

	a.db = db
	return nil
}

// Close releases database connections.
func (a *MySQLAdapter) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

// Fetch retrieves one or more records from MySQL.
func (a *MySQLAdapter) Fetch(ctx context.Context, op *adapter.Operation, params map[string]interface{}) ([]interface{}, error) {
	if a.db == nil {
		return nil, fmt.Errorf("mysql: adapter not connected")
	}

	// Replace placeholders in query with positional parameters
	query, args := a.buildQuery(op.Statement, params)

	// Prepare statement
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("mysql: failed to prepare query: %w", err)
	}
	defer stmt.Close()

	// Execute query
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("mysql: query failed: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("mysql: failed to get columns: %w", err)
	}

	// Scan results
	var results []interface{}
	for rows.Next() {
		// Create slice of interface{} for scanning
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan row
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("mysql: failed to scan row: %w", err)
		}

		// Build result map
		result := make(map[string]interface{})
		for i, col := range columns {
			result[col] = values[i]
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("mysql: row iteration error: %w", err)
	}

	// Check if we found anything
	if len(results) == 0 && !op.Multi {
		return nil, adapter.ErrNotFound
	}

	return results, nil
}

// Insert creates new records in MySQL.
func (a *MySQLAdapter) Insert(ctx context.Context, op *adapter.Operation, objects []interface{}) error {
	if a.db == nil {
		return fmt.Errorf("mysql: adapter not connected")
	}

	if len(objects) == 0 {
		return nil
	}

	// Handle bulk inserts
	if op.Bulk && len(objects) > 1 {
		return a.bulkInsert(ctx, op, objects)
	}

	// Single insert
	for _, obj := range objects {
		if err := a.singleInsert(ctx, op, obj); err != nil {
			return err
		}
	}

	return nil
}

// singleInsert handles inserting a single record.
func (a *MySQLAdapter) singleInsert(ctx context.Context, op *adapter.Operation, obj interface{}) error {
	// Extract data from object
	data, ok := obj.(map[string]interface{})
	if !ok {
		return fmt.Errorf("mysql: object must be map[string]interface{}")
	}

	// Build INSERT statement
	var fields []string
	var placeholders []string
	var values []interface{}

	for _, prop := range op.Properties {
		// Skip generated fields
		isGenerated := false
		for _, gen := range op.Generated {
			if gen.DataField == prop.DataField {
				isGenerated = true
				break
			}
		}
		if isGenerated {
			continue
		}

		if val, ok := data[prop.ObjectField]; ok {
			fields = append(fields, prop.DataField)
			placeholders = append(placeholders, "?")
			values = append(values, val)
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		op.Statement,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "))

	// Execute insert
	result, err := a.db.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("mysql: insert failed: %w", err)
	}

	// Handle generated IDs
	if len(op.Generated) > 0 {
		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("mysql: failed to get last insert ID: %w", err)
		}

		// Set generated ID back to object
		for _, gen := range op.Generated {
			data[gen.ObjectField] = lastID
		}
	}

	return nil
}

// bulkInsert handles inserting multiple records efficiently.
func (a *MySQLAdapter) bulkInsert(ctx context.Context, op *adapter.Operation, objects []interface{}) error {
	if len(objects) == 0 {
		return nil
	}

	// Extract field names from first object
	firstObj, ok := objects[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("mysql: object must be map[string]interface{}")
	}

	var fields []string
	for _, prop := range op.Properties {
		// Skip generated fields
		isGenerated := false
		for _, gen := range op.Generated {
			if gen.DataField == prop.DataField {
				isGenerated = true
				break
			}
		}
		if !isGenerated {
			if _, ok := firstObj[prop.ObjectField]; ok {
				fields = append(fields, prop.DataField)
			}
		}
	}

	// Build bulk INSERT statement
	var valueSets []string
	var values []interface{}

	for _, obj := range objects {
		data, ok := obj.(map[string]interface{})
		if !ok {
			return fmt.Errorf("mysql: object must be map[string]interface{}")
		}

		var placeholders []string
		for _, prop := range op.Properties {
			isGenerated := false
			for _, gen := range op.Generated {
				if gen.DataField == prop.DataField {
					isGenerated = true
					break
				}
			}
			if !isGenerated {
				if val, ok := data[prop.ObjectField]; ok {
					placeholders = append(placeholders, "?")
					values = append(values, val)
				}
			}
		}
		valueSets = append(valueSets, "("+strings.Join(placeholders, ", ")+")")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		op.Statement,
		strings.Join(fields, ", "),
		strings.Join(valueSets, ", "))

	// Execute bulk insert
	_, err := a.db.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("mysql: bulk insert failed: %w", err)
	}

	return nil
}

// Update modifies existing records in MySQL.
func (a *MySQLAdapter) Update(ctx context.Context, op *adapter.Operation, objects []interface{}) error {
	if a.db == nil {
		return fmt.Errorf("mysql: adapter not connected")
	}

	if len(objects) == 0 {
		return nil
	}

	// Handle each object
	for _, obj := range objects {
		if err := a.singleUpdate(ctx, op, obj); err != nil {
			return err
		}
	}

	return nil
}

// singleUpdate handles updating a single record.
func (a *MySQLAdapter) singleUpdate(ctx context.Context, op *adapter.Operation, obj interface{}) error {
	data, ok := obj.(map[string]interface{})
	if !ok {
		return fmt.Errorf("mysql: object must be map[string]interface{}")
	}

	// Build UPDATE statement
	var setClauses []string
	var values []interface{}

	for _, prop := range op.Properties {
		// Skip identifier fields
		isIdentifier := false
		for _, id := range op.Identifier {
			if id.DataField == prop.DataField {
				isIdentifier = true
				break
			}
		}
		if isIdentifier {
			continue
		}

		if val, ok := data[prop.ObjectField]; ok {
			setClauses = append(setClauses, prop.DataField+" = ?")
			values = append(values, val)
		}
	}

	// Build WHERE clause
	var whereClauses []string
	for _, id := range op.Identifier {
		if val, ok := data[id.ObjectField]; ok {
			whereClauses = append(whereClauses, id.DataField+" = ?")
			values = append(values, val)
		} else {
			return fmt.Errorf("mysql: missing identifier field: %s", id.ObjectField)
		}
	}

	// Add optimistic locking condition if present
	for _, cond := range op.Condition {
		if val, ok := data[cond.ObjectField]; ok {
			whereClauses = append(whereClauses, cond.DataField+" = ?")
			values = append(values, val)
		}
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		op.Statement,
		strings.Join(setClauses, ", "),
		strings.Join(whereClauses, " AND "))

	// Execute update
	result, err := a.db.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("mysql: update failed: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("mysql: failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return adapter.ErrNotFound
	}

	return nil
}

// Delete removes records from MySQL.
func (a *MySQLAdapter) Delete(ctx context.Context, op *adapter.Operation, identifiers []interface{}) error {
	if a.db == nil {
		return fmt.Errorf("mysql: adapter not connected")
	}

	if len(identifiers) == 0 {
		return nil
	}

	// Handle each identifier
	for _, id := range identifiers {
		if err := a.singleDelete(ctx, op, id); err != nil {
			return err
		}
	}

	return nil
}

// singleDelete handles deleting a single record.
func (a *MySQLAdapter) singleDelete(ctx context.Context, op *adapter.Operation, identifier interface{}) error {
	// Build WHERE clause
	var whereClauses []string
	var values []interface{}

	switch id := identifier.(type) {
	case map[string]interface{}:
		// Complex identifier with multiple fields
		for _, idField := range op.Identifier {
			if val, ok := id[idField.ObjectField]; ok {
				whereClauses = append(whereClauses, idField.DataField+" = ?")
				values = append(values, val)
			} else {
				return fmt.Errorf("mysql: missing identifier field: %s", idField.ObjectField)
			}
		}
	default:
		// Simple identifier (single field)
		if len(op.Identifier) != 1 {
			return fmt.Errorf("mysql: simple identifier requires exactly one identifier field")
		}
		whereClauses = append(whereClauses, op.Identifier[0].DataField+" = ?")
		values = append(values, identifier)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s",
		op.Statement,
		strings.Join(whereClauses, " AND "))

	// Execute delete
	result, err := a.db.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("mysql: delete failed: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("mysql: failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return adapter.ErrNotFound
	}

	return nil
}

// Execute runs custom SQL statements or stored procedures.
func (a *MySQLAdapter) Execute(ctx context.Context, action *adapter.Action, params map[string]interface{}) (interface{}, error) {
	if a.db == nil {
		return nil, fmt.Errorf("mysql: adapter not connected")
	}

	// Replace placeholders in statement
	query, args := a.buildQuery(action.Statement, params)

	// Determine if this is a query or exec based on Result mapping
	if action.Result != nil {
		// Execute query (SELECT, CALL with results)
		return a.executeQuery(ctx, query, args)
	}

	// Execute statement (INSERT, UPDATE, DELETE, CALL without results)
	result, err := a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("mysql: execute failed: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return map[string]interface{}{
		"rows_affected": rowsAffected,
	}, nil
}

// executeQuery executes a query and returns results.
func (a *MySQLAdapter) executeQuery(ctx context.Context, query string, args []interface{}) (interface{}, error) {
	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("mysql: query failed: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("mysql: failed to get columns: %w", err)
	}

	// Scan results
	var results []interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("mysql: failed to scan row: %w", err)
		}

		result := make(map[string]interface{})
		for i, col := range columns {
			result[col] = values[i]
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("mysql: row iteration error: %w", err)
	}

	return results, nil
}

// buildQuery replaces named placeholders with positional ones and extracts values.
func (a *MySQLAdapter) buildQuery(query string, params map[string]interface{}) (string, []interface{}) {
	var args []interface{}

	// Replace {param_name} with ? and collect values
	result := query
	for key, value := range params {
		placeholder := "{" + key + "}"
		if strings.Contains(result, placeholder) {
			result = strings.Replace(result, placeholder, "?", 1)
			args = append(args, value)
		}
	}

	return result, args
}

// Helper functions for config extraction
func getStringConfig(config map[string]interface{}, key, defaultValue string) string {
	if val, ok := config[key].(string); ok {
		return val
	}
	return defaultValue
}

func getIntConfig(config map[string]interface{}, key string, defaultValue int) int {
	if val, ok := config[key].(int); ok {
		return val
	}
	if val, ok := config[key].(float64); ok {
		return int(val)
	}
	return defaultValue
}
