package model // กำหนด package เป็น "model" ซึ่งใช้สำหรับเก็บโครงสร้างข้อมูล (structs) ที่เกี่ยวข้องกับโมเดลในโปรเจค

// กำหนดโครงสร้างข้อมูล User ซึ่งใช้แทนข้อมูลผู้ใช้ในระบบ
type User struct {
	Username string `json:"username"` // กำหนดฟิลด์ Username สำหรับเก็บชื่อผู้ใช้ และใช้ JSON tag เพื่อให้สามารถแปลงเป็น JSON ได้ในกรณีที่ส่งผ่าน API
	Password string `json:"password"` // กำหนดฟิลด์ Password สำหรับเก็บรหัสผ่าน และใช้ JSON tag เพื่อให้สามารถแปลงเป็น JSON ได้เช่นเดียวกัน
}
