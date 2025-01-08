package model

import "time"

type Departments struct {
	DepartmentId   int       `gorm:"primaryKey;autoIncrement;column:DepartmentId"`
	DepartmentName string    `gorm:"type:nvarchar(100);not null;unique"`
	Description    string    `gorm:"type:nvarchar(256);not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime;column:CreatedAt;not null"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime;column:UpdatedAt;not null"`
}

func (Departments) TableName() string {
	return "Departments"
}
