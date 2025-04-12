package planner

import (
	"sort"
	"todo-planning/internal/model"
)

// TaskSorter defines the interface for different task sorting strategies
type TaskSorter interface {
	Sort(tasks []model.Task) []model.Task
}

// DefaultTaskSorter implements the default sorting strategy (by weight)
type DefaultTaskSorter struct{}

func (s *DefaultTaskSorter) Sort(tasks []model.Task) []model.Task {
	sortedTasks := make([]model.Task, len(tasks))
	copy(sortedTasks, tasks)

	sort.Slice(sortedTasks, func(i, j int) bool {
		wi := sortedTasks[i].EstimatedDuration * sortedTasks[i].Difficulty
		wj := sortedTasks[j].EstimatedDuration * sortedTasks[j].Difficulty
		return wi > wj
	})

	return sortedTasks
}
