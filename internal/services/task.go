package services

import (
	"web-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskService struct {
	db *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{db: db}
}

func (s *TaskService) CreateTask(body string) (*models.Task, error) {
	task := &models.Task{
		ID:   uuid.New().String(),
		Body: body,
		//UserID: "none",
	}
	err := models.CreateTask(s.db, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task

	tasks, err := models.GetAllTasks(s.db)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
