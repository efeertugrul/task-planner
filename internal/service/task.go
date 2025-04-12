package service

import (
	"fmt"

	"todo-planning/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TaskService struct {
	db *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{
		db: db,
	}
}

// StoreTasks stores tasks in the database using ON CONFLICT DO NOTHING
func (s *TaskService) StoreTasks(tasks []model.Task) error {
	if len(tasks) == 0 {
		return nil
	}

	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "source"}, {Name: "external_id"}},
		DoNothing: true,
	}).Create(&tasks).Error
}

// GetTasks returns all tasks from the database
func (s *TaskService) GetTasks() ([]model.Task, error) {
	var tasks []model.Task
	if err := s.db.Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}
