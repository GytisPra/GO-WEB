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
	result := db.Where("user_id = ?", session.UserID).FirstOrCreate(&session)

	if result.Error != nil {
		return fmt.Errorf("failed to create session: %w", result.Error)
	}

	log.Println("created session: ", session.ID, "for user: ", session.UserID)

	return nil
}

func CleanupExpiredSessions(db *gorm.DB) error {
	for {
		now := time.Now()
		if err := db.Where("expires < ?", now).Delete(&Session{}).Error; err != nil {
			return fmt.Errorf("error deleting expired sessions: %w", err)
		}
		time.Sleep(1 * time.Hour)
	}
}

func GetSessionByToken(sessionToken string, db *gorm.DB) (*Session, error) {
	var session Session

	if err := db.Where("session_token = ?", sessionToken).First(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func RemoveSession(id string, db *gorm.DB) error {
	var sesion Session
	if err := db.Where("id = ?", id).Delete(&sesion).Error; err != nil {
		return err
	}

	log.Println("Removed session: ", id)

	return nil
}
