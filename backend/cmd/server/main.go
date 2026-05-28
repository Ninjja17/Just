package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/yourusername/10000hr/internal/handlers"
	"github.com/yourusername/10000hr/internal/middleware"
	"github.com/yourusername/10000hr/internal/repositories"
	"github.com/yourusername/10000hr/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient := initRedis()
	defer redisClient.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	otpRepo := repositories.NewOTPRepository(db, redisClient)
	skillRepo := repositories.NewSkillRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	goalRepo := repositories.NewGoalRepository(db)
	followerRepo := repositories.NewFollowerRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, otpRepo)
	skillService := services.NewSkillService(skillRepo)
	sessionService := services.NewSessionService(sessionRepo)
	goalService := services.NewGoalService(goalRepo)
	socialService := services.NewSocialService(userRepo, followerRepo)
	notificationService := services.NewNotificationService(notificationRepo)
	analyticsService := services.NewAnalyticsService(sessionRepo, skillRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	skillHandler := handlers.NewSkillHandler(skillService)
	sessionHandler := handlers.NewSessionHandler(sessionService)
	goalHandler := handlers.NewGoalHandler(goalService)
	socialHandler := handlers.NewSocialHandler(socialService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	// Setup router
	router := setupRouter(
		authHandler,
		skillHandler,
		sessionHandler,
		goalHandler,
		socialHandler,
		notificationHandler,
		analyticsHandler,
	)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func initDB() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Database connected successfully")
	return db, nil
}

func initRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("Redis connected successfully")
	}

	return client
}

func setupRouter(
	authHandler *handlers.AuthHandler,
	skillHandler *handlers.SkillHandler,
	sessionHandler *handlers.SessionHandler,
	goalHandler *handlers.GoalHandler,
	socialHandler *handlers.SocialHandler,
	notificationHandler *handlers.NotificationHandler,
	analyticsHandler *handlers.AnalyticsHandler,
) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
			auth.POST("/google", authHandler.GoogleLogin)
			auth.POST("/google/callback", authHandler.GoogleCallback)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Auth
			protected.POST("/auth/logout", authHandler.Logout)

			// Skills
			protected.GET("/skills", skillHandler.GetSkills)
			protected.POST("/skills", skillHandler.CreateSkill)
			protected.GET("/skills/:id", skillHandler.GetSkill)
			protected.PUT("/skills/:id", skillHandler.UpdateSkill)
			protected.DELETE("/skills/:id", skillHandler.DeleteSkill)
			protected.GET("/skills/:id/stats", skillHandler.GetSkillStats)

			// Sessions
			protected.GET("/sessions", sessionHandler.GetSessions)
			protected.POST("/sessions", sessionHandler.CreateSession)
			protected.GET("/sessions/:id", sessionHandler.GetSession)
			protected.PUT("/sessions/:id", sessionHandler.UpdateSession)
			protected.DELETE("/sessions/:id", sessionHandler.DeleteSession)
			protected.POST("/sessions/start", sessionHandler.StartTimer)
			protected.POST("/sessions/stop", sessionHandler.StopTimer)

			// Goals
			protected.GET("/goals", goalHandler.GetGoals)
			protected.POST("/goals", goalHandler.CreateGoal)
			protected.PUT("/goals/:id", goalHandler.UpdateGoal)
			protected.DELETE("/goals/:id", goalHandler.DeleteGoal)

			// Analytics
			protected.GET("/analytics/overview", analyticsHandler.GetOverview)
			protected.GET("/analytics/trends", analyticsHandler.GetTrends)
			protected.GET("/analytics/predictions", analyticsHandler.GetPredictions)

			// Social
			protected.GET("/users/:id/profile", socialHandler.GetProfile)
			protected.POST("/social/follow/:id", socialHandler.Follow)
			protected.DELETE("/social/follow/:id", socialHandler.Unfollow)
			protected.GET("/social/followers", socialHandler.GetFollowers)
			protected.GET("/social/following", socialHandler.GetFollowing)
			protected.GET("/social/leaderboard", socialHandler.GetLeaderboard)
			protected.GET("/social/feed", socialHandler.GetFeed)

			// Notifications
			protected.GET("/notifications", notificationHandler.GetNotifications)
			protected.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
			protected.DELETE("/notifications/:id", notificationHandler.DeleteNotification)
		}
	}

	return router
}
