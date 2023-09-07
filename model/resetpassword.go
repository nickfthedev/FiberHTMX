package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User struct
type ResetPassword struct {
	gorm.Model
	UserUUID uuid.UUID `gorm:"type:uuid;default:null" json:"uuid" form:"UUID"`
	Key      uuid.UUID `gorm:"type:uuid;default:null" json:"key" form:"key"`
}
