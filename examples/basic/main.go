package main

import (
	"context"
	"log"

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
	
	mapper.RegisterAdapter("mysql", func(source config.Source) (engine.Adapter, error) {
		adapter := mysql.NewMySQLAdapter()
		return adapter, nil
	})

	if err := mapper.LoadMappings("mappings/"); err != nil {
		log.Fatalf("Failed to load mappings: %v", err)
	}

	ctx := context.Background()

	// Insert a new user
	newUser := &User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	if err := mapper.Insert(ctx, "User.insert", newUser); err != nil {
		log.Fatalf("Failed to insert user: %v", err)
	}

	log.Printf("✓ Created user with ID: %d\n", newUser.ID)

	// Fetch by ID
	user := &User{}
	params := map[string]interface{}{"id": newUser.ID}

	if err := mapper.Fetch(ctx, "User.fetch_by_id", params, user); err != nil {
		log.Fatalf("Failed to fetch user: %v", err)
	}

	log.Printf("✓ Found user: %s (%s)\n", user.Name, user.Email)

	// Update user
	user.Email = "john.doe@example.com"
	if err := mapper.Update(ctx, "User.update", user); err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}

	log.Printf("✓ Updated user email to: %s\n", user.Email)

	// Delete user
	if err := mapper.Delete(ctx, "User.delete", newUser.ID); err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}

	log.Printf("✓ Deleted user ID: %d\n", newUser.ID)

	mapper.Close()
}
