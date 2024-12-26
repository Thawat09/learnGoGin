package service // กำหนดชื่อแพ็กเกจเป็น "service" ซึ่งจะมีฟังก์ชันที่เกี่ยวข้องกับการทำงานหลัก เช่น การเข้าสู่ระบบ และการลงทะเบียน

import (
	"errors"                         // นำเข้าแพ็กเกจ errors สำหรับการสร้าง error messages
	"goGin/internal/auth/model"      // นำเข้าแพ็กเกจ model เพื่อใช้ประเภทข้อมูล Users
	"goGin/internal/auth/repository" // นำเข้าแพ็กเกจ repository ซึ่งจะเชื่อมต่อกับฐานข้อมูลเพื่อจัดการข้อมูลผู้ใช้
	"goGin/internal/database"
	"time"

	"golang.org/x/crypto/bcrypt" // นำเข้า bcrypt เพื่อใช้ในการตรวจสอบรหัสผ่านที่เข้ารหัสแล้ว
)

// ฟังก์ชัน Login ใช้ในการตรวจสอบการเข้าสู่ระบบโดยใช้ชื่อผู้ใช้และรหัสผ่าน
func Login(username, password string) error {
	// หา user ใน repository โดยใช้ชื่อผู้ใช้ที่ส่งมา
	user, err := repository.FindUserByUsername(username)
	if err != nil {
		return err // ถ้าหาผู้ใช้ไม่เจอจะคืนค่าผิดพลาด (error)
	}

	// ตรวจสอบ password ว่า match กับ hash ที่เก็บไว้หรือไม่
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) // ใช้ bcrypt ตรวจสอบรหัสผ่าน
	if err != nil {                                                              // ถ้ารหัสผ่านไม่ตรงกัน
		return errors.New("invalid credentials") // คืนค่าผิดพลาดว่า credentials ไม่ถูกต้อง
	}

	return nil // คืนค่าปกติถ้าทุกอย่างถูกต้อง
}

// ฟังก์ชัน Register ใช้ในการลงทะเบียนผู้ใช้ใหม่
func Register(username, password, email string, departmentId int) error {
	// ตรวจสอบว่า Username หรือ Email เป็นค่าว่างหรือไม่
	if username == "" || email == "" || password == "" {
		return errors.New("username, email, and password must not be empty")
	}

	// ตรวจสอบว่ามีผู้ใช้ในระบบนี้อยู่แล้วหรือไม่
	_, err := repository.FindUserByUsername(username)
	if err == nil {
		return errors.New("username already exists") // ถ้ามีผู้ใช้อยู่แล้ว จะไม่สามารถลงทะเบียนใหม่ได้
	} else if err.Error() != "user not found" {
		return err
	}

	// ตรวจสอบว่า DepartmentId ที่ส่งมามีอยู่ในตาราง Departments หรือไม่
	_, err = repository.FindDepartmentById(database.DB, departmentId)
	if err != nil {
		return errors.New("invalid department ID") // ถ้าไม่พบ DepartmentId ในตาราง Departments
	}

	// เข้ารหัส password ก่อนบันทึก
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err // หากเกิดข้อผิดพลาดในการเข้ารหัส
	}

	// กำหนดเวลาให้เป็น UTC
	now := time.Now().UTC() // ใช้ UTC เพื่อหลีกเลี่ยงปัญหากับ Timezone

	// บันทึกผู้ใช้ใหม่ลงในฐานข้อมูล
	return repository.SaveUser(model.Users{
		Username:     username,
		Password:     string(hashedPassword),
		Email:        email,        // บันทึกค่า Email ด้วย
		DepartmentID: departmentId, // กำหนด DepartmentId ที่ถูกต้อง
		CreatedAt:    now,          // ใช้เวลาใน UTC
		UpdatedAt:    now,          // ใช้เวลาใน UTC
		LastLogin:    now,          // กำหนดค่า LastLogin เป็นเวลาปัจจุบัน
	})
}
