package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"eatright-backend/internal/app/config"
	"eatright-backend/internal/app/handlers"
	"eatright-backend/internal/app/middlewares"
	"eatright-backend/internal/app/repositories"
	"eatright-backend/internal/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	_ "eatright-backend/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title EatRight API
// @version 1.0
// @description Production-ready backend for EatRight food waste reduction platform
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@eatright.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Setup database
	db, err := config.SetupDatabase(cfg.Database.URL, cfg.Server.Env)
	if err != nil {
		log.Fatalf("❌ Failed to setup database: %v", err)
	}

	// Auto-migrate models (comment out after first run)
	// err = config.AutoMigrate(db,
	// 	&models.User{},
	// 	&models.Restaurant{},
	// 	&models.Listing{},
	// 	&models.Order{},
	// )
	// if err != nil {
	// 	log.Fatalf("❌ Failed to migrate database: %v", err)
	// }

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	restaurantRepo := repositories.NewRestaurantRepository(db)
	listingRepo := repositories.NewListingRepository(db)
	orderRepo := repositories.NewOrderRepository(db, listingRepo)

	// Initialize services
	authService, err := services.NewAuthService(userRepo, cfg)
	if err != nil {
		log.Fatalf("❌ Failed to create auth service: %v", err)
	}
	userService := services.NewUserService(userRepo)
	restaurantService := services.NewRestaurantService(restaurantRepo, userRepo)
	listingService := services.NewListingService(listingRepo, restaurantRepo)
	orderService := services.NewOrderService(orderRepo, listingRepo, restaurantRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService)
	listingHandler := handlers.NewListingHandler(listingService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
		AppName:      "EatRight API",
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(middlewares.Logger())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     "GET,POST,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"env":    cfg.Server.Env,
		})
	})

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Serve swagger.json directly
	app.Static("/swagger", "./docs")

	// API routes
	api := app.Group("/api")

	// Auth routes (public)
	authRoutes := api.Group("/auth")
	authRoutes.Post("/verify", authHandler.VerifyToken)

	// User routes (protected)
	userRoutes := api.Group("/users", middlewares.AuthMiddleware(cfg))
	userRoutes.Get("/me", userHandler.GetMe)

	// Restaurant routes
	restaurantRoutes := api.Group("/restaurants")
	restaurantRoutes.Get("/", restaurantHandler.GetRestaurants)       // Public
	restaurantRoutes.Get("/:id", restaurantHandler.GetRestaurantByID) // Public
	restaurantRoutes.Post("/",                                        // Protected, restaurant role only
		middlewares.AuthMiddleware(cfg),
		middlewares.RestaurantOnly(),
		restaurantHandler.CreateRestaurant,
	)

	// Listing routes
	listingRoutes := api.Group("/listings")
	listingRoutes.Get("/", listingHandler.GetListings)       // Public
	listingRoutes.Get("/:id", listingHandler.GetListingByID) // Public

	// Create listing (protected, restaurant role only)
	api.Post("/restaurants/:id/listings",
		middlewares.AuthMiddleware(cfg),
		middlewares.RestaurantOnly(),
		listingHandler.CreateListing,
	)

	// Update listing stock and status (protected, restaurant role only)
	listingRoutes.Patch("/:id/stock",
		middlewares.AuthMiddleware(cfg),
		middlewares.RestaurantOnly(),
		listingHandler.UpdateStock,
	)
	listingRoutes.Patch("/:id/status",
		middlewares.AuthMiddleware(cfg),
		middlewares.RestaurantOnly(),
		listingHandler.UpdateStatus,
	)

	// Order routes (protected)
	orderRoutes := api.Group("/orders", middlewares.AuthMiddleware(cfg))
	orderRoutes.Post("/", orderHandler.CreateOrder)
	orderRoutes.Get("/me", orderHandler.GetMyOrders)
	orderRoutes.Patch("/:id/status",
		middlewares.RestaurantOnly(),
		orderHandler.UpdateOrderStatus,
	)

	// Start server
	port := cfg.Server.Port
	log.Printf("🚀 Server starting on port %s (env: %s)", port, cfg.Server.Env)

	// Graceful shutdown
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("❌ Server shutdown error: %v", err)
	}

	log.Println("✅ Server shutdown complete")
}

// customErrorHandler handles errors globally
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}
