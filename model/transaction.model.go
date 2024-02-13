package model

import (
	"time"
)

type UserModel struct {
	Id      int    `gorm:"type:int;primary_key"`
	Name    string `gorm:"type:varchar(255)"`
	Limit   int    `gorm:"type:int"`
	Balance int    `gorm:"type:int"`
}

type Model struct {
	Id          int       `gorm:"type:int;primary_key"`
	Value       int       `gorm:"type:int"`
	Description string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time `gorm:"type:timestamp"`
	UserId      int       `gorm:"type:int"`
}
