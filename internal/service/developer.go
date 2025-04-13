package service

import (
	"fmt"
	"todo-planning/internal/model"

	"gorm.io/gorm"
)

type DeveloperService struct {
	db *gorm.DB
}

func NewDeveloperService(db *gorm.DB) *DeveloperService {
	return &DeveloperService{db: db}
}

func (s *DeveloperService) GetDevelopers() ([]model.Developer, error) {
	var developers []model.Developer
	if err := s.db.Find(&developers).Error; err != nil {
		return nil, fmt.Errorf("failed to get developers: %w", err)
	}
	return developers, nil
}
