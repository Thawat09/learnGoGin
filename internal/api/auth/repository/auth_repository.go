package repository

import (
	"errors"
	"goGin/internal/config/database"
	"goGin/internal/config/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

func FindUserByUsername(username string) (model.Users, error) {
	var user model.Users

	err := database.DB.
		Select("UserId", "username", "password", "salt", "email", "FirstName", "LastName", "DepartmentId").
		Where("username = ?", username).
		Where("status = ?", "Active").
		Preload("Department").
		Preload("UserRoles.Role").
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Users{}, errors.New("user not found")
		}
		return model.Users{}, err
	}

	return user, nil
}

func FindDepartmentById(db *gorm.DB, departmentId int) (*model.Departments, error) {
	var department model.Departments

	if err := db.First(&department, departmentId).Error; err != nil {
		return nil, err
	}
	return &department, nil
}

func SaveUser(user model.Users) error {
	if err := database.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "UNIQUE KEY") {
			return errors.New("duplicate username or email")
		}
		return err
	}
	return nil
}

func UpdateLastLogin(username string) error {
	location, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(location)

	result := database.DB.Exec("UPDATE Users SET LastLogin = ? WHERE username = ?", now, username)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found or no update occurred")
	}

	return nil
}

func LogLoginHistory(userId int, ipAddress string) error {
	location, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(location)

	result := database.DB.Exec("INSERT INTO LoginHistory (UserId, IPAddress, LoginTime) VALUES (?, ?, ?)", userId, ipAddress, now)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
