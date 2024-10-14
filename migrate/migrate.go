// migrate/migrate.go

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/suidevv/tableye-api/initializers"
	"github.com/suidevv/tableye-api/models"
)

func dropDatabase() error {
	if err := initializers.DB.Exec("DROP SCHEMA public CASCADE").Error; err != nil {
		return fmt.Errorf("failed to drop schema: %v", err)
	}
	if err := initializers.DB.Exec("CREATE SCHEMA public").Error; err != nil {
		return fmt.Errorf("failed to recreate schema: %v", err)
	}
	fmt.Println("👍 Database dropped successfully")
	return nil
}

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables: ", err)
	}
	initializers.ConnectDB(&config)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--drop" {
		if err := dropDatabase(); err != nil {
			log.Fatal("Failed to drop database: ", err)
		}
	}

	if err := initializers.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Fatal("Failed to create uuid-ossp extension: ", err)
	}

	if err := initializers.DB.AutoMigrate(
		&models.User{},
		&models.Casino{},
		&models.Game{},
		&models.Dealer{},
		&models.Player{},
		&models.GameSummary{},
		&models.Transaction{},
		&models.Admin{}, // Add the new Admin model
	); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	// Create many-to-many relationship tables and add indexes
	queries := []string{
		`CREATE TABLE IF NOT EXISTS casino_dealers (
			casino_id UUID REFERENCES casinos(id) ON DELETE CASCADE,
			dealer_id UUID REFERENCES dealers(id) ON DELETE CASCADE,
			PRIMARY KEY (casino_id, dealer_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_casino_dealers_casino_id ON casino_dealers(casino_id)`,
		`CREATE INDEX IF NOT EXISTS idx_casino_dealers_dealer_id ON casino_dealers(dealer_id)`,
		`CREATE TABLE IF NOT EXISTS casino_games (
			casino_id UUID REFERENCES casinos(id) ON DELETE CASCADE,
			game_id UUID REFERENCES games(id) ON DELETE CASCADE,
			PRIMARY KEY (casino_id, game_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_casino_games_casino_id ON casino_games(casino_id)`,
		`CREATE INDEX IF NOT EXISTS idx_casino_games_game_id ON casino_games(game_id)`,
		`CREATE TABLE IF NOT EXISTS game_players (
			game_summary_id UUID REFERENCES game_summaries(id) ON DELETE CASCADE,
			player_id UUID REFERENCES players(id) ON DELETE CASCADE,
			PRIMARY KEY (game_summary_id, player_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_game_players_game_summary_id ON game_players(game_summary_id)`,
		`CREATE INDEX IF NOT EXISTS idx_game_players_player_id ON game_players(player_id)`,
		// Additional indexes based on the model definitions
		`CREATE INDEX IF NOT EXISTS idx_dealers_user_id ON dealers(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_game_summaries_game_id ON game_summaries(game_id)`,
		`CREATE INDEX IF NOT EXISTS idx_game_summaries_casino_id ON game_summaries(casino_id)`,
		`CREATE INDEX IF NOT EXISTS idx_game_summaries_dealer_id ON game_summaries(dealer_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_game_summary_id ON transactions(game_summary_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_player_id ON transactions(player_id)`,
		// New index for the role column in the users table
		`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)`,
	}

	for _, query := range queries {
		if err := initializers.DB.Exec(query).Error; err != nil {
			log.Fatalf("Failed to execute query: %s\nError: %v", query, err)
		}
	}

	fmt.Println("👍 Migration complete")
}
