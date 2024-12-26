package main

import (
	"fmt"
	"net/http"
	"time"

	// นำเข้าชุดคำสั่งที่เกี่ยวกับเส้นทางสำหรับการจัดการ authentication
	authRoutes "goGin/internal/auth/routes"
	// นำเข้าชุดคำสั่งที่เกี่ยวข้องกับการเชื่อมต่อฐานข้อมูล
	"goGin/internal/database"
	// นำเข้าชุดคำสั่งสำหรับการใช้งาน middleware
	"goGin/internal/middleware"
	// นำเข้าชุดคำสั่งที่เกี่ยวข้องกับการจัดการ static files
	staticRoutes "goGin/internal/static/routes"
	// นำเข้าคำสั่งสำหรับการบันทึก error
	"log"
	// นำเข้าคำสั่งสำหรับการใช้งาน environment variables
	"os"

	"github.com/gin-contrib/cors" // นำเข้า CORS middleware
	"github.com/gin-gonic/gin"    // นำเข้าคลาส gin สำหรับสร้าง HTTP server
	"github.com/joho/godotenv"    // นำเข้าคลาส godotenv สำหรับโหลดค่า environment variables จากไฟล์ .env
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	// โหลดไฟล์ .env ที่เก็บค่า environment variables
	err := godotenv.Load()
	if err != nil {
		// ถ้าเกิดข้อผิดพลาดในการโหลดไฟล์ .env ให้แสดง error และหยุดการทำงานของโปรแกรม
		log.Fatalf("Error loading .env file: %v", err)
	}

	// กำหนดโหมดการทำงานของ Gin จากค่า GIN_MODE ในไฟล์ .env
	gin.SetMode(os.Getenv("GIN_MODE"))

	// สร้าง store สำหรับจำข้อมูลการจำกัด
	store := memory.NewStore()

	// กำหนด rate limit (เช่น 100 requests ต่อ 1 นาที)
	rate := limiter.Rate{
		Limit:  5,           // จำนวนคำขอสูงสุด
		Period: time.Minute, // ช่วงเวลาสำหรับจำนวนคำขอ
	}

	// สร้าง limiter จาก store และ rate ที่กำหนด
	instance := limiter.New(store, rate)

	// สร้าง Gin router ใหม่
	r := gin.New()
	// ใช้ middleware สำหรับการบันทึกข้อมูลการเข้าถึง server
	r.Use(gin.Logger())
	// ใช้ middleware สำหรับการจัดการกับ recovery เมื่อเกิดข้อผิดพลาดใน HTTP requests
	r.Use(gin.Recovery())
	// ใช้ middleware สำหรับบันทึก log ที่กำหนดเอง (เช่น log การเข้าถึง API หรือข้อผิดพลาดต่างๆ)
	r.Use(middleware.LoggingMiddleware())

	// ใช้ CORS middleware
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

	// Middleware สำหรับการจัดการ Rate Limiting
	r.Use(func(c *gin.Context) {
		ip := c.ClientIP()

		// เช็คว่า IP นี้สามารถเข้าถึง API ได้หรือไม่
		context := c.Request.Context()                 // สร้าง context จาก HTTP request
		limitContext, err := instance.Get(context, ip) // ใช้ context ที่ถูกต้อง
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		// ถ้า request เกิน limit
		if limitContext.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	})

	// เชื่อมต่อกับ SQL Server โดยใช้ค่าที่ได้จากไฟล์ .env
	sqlServer, err := database.ConnectSQLServer(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	if err != nil {
		// ถ้าเกิดข้อผิดพลาดในการเชื่อมต่อฐานข้อมูล SQL Server ให้แสดง error และหยุดการทำงาน
		log.Fatalf("Failed to connect to SQL Server: %v", err)
	}
	// เมื่อโปรแกรมทำงานเสร็จให้ปิดการเชื่อมต่อกับฐานข้อมูล SQL Server
	defer func() {
		sqlDB, err := sqlServer.DB() // รับ instance ของ SQL database
		if err != nil {
			// ถ้าไม่สามารถดึง SQL database instance ได้ให้แสดง error และออกจากฟังก์ชัน
			fmt.Println("Error retrieving SQL database instance:", err)
			return
		}
		// ปิดการเชื่อมต่อกับ SQL Server
		sqlDB.Close()
	}()

	// เชื่อมต่อกับ Redis โดยใช้ค่าที่ได้จากไฟล์ .env
	redis, err := database.ConnectRedis(
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)
	if err != nil {
		// ถ้าเกิดข้อผิดพลาดในการเชื่อมต่อ Redis ให้แสดง error และหยุดการทำงาน
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	// เมื่อโปรแกรมทำงานเสร็จให้ปิดการเชื่อมต่อกับ Redis
	defer redis.Close()

	// ลงทะเบียนเส้นทางของ Auth
	authRoutes.RegisterAuthRoutes(r)
	// ลงทะเบียนเส้นทางของ Static Files
	staticRoutes.RegisterStaticRoutes(r)

	// อ่านค่าพอร์ตจาก .env และเริ่มต้นเซิร์ฟเวอร์
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080" // ถ้าไม่พบพอร์ตใน .env กำหนดให้เป็นค่าเริ่มต้น 8080
	}

	// เริ่มต้นเซิร์ฟเวอร์ที่พอร์ต 8080
	fmt.Printf("Server is running on http://localhost:%s\n", port)
	r.Run(":" + port)
}
