package repository // กำหนดชื่อแพ็กเกจเป็น repository ซึ่งทำหน้าที่ในการจัดการข้อมูล (ในที่นี้คือข้อมูลผู้ใช้)

import "errors" // นำเข้าแพ็กเกจ errors เพื่อใช้ในการสร้างข้อผิดพลาดที่กำหนดเอง

// จำลองข้อมูลผู้ใช้ในรูปแบบของ slice (เหมือนกับฐานข้อมูลในหน่วยความจำ)
var users = []User{
	{ID: "1", Username: "john", Email: "john@example.com"}, // ข้อมูลของผู้ใช้คนแรก
	{ID: "2", Username: "jane", Email: "jane@example.com"}, // ข้อมูลของผู้ใช้คนที่สอง
}

// กำหนดโครงสร้างข้อมูลของผู้ใช้
type User struct {
	ID       string `json:"id"`       // รหัสผู้ใช้
	Username string `json:"username"` // ชื่อผู้ใช้
	Email    string `json:"email"`    // อีเมลของผู้ใช้
}

// ฟังก์ชัน FindUserByID ใช้ในการค้นหาผู้ใช้ตาม ID ที่ส่งมา
func FindUserByID(id string) (User, error) {
	// ทำการวนลูปค้นหาผู้ใช้จาก slice "users"
	for _, user := range users {
		// ถ้าพบ ID ที่ตรงกับค่าที่ค้นหา
		if user.ID == id {
			return user, nil // ส่งกลับข้อมูลผู้ใช้พร้อมกับค่าผลลัพธ์ error เป็น nil
		}
	}
	// หากไม่พบผู้ใช้ที่ตรงกับ ID ที่ค้นหา ส่งค่าผลลัพธ์เป็น User ว่างๆ และ error ว่า "user not found"
	return User{}, errors.New("user not found")
}
