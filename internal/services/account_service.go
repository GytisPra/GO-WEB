package services

import (
	"fmt"
	"web-app/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{db: db}
}

func (s *AccountService) CreateAccount(userID string,
	authType string,
	provider string,
	providerAccountId *string,
	refreshToken *string,
	accessToken *string,
	expiresAt *float64,
	tokenType *string,
	scope *string,
	idToken *string,
	sessionState *string,
	refreshTokenExpiresIn *float64,
	password *string,
) (*models.Account, error) {
	var account *models.Account

	switch provider {
	case "local":
		if password == nil {
			return nil, fmt.Errorf("password must be provided for local accounts")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		hashedPasswordString := string(hashedPassword)

		account = &models.Account{
			ID:           uuid.New().String(),
			UserID:       userID,
			Type:         authType,
			Provider:     provider,
			PasswordHash: &hashedPasswordString,
		}

		err = models.CreateAccount(s.db, account)
		if err != nil {
			return nil, err
		}

	case "discord":
		account = &models.Account{
			ID:                    uuid.New().String(),
			UserID:                userID,
			Type:                  authType,
			Provider:              provider,
			ProviderAccountId:     providerAccountId,
			RefreshToken:          refreshToken,
			AccessToken:           accessToken,
			ExpiresAt:             expiresAt,
			TokenType:             tokenType,
			Scope:                 scope,
			IdToken:               idToken,
			SessionState:          sessionState,
			RefreshTokenExpiresIn: refreshTokenExpiresIn,
		}
		err := models.CreateAccount(s.db, account)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	return account, nil
}

func (s *AccountService) GetLocalAccountByUserID(ID string) (*models.Account, error) {
	account, err := models.GetLocalAccountByUserID(ID, s.db)
	if err != nil {
		return nil, err
	}

	return account, nil
}
