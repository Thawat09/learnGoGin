package model

type UserRoles struct {
	UserId int   `gorm:"primaryKey;autoIncrement;column:UserId"`
	RoleId int   `gorm:"column:RoleId;not null"`
	Role   Roles `gorm:"foreignKey:RoleId;references:RoleId"`
}

func (UserRoles) TableName() string {
	return "UserRoles"
}
