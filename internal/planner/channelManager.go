package planner

import (
	"todo-planning/internal/logger"
	"todo-planning/internal/model"
)

// max hours a developer can work in a week
const MaxHoursPerWeek = 45

type devState struct {
	Developer model.Developer
	WeekLoads map[int]float64 // week -> hours
}

// ChannelManager defines the interface for managing channel operations
type ChannelManager interface {
	// Core operations
	SendDevelopers(developers []model.Developer)
	SendTask(task model.Task)
	ReceiveAssignments() []model.Assignment
	HandleChannels()

	// Channel access for handlers
	GetTaskChannel() <-chan model.Task
	GetDeveloperChannel() <-chan []model.Developer
	GetAssignmentsChannel() chan<- []model.Assignment
	GetDoneChannel() chan bool
}

// DefaultChannelManager implements the default channel management strategy
type DefaultChannelManager struct {
	developerChannel   chan []model.Developer
	taskChannel        chan model.Task
	assignmentsChannel chan []model.Assignment
	doneChannel        chan bool
}

func NewDefaultChannelManager() *DefaultChannelManager {
	return &DefaultChannelManager{
		developerChannel:   make(chan []model.Developer),
		taskChannel:        make(chan model.Task),
		assignmentsChannel: make(chan []model.Assignment),
		doneChannel:        make(chan bool),
	}
}

func (cm *DefaultChannelManager) SendDevelopers(developers []model.Developer) {
	cm.developerChannel <- developers
}

func (cm *DefaultChannelManager) SendTask(task model.Task) {
	cm.taskChannel <- task
}

func (cm *DefaultChannelManager) ReceiveAssignments() []model.Assignment {
	return <-cm.assignmentsChannel
}

func (cm *DefaultChannelManager) GetTaskChannel() <-chan model.Task {
	return cm.taskChannel
}

func (cm *DefaultChannelManager) GetDeveloperChannel() <-chan []model.Developer {
	return cm.developerChannel
}

func (cm *DefaultChannelManager) GetAssignmentsChannel() chan<- []model.Assignment {
	return cm.assignmentsChannel
}

func (cm *DefaultChannelManager) GetDoneChannel() chan bool {
	return cm.doneChannel
}

func (cm *DefaultChannelManager) HandleChannels() {
	var (
		taskBuffer   = make([]model.Task, 0)
		taskAssigner *TaskAssigner
	)

	for {
		select {
		case currentTask := <-cm.taskChannel:
			if currentTask.ID > 0 {
				logger.Info("Received task", currentTask)

				assignments := make([]model.Assignment, 0)

				if taskAssigner == nil {
					taskBuffer = append(taskBuffer, currentTask)
				} else {
					for len(taskBuffer) > 0 {
						bTask := taskBuffer[0]
						taskBuffer = taskBuffer[1:]
						if assignment := taskAssigner.AssignTask(bTask); assignment != nil {
							assignments = append(assignments, *assignment)
						} else {
							logger.Info("Assignment can't be made to any Developer for task: ", bTask.Source+"-"+bTask.ExternalID, " consider splitting it into 2 issues")
						}
					}

					if assignment := taskAssigner.AssignTask(currentTask); assignment != nil {
						assignments = append(assignments, *assignment)
					} else {
						logger.Info("Assignment can't be made to any Developer for task: ", currentTask.Source+"-"+currentTask.ExternalID, " consider splitting it into 2 issues")
					}
				}

				logger.Info("Sending assignments")
				cm.assignmentsChannel <- assignments
			}
		case developers := <-cm.developerChannel:
			logger.Info("Received developers", developers)
			taskAssigner = NewTaskAssigner(developers)
		case <-cm.doneChannel:
			return
		}
	}
}
