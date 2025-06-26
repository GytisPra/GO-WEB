package models

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID        string    `gorm:"type:uuid;primaryKey;" json:"id"`
	Body      string    `gorm:"not null" json:"body"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (t *Task) Validate() error {
	if t.Body == "" {
		return errors.New("body cannot be empty")
	}
	return nil
}

func CreateTask(db *gorm.DB, task *Task) error {
	if err := task.Validate(); err != nil {
		return err
	}
	result := db.Create(&task)

	if result.Error != nil {
		log.Println("Failed to insert data: ", result.Error)
		return result.Error
	}

	log.Println("Created new task: ", task.ID, "By user: ", task.UserID)
	return nil
}

func GetAllTasks(db *gorm.DB) ([]Task, error) {
	var tasks []Task

	result := db.Find(&tasks)
	if result.Error != nil {
		log.Println("Failed to fetch tasks: ", result.Error)
		return nil, result.Error
	}

	log.Println("Fetched all tasks")
	return tasks, nil
}

func GetUserTasks(userId string, db *gorm.DB) ([]Task, error) {
	var tasks []Task

	result := db.Where("user_id = ?", userId).Find(&tasks)
	if result.Error != nil {
		log.Println("Failed to fetch tasks: ", result.Error)
		return nil, result.Error
	}

	log.Println("Fetched all tasks for user: ", userId)
	return tasks, nil
}
