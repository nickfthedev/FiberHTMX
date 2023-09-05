package model

import "gorm.io/gorm"

// User struct
type User struct {
	gorm.Model
	Name       string `gorm:"not null; default:null" json:"name" form:"name" binding:"required"`
	Username   string `gorm:"unique_index;not null;unique;default:null" json:"username" form:"username" binding:"required"`
	Email      string `gorm:"unique_index;not null;unique;default:null" json:"email" form:"email" binding:"required"`
	Password   string `gorm:"not null; default:null" json:"password" form:"password" binding:"required"`
	Verified   bool   `gorm:"default:false"`
	Fullname   string `json:"fullname" form:"fullname"`
	AvatarPath string `json:"avatarpath" form:"avaterpath"`
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
	Fullname string `json:"fullname" binding:"required" form:"fullname"`
}

// Create fuser from register form
type RegisterUser struct {
	Name string `gorm:"not null; default:null" json:"name" form:"name" binding:"required"`
	//Username        string `gorm:"unique_index;not null;unique" json:"username" form:"username" binding:"required"`
	Email           string `gorm:"unique_index;not null;unique" json:"email" form:"email" binding:"required"`
	Password        string `gorm:"not null" json:"password" binding:"required" form:"password"`
	ConfirmPassword string `json:"confirmpassword" binding:"required" form:"confirmpassword"`
	Fullname        string `json:"fullname" binding:"required" form:"fullname"`
}

/////
///// Not used for now
/////
type UserAPI struct {
	ID         int
	Email      string `json:"email"`
	Username   string `json:"username"`
	Fullname   string `json:"fullname"`
	AvatarPath string `json:"avatarpath"`
}

type UserChangePassword struct {
	ID             int
	OldPassword    string `json:"oldpassword"`
	NewPassword    string `json:"newpassword"`
	RepeatPassword string `json:"repeatpassword"`
}

type UserUpdate struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	Fullname   string `json:"fullname"`
	AvatarPath string `json:"avatarpath"`
}
