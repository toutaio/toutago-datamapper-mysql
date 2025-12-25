package mysql

import (
	"context"
	"testing"

	"github.com/toutago/toutago-datamapper/adapter"
)

func TestMySQLAdapter_Name(t *testing.T) {
	a := NewMySQLAdapter()
	if a.Name() != "mysql" {
		t.Errorf("expected adapter name 'mysql', got '%s'", a.Name())
	}
}

func TestMySQLAdapter_ConnectInvalidConfig(t *testing.T) {
	a := NewMySQLAdapter()
	ctx := context.Background()

	// Test with invalid connection - should fail ping
	config := map[string]interface{}{
		"host":     "invalid-host",
		"port":     3306,
		"user":     "root",
		"password": "",
		"database": "test",
	}

	err := a.Connect(ctx, config)
	if err == nil {
		t.Error("expected connection error to invalid host, got nil")
		a.Close()
	}
}

func TestMySQLAdapter_BuildQuery(t *testing.T) {
	a := NewMySQLAdapter()

	tests := []struct {
		name     string
		query    string
		params   map[string]interface{}
		expected string
		argCount int
	}{
		{
			name:     "Single parameter",
			query:    "SELECT * FROM users WHERE id = {id}",
			params:   map[string]interface{}{"id": 123},
			expected: "SELECT * FROM users WHERE id = ?",
			argCount: 1,
		},
		{
			name:     "Multiple parameters",
			query:    "SELECT * FROM users WHERE name = {name} AND email = {email}",
			params:   map[string]interface{}{"name": "John", "email": "john@example.com"},
			expected: "SELECT * FROM users WHERE name = ? AND email = ?",
			argCount: 2,
		},
		{
			name:     "No parameters",
			query:    "SELECT * FROM users",
			params:   map[string]interface{}{},
			expected: "SELECT * FROM users",
			argCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, args := a.buildQuery(tt.query, tt.params)
			
			if result != tt.expected {
				t.Errorf("expected query '%s', got '%s'", tt.expected, result)
			}
			
			if len(args) != tt.argCount {
				t.Errorf("expected %d arguments, got %d", tt.argCount, len(args))
			}
		})
	}
}

func TestMySQLAdapter_ConfigHelpers(t *testing.T) {
	config := map[string]interface{}{
		"string_val": "test",
		"int_val":    42,
		"float_val":  3.14,
	}

	// Test getStringConfig
	if val := getStringConfig(config, "string_val", "default"); val != "test" {
		t.Errorf("expected 'test', got '%s'", val)
	}
	
	if val := getStringConfig(config, "missing", "default"); val != "default" {
		t.Errorf("expected 'default', got '%s'", val)
	}

	// Test getIntConfig
	if val := getIntConfig(config, "int_val", 0); val != 42 {
		t.Errorf("expected 42, got %d", val)
	}
	
	if val := getIntConfig(config, "float_val", 0); val != 3 {
		t.Errorf("expected 3, got %d", val)
	}
	
	if val := getIntConfig(config, "missing", 99); val != 99 {
		t.Errorf("expected 99, got %d", val)
	}
}

func TestMySQLAdapter_NotConnectedErrors(t *testing.T) {
	a := NewMySQLAdapter()
	ctx := context.Background()
	
	op := &adapter.Operation{
		Type:      adapter.OpFetch,
		Statement: "users",
	}

	// Test Fetch without connection
	_, err := a.Fetch(ctx, op, nil)
	if err == nil {
		t.Error("expected error when fetching without connection")
	}

	// Test Insert without connection
	err = a.Insert(ctx, op, []interface{}{})
	if err == nil {
		t.Error("expected error when inserting without connection")
	}

	// Test Update without connection
	err = a.Update(ctx, op, []interface{}{})
	if err == nil {
		t.Error("expected error when updating without connection")
	}

	// Test Delete without connection
	err = a.Delete(ctx, op, []interface{}{})
	if err == nil {
		t.Error("expected error when deleting without connection")
	}

	// Test Execute without connection
	action := &adapter.Action{
		Statement: "SELECT 1",
	}
	_, err = a.Execute(ctx, action, nil)
	if err == nil {
		t.Error("expected error when executing without connection")
	}
}
