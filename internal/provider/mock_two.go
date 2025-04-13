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

// MockTwoClient implements the Client interface for mock-two API
type MockTwoClient struct {
	url string
	http.Client
}

func NewMockTwoClient(url string) *MockTwoClient {
	return &MockTwoClient{
		url: url,
		Client: http.Client{
			Timeout: time.Second * 15,
		},
	}
}

func (mtc *MockTwoClient) FetchTasks() ([]model.Task, error) {
	var tasks []*MockTwoTask

	resp, err := mtc.Client.Get(mtc.url)
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

// MockTwoTask represents the task structure from the mock-two provider
type MockTwoTask struct {
	ID     uint    `json:"id"`
	Zorluk float64 `json:"zorluk"`
	Sure   float64 `json:"sure"` // in hours
}

func (mt *MockTwoTask) ToTask() model.Task {
	now := time.Now()

	return model.Task{
		ExternalID:        strconv.Itoa(int(mt.ID)),
		EstimatedDuration: mt.Sure,
		Name:              utility.ToPointer(fmt.Sprintf("Mock Two Task %d", mt.ID)),
		Difficulty:        mt.Zorluk,
		Source:            "mock-two",
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}
