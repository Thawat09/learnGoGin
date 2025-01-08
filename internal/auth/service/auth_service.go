package service // กำหนดชื่อแพ็กเกจเป็น "service" ซึ่งจะมีฟังก์ชันที่เกี่ยวข้องกับการทำงานหลัก เช่น การเข้าสู่ระบบ และการลงทะเบียน

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"errors"                         // นำเข้าแพ็กเกจ errors สำหรับการสร้าง error messages
	"goGin/internal/auth/model"      // นำเข้าแพ็กเกจ model เพื่อใช้ประเภทข้อมูล Users
	"goGin/internal/auth/repository" // นำเข้าแพ็กเกจ repository ซึ่งจะเชื่อมต่อกับฐานข้อมูลเพื่อจัดการข้อมูลผู้ใช้
	"goGin/internal/database"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt" // นำเข้า bcrypt เพื่อใช้ในการตรวจสอบรหัสผ่านที่เข้ารหัสแล้ว
)

type Claims struct {
	UserId         int    `json:"userId"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	RoleId         int    `json:"roleId"`
	RoleName       string `json:"roleName"`
	DepartmentId   int    `json:"departmentId"`
	DepartmentName string `json:"departmentName"`
	jwt.RegisteredClaims
}

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

// ฟังก์ชันแปลงข้อความด้วย AES โดยใช้ SECRETKEYDATA จาก environment variables
func Encrypt(message string) (string, error) {
	// ดึงค่า SECRETKEYDATA จาก environment variables
	secretKey := os.Getenv("SECRETKEYDATA")
	if secretKey == "" {
		return "", errors.New("missing SECRETKEYDATA in environment variables")
	}

	// ใช้ SHA256 เพื่อให้แน่ใจว่า secretKey มีขนาดที่เหมาะสม
	hash := sha256.Sum256([]byte(secretKey))

	// สร้าง block cipher ด้วย AES
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", err
	}

	// สร้าง IV (Initialization Vector) สำหรับการเข้ารหัส
	iv := []byte("123456789012") // คุณอาจใช้ค่า IV แบบคงที่หรือสุ่มขึ้นมา (ควรสุ่มในกรณีใช้งานจริง)

	// สร้าง GCM (Galois/Counter Mode) สำหรับ AES
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// เข้ารหัสข้อความ
	ciphertext := gcm.Seal(nil, iv, []byte(message), nil)

	// เข้ารหัสข้อความเป็น base64 เพื่อให้ง่ายต่อการส่งใน JSON
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encoded, nil
}

func Decrypt(encodedMessage string) (string, error) {
	secretKey := os.Getenv("SECRETKEYDATA")
	if secretKey == "" {
		return "", errors.New("missing SECRETKEYDATA in environment variables")
	}

	hash := sha256.Sum256([]byte(secretKey))

	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Decode base64-encoded message
	ciphertext, err := base64.StdEncoding.DecodeString(encodedMessage)
	if err != nil {
		return "", err
	}

	// Use the same IV (must match the one used during encryption)
	iv := []byte("123456789012")
	if len(ciphertext) < len(iv) {
		return "", errors.New("ciphertext too short")
	}

	// Decrypt the message
	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func CreateAccessToken(user *Claims) (string, error) {
	// ดึงค่า Secret Key สำหรับ Access Token จาก environment variables
	secretKey := os.Getenv("SECRETTOKENKEY")
	if secretKey == "" {
		return "", errors.New("missing SECRETTOKENKEY in environment variables")
	}

	// สร้าง token พร้อม claims ที่ต้องการ
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserId:         user.UserId,
		Username:       user.Username,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		RoleId:         user.RoleId,
		RoleName:       user.RoleName,
		DepartmentId:   user.DepartmentId,
		DepartmentName: user.DepartmentName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)), // กำหนดเวลาหมดอายุ 30 นาที
			IssuedAt:  jwt.NewNumericDate(time.Now()),                       // เวลาที่ token ถูกสร้าง
		},
	})

	// เซ็น token ด้วย Secret Key
	return token.SignedString([]byte(secretKey))
}

func CreateRefreshToken(user *Claims) (string, error) {
	// ดึงค่า Secret Key สำหรับ Refresh Token จาก environment variables
	secretKey := os.Getenv("SECRETREFRESHTOKENKEY")
	if secretKey == "" {
		return "", errors.New("missing SECRETREFRESHTOKENKEY in environment variables")
	}

	// สร้าง token พร้อม claims ที่ต้องการ
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.UserId,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(), // กำหนดเวลาหมดอายุ 7 วัน
	})

	// เซ็น token ด้วย Secret Key
	return token.SignedString([]byte(secretKey))
}
