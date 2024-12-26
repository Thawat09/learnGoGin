package database // กำหนดชื่อแพ็กเกจเป็น "database" ซึ่งจะใช้ในการเชื่อมต่อกับฐานข้อมูลต่าง ๆ เช่น Redis

import (
	"context" // นำเข้าแพ็กเกจ context สำหรับการจัดการบริบทการทำงานในโค้ดที่รองรับการทำงานแบบ concurrent
	"fmt"     // นำเข้าแพ็กเกจ fmt สำหรับการจัดการกับข้อความ เช่น การสร้างสตริง
	"log"     // นำเข้าแพ็กเกจ log สำหรับการบันทึกข้อความต่าง ๆ เช่น ข้อความ error

	"github.com/go-redis/redis/v8" // นำเข้าแพ็กเกจ redis สำหรับการทำงานกับ Redis
)

var ctx = context.Background() // กำหนด context ทั่วไปที่ใช้สำหรับการทำงานกับ Redis

// ฟังก์ชัน ConnectRedis เชื่อมต่อกับ Redis โดยใช้ host และ port
func ConnectRedis(host, port string) (*redis.Client, error) {
	// สร้าง client สำหรับเชื่อมต่อกับ Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port), // กำหนดที่อยู่ของ Redis server จาก host และ port ที่ได้รับ
	})

	// ตรวจสอบการเชื่อมต่อกับ Redis
	_, err := rdb.Ping(ctx).Result() // ใช้คำสั่ง Ping เพื่อตรวจสอบการเชื่อมต่อ
	if err != nil {
		// ถ้ามีข้อผิดพลาดในการเชื่อมต่อ Redis
		log.Printf("Redis connection error: %v", err) // บันทึกข้อความ error ลงใน log
		return nil, err                               // คืนค่า error
	}
	log.Println("Redis connected successfully") // ถ้าการเชื่อมต่อสำเร็จ ให้แสดงข้อความว่าเชื่อมต่อ Redis สำเร็จ
	return rdb, nil                             // คืนค่า client ของ Redis
}
