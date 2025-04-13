package planner

import (
	"testing"
	"time"

	"todo-planning/internal/model"
)

func TestDefaultChannelManager(t *testing.T) {
	manager := NewDefaultChannelManager()

	t.Run("send and receive developers", func(t *testing.T) {
		developers := []model.Developer{
			{ID: 1},
			{ID: 2},
		}

		go func() {
			manager.SendDevelopers(developers)
		}()

		select {
		case received := <-manager.GetDeveloperChannel():
			if len(received) != len(developers) {
				t.Errorf("expected %d developers, got %d", len(developers), len(received))
			}
			for i := range received {
				if received[i].ID != developers[i].ID {
					t.Errorf("developer %d: expected ID %d, got %d", i, developers[i].ID, received[i].ID)
				}
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("timeout waiting for developers")
		}
	})

	t.Run("send and receive task", func(t *testing.T) {
		task := model.Task{ID: 1}

		go func() {
			manager.SendTask(task)
		}()

		select {
		case received := <-manager.GetTaskChannel():
			if received.ID != task.ID {
				t.Errorf("expected task ID %d, got %d", task.ID, received.ID)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("timeout waiting for task")
		}
	})

	t.Run("send and receive assignments", func(t *testing.T) {
		assignments := []model.Assignment{
			{TaskID: 1, DeveloperID: 1, WeekNumber: 1},
			{TaskID: 2, DeveloperID: 2, WeekNumber: 1},
		}

		go func() {
			manager.GetAssignmentsChannel() <- assignments
		}()

		received := manager.ReceiveAssignments()
		if len(received) != len(assignments) {
			t.Errorf("expected %d assignments, got %d", len(assignments), len(received))
		}
		for i := range received {
			if received[i].TaskID != assignments[i].TaskID {
				t.Errorf("assignment %d: expected task ID %d, got %d", i, assignments[i].TaskID, received[i].TaskID)
			}
		}
	})
}
