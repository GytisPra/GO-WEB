package models

import (
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
	result := db.Create(&user)

	if result.Error != nil {
		log.Println("Failed to insert data: ", result.Error)
		return result.Error
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
		return nil, err
	}

	log.Println("Got user by ID: ", user.ID)
	return &user, nil
}
