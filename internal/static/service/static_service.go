package service // กำหนดชื่อแพ็กเกจเป็น service ซึ่งใช้สำหรับการจัดการฟังก์ชันหลักที่เกี่ยวกับบริการหรือการประมวลผลข้อมูล

import "goGin/internal/static/repository" // นำเข้าแพ็กเกจ repository ที่ใช้สำหรับการทำงานกับข้อมูล เช่น การดึงข้อมูลผู้ใช้จากฐานข้อมูล

// ฟังก์ชัน GetUserByID ใช้เพื่อดึงข้อมูลผู้ใช้ตาม ID
func GetUserByID(id string) (repository.User, error) {
	return repository.FindUserByID(id) // เรียกฟังก์ชัน FindUserByID จาก repository เพื่อตรวจสอบและดึงข้อมูลผู้ใช้ตาม ID ที่ส่งเข้ามา
}
