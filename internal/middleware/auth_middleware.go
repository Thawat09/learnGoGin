package middleware // กำหนดชื่อแพ็กเกจเป็น "middleware" เพื่อเก็บโค้ดที่เกี่ยวข้องกับ middleware ใน Gin

import (
	"fmt"      // นำเข้าแพ็กเกจ fmt สำหรับการจัดการข้อความ เช่น การพิมพ์ข้อความลงใน console
	"net/http" // นำเข้าแพ็กเกจ http เพื่อใช้ในการกำหนดสถานะ HTTP และการจัดการการตอบกลับ HTTP

	"github.com/gin-gonic/gin" // นำเข้า Gin framework สำหรับการสร้าง Web application
)

// ฟังก์ชัน LoggingMiddleware สำหรับการล็อกข้อมูลการเข้าถึง API
func LoggingMiddleware() gin.HandlerFunc {
	// คืนค่าฟังก์ชัน HandlerFunc ซึ่งเป็นประเภทของ middleware ใน Gin
	return func(c *gin.Context) {
		// ก่อนที่จะไปที่ handler จริง ๆ ให้ middleware นี้ทำงานก่อน (จะอยู่ในขั้นตอนของ request)
		// อาจจะใช้สำหรับการล็อกข้อมูลที่เกี่ยวข้องกับ request
		c.Next() // เรียก c.Next() เพื่อให้ Gin ไปที่ handler ถัดไป

		// หลังจากที่ไปถึง handler แล้ว (หลังจากที่ response ถูกส่งกลับ)
		// สถานะของ response (เช่น 200, 404, 500) จะถูกบันทึก
		status := c.Writer.Status()                 // ดึงสถานะ HTTP ของ response จาก gin.Writer
		fmt.Printf("Response Status: %d\n", status) // พิมพ์สถานะของ response ลงใน console
	}
}

// ฟังก์ชัน AuthMiddleware สำหรับการตรวจสอบ token ใน header ของ request
func AuthMiddleware() gin.HandlerFunc {
	// คืนค่าฟังก์ชัน HandlerFunc ซึ่งเป็น middleware ที่ใช้ตรวจสอบการยืนยันตัวตน
	return func(c *gin.Context) {
		// ตรวจสอบค่า Authorization header ว่ามีค่าเป็น "valid-token" หรือไม่
		token := c.GetHeader("Authorization") // ดึงค่า Authorization header จาก request
		if token != "valid-token" {           // ถ้าค่า token ไม่ตรงกับ "valid-token"
			// ส่ง response กลับไปยัง client ด้วยสถานะ HTTP 401 Unauthorized และข้อความ error
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort() // หยุดการประมวลผลของ request และไม่ให้ไปที่ handler ถัดไป
			return
		}
		// ถ้า token ถูกต้อง ให้ดำเนินการต่อไปยัง handler ถัดไป
		c.Next()
	}
}
