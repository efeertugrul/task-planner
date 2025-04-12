package model

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	ExternalID        string         `gorm:"uniqueIndex:idx_source_external_id" json:"external_id"`
	Name              *string        `json:"name"`
	Difficulty        float64        `json:"difficulty"`
	EstimatedDuration float64        `json:"estimated_duration"`
	Source            string         `gorm:"uniqueIndex:idx_source_external_id" json:"source"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Assignment        *Assignment    `gorm:"foreignKey:TaskID" json:"assignment,omitempty"`
}

type Developer struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `json:"name"`
	Productivity float64        `json:"productivity"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Assignments  []Assignment   `gorm:"foreignKey:DeveloperID" json:"assignments,omitempty"`
}

type Assignment struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	DeveloperID     uint           `json:"developer_id"`
	TaskID          uint           `json:"task_id"`
	WeekNumber      int            `json:"week_number"`
	CalculatedHours float64        `json:"calculated_hours"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Developer       Developer      `gorm:"foreignKey:DeveloperID" json:"developer,omitempty"`
	Task            Task           `gorm:"foreignKey:TaskID" json:"task,omitempty"`
}
