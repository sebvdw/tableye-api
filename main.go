package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/suidevv/tableye-api/controllers"
	_ "github.com/suidevv/tableye-api/docs" // Import the docs package
	"github.com/suidevv/tableye-api/initializers"
	"github.com/suidevv/tableye-api/routes"
)

// @title           Tableye API
// @version         1.0
// @description     A REST API for Tableye application
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
var (
	server                 *gin.Engine
	AuthController         controllers.AuthController
	AuthRouteController    routes.AuthRouteController
	UserController         controllers.UserController
	UserRouteController    routes.UserRouteController
	PostController         controllers.PostController
	PostRouteController    routes.PostRouteController
	CommentController      controllers.CommentController
	CommentRouteController routes.CommentRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)
	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routes.NewRouteUserController(UserController)
	PostController = controllers.NewPostController(initializers.DB)
	PostRouteController = routes.NewRoutePostController(PostController)
	CommentController = controllers.NewCommentController(initializers.DB)
	CommentRouteController = routes.NewRouteCommentController(CommentController)

	server = gin.Default()
}

// @Summary Health check endpoint
// @Description Get API health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /healthchecker [get]
func healthHandler(ctx *gin.Context) {
	message := "Welcome to Golang with Gorm and Postgres"
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		config.ClientOrigin,
	}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))
	server.StaticFile("", "templates/index.html")
	router := server.Group("/api")

	// Health check endpoint
	router.GET("/healthchecker", healthHandler)

	// Swagger documentation endpoint - will serve Swagger UI directly
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	PostRouteController.PostRoute(router)
	CommentRouteController.CommentRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
