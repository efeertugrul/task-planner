package planner

import (
	"math"
	"testing"

	"todo-planning/internal/model"
)

func TestCalculateTaskEffort(t *testing.T) {
	tests := []struct {
		name     string
		task     model.Task
		expected float64
	}{
		{
			name:     "basic calculation",
			task:     model.Task{Difficulty: 2, EstimatedDuration: 3},
			expected: 6.0,
		},
		{
			name:     "zero values",
			task:     model.Task{Difficulty: 0, EstimatedDuration: 0},
			expected: 0.0,
		},
		{
			name:     "large values",
			task:     model.Task{Difficulty: 10, EstimatedDuration: 10},
			expected: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTaskEffort(tt.task)
			if result != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestCalculateHoursNeeded(t *testing.T) {
	tests := []struct {
		name       string
		taskEffort float64
		developer  model.Developer
		expected   float64
	}{
		{
			name:       "basic calculation",
			taskEffort: 10.0,
			developer:  model.Developer{Productivity: 2},
			expected:   5.0,
		},
		{
			name:       "zero productivity",
			taskEffort: 10.0,
			developer:  model.Developer{Productivity: 0},
			expected:   math.MaxFloat64,
		},
		{
			name:       "high productivity",
			taskEffort: 100.0,
			developer:  model.Developer{Productivity: 10},
			expected:   10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateHoursNeeded(tt.taskEffort, tt.developer)
			if result != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestTaskAssigner_FindBestFit(t *testing.T) {
	developers := []model.Developer{
		{ID: 1, Productivity: 2},
		{ID: 2, Productivity: 3},
	}

	tests := []struct {
		name          string
		task          model.Task
		devStates     []*devState
		expectedDevID uint
		expectedWeek  int
		expectedHours float64
	}{
		{
			name: "first developer gets task",
			task: model.Task{ID: 1, Difficulty: 2, EstimatedDuration: 3},
			devStates: []*devState{
				{
					Developer: developers[0],
					WeekLoads: make(map[int]float64),
				},
				{
					Developer: developers[1],
					WeekLoads: map[int]float64{1: 45},
				},
			},
			expectedDevID: 1,
			expectedWeek:  1,
			expectedHours: 3.0,
		},
		{
			name:          "second developer gets task due to higher productivity",
			task:          model.Task{ID: 2, Difficulty: 3, EstimatedDuration: 2},
			expectedDevID: 2,
			expectedWeek:  1,
			expectedHours: 2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskAssigner := NewTaskAssigner(developers)

			if len(tt.devStates) > 0 {
				taskAssigner.devStates = tt.devStates
			}

			devState, week, hours := taskAssigner.FindBestFit(tt.task)

			if devState.Developer.ID != tt.expectedDevID {
				t.Errorf("expected developer ID %d, got %d", tt.expectedDevID, devState.Developer.ID)

				return
			}
			if week != tt.expectedWeek {
				t.Errorf("expected week %d, got %d", tt.expectedWeek, week)

				return
			}
			if hours != tt.expectedHours {
				t.Errorf("expected hours %f, got %f", tt.expectedHours, hours)

				return
			}
		})
	}
}

func TestTaskAssigner_AssignTask(t *testing.T) {
	developers := []model.Developer{
		{ID: 1, Productivity: 2},
	}
	taskAssigner := NewTaskAssigner(developers)

	task := model.Task{ID: 1, Difficulty: 2, EstimatedDuration: 3}
	assignment := taskAssigner.AssignTask(task)

	if assignment.TaskID != task.ID {
		t.Errorf("expected task ID %d, got %d", task.ID, assignment.TaskID)
	}
	if assignment.DeveloperID != developers[0].ID {
		t.Errorf("expected developer ID %d, got %d", developers[0].ID, assignment.DeveloperID)
	}
	if assignment.WeekNumber != 1 {
		t.Errorf("expected week 1, got %d", assignment.WeekNumber)
	}
	if assignment.CalculatedHours != 3.0 {
		t.Errorf("expected hours 3.0, got %f", assignment.CalculatedHours)
	}
}
