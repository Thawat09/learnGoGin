package service // กำหนดชื่อแพ็กเกจเป็น "service" ซึ่งจะมีฟังก์ชันที่เกี่ยวข้องกับการทำงานหลัก เช่น การเข้าสู่ระบบ และการลงทะเบียน

import (
	"errors"                         // นำเข้าแพ็กเกจ errors สำหรับการสร้าง error messages
	"fmt"                            // นำเข้าแพ็กเกจ fmt สำหรับการพิมพ์ข้อความ
	"goGin/internal/auth/model"      // นำเข้าแพ็กเกจ model เพื่อใช้ประเภทข้อมูล User
	"goGin/internal/auth/repository" // นำเข้าแพ็กเกจ repository ซึ่งจะเชื่อมต่อกับฐานข้อมูลเพื่อจัดการข้อมูลผู้ใช้
)

// ฟังก์ชัน Login ใช้ในการตรวจสอบการเข้าสู่ระบบโดยใช้ชื่อผู้ใช้และรหัสผ่าน
func Login(username, password string) error {
	// หา user ใน repository โดยใช้ชื่อผู้ใช้ที่ส่งมา
	user, err := repository.FindUserByUsername(username)
	if err != nil {
		return err // ถ้าหาผู้ใช้ไม่เจอจะคืนค่าผิดพลาด (error)
	}

	// ตรวจสอบ password ว่า match กับ hash ที่เก็บไว้หรือไม่
	if user.Password != password { // ในกรณีนี้ยังไม่ใช้การเข้ารหัส password จริงๆ
		fmt.Println("Passwords do not match")    // แสดงข้อความถ้ารหัสผ่านไม่ตรงกัน
		return errors.New("invalid credentials") // คืนค่า error ว่า credentials ไม่ถูกต้อง
	}
	return nil // คืนค่าปกติถ้าทุกอย่างถูกต้อง
}

// ฟังก์ชัน Register ใช้ในการลงทะเบียนผู้ใช้ใหม่
func Register(username, password string) error {
	// เข้ารหัส password (ในโค้ดนี้ยังไม่ได้เข้ารหัสจริงๆ ต้องมีการเข้ารหัสในอนาคต)
	// บันทึกผู้ใช้ใหม่ลงใน repository
	return repository.SaveUser(model.User{Username: username, Password: password})
}
