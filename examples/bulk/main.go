package main

import (
	"context"
	"log"

	"github.com/toutaio/toutago-datamapper/adapter"
	"github.com/toutaio/toutago-datamapper/config"
	"github.com/toutaio/toutago-datamapper/engine"
	mysql "github.com/toutaio/toutago-datamapper-mysql"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// For this example to work, you need valid configuration files
	// See sources.yaml and mappings/users.yaml for examples
	
	mapper, err := engine.NewMapper("sources.yaml")
	if err != nil {
		log.Fatalf("Failed to create mapper: %v", err)
	}
	
	mapper.RegisterAdapter("mysql", func(source config.Source) (adapter.Adapter, error) {
		adapter := mysql.NewMySQLAdapter()
		return adapter, nil
	})

	ctx := context.Background()

	// Create a batch of users
	users := []interface{}{
		&User{Name: "Alice Johnson", Email: "alice@example.com"},
		&User{Name: "Bob Smith", Email: "bob@example.com"},
		&User{Name: "Carol White", Email: "carol@example.com"},
		&User{Name: "David Brown", Email: "david@example.com"},
		&User{Name: "Eve Davis", Email: "eve@example.com"},
	}

	log.Printf("Inserting %d users in bulk...\n", len(users))

	// Bulk insert - single database roundtrip
	if err := mapper.Insert(ctx, "User.bulk_insert", users); err != nil {
		log.Fatalf("Failed to bulk insert: %v", err)
	}

	log.Printf("✓ Successfully inserted %d users in a single operation\n", len(users))

	// Fetch all inserted users
	var fetchedUsers []*User
	if err := mapper.FetchMulti(ctx, "User.fetch_all", nil, &fetchedUsers); err != nil {
		log.Fatalf("Failed to fetch users: %v", err)
	}

	log.Printf("✓ Found %d users in database:\n", len(fetchedUsers))
	for _, u := range fetchedUsers {
		log.Printf("  - ID: %d, Name: %s, Email: %s\n", u.ID, u.Name, u.Email)
	}

	mapper.Close()
}
