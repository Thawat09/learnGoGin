package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	authRoutes "goGin/internal/api/auth/routes"
	chackRoutes "goGin/internal/api/check/routes"
	"goGin/internal/config/database"
	"goGin/internal/middleware"
	staticRoutes "goGin/internal/api/static/routes"
	tokenRoutes "goGin/internal/api/token/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	gin.SetMode(os.Getenv("GIN_MODE"))

	store := memory.NewStore()
	rate := limiter.Rate{
		Limit:  30,
		Period: time.Minute,
	}
	instance := limiter.New(store, rate)

	r := gin.New()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=86400")
		c.Header("Pragma", "cache")
		c.Header("Expires", time.Now().Add(24*time.Hour).Format(time.RFC1123))
		c.Next()
	})

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
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "Internal Server Error",
				"status": http.StatusInternalServerError,
			})
			c.Abort()
			return
		}

		if limitContext.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":  "Too many requests",
				"status": http.StatusTooManyRequests,
			})
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
	} else {
		fmt.Println("SQL Server connected successfully")
	}
	defer func() {
		sqlDB, err := sqlServer.DB()
		if err != nil {
			fmt.Println("Error retrieving SQL database instance:", err)
			return
		}
		sqlDB.Close()
	}()

	redisClient, err := database.ConnectRedis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	} else {
		fmt.Println("Redis connected successfully")
	}
	defer redisClient.Close()

	database.SetRedisClient(redisClient)

	apiV1 := r.Group("/api/v1")
	{
		authRoutes.RegisterAuthRoutes(apiV1)
		chackRoutes.CheckRoutes(apiV1)
		tokenRoutes.TokenRoutes(apiV1)
		staticRoutes.StaticRoutes(apiV1)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "Endpoint not found",
			"status": http.StatusNotFound,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on http://localhost:%s\n", port)

	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout:    30 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
