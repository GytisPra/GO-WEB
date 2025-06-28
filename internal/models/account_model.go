package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Account struct {
	ID                    string   `gorm:"type:uuid;primaryKey;" json:"id"`
	UserID                string   `gorm:"type:uuid;not null;index" json:"user_id"`
	User                  User     `gorm:"foreignKey:UserID" json:"-"`
	Type                  string   `gorm:"not null" json:"type"`
	Provider              string   `gorm:"not null;uniqueIndex" json:"provider"`
	ProviderAccountId     *string  `gorm:"uniqueIndex" json:"provider_account_id"`
	RefreshToken          *string  `gorm:"" json:"refresh_token,omitempty"`
	AccessToken           *string  `gorm:"" json:"access_token,omitempty"`
	ExpiresAt             *float64 `gorm:"" json:"expires_at,omitempty"`
	TokenType             *string  `gorm:"" json:"token_type,omitempty"`
	Scope                 *string  `gorm:"" json:"scope,omitempty"`
	IdToken               *string  `gorm:"" json:"id_token,omitempty"`
	SessionState          *string  `gorm:"" json:"session_state,omitempty"`
	RefreshTokenExpiresIn *float64 `gorm:"" json:"refresh_token_Expires_in,omitempty"`

	// for local login only
	PasswordHash *string `gorm:"" json:"password_hash,omitempty"`
}

func (a *Account) Validate() error {
	if a.UserID == "" {
		return errors.New("userID cannot be empty")
	}
	if a.Type == "" {
		return errors.New("token type cannot be empty")
	}
	if a.Provider == "" {
		return errors.New("provider cannot be empty")
	}
	return nil
}

func CreateAccount(db *gorm.DB, account *Account) error {
	if err := account.Validate(); err != nil {
		return fmt.Errorf("new account validation failed: %w", err)
	}

	result := db.Where("user_id = ?", account.UserID).FirstOrCreate(account)
	if result.Error != nil {
		return fmt.Errorf("failed to create new account: %w", result.Error)
	}

	return nil
}

func GetLocalAccountByUserID(userID string, db *gorm.DB) (*Account, error) {
	var account Account

	result := db.Where("user_id = ? AND provider = ? AND password_hash IS NOT NULL", userID, "local").First(&account)

	if result.Error != nil {
		return nil, fmt.Errorf("error finding account by user ID: %w", result.Error)
	}

	return &account, nil
}
