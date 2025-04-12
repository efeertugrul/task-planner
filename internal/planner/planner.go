package planner

import (
	"todo-planning/internal/logger"
	"todo-planning/internal/model"

	"gorm.io/gorm"
)

// Service interfaces for dependency injection
type TaskService interface {
	GetTasks() ([]model.Task, error)
}

type DeveloperService interface {
	GetDevelopers() ([]model.Developer, error)
}

type AssignmentService interface {
	// Add methods as needed
}

type Planner struct {
	taskService       TaskService
	developerService  DeveloperService
	assignmentService AssignmentService
	taskSorter        TaskSorter
	channelManager    ChannelManager
}

type PlanningOptions struct {
	DB                *gorm.DB
	SaveAssignments   bool
	TaskService       TaskService
	DeveloperService  DeveloperService
	AssignmentService AssignmentService
	TaskSorter        TaskSorter
	ChannelManager    ChannelManager
}

func NewPlanner(options PlanningOptions) *Planner {

	taskSorter := options.TaskSorter
	if taskSorter == nil {
		taskSorter = &DefaultTaskSorter{}
	}

	channelManager := options.ChannelManager
	if channelManager == nil {
		channelManager = NewDefaultChannelManager()
	}

	return &Planner{
		taskService:       options.TaskService,
		developerService:  options.DeveloperService,
		assignmentService: options.AssignmentService,
		taskSorter:        taskSorter,
		channelManager:    channelManager,
	}
}

func (p *Planner) RunRoutines() {
	go p.channelManager.HandleChannels()
}

func (p *Planner) Plan() []model.Assignment {
	// Fetch developers first
	developers, err := p.developerService.GetDevelopers()
	if err != nil {
		logger.Error(err)
		return nil
	}
	p.channelManager.SendDevelopers(developers)

	// Fetch and sort tasks
	tasks, err := p.taskService.GetTasks()
	if err != nil {
		logger.Error(err)
		return nil
	}

	// Sort tasks using the configured sorter
	sortedTasks := p.taskSorter.Sort(tasks)

	// Send tasks in batches
	var assignments []model.Assignment
	for _, task := range sortedTasks {
		p.channelManager.SendTask(task)
		currentAssignments := p.channelManager.ReceiveAssignments()
		assignments = append(assignments, currentAssignments...)
	}

	return assignments
}
