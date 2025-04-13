package planner

import (
	"math"
	"todo-planning/internal/logger"
	"todo-planning/internal/model"
)

// TaskAssigner handles the core task assignment logic
type TaskAssigner struct {
	developers []model.Developer
	devStates  []*devState
}

func NewTaskAssigner(developers []model.Developer) *TaskAssigner {
	devStates := make([]*devState, 0, len(developers))
	for _, dev := range developers {
		devStates = append(devStates, &devState{
			Developer: dev,
			WeekLoads: make(map[int]float64, 10),
		})
	}
	return &TaskAssigner{
		developers: developers,
		devStates:  devStates,
	}
}

func (ta *TaskAssigner) AssignTask(task model.Task) *model.Assignment {
	bestDev, week, hours := ta.FindBestFit(task)
	if bestDev == nil {
		return nil
	}

	return &model.Assignment{
		TaskID:          task.ID,
		DeveloperID:     bestDev.Developer.ID,
		WeekNumber:      week,
		CalculatedHours: hours,
		Task:            task,
		Developer:       bestDev.Developer,
	}
}

// FindBestFit finds the best developer and week for a task
func (ta *TaskAssigner) FindBestFit(task model.Task) (*devState, int, float64) {
	var (
		bestDev  *devState
		bestWeek = 1
		minTotal = math.MaxInt64
	)

	taskEffort := CalculateTaskEffort(task)

	for _, dev := range ta.devStates {
		hoursNeeded := CalculateHoursNeeded(taskEffort, dev.Developer)
		if hoursNeeded > MaxHoursPerWeek {
			continue
		}
		// Find the first week where the task can fit
		week := 1
		for {
			if dev.WeekLoads[week]+hoursNeeded <= MaxHoursPerWeek {
				break
			}
			week++
		}
		changeDev := false
		if week < minTotal {
			changeDev = true
		} else if week == minTotal {
			var (
				currentDevsWorkload float64
				bestDevWorkload     float64
				ok                  bool
			)

			if currentDevsWorkload, ok = dev.WeekLoads[week]; !ok {
				currentDevsWorkload = 0
			}

			if bestDevWorkload, ok = bestDev.WeekLoads[bestWeek]; !ok {
				bestDevWorkload = 0
			}

			if currentDevsWorkload+hoursNeeded < bestDevWorkload {
				changeDev = true
			}
		}

		if changeDev {
			if bestDev != nil {
				logger.Info("bestDev is changing from:  ", *bestDev, "to ", *dev)
				bestDev.WeekLoads[bestWeek] -= CalculateHoursNeeded(taskEffort, bestDev.Developer)
			}

			bestDev = dev
			bestDev.WeekLoads[week] += hoursNeeded
			bestWeek = week
			minTotal = week
		}

		// if week < minTotal || (week == minTotal && dev.WeekLoads[week]+hoursNeeded < bestDev.WeekLoads[bestWeek]) {
		// 	if bestDev != nil {
		// 		logger.Info("bestDev is changing from:  ", *bestDev, "to ", *dev)
		// 		bestDev.WeekLoads[bestWeek] -= CalculateHoursNeeded(taskEffort, bestDev.Developer)
		// 	}

		// 	bestDev = dev
		// 	bestDev.WeekLoads[week] += hoursNeeded
		// 	bestWeek = week
		// 	minTotal = week
		// }
	}

	if bestDev == nil {
		return nil, 0, 0
	}
	// Update the selected developer's workload

	hoursNeeded := CalculateHoursNeeded(taskEffort, bestDev.Developer)

	return bestDev, bestWeek, hoursNeeded
}

// CalculateTaskEffort calculates the effort needed for a task
func CalculateTaskEffort(task model.Task) float64 {
	return float64(task.Difficulty * task.EstimatedDuration)
}

// CalculateHoursNeeded calculates the hours needed for a task based on developer productivity
func CalculateHoursNeeded(taskEffort float64, developer model.Developer) float64 {
	if developer.Productivity == 0 {
		return math.MaxFloat64
	}

	return taskEffort / float64(developer.Productivity)
}
