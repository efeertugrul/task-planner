package planner

import (
	"errors"
	"testing"

	"todo-planning/internal/model"
)

type mockTaskService struct {
	tasks []model.Task
	err   error
}

func (m *mockTaskService) GetTasks() ([]model.Task, error) {
	return m.tasks, m.err
}

type mockDeveloperService struct {
	developers []model.Developer
	err        error
}

func (m *mockDeveloperService) GetDevelopers() ([]model.Developer, error) {
	return m.developers, m.err
}

func TestPlanner_Plan(t *testing.T) {
	tests := []struct {
		name          string
		tasks         []model.Task
		developers    []model.Developer
		taskErr       error
		developerErr  error
		expectedCount int
		expectError   bool
	}{
		{
			name: "successful planning",
			tasks: []model.Task{
				{ID: 1, Difficulty: 2, EstimatedDuration: 3},
				{ID: 2, Difficulty: 1, EstimatedDuration: 4},
			},
			developers: []model.Developer{
				{ID: 1, Productivity: 2},
				{ID: 2, Productivity: 3},
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:          "no tasks",
			tasks:         []model.Task{},
			developers:    []model.Developer{{ID: 1}},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "no developers",
			tasks:         []model.Task{{ID: 1}},
			developers:    []model.Developer{},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "task service error",
			taskErr:       errors.New("task service error"),
			developers:    []model.Developer{{ID: 1}},
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "developer service error",
			tasks:         []model.Task{{ID: 1}},
			developerErr:  errors.New("developer service error"),
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			print(tt.name)
			// Create mock services
			taskService := &mockTaskService{
				tasks: tt.tasks,
				err:   tt.taskErr,
			}
			developerService := &mockDeveloperService{
				developers: tt.developers,
				err:        tt.developerErr,
			}

			// Create planner with mock services
			planner := newPlanner(PlanningOptions{
				TaskService:       taskService,
				DeveloperService:  developerService,
				ChannelManager:    NewDefaultChannelManager(),
				SaveAssignments:   false,
				AssignmentService: nil,
				TaskSorter:        nil,
			})

			// Run planning
			assignments, err := planner.Plan()

			// Verify results
			if err != nil && !tt.expectError {
				t.Error("expected  error")

				return
			}

			if len(assignments) != tt.expectedCount {
				t.Errorf("expected %d assignments, got %d", tt.expectedCount, len(assignments))
			}

			// Verify assignments are valid
			for _, assignment := range assignments {
				if assignment.TaskID == 0 {
					t.Error("assignment has invalid task ID")
				}
				if assignment.DeveloperID == 0 {
					t.Error("assignment has invalid developer ID")
				}
				if assignment.WeekNumber < 1 {
					t.Error("assignment has invalid week number")
				}
				if assignment.CalculatedHours <= 0 {
					t.Error("assignment has invalid calculated hours")
				}
			}
		})
	}
}
