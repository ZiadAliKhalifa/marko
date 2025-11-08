package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/marko/backend/internal/auth"
	"github.com/marko/backend/internal/config"
	"github.com/marko/backend/internal/db"
	"github.com/marko/backend/internal/groups"
	"github.com/marko/backend/internal/locations"
	"github.com/marko/backend/internal/notifications"
)

func main() {
	// Configure logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Set log level based on environment
	if cfg.Environment == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().Str("environment", cfg.Environment).Msg("Starting Marko API")

	// Initialize database
	database, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	// Initialize services
	notificationService := notifications.NewService(cfg.ExpoPushToken)

	// Create Gin router
	router := gin.New()
	
	// Add middleware
	router.Use(gin.Recovery())
	router.Use(zerologMiddleware())
	
	// Add CORS middleware for development
	if cfg.Environment == "development" {
		router.Use(corsMiddleware())
	}

	// Health check endpoint (no auth required)
	router.GET("/healthz", healthCheckHandler(database))

	// API routes
	api := router.Group("/api/v1")
	
	// Auth middleware
	authMiddleware := auth.AuthMiddleware(cfg.SupabaseJWTSecret)

	// Initialize handlers
	groupsHandler := groups.NewHandler(database)
	locationsHandler := locations.NewHandler(database, notificationService)
	notificationsHandler := notifications.NewHandler(database)

	// Register routes
	groupsHandler.RegisterRoutes(api, authMiddleware)
	locationsHandler.RegisterRoutes(api, authMiddleware)
	notificationsHandler.RegisterRoutes(api, authMiddleware)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Str("port", cfg.Port).Msg("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}

// healthCheckHandler returns a health check handler
func healthCheckHandler(database *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check database health
		if err := database.Health(c.Request.Context()); err != nil {
			log.Error().Err(err).Msg("Database health check failed")
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "Database connection failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"timestamp": time.Now().Unix(),
		})
	}
}

// zerologMiddleware adds zerolog logging to Gin
func zerologMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Info().
			Str("client_ip", clientIP).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Int("body_size", bodySize).
			Msg("Request completed")
	}
}

// corsMiddleware adds CORS headers for development
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}