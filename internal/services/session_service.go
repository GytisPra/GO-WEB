package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"
	"web-app/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	appErrors "web-app/pkg/errors"
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

func (s *SessionService) CreateSession(userID string) (*models.Session, error) {
	token, err := generateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	log.Println("Generated session token:", token)

	session := &models.Session{
		ID:           uuid.New().String(),
		SessionToken: token,
		Expires:      time.Now().Add(24 * time.Hour),
		UserID:       userID,
	}
	if err := models.CreateSession(s.db, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) CleanupExpiredSessions() {
	go models.CleanupExpiredSessions(s.db)
}

func (s *SessionService) GetSessionByToken(sessionToken string) (*models.Session, error) {
	session, err := models.GetSessionByToken(sessionToken, s.db)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) LoginWithLocal(email, password string) (*models.Session, error) {
	user, err := models.GetUserByEmail(s.db, email)
	if err != nil {
		return nil, appErrors.ErrInvalidCredentials
	}

	account, err := models.GetLocalAccountByUserID(user.ID, s.db)
	if err != nil || account.PasswordHash == nil {
		return nil, appErrors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*account.PasswordHash), []byte(password)); err != nil {
		return nil, appErrors.ErrInvalidCredentials
	}

	return s.CreateSession(user.ID)
}

func (s *SessionService) IsUserLoggedIn(sessionToken string) (bool, *models.User, error) {
	session, err := s.GetSessionByToken(sessionToken)
	if err != nil {
		return false, nil, err
	}

	if session.Expires.Before(time.Now()) {
		return false, nil, nil
	}

	user, err := models.GetUserById(s.db, session.UserID)
	if err != nil {
		return false, nil, err
	}

	return true, user, nil
}

func (s *SessionService) LogUserOut(sessionToken string) error {
	session, err := models.GetSessionByToken(sessionToken, s.db)
	if err != nil {
		return err
	}

	if err = models.RemoveSession(session.ID, s.db); err != nil {
		return err
	}

	return nil
}

func generateSessionToken() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
