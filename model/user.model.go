package model

import "gorm.io/gorm"

// User struct
type User struct {
	gorm.Model
	Username   string `gorm:"unique_index;not null;unique;default:null" json:"username" form:"username" binding:"required"`
	Email      string `gorm:"unique_index;not null;unique;default:null" json:"email" form:"email" binding:"required"`
	Password   string `gorm:"not null; default:null" json:"password" form:"password" binding:"required"`
	Fullname   string `json:"fullname" form:"fullname"`
	AvatarPath string `json:"avatarpath" form:"avaterpath"`
}

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

type CreateUser struct {
	Username string `gorm:"unique_index;not null;unique" json:"username" form:"username" binding:"required"`
	Email    string `gorm:"unique_index;not null;unique" json:"email" form:"email" binding:"required"`
	Password string `gorm:"not null" json:"password" binding:"required" form:"password"`
	Fullname string `json:"fullname" binding:"required" form:"fullname"`
}
type RegisterUser struct {
	Username        string `gorm:"unique_index;not null;unique" json:"username" form:"username" binding:"required"`
	Email           string `gorm:"unique_index;not null;unique" json:"email" form:"email" binding:"required"`
	Password        string `gorm:"not null" json:"password" binding:"required" form:"password"`
	ConfirmPassword string `json:"confirmpassword" binding:"required" form:"confirmpassword"`
	Fullname        string `json:"fullname" binding:"required" form:"fullname"`
}
type UserLoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
