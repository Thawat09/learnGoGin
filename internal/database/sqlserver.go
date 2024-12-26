package database // กำหนดชื่อแพ็กเกจเป็น "database" ซึ่งใช้ในการเชื่อมต่อกับฐานข้อมูล SQL Server

import (
	"fmt" // นำเข้าแพ็กเกจ fmt สำหรับการจัดการกับข้อความ เช่น การสร้างสตริง

	"gorm.io/driver/sqlserver" // นำเข้าแพ็กเกจ sqlserver จาก GORM เพื่อเชื่อมต่อกับฐานข้อมูล SQL Server
	"gorm.io/gorm"             // นำเข้า GORM ซึ่งเป็น ORM ที่ใช้ในการเชื่อมต่อกับฐานข้อมูล
)

// กำหนดตัวแปร DB สำหรับเก็บการเชื่อมต่อฐานข้อมูล
var DB *gorm.DB

// ฟังก์ชัน ConnectSQLServer ใช้เชื่อมต่อฐานข้อมูล SQL Server
func ConnectSQLServer(host, port, user, password string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=testITD", user, password, host, port)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQL Server: %w", err)
	}

	DB = db // เก็บการเชื่อมต่อไว้ในตัวแปร DB
	return db, nil
}
