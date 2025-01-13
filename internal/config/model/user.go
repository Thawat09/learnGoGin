package model

import "time"

type Users struct {
	UserID       int       `gorm:"primaryKey;autoIncrement;column:UserId"`
	Username     string    `gorm:"type:nvarchar(50);not null;unique"`
	Password     string    `gorm:"type:nvarchar(256);not null"`
	Salt         string    `gorm:"type:nvarchar(256);not null"`
	Email        string    `gorm:"type:nvarchar(100);not null"`
	FirstName    string    `gorm:"column:FirstName;type:nvarchar(50);not null"`
	LastName     string    `gorm:"column:LastName;type:nvarchar(50);not null"`
	DepartmentID int       `gorm:"column:DepartmentId;not null"`
	Status       string    `gorm:"type:nvarchar(20);default:'Active';not null"`
	LastLogin    time.Time `gorm:"column:LastLogin;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime;column:CreatedAt;not null"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;column:UpdatedAt;not null"`

	Department Departments `gorm:"foreignKey:DepartmentID;references:DepartmentId"`
	UserRoles  []UserRoles `gorm:"foreignKey:UserId;references:UserID"`
}

func (Users) TableName() string {
	return "Users"
}
