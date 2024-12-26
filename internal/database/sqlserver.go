package database // กำหนดชื่อแพ็กเกจเป็น "database" ซึ่งใช้ในการเชื่อมต่อกับฐานข้อมูล SQL Server

import (
	"fmt" // นำเข้าแพ็กเกจ fmt สำหรับการจัดการกับข้อความ เช่น การสร้างสตริง

	"gorm.io/driver/sqlserver" // นำเข้าแพ็กเกจ sqlserver จาก GORM เพื่อเชื่อมต่อกับฐานข้อมูล SQL Server
	"gorm.io/gorm"             // นำเข้า GORM ซึ่งเป็น ORM ที่ใช้ในการเชื่อมต่อกับฐานข้อมูล
)

// ฟังก์ชัน ConnectSQLServer เชื่อมต่อกับ SQL Server โดยรับพารามิเตอร์ host, port, user, และ password
func ConnectSQLServer(host, port, user, password string) (*gorm.DB, error) {
	// สร้าง Data Source Name (DSN) สำหรับการเชื่อมต่อกับ SQL Server โดยใช้ข้อมูล host, port, user, และ password
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=testITD", user, password, host, port)

	// เชื่อมต่อกับ SQL Server โดยใช้ GORM และ DSN ที่ได้สร้าง
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		// ถ้ามีข้อผิดพลาดในการเชื่อมต่อ จะคืนค่า error ที่บอกข้อผิดพลาดในการเชื่อมต่อ
		return nil, fmt.Errorf("failed to connect to SQL Server: %w", err)
	}

	// ตรวจสอบการเชื่อมต่อ SQL Server
	sqlDB, err := db.DB() // ใช้ db.DB() เพื่อดึงตัวเชื่อมต่อ SQL database
	if err != nil {
		// ถ้าดึงการเชื่อมต่อ SQL database ไม่สำเร็จ จะคืนค่า error
		return nil, fmt.Errorf("failed to retrieve SQL database instance: %w", err)
	}

	// ทดสอบการเชื่อมต่อกับ SQL Server โดยการ Ping
	if err := sqlDB.Ping(); err != nil {
		// ถ้าการ Ping ไม่สำเร็จ แสดงว่าไม่สามารถเชื่อมต่อกับฐานข้อมูลได้
		return nil, fmt.Errorf("failed to ping SQL Server: %w", err)
	}

	// ถ้าทุกอย่างทำงานถูกต้อง จะคืนค่า db ซึ่งเป็นการเชื่อมต่อกับฐานข้อมูล SQL Server
	return db, nil
}
