package main

import (
	"fmt"
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

	"github.com/gin-gonic/gin" // นำเข้าคลาส gin สำหรับสร้าง HTTP server
	"github.com/joho/godotenv" // นำเข้าคลาส godotenv สำหรับโหลดค่า environment variables จากไฟล์ .env
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

	// สร้าง Gin router ใหม่
	r := gin.New()
	// ใช้ middleware สำหรับการบันทึกข้อมูลการเข้าถึง server
	r.Use(gin.Logger())
	// ใช้ middleware สำหรับการจัดการกับ recovery เมื่อเกิดข้อผิดพลาดใน HTTP requests
	r.Use(gin.Recovery())
	// ใช้ middleware สำหรับบันทึก log ที่กำหนดเอง (เช่น log การเข้าถึง API หรือข้อผิดพลาดต่างๆ)
	r.Use(middleware.LoggingMiddleware())

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
