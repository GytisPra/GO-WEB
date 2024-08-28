package services

import (
	"time"
	"web-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

func (s *SessionService) CreateSession(userID string) (*models.Session, error) {
	var session *models.Session

	result := s.db.Where("user_id = ?", userID).First(&session)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			session := &models.Session{
				ID:           uuid.New().String(),
				SessionToken: uuid.New().String(),
				Expires:      time.Now().Add(24 * time.Hour),
				UserID:       userID,
			}
			err := models.CreateSession(s.db, session)
			if err != nil {
				return nil, err
			}
			return session, nil
		}
		return nil, result.Error
	}

	return session, nil
}

func (s *SessionService) CleanupExpiredSessions() {
	go models.CleanupExpiredSessions(s.db)
}

func (s *SessionService) IsUserLoggedIn(sessionToken string) (bool, *models.User, error) {
	session, err := models.GetSessionByToken(sessionToken, s.db)
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
