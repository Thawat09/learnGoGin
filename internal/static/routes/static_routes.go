package routes // กำหนดชื่อแพ็กเกจเป็น routes ซึ่งใช้สำหรับการตั้งค่าเส้นทาง (routes) ในแอปพลิเคชัน

import (
	"goGin/internal/middleware"     // นำเข้าแพ็กเกจ middleware ที่ใช้สำหรับการจัดการกลาง เช่น การตรวจสอบการเข้าถึง
	"goGin/internal/static/handler" // นำเข้า handler จากแพ็กเกจ static ที่ใช้ในการจัดการคำขอ (request) สำหรับ static routes

	"github.com/gin-gonic/gin" // นำเข้า Gin ซึ่งเป็น web framework ที่ใช้ในการจัดการ HTTP requests
)

func RegisterStaticRoutes(r *gin.Engine) { // ฟังก์ชันที่ใช้ในการลงทะเบียน static routes
	static := r.Group("/static")            // สร้างกลุ่มเส้นทางใหม่ที่เริ่มต้นด้วย "/static"
	static.Use(middleware.AuthMiddleware()) // ใช้ middleware สำหรับการตรวจสอบการเข้าถึง (เช่น การตรวจสอบ token)
	{
		static.GET("/:id", handler.GetUser) // ลงทะเบียนเส้นทาง GET "/static/:id" ซึ่งจะเรียกฟังก์ชัน handler.GetUser เมื่อได้รับคำขอ
	}
}
