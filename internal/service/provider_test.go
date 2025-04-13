package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"todo-planning/internal/model"
	"todo-planning/internal/provider"
	"todo-planning/internal/utility"
)

type mockProvider struct {
	tasks []model.Task
	err   error
}

func (m *mockProvider) FetchTasks() ([]model.Task, error) {
	return m.tasks, m.err
}

type mockProviderServer struct {
}

func (m *mockProviderServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"tasks": [{"id": "1", "name": "Task 1", "difficulty": 1, "estimated_duration": 1}, {"id": "2", "name": "Task 2", "difficulty": 2, "estimated_duration": 2}]}`))
}

func serveMockProvider2(port string) {
	handler := &mockProviderServer{}
	http.ListenAndServe("0.0.0.0:"+port, handler)

}

type mockProviderClient2 struct {
	port string
}

func (m *mockProviderClient2) FetchTasks() ([]model.Task, error) {
	response, err := http.Get("http://0.0.0.0:" + m.port)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	type MockStruct struct {
		ID                string  `json:"id"`
		Name              string  `json:"name"`
		Difficulty        float64 `json:"difficulty"`
		EstimatedDuration float64 `json:"estimated_duration"`
	}

	type Response struct {
		Tasks []MockStruct `json:"tasks"`
	}

	var mockcliTasks Response

	err = json.NewDecoder(response.Body).Decode(&mockcliTasks)
	if err != nil {
		return nil, err
	}

	var tasks []model.Task
	for _, task := range mockcliTasks.Tasks {
		tasks = append(tasks, model.Task{
			ExternalID:        task.ID,
			Name:              utility.ToPointer(task.Name),
			Difficulty:        task.Difficulty,
			EstimatedDuration: task.EstimatedDuration,
		})
	}

	return tasks, nil
}

func TestProviderService_FetchTasksFromProviders(t *testing.T) {
	tests := []struct {
		name      string
		providers []provider.Provider
		wantTasks int
		wantErr   bool
	}{
		{
			name: "successful fetch from single provider",
			providers: []provider.Provider{
				&mockProvider{
					tasks: []model.Task{
						{ExternalID: "1", Name: utility.ToPointer("Task 1")},
						{ExternalID: "2", Name: utility.ToPointer("Task 2")},
					},
					err: nil,
				},
			},
			wantTasks: 2,
			wantErr:   false,
		},
		{
			name: "successful fetch from multiple providers",
			providers: []provider.Provider{
				&mockProvider{
					tasks: []model.Task{
						{ExternalID: "1", Name: utility.ToPointer("Task 1")},
					},
					err: nil,
				},
				&mockProvider{
					tasks: []model.Task{
						{ExternalID: "2", Name: utility.ToPointer("Task 2")},
						{ExternalID: "3", Name: utility.ToPointer("Task 3")},
					},
					err: nil,
				},
			},
			wantTasks: 3,
			wantErr:   false,
		},
		{
			name: "provider error continues to next provider",
			providers: []provider.Provider{
				&mockProvider{
					tasks: nil,
					err:   errors.New("provider error"),
				},
				&mockProvider{
					tasks: []model.Task{
						{ExternalID: "1", Name: utility.ToPointer("Task 1")},
					},
					err: nil,
				},
			},
			wantTasks: 1,
			wantErr:   false,
		},
		{
			name: "multiple providers with mock provider client 2",
			providers: []provider.Provider{
				&mockProviderClient2{port: "8081"},
				&mockProviderClient2{port: "8082"},
			},
			wantTasks: 4,
			wantErr:   false,
		},
	}

	go serveMockProvider2("8081")
	go serveMockProvider2("8082")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new service with mock providers
			service := &ProviderService{
				providers: make([]provider.Provider, len(tt.providers)),
			}

			copy(service.providers, tt.providers)

			got, err := service.FetchTasksFromProviders()
			if (err != nil) != tt.wantErr {
				t.Errorf("ProviderService.FetchTasksFromProviders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.wantTasks {
				t.Errorf("ProviderService.FetchTasksFromProviders() got = %v tasks, want %v tasks", len(got), tt.wantTasks)
			}
		})
	}
}
