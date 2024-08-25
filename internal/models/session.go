package models

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	SessionToken string    `gorm:"not null;uniqueIndex" json:"session_token"`
	Expires      time.Time `gorm:"not null" json:"expires"`
	UserID       string    `gorm:"type:uuid;not null;index" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"-"`
}

func CreateSession(db *gorm.DB, session *Session) error {
	result := db.Create(&session)

	if result.Error != nil {
		log.Println("Failed to insert data: ", result.Error)
		return result.Error
	}

	return nil
}

func CleanupExpiredSessions(db *gorm.DB) {
	for {
		now := time.Now()
		if err := db.Where("expires < ?", now).Delete(&Session{}).Error; err != nil {
			fmt.Println("Error deleting expired sessions:", err)
		}
		time.Sleep(1 * time.Hour)
	}
}
