package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/suidevv/tableye-api/initializers"
	"github.com/suidevv/tableye-api/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	numUsers         = 40 // Adjusted to accommodate admin, casino owners, and dealers
	numCasinos       = 10
	numGames         = 5
	numDealers       = 30
	numPlayers       = 50
	numGameSummaries = 100
	numTransactions  = 500
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables: ", err)
	}
	initializers.ConnectDB(&config)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	db := initializers.DB

	clearTables(db)

	users := seedUsers(db)
	casinos := seedCasinos(db)
	games := seedGames(db)
	dealers := seedDealers(db, users[numCasinos:numCasinos+numDealers])
	players := seedPlayers(db)
	gameSummaries := seedGameSummaries(db, games, casinos, dealers)
	seedTransactions(db, gameSummaries, players)

	createRelationships(db, casinos, dealers, games, gameSummaries, players)

	fmt.Println("Seeding completed successfully!")
}

func clearTables(db *gorm.DB) {
	tables := []string{"game_players", "casino_dealers", "casino_games", "transactions", "game_summaries", "players", "dealers", "games", "casinos", "users"}
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			log.Fatalf("Failed to clear table %s: %v", table, err)
		}
	}
}

func seedUsers(db *gorm.DB) []models.User {
	users := make([]models.User, numUsers)
	roles := []string{"admin", "casino_owner", "dealer"}
	for i := 0; i < numUsers; i++ {
		users[i] = models.User{
			ID:        uuid.New(),
			Name:      fmt.Sprintf("User %d", i+1),
			Email:     fmt.Sprintf("user%d@example.com", i+1),
			Password:  hashPassword(fmt.Sprintf("password%d", i+1)),
			Role:      roles[i%len(roles)],
			Provider:  "local",
			Verified:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	if err := db.Create(&users).Error; err != nil {
		log.Fatalf("Failed to create users: %v", err)
	}
	return users
}

func seedCasinos(db *gorm.DB) []models.Casino {
	casinos := make([]models.Casino, numCasinos)
	for i := 0; i < numCasinos; i++ {
		casinos[i] = models.Casino{
			ID:            uuid.New(),
			Name:          fmt.Sprintf("Casino %d", i+1),
			Location:      fmt.Sprintf("City %d", i+1),
			LicenseNumber: fmt.Sprintf("LN%05d", i+1),
			Description:   fmt.Sprintf("Description for Casino %d", i+1),
			OpeningHours:  "24/7",
			Website:       fmt.Sprintf("https://casino%d.com", i+1),
			PhoneNumber:   fmt.Sprintf("+1-123-456-%04d", i+1),
			MaxCapacity:   500 + rand.Intn(1500),
			Status:        "Active",
			Rating:        4.0 + rand.Float32(),
		}
	}
	if err := db.Create(&casinos).Error; err != nil {
		log.Fatalf("Failed to create casinos: %v", err)
	}
	return casinos
}

func seedGames(db *gorm.DB) []models.Game {
	gameTypes := []string{"Poker", "Blackjack", "Roulette", "Slots", "Baccarat"}
	games := make([]models.Game, numGames)
	for i := 0; i < numGames; i++ {
		games[i] = models.Game{
			ID:          uuid.New(),
			Name:        fmt.Sprintf("%s %d", gameTypes[i], i+1),
			Type:        gameTypes[i],
			Description: fmt.Sprintf("Description for %s %d", gameTypes[i], i+1),
			MaxPlayers:  4 + rand.Intn(8),
			MinPlayers:  1 + rand.Intn(3),
			MinBet:      float64(5 + rand.Intn(20)),
			MaxBet:      float64(100 + rand.Intn(900)),
		}
	}
	if err := db.Create(&games).Error; err != nil {
		log.Fatalf("Failed to create games: %v", err)
	}
	return games
}

func seedDealers(db *gorm.DB, dealerUsers []models.User) []models.Dealer {
	dealers := make([]models.Dealer, numDealers)
	for i := 0; i < numDealers; i++ {
		dealers[i] = models.Dealer{
			ID:         uuid.New(),
			UserID:     dealerUsers[i].ID,
			DealerCode: fmt.Sprintf("D%04d", i+1),
			Status:     "Active",
			GamesDealt: rand.Intn(200),
			Rating:     4.0 + rand.Float32(),
		}
	}
	if err := db.Create(&dealers).Error; err != nil {
		log.Fatalf("Failed to create dealers: %v", err)
	}
	return dealers
}

func seedPlayers(db *gorm.DB) []models.Player {
	ranks := []string{"Bronze", "Silver", "Gold", "Platinum", "Diamond"}
	players := make([]models.Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = models.Player{
			ID:            uuid.New(),
			Nickname:      fmt.Sprintf("Player%d", i+1),
			TotalWinnings: float64(rand.Intn(10000)),
			Rank:          ranks[rand.Intn(len(ranks))],
			Status:        "Active",
		}
	}
	if err := db.Create(&players).Error; err != nil {
		log.Fatalf("Failed to create players: %v", err)
	}
	return players
}

func seedGameSummaries(db *gorm.DB, games []models.Game, casinos []models.Casino, dealers []models.Dealer) []models.GameSummary {
	gameSummaries := make([]models.GameSummary, numGameSummaries)
	for i := 0; i < numGameSummaries; i++ {
		startTime := time.Now().Add(time.Duration(-rand.Intn(30)) * 24 * time.Hour)
		endTime := startTime.Add(time.Duration(rand.Intn(4)+1) * time.Hour)
		gameSummaries[i] = models.GameSummary{
			ID:           uuid.New(),
			GameID:       games[rand.Intn(len(games))].ID,
			CasinoID:     casinos[rand.Intn(len(casinos))].ID,
			StartTime:    startTime,
			EndTime:      endTime,
			DealerID:     dealers[rand.Intn(len(dealers))].ID,
			TotalPot:     float64(100 + rand.Intn(10000)),
			Status:       []string{"Completed", "In Progress"}[rand.Intn(2)],
			RoundsPlayed: rand.Intn(50),
			HighestBet:   float64(50 + rand.Intn(950)),
		}
	}
	if err := db.Create(&gameSummaries).Error; err != nil {
		log.Fatalf("Failed to create game summaries: %v", err)
	}
	return gameSummaries
}

func seedTransactions(db *gorm.DB, gameSummaries []models.GameSummary, players []models.Player) {
	transactions := make([]models.Transaction, numTransactions)
	for i := 0; i < numTransactions; i++ {
		gameSummary := gameSummaries[rand.Intn(len(gameSummaries))]
		player := players[rand.Intn(len(players))]
		transactionTime := gameSummary.StartTime.Add(time.Duration(rand.Intn(int(gameSummary.EndTime.Sub(gameSummary.StartTime)))))
		transactionType := []string{"bet", "win"}[rand.Intn(2)]
		amount := float64(10 + rand.Intn(990))
		if transactionType == "bet" {
			amount = -amount
		}
		transactions[i] = models.Transaction{
			ID:            uuid.New(),
			GameSummaryID: gameSummary.ID,
			PlayerID:      player.ID,
			Amount:        amount,
			Type:          transactionType,
			Outcome:       map[string]string{"bet": "loss", "win": "win"}[transactionType],
			CreatedAt:     transactionTime,
			UpdatedAt:     transactionTime,
		}
	}
	if err := db.Create(&transactions).Error; err != nil {
		log.Fatalf("Failed to create transactions: %v", err)
	}
}

func createRelationships(db *gorm.DB, casinos []models.Casino, dealers []models.Dealer, games []models.Game, gameSummaries []models.GameSummary, players []models.Player) {
	// Casino - Dealers
	casinoDealerMap := make(map[string]map[string]bool)
	for _, casino := range casinos {
		casinoDealerMap[casino.ID.String()] = make(map[string]bool)
		numDealers := rand.Intn(4) + 2
		dealersAdded := 0
		for dealersAdded < numDealers {
			dealer := dealers[rand.Intn(len(dealers))]
			if !casinoDealerMap[casino.ID.String()][dealer.ID.String()] {
				if err := db.Exec("INSERT INTO casino_dealers (casino_id, dealer_id) VALUES (?, ?)", casino.ID, dealer.ID).Error; err != nil {
					log.Printf("Failed to create casino-dealer relationship: %v", err)
					continue
				}
				casinoDealerMap[casino.ID.String()][dealer.ID.String()] = true
				dealersAdded++
			}
		}
	}

	// Casino - Games
	casinoGameMap := make(map[string]map[string]bool)
	for _, casino := range casinos {
		casinoGameMap[casino.ID.String()] = make(map[string]bool)
		numGames := rand.Intn(3) + 2
		gamesAdded := 0
		for gamesAdded < numGames {
			game := games[rand.Intn(len(games))]
			if !casinoGameMap[casino.ID.String()][game.ID.String()] {
				if err := db.Exec("INSERT INTO casino_games (casino_id, game_id) VALUES (?, ?)", casino.ID, game.ID).Error; err != nil {
					log.Printf("Failed to create casino-game relationship: %v", err)
					continue
				}
				casinoGameMap[casino.ID.String()][game.ID.String()] = true
				gamesAdded++
			}
		}
	}

	// Game Summary - Players
	gameSummaryPlayerMap := make(map[string]map[string]bool)
	for _, gameSummary := range gameSummaries {
		gameSummaryPlayerMap[gameSummary.ID.String()] = make(map[string]bool)
		numPlayers := rand.Intn(6) + 2
		playersAdded := 0
		for playersAdded < numPlayers {
			player := players[rand.Intn(len(players))]
			if !gameSummaryPlayerMap[gameSummary.ID.String()][player.ID.String()] {
				if err := db.Exec("INSERT INTO game_players (game_summary_id, player_id) VALUES (?, ?)", gameSummary.ID, player.ID).Error; err != nil {
					log.Printf("Failed to create game summary-player relationship: %v", err)
					continue
				}
				gameSummaryPlayerMap[gameSummary.ID.String()][player.ID.String()] = true
				playersAdded++
			}
		}
	}
}

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Could not hash password", err)
	}
	return string(hashedPassword)
}
