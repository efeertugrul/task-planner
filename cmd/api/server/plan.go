package server

import (
	"fmt"
	"net/http"
	"todo-planning/internal/model"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetPlan(c *gin.Context) {
	// Get the plan
	assignments, err := s.planner.Plan()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create plan",
		})

		return
	}

	if len(assignments) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"assignments": [][]model.AssignmentResponse{},
			"totalHours":  0,
		})
		return
	}

	developerAssignments := make(map[uint][]model.AssignmentResponse)

	// Convert assignments to response format
	response := struct {
		Assignments [][]model.AssignmentResponse `json:"assignments"`
		TotalHours  float64                      `json:"totalHours"`
		TotalWeeks  int                          `json:"totalWeeks"`
	}{
		Assignments: make([][]model.AssignmentResponse, 0),
	}

	for _, assignment := range assignments {
		taskName := fmt.Sprintf("Task %s - %s", assignment.Task.Source, assignment.Task.ExternalID)
		if assignment.Task.Name != nil {
			taskName = *assignment.Task.Name
		}

		developerAssignments[assignment.DeveloperID] = append(developerAssignments[assignment.DeveloperID], model.AssignmentResponse{
			TaskName:        taskName,
			WeekNumber:      assignment.WeekNumber,
			CalculatedHours: assignment.CalculatedHours,
			Task:            assignment.Task,
			Developer:       assignment.Developer,
		})

		if assignment.WeekNumber > response.TotalWeeks {
			response.TotalWeeks = assignment.WeekNumber
		}

		response.TotalHours += assignment.CalculatedHours
	}

	for _, assignments := range developerAssignments {
		response.Assignments = append(response.Assignments, assignments)
	}
	c.Header("Status", "200")
	c.JSON(http.StatusOK, response)
}
