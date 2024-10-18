package integration

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/suidevv/tableye-api/controllers"
	"github.com/suidevv/tableye-api/initializers"
	"github.com/suidevv/tableye-api/middleware"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var testRouter *gin.Engine

func init() {
	// Load test configuration
	config, err := initializers.LoadConfig("../..")
	if err != nil {
		log.Fatal("Could not load config:", err)
	}

	// Connect to test database
	initializers.ConnectDB(&config)
	testDB = initializers.DB // Assuming ConnectDB sets a global DB variable

	if testDB == nil {
		log.Fatal("Failed to connect to the test database")
	}

	// Set up the test router
	gin.SetMode(gin.TestMode)
	testRouter = gin.Default()
	setupTestRoutes(testRouter, testDB)
}

func setupTestRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize controllers
	authController := controllers.NewAuthController(db)
	userController := controllers.NewUserController(db)
	casinoController := controllers.NewCasinoController(db)
	gameController := controllers.NewGameController(db)
	playerController := controllers.NewPlayerController(db)
	dealerController := controllers.NewDealerController(db)
	gameSummaryController := controllers.NewGameSummaryController(db)
	transactionController := controllers.NewTransactionController(db)
	adminController := controllers.NewAdminController(db)

	// Setup routes
	api := router.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authController.SignUpUser)
		auth.POST("/login", authController.SignInUser)
		auth.GET("/refresh", authController.RefreshAccessToken)
		auth.GET("/logout", middleware.DeserializeUser(), authController.LogoutUser)
	}

	// User routes
	users := api.Group("/users")
	users.Use(middleware.DeserializeUser())
	{
		users.GET("/me", middleware.AuthorizeRoles("admin"), userController.GetMe)
	}

	// Casino routes
	casinos := api.Group("/casinos")
	casinos.Use(middleware.DeserializeUser())
	{
		casinos.POST("/", middleware.AuthorizeRoles("admin"), casinoController.CreateCasino)
		casinos.GET("/", middleware.AuthorizeRoles("admin", "dealer"), casinoController.FindCasinos)
		casinos.GET("/:casinoId", middleware.AuthorizeRoles("admin", "dealer"), casinoController.FindCasinoById)
		casinos.PUT("/:casinoId", middleware.AuthorizeRoles("admin"), casinoController.UpdateCasino)
		casinos.DELETE("/:casinoId", middleware.AuthorizeRoles("admin"), casinoController.DeleteCasino)
	}

	// Game routes
	games := api.Group("/games")
	games.Use(middleware.DeserializeUser())
	{
		games.POST("/", middleware.AuthorizeRoles("admin"), gameController.CreateGame)
		games.GET("/", middleware.AuthorizeRoles("admin", "dealer"), gameController.FindGames)
		games.GET("/:gameId", middleware.AuthorizeRoles("admin", "dealer"), gameController.FindGameById)
		games.PUT("/:gameId", middleware.AuthorizeRoles("admin"), gameController.UpdateGame)
		games.DELETE("/:gameId", middleware.AuthorizeRoles("admin"), gameController.DeleteGame)
	}

	// Player routes
	players := api.Group("/players")
	players.Use(middleware.DeserializeUser())
	{
		players.POST("/", middleware.AuthorizeRoles("admin"), playerController.CreatePlayer)
		players.GET("/", middleware.AuthorizeRoles("admin", "dealer"), playerController.FindPlayers)
		players.GET("/:playerId", middleware.AuthorizeRoles("admin", "dealer"), playerController.FindPlayerById)
		players.PUT("/:playerId", middleware.AuthorizeRoles("admin"), playerController.UpdatePlayer)
		players.DELETE("/:playerId", middleware.AuthorizeRoles("admin"), playerController.DeletePlayer)
		players.GET("/:playerId/stats", middleware.AuthorizeRoles("admin"), playerController.FindPlayerStats)
	}

	// Dealer routes
	dealers := api.Group("/dealers")
	dealers.Use(middleware.DeserializeUser())
	{
		dealers.POST("/", middleware.AuthorizeRoles("admin"), dealerController.CreateDealer)
		dealers.GET("/", middleware.AuthorizeRoles("admin", "dealer"), dealerController.FindDealers)
		dealers.GET("/:dealerId", middleware.AuthorizeRoles("admin", "dealer"), dealerController.FindDealerById)
		dealers.PUT("/:dealerId", middleware.AuthorizeRoles("admin"), dealerController.UpdateDealer)
		dealers.DELETE("/:dealerId", middleware.AuthorizeRoles("admin"), dealerController.DeleteDealer)
	}

	// Game Summary routes
	gameSummaries := api.Group("/game-summaries")
	gameSummaries.Use(middleware.DeserializeUser())
	{
		gameSummaries.POST("/", middleware.AuthorizeRoles("admin", "dealer"), gameSummaryController.CreateGameSummary)
		gameSummaries.GET("/", middleware.AuthorizeRoles("admin"), gameSummaryController.FindGameSummaries)
		gameSummaries.GET("/:gameSummaryId", middleware.AuthorizeRoles("admin"), gameSummaryController.FindGameSummaryById)
		gameSummaries.PUT("/:gameSummaryId", middleware.AuthorizeRoles("admin", "dealer"), gameSummaryController.UpdateGameSummary)
		gameSummaries.DELETE("/:gameSummaryId", middleware.AuthorizeRoles("admin"), gameSummaryController.DeleteGameSummary)
	}

	// Transaction routes
	transactions := api.Group("/transactions")
	transactions.Use(middleware.DeserializeUser())
	{
		transactions.POST("/", middleware.AuthorizeRoles("admin", "dealer"), transactionController.CreateTransaction)
		transactions.GET("/", middleware.AuthorizeRoles("admin", "dealer"), transactionController.FindTransactions)
		transactions.GET("/:transactionId", middleware.AuthorizeRoles("admin", "dealer"), transactionController.FindTransactionById)
		transactions.PUT("/:transactionId", middleware.AuthorizeRoles("admin"), transactionController.UpdateTransaction)
		transactions.DELETE("/:transactionId", middleware.AuthorizeRoles("admin"), transactionController.DeleteTransaction)
	}

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"))
	{
		admin.POST("/assign-admin", adminController.AssignAdminRole)
	}
}

func GetTestDB() *gorm.DB {
	return testDB
}

func GetTestRouter() *gin.Engine {
	return testRouter
}
