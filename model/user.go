package model

import (
	"time"

	"gorm.io/gorm"
)

type UserType string

// Enum values for UserStatus
const (
	UserGenPod UserType = "gen_pod"
	UserSubPod UserType = "sub_pod"
	UserBrigad UserType = "brigad"
	UserQA     UserType = "qa"
)

type User struct {
	ID        string `gorm:"primaryKey" json:"id"`
	Username  string `json:"username" gorm:"unique"`
	Password  string `json:"password"`
	Role      string
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Task struct {
	ID          string  `gorm:"primaryKey" json:"id"`
	OrderID     int     `gorm:"autoIncrement"`
	Name        string  `json:"name"`
	Responsible *string `json:"responsible"`
	User        User    `gorm:"foreignKey:Responsible;references:ID"`
	IsDone      bool    `gorm:"default:false"`
	IsChecked   bool    `gorm:"default:false"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
