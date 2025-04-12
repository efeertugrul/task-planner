package provider

import (
	"todo-planning/internal/model"
)

// Provider defines the interface for API clients
type Provider interface {
	FetchTasks() ([]model.Task, error)
}
