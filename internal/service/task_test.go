package service

import (
	"testing"
	"time"

	"todo-planning/internal/model"
	"todo-planning/internal/utility"
)

func setupTaskTest(t *testing.T) (*TaskService, func()) {
	db := utility.GetTestDB()
	utility.AutoMigrate(&model.Task{})

	service := NewTaskService(db)

	// Return cleanup function
	cleanup := func() {
		utility.ClearTables()
		utility.CloseTestDB()
	}

	return service, cleanup
}

func TestTaskService_StoreTasks(t *testing.T) {
	service, cleanup := setupTaskTest(t)
	defer cleanup()

	tests := []struct {
		name    string
		tasks   []model.Task
		wantErr bool
	}{
		{
			name: "store single task",
			tasks: []model.Task{
				{
					ExternalID:        "1",
					Name:              utility.ToPointer("Test Task 1"),
					Difficulty:        3.0,
					EstimatedDuration: 2.0,
					Source:            "test",
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "store multiple tasks",
			tasks: []model.Task{
				{
					ExternalID:        "2",
					Name:              utility.ToPointer("Test Task 2"),
					Difficulty:        2.0,
					EstimatedDuration: 1.0,
					Source:            "test",
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				},
				{
					ExternalID:        "3",
					Name:              utility.ToPointer("Test Task 3"),
					Difficulty:        4.0,
					EstimatedDuration: 3.0,
					Source:            "test",
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "store duplicate task",
			tasks: []model.Task{
				{
					ExternalID:        "1",
					Name:              utility.ToPointer("Test Task 1"),
					Difficulty:        3.0,
					EstimatedDuration: 2.0,
					Source:            "test",
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				},
			},
			wantErr: false, // Should not error due to ON CONFLICT DO NOTHING
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := service.StoreTasks(tt.tasks); (err != nil) != tt.wantErr {
				t.Errorf("TaskService.StoreTasks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskService_GetTasks(t *testing.T) {
	service, cleanup := setupTaskTest(t)
	defer cleanup()

	// First store some tasks
	tasks := []model.Task{
		{
			ExternalID:        "1",
			Name:              utility.ToPointer("Test Task 1"),
			Difficulty:        3.0,
			EstimatedDuration: 2.0,
			Source:            "test",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			ExternalID:        "2",
			Name:              utility.ToPointer("Test Task 2"),
			Difficulty:        2.0,
			EstimatedDuration: 1.0,
			Source:            "test",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
	}

	if err := service.StoreTasks(tasks); err != nil {
		t.Fatalf("Failed to store tasks: %v", err)
	}

	// Test getting tasks
	got, err := service.GetTasks()
	if err != nil {
		t.Errorf("TaskService.GetTasks() error = %v", err)
		return
	}

	if len(got) != len(tasks) {
		t.Errorf("TaskService.GetTasks() got = %v tasks, want %v tasks", len(got), len(tasks))
	}
}
