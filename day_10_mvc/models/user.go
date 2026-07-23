package models

import "gorm.io/gorm"

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	UserName string `json:"username"`
	Password string `json:"-"`
}

func GetUserByUserName(db *gorm.DB, username string) (User, error) {
	var user User
	result := db.Where("user_name = ?", username).First(&user)
	return user, result.Error
}

func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}
