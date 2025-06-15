package main

import (
	"log"
	"os"
	"time"

	_ "github.com/DMaryanskiy/bookshare-api/docs" // swag init output
	"github.com/DMaryanskiy/bookshare-api/internal/admin"
	"github.com/DMaryanskiy/bookshare-api/internal/auth"
	"github.com/DMaryanskiy/bookshare-api/internal/books"
	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/middleware"
	"github.com/DMaryanskiy/bookshare-api/internal/task/distributor"
	"github.com/DMaryanskiy/bookshare-api/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	debug_mode := os.Getenv("DEBUG")
	if debug_mode != "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Failed to load env:", err)
		}
	}

	db.InitDB()

	redisAddr := os.Getenv("REDIS_ADDR")
	taskDist := distributor.NewTaskDistributor(redisAddr)
	tokenStore := auth.NewTokenStore(redisAddr)

	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})

	// An example of rate limiter
	rateLimiter := middleware.NewRateLimiter(redisClient, map[string]middleware.RateLimitRule{
		"api/v1/login": {
			Limit:  5,
			Window: time.Minute,
		},
		"api/v1/books": {
			Limit:  100,
			Window: time.Minute,
		},
	})

	r := gin.Default()
	r.Use(rateLimiter.Middleware())

	userHandler := user.NewHandler(taskDist, tokenStore)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(ctx *gin.Context) { // Check status handler
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	// Group: Public routes
	public := r.Group("/api/v1")
	public.POST("/register", userHandler.RegisterUser)
	public.GET("/verify", userHandler.VerifyEmail)
	public.POST("/login", userHandler.LoginUser)
	public.POST("/refresh", userHandler.RefreshToken)

	// Group: Authenticated user routes
	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuthMiddleware())
	auth.GET("/me", userHandler.GetMe)
	auth.POST("/logout", userHandler.Logout)

	bookHandler := books.NewHandler()

	// Group: Books CRUD
	booksGroup := r.Group("/api/v1/books")
	booksGroup.Use(middleware.JWTAuthMiddleware())

	booksGroup.POST("", bookHandler.CreateBook)
	booksGroup.GET("", bookHandler.ListBooks)
	booksGroup.GET("/:id", bookHandler.GetBook)
	booksGroup.PUT("/:id", bookHandler.UpdateBook)
	booksGroup.DELETE("/:id", bookHandler.DeleteBook)

	adminHandler := admin.NewHandler()

	// Group: Admin handler
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.JWTAuthMiddleware(), middleware.AdminOnly())

	adminGroup.GET("/users", adminHandler.ListUsers)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(r.Run(":" + port))
}
