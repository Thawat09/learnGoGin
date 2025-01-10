package service

import (
	"errors"
	"goGin/internal/auth/repository"
	"goGin/internal/database"
	"goGin/internal/model"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
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

func Login(username, password string, ipAddress string) (*Claims, error) {
	user, err := repository.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	go func() {
		err := repository.UpdateLastLogin(username)
		if err != nil {
			log.Println("Failed to update last login:", err)
		}
	}()

	go func() {
		err := repository.LogLoginHistory(user.UserID, ipAddress)
		if err != nil {
			log.Println("Failed to log login history:", err)
		}
	}()

	var roleName string
	var roleId int
	if len(user.UserRoles) > 0 {
		roleName = user.UserRoles[0].Role.RoleName
		roleId = user.UserRoles[0].RoleId
	}
	departmentName := user.Department.DepartmentName

	data := &Claims{
		UserId:         user.UserID,
		Username:       user.Username,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		RoleId:         roleId,
		RoleName:       roleName,
		DepartmentId:   user.DepartmentID,
		DepartmentName: departmentName,
	}

	return data, nil
}

func Register(username, password, email string, departmentId int) error {
	if username == "" || email == "" || password == "" {
		return errors.New("username, email, and password must not be empty")
	}

	_, err := repository.FindUserByUsername(username)

	if err == nil {
		return errors.New("username already exists")
	} else if err.Error() != "user not found" {
		return err
	}

	_, err = repository.FindDepartmentById(database.DB, departmentId)

	if err != nil {
		return errors.New("invalid department ID")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	location, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(location)

	return repository.SaveUser(model.Users{
		Username:     username,
		Password:     string(hashedPassword),
		Email:        email,
		DepartmentID: departmentId,
		CreatedAt:    now,
		UpdatedAt:    now,
		LastLogin:    now,
	})
}

func CreateAccessToken(user *Claims) (string, error) {
	secretKey := os.Getenv("SECRETTOKENKEY")

	if secretKey == "" {
		return "", errors.New("missing SECRETTOKENKEY in environment variables")
	}

	claims := jwt.MapClaims{
		"UserId":         user.UserId,
		"Username":       user.Username,
		"Email":          user.Email,
		"FirstName":      user.FirstName,
		"LastName":       user.LastName,
		"RoleId":         user.RoleId,
		"RoleName":       user.RoleName,
		"DepartmentId":   user.DepartmentId,
		"DepartmentName": user.DepartmentName,
		"Exp":            time.Now().Add(1 * time.Hour).Unix(), // time.Now().Add(1 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

func CreateRefreshToken(user *Claims) (string, error) {
	secretKey := os.Getenv("SECRETREFRESHTOKENKEY")

	if secretKey == "" {
		return "", errors.New("missing SECRETREFRESHTOKENKEY in environment variables")
	}

	claims := jwt.MapClaims{
		"UserId":         user.UserId,
		"Username":       user.Username,
		"Email":          user.Email,
		"FirstName":      user.FirstName,
		"LastName":       user.LastName,
		"RoleId":         user.RoleId,
		"RoleName":       user.RoleName,
		"DepartmentId":   user.DepartmentId,
		"DepartmentName": user.DepartmentName,
		"Exp":            time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}
