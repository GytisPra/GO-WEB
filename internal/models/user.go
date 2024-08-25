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

	return nil
}
