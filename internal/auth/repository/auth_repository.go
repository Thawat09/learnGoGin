package repository // กำหนดชื่อแพ็กเกจเป็น "repository" ซึ่งใช้สำหรับการจัดการฐานข้อมูลหรือที่เก็บข้อมูลของแอปพลิเคชัน

import (
	"errors"                    // นำเข้าแพ็กเกจ errors เพื่อใช้ในการสร้างข้อผิดพลาด
	"goGin/internal/auth/model" // นำเข้าแพ็กเกจ model ที่มีโครงสร้างข้อมูล User ซึ่งใช้สำหรับการจัดการข้อมูลผู้ใช้
)

// จำลองฐานข้อมูล (คุณสามารถเชื่อมต่อกับฐานข้อมูลจริงได้ในอนาคต)
var users = []model.User{ // สร้างตัวแปร users ซึ่งเป็นลิสต์ของผู้ใช้จำลอง
	{Username: "admin", Password: "admin123"}, // ผู้ใช้ตัวอย่างที่มีชื่อ "admin" และรหัสผ่าน "admin123"
}

func FindUserByUsername(username string) (model.User, error) { // ฟังก์ชันเพื่อค้นหาผู้ใช้จากชื่อผู้ใช้
	for _, user := range users { // วนลูปผ่านรายการ users
		if user.Username == username { // หากชื่อผู้ใช้ตรงกับที่ค้นหา
			return user, nil // คืนค่าผู้ใช้ที่พบ พร้อมกับค่า error เป็น nil (หมายถึงไม่พบข้อผิดพลาด)
		}
	}
	return model.User{}, errors.New("user not found") // หากไม่พบผู้ใช้ให้คืนค่าเป็น User เปล่า และสร้างข้อผิดพลาดว่า "user not found"
}

func SaveUser(user model.User) error { // ฟังก์ชันสำหรับบันทึกผู้ใช้ใหม่
	users = append(users, user) // เพิ่มผู้ใช้ใหม่เข้าไปในลิสต์ users
	return nil                  // คืนค่า nil เพราะไม่มีข้อผิดพลาด
}
