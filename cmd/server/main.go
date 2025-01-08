package main

import (
	"fmt"
	authRoutes "goGin/internal/auth/routes"
	"goGin/internal/database"
	"goGin/internal/middleware"
	staticRoutes "goGin/internal/static/routes"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	gin.SetMode(os.Getenv("GIN_MODE"))
	store := memory.NewStore()

	rate := limiter.Rate{
		Limit:  5,
		Period: time.Minute,
	}

	instance := limiter.New(store, rate)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())

	corsOptions := cors.Config{
		AllowOrigins: []string{
			"http://4.156.59.104",
			"http://localhost",
			"http://localhost:8081",
			"http://192.10.51.7",
			"http://192.10.51.7:8081",
			"http://frontend.thawat.site",
			"https://frontend.thawat.site",
		},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Content-Type", "Authorization", "Accept", "X-Requested-With", "Cache-Control",
			"X-CSRF-Token", "X-User-IP", "UserAgent", "UserIP", "UserOS", "Token", "HostName",
			"City", "Region", "Country", "Loc", "Org", "Postal", "Timezone",
		},
	}

	r.Use(cors.New(corsOptions))

	r.Use(func(c *gin.Context) {
		ip := c.ClientIP()
		context := c.Request.Context()
		limitContext, err := instance.Get(context, ip)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		if limitContext.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	})

	sqlServer, err := database.ConnectSQLServer(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)

	if err != nil {
		log.Fatalf("Failed to connect to SQL Server: %v", err)
	}

	defer func() {
		sqlDB, err := sqlServer.DB()

		if err != nil {
			fmt.Println("Error retrieving SQL database instance:", err)
			return
		}

		sqlDB.Close()
	}()

	redis, err := database.ConnectRedis(
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)

	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	defer redis.Close()

	authRoutes.RegisterAuthRoutes(r)
	staticRoutes.RegisterStaticRoutes(r)
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on http://localhost:%s\n", port)
	r.Run(":" + port)
}
