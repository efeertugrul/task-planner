package service

import (
	"todo-planning/internal/model"

	"gorm.io/gorm"
)

type AssignmentService struct {
	db *gorm.DB
}

func NewAssignmentService(db *gorm.DB) *AssignmentService {
	return &AssignmentService{db: db}
}

func (s *AssignmentService) CreateAssignment(assignment *model.Assignment) error {
	return s.db.Create(assignment).Error
}

func (s *AssignmentService) GetAssignments() ([]model.Assignment, error) {
	var assignments []model.Assignment
	return assignments, s.db.Find(&assignments).Error
}

func (s *AssignmentService) CreateBatchAssignments(assignments []model.Assignment) error {
	return s.db.CreateInBatches(&assignments, 100).Error
}
