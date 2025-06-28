package services

import (
	"errors"
	"web-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"

	appErrors "web-app/pkg/errors"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(name string, email string) (*models.User, error) {
	user := &models.User{
		ID:    uuid.New().String(),
		Name:  &name,
		Email: &email,
	}

	err := models.CreateUser(s.db, user)
	if err != nil {
		if errors.Is(err, appErrors.ErrUserAlreadyExits) {
			return nil, appErrors.ErrUserAlreadyExits
		}
		return nil, err
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

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := models.GetUserByEmail(s.db, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteUserByID(id string) error {
	err := models.DeleteUserByID(s.db, id)
	if err != nil {
		return err
	}

	return nil
}
