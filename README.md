# learnGoGin

Learn Golang Gin Framework # คำอธิบายเกี่ยวกับโปรเจกต์นี้ ซึ่งเป็นการเรียนรู้การใช้ Gin Framework บนภาษา Go

# Structure Folder Project # โครงสร้างของโปรเจกต์

goGin/ # โฟลเดอร์หลักของโปรเจกต์
│
├── cmd/ # โฟลเดอร์ที่เก็บโค้ดที่ใช้ในการรันแอปพลิเคชัน
│ └── server/ # โฟลเดอร์ที่เก็บไฟล์การตั้งค่าและการรัน server
│ └── main.go # จุดเริ่มต้นของโปรเจกต์ (entry point) ซึ่งจะทำการเริ่มการทำงานของเซิร์ฟเวอร์
│
├── internal/ # โฟลเดอร์ที่เก็บโค้ดของแอปพลิเคชันทั้งหมด แบ่งตามฟังก์ชัน
│ ├── auth/ # แพ็กเกจที่เกี่ยวข้องกับการจัดการการยืนยันตัวตน (Authentication)
│ │ ├── handler/ # โฟลเดอร์สำหรับจัดการกับ HTTP request และ response สำหรับการยืนยันตัวตน
│ │ │ └── handler.go # ไฟล์ที่จัดการกับ request สำหรับการ login หรือ register
│ │ ├── service/ # แพ็กเกจที่เก็บ logic ต่างๆ เช่น การตรวจสอบการ login
│ │ │ └── auth_service.go # ฟังก์ชันที่เกี่ยวข้องกับ logic การยืนยันตัวตน เช่น การ login, register
│ │ ├── repository/ # แพ็กเกจที่จัดการข้อมูลในฐานข้อมูล
│ │ │ └── auth_repository.go # จัดการข้อมูลของผู้ใช้ที่เกี่ยวข้องกับการยืนยันตัวตน
│ │ └── model/ # แพ็กเกจที่เก็บโครงสร้างของข้อมูล (model)
│ │ └── user.go # โมเดลของผู้ใช้ที่เกี่ยวข้องกับการยืนยันตัวตน
│ │
│ ├── static/ # แพ็กเกจที่จัดการข้อมูลของผู้ใช้
│ │ ├── handler/ # โฟลเดอร์สำหรับจัดการกับ HTTP request และ response สำหรับข้อมูลผู้ใช้
│ │ │ └── handler.go # ไฟล์ที่จัดการกับ request สำหรับการดึงข้อมูลหรืออัปเดตข้อมูลผู้ใช้
│ │ ├── service/ # แพ็กเกจที่เก็บ logic ต่างๆ เช่น การดึงข้อมูลผู้ใช้
│ │ │ └── static_service.go # ฟังก์ชันที่เกี่ยวข้องกับ logic การจัดการข้อมูลผู้ใช้
│ │ ├── repository/ # แพ็กเกจที่จัดการข้อมูลในฐานข้อมูล
│ │ │ └── static_repository.go # จัดการข้อมูลของผู้ใช้ในฐานข้อมูล
│ │ └── model/ # แพ็กเกจที่เก็บโครงสร้างของข้อมูล (model)
│ │
│ └── middleware/ # แพ็กเกจที่เก็บ middleware ที่ใช้ในการจัดการการทำงานของเซิร์ฟเวอร์
│ ├── auth_middleware.go # Middleware สำหรับตรวจสอบการเข้าถึง API เช่น การตรวจสอบ token
│ └── logging_middleware.go # Middleware สำหรับ log ข้อมูล request และ response
│
├── go.mod # ไฟล์ที่เก็บข้อมูลเกี่ยวกับ dependency ของโปรเจกต์
└── go.sum # ไฟล์ที่เก็บข้อมูล checksum ของ dependencies ที่ใช้

# Install # คำแนะนำในการติดตั้ง dependencies ของโปรเจกต์

1. go mod init goGin # สร้างไฟล์ go.mod ใหม่ เพื่อระบุชื่อของโปรเจกต์
2. go get -u github.com/gin-gonic/gin # ติดตั้ง Gin framework ที่ใช้ในการพัฒนาเว็บแอปพลิเคชัน
3. go get gorm.io/gorm # ติดตั้ง GORM ORM เพื่อใช้งานกับฐานข้อมูล
4. go get gorm.io/driver/sqlserver # ติดตั้ง driver สำหรับการเชื่อมต่อกับฐานข้อมูล SQL Server
5. go get github.com/denisenkom/go-mssqldb # ติดตั้ง driver สำหรับการเชื่อมต่อ SQL Server จาก Go
6. go get github.com/redis/go-redis/v9 # ติดตั้ง Redis client ที่ใช้ในเวอร์ชัน 8
7. go get github.com/joho/godotenv # ติดตั้งไลบรารีสำหรับโหลดค่าจากไฟล์ .env
8. go get -u golang.org/x/crypto/bcrypt # ติดตั้งไลบรารีสำหรับตรวจสอบว่ารหัสผ่าน
9. get github.com/ulule/limiter/v3 # ติดตั้ง rate limit
10. go get github.com/ulule/limiter/v3/drivers/store/memory # ติดตั้ง rate limit
11. go get github.com/gin-contrib/cors # ติดตั้ง cors
12. go get github.com/gin-contrib/gzip

# Delete Dependencies # วิธีลบ dependencies ที่ไม่ใช้งานออกจาก go.mod
1. go mod tidy # ลบ dependencies ที่ไม่ได้ใช้ออกจากโปรเจกต์

# Run # วิธีการรันโปรเจกต์
1. go run cmd/server/main.go # รันเซิร์ฟเวอร์โดยใช้ไฟล์ main.go ซึ่งเป็นจุดเริ่มต้นของโปรเจกต์
