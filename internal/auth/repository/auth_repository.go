package repository

import (
	"errors"
	"goGin/internal/auth/model"
	"goGin/internal/database"
	"strings"

	"gorm.io/gorm"
)

func FindUserByUsername(username string) (model.Users, error) {
	var user model.Users
	err := database.DB.Where("username = ?", username).First(&user).Error

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
