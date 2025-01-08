package auth // กำหนดชื่อแพ็กเกจเป็น "auth" ซึ่งจะใช้สำหรับจัดการเส้นทาง (routes) ที่เกี่ยวข้องกับการยืนยันตัวตน (authentication)

import (
	"goGin/internal/auth/handler" // นำเข้าแพ็กเกจ handler ซึ่งมีฟังก์ชันที่จัดการกับการทำงานของการเข้าสู่ระบบ (login) และการลงทะเบียน (register)

	"github.com/gin-gonic/gin" // นำเข้า Gin framework ซึ่งใช้สำหรับสร้าง Web API
)

// RegisterAuthRoutes เป็นฟังก์ชันที่ใช้ลงทะเบียนเส้นทาง (routes) ที่เกี่ยวข้องกับการยืนยันตัวตน
func RegisterAuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth") // สร้างกลุ่มเส้นทาง (route group) ที่มีพาธเริ่มต้นเป็น "/auth"
	// ไม่จำเป็นต้องใช้ AuthMiddleware ที่นี่ เพราะ login และ register เป็น public
	{
		auth.POST("/login", handler.Login)            // ลงทะเบียนเส้นทาง POST "/auth/login" ให้เรียกใช้ฟังก์ชัน Login ใน handler
		auth.POST("/register", handler.Register)      // ลงทะเบียนเส้นทาง POST "/auth/register" ให้เรียกใช้ฟังก์ชัน Register ใน handler
		auth.POST("/encrypt", handler.EncryptMessage) // route สำหรับการเข้ารหัสข้อความ
	}
}
