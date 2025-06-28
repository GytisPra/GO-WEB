package models

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            string     `gorm:"type:uuid;primaryKey;" json:"id"`
	Name          *string    `gorm:"uniqueIndex" json:"name"`
	Email         *string    `gorm:"uniqueIndex" json:"email"`
	EmailVerified *time.Time `gorm:"uniqueIndex" json:"email_verified"`
	Image         *string    `gorm:"uniqueIndex" json:"image"`
	Accounts      []Account  `gorm:"foreignKey:UserID;" json:"accounts"`
	Sessions      []Session  `gorm:"foreignKey:UserID;" json:"sessions"`
	Tasks         []Task     `gorm:"foreignKey:UserID;" json:"tasks"`
}

func CreateUser(db *gorm.DB, user *User) error {
	result := db.Where("id = ? AND email = ?", user.ID, user.Email).FirstOrCreate(user)

	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}

	log.Println("Created new user: ", user.ID)
	return nil
}

func GetUserById(db *gorm.DB, id string) (*User, error) {
	var user User

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	log.Println("Got user by ID: ", user.ID)
	return &user, nil
}

func DeleteUserByID(db *gorm.DB, id string) error {
	if err := db.Where("id = ?", id).Delete(&User{}).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	log.Println("Deleted user: ", id)
	return nil
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User

	result := db.Where("email = ?", email).Find(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", result.Error)
	}

	return &user, nil
}
