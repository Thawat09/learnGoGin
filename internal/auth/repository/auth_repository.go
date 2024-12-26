package repository // กำหนดชื่อแพ็กเกจเป็น "repository" ซึ่งใช้สำหรับการจัดการฐานข้อมูลหรือที่เก็บข้อมูลของแอปพลิเคชัน

import (
	"errors"                    // นำเข้าแพ็กเกจ errors เพื่อใช้ในการสร้างข้อผิดพลาด
	"goGin/internal/auth/model" // นำเข้าแพ็กเกจ model ที่มีโครงสร้างข้อมูล Users ซึ่งใช้สำหรับการจัดการข้อมูลผู้ใช้
	"goGin/internal/database"   // นำเข้าแพ็กเกจ database เพื่อใช้ในการเชื่อมต่อฐานข้อมูล

	"strings" // นำเข้า strings เพื่อใช้ในการตรวจสอบข้อผิดพลาดที่เกี่ยวกับข้อมูลซ้ำ

	"gorm.io/gorm" // นำเข้า GORM สำหรับจัดการฐานข้อมูล
)

// ฟังก์ชัน FindUserByUsername ค้นหาผู้ใช้ในฐานข้อมูลตามชื่อผู้ใช้
func FindUserByUsername(username string) (model.Users, error) {
	var user model.Users
	// ใช้ GORM สำหรับค้นหาผู้ใช้ในฐานข้อมูล
	err := database.DB.Where("username = ?", username).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Users{}, errors.New("user not found") // ส่งข้อความที่เข้าใจง่ายเมื่อไม่พบผู้ใช้
		}
		return model.Users{}, err
	}

	return user, nil
}

// ฟังก์ชัน FindDepartmentById ค้นหาข้อมูลแผนกจาก DepartmentId
func FindDepartmentById(db *gorm.DB, departmentId int) (*model.Departments, error) {
	var department model.Departments
	// ค้นหาข้อมูลแผนกในฐานข้อมูลโดยใช้ DepartmentId
	if err := db.First(&department, departmentId).Error; err != nil {
		return nil, err // หากไม่พบข้อมูลหรือเกิดข้อผิดพลาดในการค้นหา
	}
	return &department, nil // คืนค่าข้อมูลแผนกที่พบ
}

// ฟังก์ชันบันทึกผู้ใช้ใหม่ลงในฐานข้อมูล
func SaveUser(user model.Users) error {
	if err := database.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "UNIQUE KEY") {
			return errors.New("duplicate username or email") // แจ้งข้อผิดพลาดเฉพาะกรณีข้อมูลซ้ำ
		}
		return err // ส่งต่อข้อผิดพลาดที่เกิดขึ้นอื่น ๆ
	}
	return nil
}
