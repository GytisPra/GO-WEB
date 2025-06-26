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

func (s *TaskService) CreateTask(body string, userId string) (*models.Task, error) {
	task := &models.Task{
		ID:     uuid.New().String(),
		Body:   body,
		UserID: userId,
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

func (s *TaskService) GetUserTasks(userId string) ([]models.Task, error) {
	var tasks []models.Task

	tasks, err := models.GetUserTasks(userId, s.db)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) GetTaskById(taskId string) (*models.Task, error) {
	var task *models.Task

	task, err := models.GetTaskById(taskId, s.db)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) UpdateTask(userId string, taskId string, newTaskBody string) error {
	err := models.UpdateTask(userId, taskId, newTaskBody, s.db)
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskService) DeleteTask(ID string) error {
	err := models.DeleteTask(ID, s.db)
	if err != nil {
		return err
	}

	return nil
}
