package models

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:text;not null" json:"body"`
	Password  string    `gorm:"type:text;not null" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (t *User) Validate() error {
	if t.Username == "" {
		return errors.New("body cannot be empty")
	}

	if t.Password == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}

func CreateUser(db *gorm.DB, user *User) (int64, error) {
	if err := user.Validate(); err != nil {
		return 0, err
	}
	result := db.Create(&user)

	if result.Error != nil {
		log.Println("Failed to insert data: ", result.Error)
		return 0, result.Error
	}

	return user.ID, nil
}
