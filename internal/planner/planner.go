package planner

import (
	"fmt"
	"sync"
	"time"
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
	mu                sync.Mutex
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

var planner *Planner

func newPlanner(options PlanningOptions) *Planner {
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

func NewPlanner(options PlanningOptions) *Planner {
	if planner == nil {
		planner = newPlanner(options)
	}

	return planner
}

func (p *Planner) RunRoutines() {
	go p.channelManager.HandleChannels()
}

func (p *Planner) Plan() ([]model.Assignment, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.RunRoutines()

	assignments, err := p.plan()
	p.Stop()

	return assignments, err
}

func (p *Planner) plan() ([]model.Assignment, error) {
	// Fetch developers first
	developers, err := p.developerService.GetDevelopers()
	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to get developers: %w", err)
	}
	p.channelManager.SendDevelopers(developers)

	// Fetch and sort tasks
	tasks, err := p.taskService.GetTasks()
	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	if len(tasks) == 0 {
		return nil, nil
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

	return assignments, nil
}

func (p *Planner) Stop() {
	p.channelManager.GetDoneChannel() <- true
	time.Sleep(100 * time.Millisecond)
}
