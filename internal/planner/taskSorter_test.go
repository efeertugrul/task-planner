package planner

import (
	"testing"

	"todo-planning/internal/model"
)

func TestDefaultTaskSorter_Sort(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []model.Task
		expected []model.Task
	}{
		{
			name: "sort by weight (duration * difficulty)",
			tasks: []model.Task{
				{ID: 1, EstimatedDuration: 2, Difficulty: 3}, // weight: 6
				{ID: 2, EstimatedDuration: 4, Difficulty: 1}, // weight: 4
				{ID: 3, EstimatedDuration: 3, Difficulty: 2}, // weight: 6
			},
			expected: []model.Task{
				{ID: 1, EstimatedDuration: 2, Difficulty: 3}, // weight: 6
				{ID: 3, EstimatedDuration: 3, Difficulty: 2}, // weight: 6
				{ID: 2, EstimatedDuration: 4, Difficulty: 1}, // weight: 4
			},
		},
		{
			name:     "empty tasks",
			tasks:    []model.Task{},
			expected: []model.Task{},
		},
		{
			name: "single task",
			tasks: []model.Task{
				{ID: 1, EstimatedDuration: 2, Difficulty: 3},
			},
			expected: []model.Task{
				{ID: 1, EstimatedDuration: 2, Difficulty: 3},
			},
		},
	}

	sorter := &DefaultTaskSorter{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sorter.Sort(tt.tasks)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d tasks, got %d", len(tt.expected), len(result))
				return
			}

			for i := range result {
				if result[i].ID != tt.expected[i].ID {
					t.Errorf("task %d: expected ID %d, got %d", i, tt.expected[i].ID, result[i].ID)
				}
				if result[i].EstimatedDuration != tt.expected[i].EstimatedDuration {
					t.Errorf("task %d: expected duration %.2f, got %.2f", i, tt.expected[i].EstimatedDuration, result[i].EstimatedDuration)
				}
				if result[i].Difficulty != tt.expected[i].Difficulty {
					t.Errorf("task %d: expected difficulty %.2f, got %.2f", i, tt.expected[i].Difficulty, result[i].Difficulty)
				}
			}
		})
	}
}
