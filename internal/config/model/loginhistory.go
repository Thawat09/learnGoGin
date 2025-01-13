package model

import "time"

type LoginHistory struct {
	LoginId   int       `gorm:"primaryKey;autoIncrement;column:LoginId"`
	UserId    int       `gorm:"column:UserId;not null"`
	LoginTime time.Time `gorm:"column:LoginTime;not null"`
	IPAddress string    `gorm:"column:IPAddress;type:nvarchar(50);not null"`
}

func (LoginHistory) TableName() string {
	return "LoginHistory"
}
