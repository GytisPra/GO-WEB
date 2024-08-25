package services

import (
	"web-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AcountService struct {
	db *gorm.DB
}

func NewAcountService(db *gorm.DB) *AcountService {
	return &AcountService{db: db}
}

func (s *AcountService) CreateAccount(userID string,
	authType string,
	provider string,
	providerAccountId string,
	refreshToken *string,
	accessToken *string,
	expiresAt *float64,
	tokenType *string,
	scope *string,
	idToken *string,
	sessionState *string,
	refreshTokenExpiresIn *float64) (*models.Account, error) {
	var account *models.Account

	result := s.db.Where("user_id = ?", userID).First(&account)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			account := &models.Account{
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
			return account, nil
		}
		return nil, result.Error
	}

	return account, nil
}
