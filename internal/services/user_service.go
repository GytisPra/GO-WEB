package services

import (
	"web-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(name string, email string) (*models.User, error) {
	var user *models.User

	result := s.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			user := &models.User{
				ID:    uuid.New().String(),
				Name:  &name,
				Email: &email,
			}
			err := models.CreateUser(s.db, user)
			if err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, result.Error
	}

	return user, nil
}

func (s *UserService) GetUserById(id string) (*models.User, error) {
	user, err := models.GetUserById(s.db, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
