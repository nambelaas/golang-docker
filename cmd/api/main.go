package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nambelaas/golang-docker/configs"
	appinit "github.com/nambelaas/golang-docker/internal/init"
	"github.com/nambelaas/golang-docker/internal/handler"
	"github.com/nambelaas/golang-docker/internal/repository"
)

func main() {
	// Load configuration
	cfg := configs.New()
	log.Println("✓ Configuration loaded")

	// Initialize database
	db, err := appinit.InitDB(cfg)
	if err != nil {
		log.Fatalf("✗ Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	productRepo := repository.NewProductRepository(db)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productRepo)

	// Setup router
	router := setupRouter(productHandler)

	// Start server
	log.Printf("✓ Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("✗ Failed to start server: %v", err)
	}
}

// setupRouter configures all routes
func setupRouter(productHandler *handler.ProductHandler) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Product routes
	api := router.Group("/products")
	{
		api.GET("", productHandler.GetAll)
		api.GET("/:id", productHandler.GetByID)
		api.POST("", productHandler.Create)
		api.PUT("/:id", productHandler.Update)
		api.DELETE("/:id", productHandler.Delete)
	}

	return router
}