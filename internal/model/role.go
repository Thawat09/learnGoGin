package model

import "time"

type Roles struct {
	RoleId    int       `gorm:"primaryKey;autoIncrement;column:RoleId"`
	RoleName  string    `gorm:"type:nvarchar(50);not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:CreatedAt;not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:UpdatedAt;not null"`
}

func (Roles) TableName() string {
	return "Roles"
}
