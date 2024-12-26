package model

import "time"

// โครงสร้างข้อมูล Departments
type Departments struct {
	DepartmentId   int       `gorm:"primaryKey;autoIncrement;column:DepartmentId"` // กำหนดให้ DepartmentId เป็น primary key
	DepartmentName string    `gorm:"type:nvarchar(100);not null;unique"`           // ชื่อแผนกต้องไม่ซ้ำกัน
	Description    string    `gorm:"type:nvarchar(256);not null"`                  // รายละเอียด
	CreatedAt      time.Time `gorm:"autoCreateTime;column:CreatedAt;not null"`     // วันที่สร้าง
	UpdatedAt      time.Time `gorm:"autoUpdateTime;column:UpdatedAt;not null"`     // วันที่อัพเดท
}

func (Departments) TableName() string {
	return "Departments" // กำหนดชื่อตารางในฐานข้อมูล
}
