package models

import (
	"gorm.io/gorm"
)

type Meso struct {
	gorm.Model `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	User       *User  `json:"-"`
	UserUUID   string ``
	UserID     uint   `gorm:"index:idx_user_id"`
	UUID       string `gorm:"index:idx_meso_uuid,unique"`
	Name       string `gorm:"not null"`

	Weeks *[]Week `gorm:"type:jsonb;serializer:json" validate:"required"`
}

type Week struct {
	gorm.Model `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	MesoID     uint `gorm:"index:idx_meso_id"`

	Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday *Day `gorm:"foreignKey:WeekID;constraint:OnDelete:CASCADE" validate:"required"`
}

type Day struct {
	gorm.Model `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	WeekID     uint    `gorm:"index:idx_meso_id"`
	Lifts      *[]Lift `gorm:"foreignKey:DayID;constraint:OnDelete:CASCADE" validate:"required"`
}

type Lift struct {
	gorm.Model `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	DayID      uint `gorm:"index:idx_week_id"`

	// Info for a lift
	Exercise string  `json:"exercise" validate:"required"`
	Sets     int     `json:"sets"`
	Weight   float32 `json:"weight"`
	Reps     int     `json:"reps"`
	Pump     int     `json:"pump"`
	Soreness int     `json:"soreness"`
}
