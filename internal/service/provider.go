package service

import (
	"fmt"

	"todo-planning/internal/config"
	"todo-planning/internal/logger"
	"todo-planning/internal/model"
	"todo-planning/internal/provider"
)

type ProviderService struct {
	providers []provider.Provider
}

func NewProviderService() *ProviderService {
	return &ProviderService{
		providers: InitProviders(),
	}
}

// FetchTasksFromProviders fetches tasks from all providers
func (s *ProviderService) FetchTasksFromProviders() ([]model.Task, error) {
	var allTasks []model.Task

	for _, p := range s.providers {
		tasks, err := p.FetchTasks()
		if err != nil {
			logger.Error(fmt.Errorf("failed to fetch tasks from provider: %w", err))
			continue
		}

		allTasks = append(allTasks, tasks...)
	}

	return allTasks, nil
}

// InitProviders returns a slice of Provider instances
// this behavior is not ideal, but it's a simple example
// to demonstrate the concept
func InitProviders() []provider.Provider {
	config, err := config.Load()
	if err != nil {
		logger.Error(err)
	}

	return []provider.Provider{
		provider.NewMockOneClient(config.ProviderConfig.MockOne.Url),
		provider.NewMockTwoClient(config.ProviderConfig.MockTwo.Url),
	}
}
