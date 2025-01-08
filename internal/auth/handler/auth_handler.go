package handler

import (
	"fmt"
	"goGin/internal/auth/service" // นำเข้าชุดฟังก์ชันที่เกี่ยวกับการตรวจสอบและลงทะเบียนผู้ใช้
	"net/http"                    // ใช้สำหรับการกำหนดสถานะ HTTP เช่น 200, 400, 401 เป็นต้น

	"github.com/gin-gonic/gin" // นำเข้า Gin framework สำหรับสร้าง Web API
)

// ฟังก์ชัน Login รับค่าจากการร้องขอ HTTP เพื่อเข้าสู่ระบบ
func Login(c *gin.Context) {
	// กำหนดโครงสร้างของข้อมูลที่ต้องการรับจากผู้ใช้ (username และ password)
	var loginReq struct {
		Username string `json:"username"` // รับค่าชื่อผู้ใช้
		Password string `json:"password"` // รับค่ารหัสผ่าน
	}

	// ตรวจสอบว่าผู้ใช้ส่งข้อมูลในรูปแบบ JSON และเชื่อมต่อข้อมูลกับตัวแปร loginReq หรือไม่
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		// ถ้าเกิดข้อผิดพลาดในการแปลงข้อมูล ให้ส่งกลับ error message ด้วยรหัสสถานะ 400
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fmt.Println("data: ", loginReq)
	fmt.Println("username: ", loginReq.Username)
	fmt.Println("password: ", loginReq.Password)

	// ถอดรหัส Username
	decryptedUsername, err := service.Decrypt(loginReq.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt username"})
		return
	}

	// ถอดรหัส Password
	decryptedPassword, err := service.Decrypt(loginReq.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt password"})
		return
	}

	// ตรวจสอบชื่อผู้ใช้และรหัสผ่านโดยการเรียกใช้ฟังก์ชัน Login จาก service
	if err := service.Login(decryptedUsername, decryptedPassword); err != nil {
		// ถ้าเกิดข้อผิดพลาดในการตรวจสอบข้อมูลผู้ใช้ ให้ส่งกลับ error message ด้วยรหัสสถานะ 401
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// ถ้าผ่านการตรวจสอบ ให้ส่งกลับข้อความว่าเข้าสู่ระบบสำเร็จ พร้อมกับรหัสสถานะ 200
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// ฟังก์ชัน Register รับค่าจากการร้องขอ HTTP เพื่อสมัครสมาชิกใหม่
func Register(c *gin.Context) {
	// กำหนดโครงสร้างของข้อมูลที่ต้องการรับจากผู้ใช้ (username, password และ departmentId)
	var userReq struct {
		Username     string `json:"username"`     // รับค่าชื่อผู้ใช้
		Password     string `json:"password"`     // รับค่ารหัสผ่าน
		Email        string `json:"email"`        // รับค่าอีเมล
		DepartmentId int    `json:"departmentId"` // รับค่า DepartmentId
	}

	// ตรวจสอบว่าผู้ใช้ส่งข้อมูลในรูปแบบ JSON และเชื่อมต่อข้อมูลกับตัวแปร userReq หรือไม่
	if err := c.ShouldBindJSON(&userReq); err != nil {
		// ถ้าเกิดข้อผิดพลาดในการแปลงข้อมูล ให้ส่งกลับ error message ด้วยรหัสสถานะ 400
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ตรวจสอบว่า Username, Email, และ Password ไม่เป็นค่าว่าง
	if userReq.Username == "" || userReq.Email == "" || userReq.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password must not be empty"})
		return
	}

	// สมัครสมาชิกใหม่ โดยการเรียกใช้ฟังก์ชัน Register จาก service
	if err := service.Register(userReq.Username, userReq.Password, userReq.Email, userReq.DepartmentId); err != nil {
		// ถ้าเกิดข้อผิดพลาดในการสมัครสมาชิก ให้ส่งกลับ error message ด้วยรหัสสถานะ 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register"})
		return
	}

	// ถ้าการสมัครสมาชิกสำเร็จ ให้ส่งกลับข้อความว่าสมัครสมาชิกสำเร็จ พร้อมกับรหัสสถานะ 200
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func EncryptMessage(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ทำการเข้ารหัสข้อความ
	encryptedUsername, err := service.Encrypt(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt message"})
		return
	}

	encryptedPassword, err := service.Encrypt(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"encryptedUsername": encryptedUsername,
		"encryptedPassword": encryptedPassword,
	})
}
