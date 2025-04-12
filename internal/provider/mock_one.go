package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"todo-planning/internal/logger"
	"todo-planning/internal/model"
	"todo-planning/internal/utility"
)

func NewMockOneClient(url string) *MockOneClient {
	return &MockOneClient{
		url: url,
		Client: http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type MockOneClient struct {
	url string
	http.Client
}

func (moc *MockOneClient) FetchTasks() ([]model.Task, error) {
	var tasks []*MockOneTask

	resp, err := moc.Client.Get(moc.url)
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var b []byte

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	if err := json.Unmarshal(b, &tasks); err != nil {
		logger.Error(err)
		return nil, err
	}

	var result []model.Task

	for i := range tasks {
		result = append(result, tasks[i].ToTask())
	}

	return result, nil
}

// MockOneTask represents the task structure from the mock-one provider
type MockOneTask struct {
	ID                uint    `json:"id"`
	Value             float64 `json:"value"`
	EstimatedDuration float64 `json:"estimated_duration"` // in hours
}

func (mot *MockOneTask) ToTask() model.Task {
	now := time.Now()
	return model.Task{
		ExternalID:        strconv.Itoa(int(mot.ID)),
		Difficulty:        mot.Value,
		EstimatedDuration: mot.EstimatedDuration,
		Name:              utility.ToPointer(fmt.Sprintf("Mock One Task %d", mot.ID)),
		Source:            "mock-one",
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}
