package models

import (
	"gorm.io/gorm"
)

type User struct {
	// Basic User Info
	gorm.Model `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	UUID       string `gorm:"index:idx_user_uuid,unique"`
	Username   string `json:"username" gorm:"uniqueIndex;not null"`
	Password   string `json:"password" gorm:"not null"`

	// User tracking/exercise info
	Mesos []*Meso `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
