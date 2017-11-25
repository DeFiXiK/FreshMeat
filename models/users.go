package models

import "time"
import "fmt"
import "crypto/sha1"
import "github.com/jinzhu/gorm"

type User struct {
	ID           uint `gorm:"primary_key"`
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

func HashPwd(password string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}

func (u *User) Create(db *gorm.DB) error {
	err := db.Create(&u).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUserByName(db *gorm.DB, username string) (*User, error) {
	user := &User{}
	err := db.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	user := &User{}
	err := db.Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) CheckPassword(password string) bool {
	if u.PasswordHash != HashPwd(password) {
		return false
	}
	return true
}

func (u *User) Save(db *gorm.DB) error {
	err := db.Save(u).Error
	if err != nil {
		return err
	}
	return nil
}
