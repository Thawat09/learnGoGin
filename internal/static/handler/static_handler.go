package handler // กำหนดชื่อแพ็กเกจเป็น "handler" ซึ่งใช้สำหรับเก็บโค้ดที่จัดการ request และ response ของ API

import (
	"goGin/internal/static/service" // นำเข้าฟังก์ชันจาก service package เพื่อใช้ในการดึงข้อมูลผู้ใช้จาก service
	"net/http"                      // นำเข้าแพ็กเกจ http เพื่อใช้ในการตั้งค่าสถานะ HTTP และการจัดการ response

	"github.com/gin-gonic/gin" // นำเข้า Gin framework สำหรับสร้าง Web application
)

// ฟังก์ชัน GetUser สำหรับดึงข้อมูลผู้ใช้ตาม ID
func GetUser(c *gin.Context) {
	// ดึงค่าจาก URL parameter ที่ชื่อว่า "id"
	userID := c.Param("id")

	// เรียกใช้ฟังก์ชันจาก service เพื่อดึงข้อมูลผู้ใช้ตาม userID ที่ได้มา
	user, err := service.GetUserByID(userID)

	// ถ้าเกิดข้อผิดพลาดในการดึงข้อมูลผู้ใช้ (เช่น ผู้ใช้ไม่พบ)
	if err != nil {
		// ส่ง response กลับไปยัง client ด้วยสถานะ HTTP 404 Not Found และข้อความ error
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return // หยุดการทำงานของฟังก์ชันนี้
	}

	// ถ้าดึงข้อมูลผู้ใช้สำเร็จ ส่งข้อมูลผู้ใช้กลับไปในรูปแบบ JSON พร้อมกับสถานะ HTTP 200 OK
	c.JSON(http.StatusOK, user)
}
