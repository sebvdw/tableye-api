package main

import (
	"fmt"
	"log"
	"os"

	"github.com/suidevv/tableye-api/initializers"
	"github.com/suidevv/tableye-api/models"
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

	err := initializers.DB.AutoMigrate(
		&models.User{},
		&models.Casino{},
		&models.Game{},
		&models.Dealer{},
		&models.Player{},
		&models.GameSummary{},
		&models.Transaction{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	// Create many-to-many relationship tables
	err = initializers.DB.Exec(`
		CREATE TABLE IF NOT EXISTS casino_dealers (
			casino_id UUID REFERENCES casinos(id),
			dealer_id UUID REFERENCES dealers(id),
			PRIMARY KEY (casino_id, dealer_id)
		)
	`).Error
	if err != nil {
		log.Fatal("Failed to create casino_dealers table: ", err)
	}

	err = initializers.DB.Exec(`
		CREATE TABLE IF NOT EXISTS casino_games (
			casino_id UUID REFERENCES casinos(id),
			game_id UUID REFERENCES games(id),
			PRIMARY KEY (casino_id, game_id)
		)
	`).Error
	if err != nil {
		log.Fatal("Failed to create casino_games table: ", err)
	}

	err = initializers.DB.Exec(`
		CREATE TABLE IF NOT EXISTS game_players (
			game_summary_id UUID REFERENCES game_summaries(id),
			player_id UUID REFERENCES players(id),
			PRIMARY KEY (game_summary_id, player_id)
		)
	`).Error
	if err != nil {
		log.Fatal("Failed to create game_players table: ", err)
	}

	fmt.Println("👍 Migration complete")
}
