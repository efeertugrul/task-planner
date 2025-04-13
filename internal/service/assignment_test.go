package service

import (
	"testing"
	"time"

	"todo-planning/internal/model"
	"todo-planning/internal/utility"
)

func setupAssignmentTest(t *testing.T) (*AssignmentService, func()) {
	db := utility.GetTestDB()
	utility.AutoMigrate(&model.Assignment{})

	service := NewAssignmentService(db)

	// Return cleanup function
	cleanup := func() {
		utility.ClearTables()
		utility.CloseTestDB()
	}

	return service, cleanup
}

func TestAssignmentService_CreateAssignment(t *testing.T) {

	tests := []struct {
		name       string
		assignment *model.Assignment
		wantErr    bool
	}{
		{
			name: "create valid assignment",
			assignment: &model.Assignment{
				TaskID:          1,
				DeveloperID:     1,
				WeekNumber:      1,
				CalculatedHours: 8.0,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			wantErr: false,
		},
		{
			name: "create assignment with zero values",
			assignment: &model.Assignment{
				TaskID:          0,
				DeveloperID:     0,
				WeekNumber:      0,
				CalculatedHours: 0,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			wantErr: false,
		},
		{
			name: "create assignment with negative hours",
			assignment: &model.Assignment{
				TaskID:          1,
				DeveloperID:     1,
				WeekNumber:      1,
				CalculatedHours: -8.0,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Create a new variable for the closure
		t.Run(tt.name, func(t *testing.T) {

			// Create a new service instance for each subtest
			service, cleanup := setupAssignmentTest(t)
			defer cleanup()

			err := service.CreateAssignment(tt.assignment)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignmentService.CreateAssignment() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify the assignment was created
				var count int64
				if err := service.db.Model(&model.Assignment{}).Count(&count).Error; err != nil {
					t.Errorf("Failed to count assignments: %v", err)
				}
				if count != 1 {
					t.Errorf("Expected 1 assignment, got %d", count)
				}
			}
		})
	}
}

func TestAssignmentService_GetAssignments(t *testing.T) {

	service, cleanup := setupAssignmentTest(t)
	defer cleanup()

	// Create test assignments
	assignments := []model.Assignment{
		{
			TaskID:          1,
			DeveloperID:     1,
			WeekNumber:      1,
			CalculatedHours: 8.0,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			TaskID:          2,
			DeveloperID:     2,
			WeekNumber:      1,
			CalculatedHours: 6.0,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	// Create assignments
	for _, assignment := range assignments {
		if err := service.CreateAssignment(&assignment); err != nil {
			t.Fatalf("Failed to create test assignment: %v", err)
		}
	}

	// Test getting assignments
	got, err := service.GetAssignments()
	if err != nil {
		t.Errorf("AssignmentService.GetAssignments() error = %v", err)
		return
	}

	if len(got) != len(assignments) {
		t.Errorf("AssignmentService.GetAssignments() got = %v assignments, want %v assignments", len(got), len(assignments))
		return
	}

	// Verify the assignments' data
	for i, assignment := range got {
		if assignment.TaskID != assignments[i].TaskID {
			t.Errorf("AssignmentService.GetAssignments() got.TaskID = %v, want %v", assignment.TaskID, assignments[i].TaskID)
		}
		if assignment.DeveloperID != assignments[i].DeveloperID {
			t.Errorf("AssignmentService.GetAssignments() got.DeveloperID = %v, want %v", assignment.DeveloperID, assignments[i].DeveloperID)
		}
		if assignment.WeekNumber != assignments[i].WeekNumber {
			t.Errorf("AssignmentService.GetAssignments() got.WeekNumber = %v, want %v", assignment.WeekNumber, assignments[i].WeekNumber)
		}
		if assignment.CalculatedHours != assignments[i].CalculatedHours {
			t.Errorf("AssignmentService.GetAssignments() got.CalculatedHours = %v, want %v", assignment.CalculatedHours, assignments[i].CalculatedHours)
		}
	}
}

func TestAssignmentService_CreateBatchAssignments(t *testing.T) {

	tests := []struct {
		name         string
		assignments  []model.Assignment
		wantErr      bool
		expectedRows int64
	}{
		{
			name: "create batch of valid assignments",
			assignments: []model.Assignment{
				{
					TaskID:          1,
					DeveloperID:     1,
					WeekNumber:      1,
					CalculatedHours: 8.0,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
				{
					TaskID:          2,
					DeveloperID:     2,
					WeekNumber:      1,
					CalculatedHours: 6.0,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
			wantErr:      false,
			expectedRows: 2,
		},
		{
			name: "create large batch of assignments",
			assignments: func() []model.Assignment {
				var assignments []model.Assignment
				for i := 0; i < 150; i++ {
					assignments = append(assignments, model.Assignment{
						TaskID:          uint(i + 1),
						DeveloperID:     uint(i + 1),
						WeekNumber:      1,
						CalculatedHours: 8.0,
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					})
				}
				return assignments
			}(),
			wantErr:      false,
			expectedRows: 150,
		},
		{
			name: "create batch with invalid assignments",
			assignments: []model.Assignment{
				{
					TaskID:          0, // Invalid TaskID
					DeveloperID:     1,
					WeekNumber:      1,
					CalculatedHours: 8.0,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
			wantErr:      false,
			expectedRows: 1,
		},
	}

	for _, tt := range tests {
		tt := tt // Create a new variable for the closure
		t.Run(tt.name, func(t *testing.T) {

			// Create a new service instance for each subtest
			service, cleanup := setupAssignmentTest(t)
			defer cleanup()

			err := service.CreateBatchAssignments(tt.assignments)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignmentService.CreateBatchAssignments() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify the assignments were created
				var count int64
				if err := service.db.Model(&model.Assignment{}).Count(&count).Error; err != nil {
					t.Errorf("Failed to count assignments: %v", err)
				}
				if count != tt.expectedRows {
					t.Errorf("Expected %d assignments, got %d", tt.expectedRows, count)
				}
			}
		})
	}
}
