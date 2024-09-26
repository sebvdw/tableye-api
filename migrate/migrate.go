package main

import (
	"fmt"
	"log"
	"os"

	"github.com/suidevv/golang-tableye/initializers"
	"github.com/suidevv/golang-tableye/models"
)

func dropDatabase() error {
	// This works for PostgreSQL
	err := initializers.DB.Exec("DROP SCHEMA public CASCADE").Error
	if err != nil {
		return fmt.Errorf("failed to drop schema: %v", err)
	}

	err = initializers.DB.Exec("CREATE SCHEMA public").Error
	if err != nil {
		return fmt.Errorf("failed to recreate schema: %v", err)
	}

	fmt.Println("👍 Database dropped successfully")
	return nil
}

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	// Check if the --drop flag is provided
	if len(os.Args) > 1 && os.Args[1] == "--drop" {
		err := dropDatabase()
		if err != nil {
			log.Fatal("Failed to drop database: ", err)
		}
	}
	initializers.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	initializers.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	fmt.Println("👍 Migration complete")
}
