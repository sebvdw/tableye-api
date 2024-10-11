package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/suidevv/tableye-api/initializers"
	"github.com/suidevv/tableye-api/models"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func clearTables() {
	tables := []string{
		"game_players",
		"casino_dealers",
		"casino_games",
		"game_summaries",
		"players",
		"dealers",
		"games",
		"casinos",
		"users",
	}

	for _, table := range tables {
		result := initializers.DB.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if result.Error != nil {
			log.Fatalf("Failed to clear table %s: %v", table, result.Error)
		}
		fmt.Printf("Cleared table: %s\n", table)
	}

	fmt.Println("All tables cleared successfully")
}

func main() {
	clearTables()

	// Seed Users
	users := []models.User{
		{
			ID:        uuid.New(),
			Name:      "John Doe",
			Email:     "john@example.com",
			Password:  hashPassword("password123"),
			Role:      "admin",
			Provider:  "local",
			Verified:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Jane Smith",
			Email:     "jane@example.com",
			Password:  hashPassword("password456"),
			Role:      "user",
			Provider:  "local",
			Verified:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range users {
		if err := initializers.DB.Create(&user).Error; err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
	}

	// Seed Casinos
	casinos := []models.Casino{
		{
			ID:            uuid.New(),
			Name:          "Golden Nugget",
			Location:      "Las Vegas",
			LicenseNumber: "LV12345",
			Description:   "A luxurious casino in the heart of Las Vegas",
			OpeningHours:  "24/7",
			Website:       "https://goldennugget.com",
			PhoneNumber:   "+1-123-456-7890",
			MaxCapacity:   1000,
			Status:        "Active",
			Rating:        4.5,
			OwnerID:       users[0].ID,
		},
		{
			ID:            uuid.New(),
			Name:          "Royal Flush",
			Location:      "Atlantic City",
			LicenseNumber: "AC67890",
			Description:   "Experience the royal treatment at Atlantic City's finest",
			OpeningHours:  "10:00 AM - 4:00 AM",
			Website:       "https://royalflush.com",
			PhoneNumber:   "+1-987-654-3210",
			MaxCapacity:   800,
			Status:        "Active",
			Rating:        4.2,
			OwnerID:       users[1].ID,
		},
	}

	for _, casino := range casinos {
		if err := initializers.DB.Create(&casino).Error; err != nil {
			log.Fatalf("Failed to create casino: %v", err)
		}
	}

	// Seed Games
	games := []models.Game{
		{
			ID:          uuid.New(),
			Name:        "Texas Hold'em",
			Type:        "Poker",
			Description: "The most popular variant of poker",
			MaxPlayers:  10,
			MinPlayers:  2,
			MinBet:      5.0,
			MaxBet:      500.0,
		},
		{
			ID:          uuid.New(),
			Name:        "Blackjack",
			Type:        "Card Game",
			Description: "Try to beat the dealer without going over 21",
			MaxPlayers:  7,
			MinPlayers:  1,
			MinBet:      10.0,
			MaxBet:      1000.0,
		},
	}

	for _, game := range games {
		if err := initializers.DB.Create(&game).Error; err != nil {
			log.Fatalf("Failed to create game: %v", err)
		}
	}

	// Seed Dealers
	dealers := []models.Dealer{
		{
			ID:         uuid.New(),
			UserID:     users[0].ID,
			DealerCode: "D001",
			Status:     "Active",
			GamesDealt: 100,
			Rating:     4.8,
		},
		{
			ID:         uuid.New(),
			UserID:     users[1].ID,
			DealerCode: "D002",
			Status:     "Active",
			GamesDealt: 75,
			Rating:     4.6,
		},
	}

	for _, dealer := range dealers {
		if err := initializers.DB.Create(&dealer).Error; err != nil {
			log.Fatalf("Failed to create dealer: %v", err)
		}
	}

	// Seed Players
	players := []models.Player{
		{
			ID:            uuid.New(),
			UserID:        users[0].ID,
			Nickname:      "Lucky John",
			TotalWinnings: 5000.0,
			GamesPlayed:   50,
			Rank:          "Gold",
			Status:        "Active",
		},
		{
			ID:            uuid.New(),
			UserID:        users[1].ID,
			Nickname:      "Queen Jane",
			TotalWinnings: 3500.0,
			GamesPlayed:   35,
			Rank:          "Silver",
			Status:        "Active",
		},
	}

	for _, player := range players {
		if err := initializers.DB.Create(&player).Error; err != nil {
			log.Fatalf("Failed to create player: %v", err)
		}
	}

	// Seed Game Summaries
	gameSummaries := []models.GameSummary{
		{
			ID:           uuid.New(),
			GameID:       games[0].ID,
			CasinoID:     casinos[0].ID,
			StartTime:    time.Now().Add(-2 * time.Hour),
			EndTime:      time.Now().Add(-30 * time.Minute),
			DealerID:     dealers[0].ID,
			TotalPot:     1000.0,
			Status:       "Completed",
			RoundsPlayed: 30,
			HighestBet:   100.0,
		},
		{
			ID:           uuid.New(),
			GameID:       games[1].ID,
			CasinoID:     casinos[1].ID,
			StartTime:    time.Now().Add(-1 * time.Hour),
			DealerID:     dealers[1].ID,
			TotalPot:     500.0,
			Status:       "In Progress",
			RoundsPlayed: 15,
			HighestBet:   50.0,
		},
	}

	for _, gameSummary := range gameSummaries {
		if err := initializers.DB.Create(&gameSummary).Error; err != nil {
			log.Fatalf("Failed to create game summary: %v", err)
		}
	}

	// Add relationships
	// Casino - Dealers
	for i, casino := range casinos {
		if err := initializers.DB.Exec("INSERT INTO casino_dealers (casino_id, dealer_id) VALUES (?, ?)", casino.ID, dealers[i].ID).Error; err != nil {
			log.Fatalf("Failed to create casino-dealer relationship: %v", err)
		}
	}

	// Casino - Games
	for i, casino := range casinos {
		if err := initializers.DB.Exec("INSERT INTO casino_games (casino_id, game_id) VALUES (?, ?)", casino.ID, games[i].ID).Error; err != nil {
			log.Fatalf("Failed to create casino-game relationship: %v", err)
		}
	}

	// Game Summary - Players
	for _, gameSummary := range gameSummaries {
		for _, player := range players {
			if err := initializers.DB.Exec("INSERT INTO game_players (game_summary_id, player_id) VALUES (?, ?)", gameSummary.ID, player.ID).Error; err != nil {
				log.Fatalf("Failed to create game summary-player relationship: %v", err)
			}
		}
	}

	fmt.Println("Seeding completed successfully!")
}

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Could not hash password", err)
	}
	return string(hashedPassword)
}
