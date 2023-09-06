package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User struct
type User struct {
	gorm.Model
	UUID       uuid.UUID `gorm:"type:uuid;default:null" json:"uuid"`
	Name       string    `gorm:"not null; default:null" json:"name" form:"name" binding:"required"`
	Username   string    `gorm:"unique_index;not null;unique;default:null" json:"username" form:"username" binding:"required"`
	Email      string    `gorm:"unique_index;not null;unique;default:null" json:"email" form:"email" binding:"required"`
	Password   string    `gorm:"not null; default:null" json:"password" form:"password" binding:"required"`
	Verified   bool      `gorm:"default:false"`
	AvatarPath string    `json:"avatarpath" form:"avaterpath"`
}

// Like Model User but without Password Hash
type UserSafe struct {
	gorm.Model
	UUID       uuid.UUID `gorm:"type:uuid;default:null" json:"uuid"`
	Name       string    `gorm:"not null; default:null" json:"name" form:"name" binding:"required"`
	Username   string    `gorm:"unique_index;not null;unique;default:null" json:"username" form:"username" binding:"required"`
	Email      string    `gorm:"unique_index;not null;unique;default:null" json:"email" form:"email" binding:"required"`
	Verified   bool      `gorm:"default:false"`
	AvatarPath string    `json:"avatarpath" form:"avaterpath"`
}

type UserChangePassword struct {
	ID              int
	OldPassword     string `json:"oldpassword" form:"oldpassword"`
	NewPassword     string `json:"newpassword" form:"newpassword"`
	ConfirmPassword string `json:"confirmpassword" form:"confirmpassword"`
}

// Input for Login
type UserLoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Create manually users
type CreateUser struct {
	Name string `gorm:"not null; default:null" json:"name" form:"name" binding:"required"`
	//Username string `gorm:"unique_index;not null;unique" json:"username" form:"username" binding:"required"`
	Email    string `gorm:"unique_index;not null;unique" json:"email" form:"email" binding:"required"`
	Password string `gorm:"not null" json:"password" binding:"required" form:"password"`
}

// Create fuser from register form
type RegisterUser struct {
	Name string `gorm:"not null; default:null" json:"name" form:"name" binding:"required"`
	//Username        string `gorm:"unique_index;not null;unique" json:"username" form:"username" binding:"required"`
	Email           string `gorm:"unique_index;not null;unique" json:"email" form:"email" binding:"required"`
	Password        string `gorm:"not null" json:"password" binding:"required" form:"password"`
	ConfirmPassword string `json:"confirmpassword" binding:"required" form:"confirmpassword"`
}

// ///
// /// Not used for now
// ///
type UserAPI struct {
	ID         int
	Email      string `json:"email"`
	Username   string `json:"username"`
	AvatarPath string `json:"avatarpath"`
}

type UserUpdate struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	AvatarPath string `json:"avatarpath"`
}
