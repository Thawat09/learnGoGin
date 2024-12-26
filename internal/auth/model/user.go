package model

import "time"

// โครงสร้างข้อมูล Users ซึ่งใช้แทนข้อมูลในฐานข้อมูล
type Users struct {
	UserID       int       `gorm:"primaryKey;autoIncrement;column:UserId"`      // กำหนดให้ UserId เป็น primary key
	Username     string    `gorm:"type:nvarchar(50);not null;unique"`           // ชื่อผู้ใช้ต้องไม่ซ้ำกัน
	Password     string    `gorm:"type:nvarchar(256);not null"`                 // รหัสผ่าน
	Salt         string    `gorm:"type:nvarchar(256);not null"`                 // Salt สำหรับการเข้ารหัสรหัสผ่าน
	Email        string    `gorm:"type:nvarchar(100);not null"`                 // อีเมล
	FirstName    string    `gorm:"column:FirstName;type:nvarchar(50);not null"` // ชื่อจริง
	LastName     string    `gorm:"column:LastName;type:nvarchar(50);not null"`  // นามสกุล
	DepartmentID int       `gorm:"column:DepartmentId;not null"`                // รหัสแผนก
	Status       string    `gorm:"type:nvarchar(20);default:'Active';not null"` // สถานะของผู้ใช้
	LastLogin    time.Time `gorm:"column:LastLogin;not null"`                   // เวลาที่ผู้ใช้ล็อกอินครั้งล่าสุด
	CreatedAt    time.Time `gorm:"autoCreateTime;column:CreatedAt;not null"`    // วันที่สร้าง
	UpdatedAt    time.Time `gorm:"autoUpdateTime;column:UpdatedAt;not null"`    // วันที่อัพเดท
}

func (Users) TableName() string {
	return "Users" // กำหนดชื่อตารางในฐานข้อมูล
}
