package models

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Body      string    `gorm:"type:text;not null" json:"body"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (t *Task) Validate() error {
	if t.Body == "" {
		return errors.New("body cannot be empty")
	}
	return nil
}

func CreateTask(db *gorm.DB, task *Task) (int64, error) {
	if err := task.Validate(); err != nil {
		return 0, err
	}
	result := db.Create(&task)

	if result.Error != nil {
		log.Println("Failed to insert data: ", result.Error)
		return 0, result.Error
	}

	return task.ID, nil
}

func GetAllTasks(db *gorm.DB) ([]Task, error) {
	var tasks []Task

	result := db.Find(&tasks)
	if result.Error != nil {
		log.Println("Failed to fetch tasks: ", result.Error)
		return nil, result.Error
	}

	return tasks, nil
}
