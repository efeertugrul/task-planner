package service

import (
	"testing"
	"time"

	"todo-planning/internal/model"
	"todo-planning/internal/utility"
)

func setupDeveloperTest(t *testing.T) (*DeveloperService, func()) {
	db := utility.GetTestDB()
	utility.AutoMigrate(&model.Developer{})

	service := NewDeveloperService(db)

	// Return cleanup function
	cleanup := func() {
		utility.ClearTables()
		utility.CloseTestDB()
	}

	return service, cleanup
}

func TestDeveloperService_GetDevelopers(t *testing.T) {
	service, cleanup := setupDeveloperTest(t)
	defer cleanup()

	// First store some developers
	developers := []model.Developer{
		{
			Name:         "Developer 1",
			Productivity: 1.0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Name:         "Developer 2",
			Productivity: 2.0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	if err := service.db.Create(&developers).Error; err != nil {
		t.Fatalf("Failed to create developers: %v", err)
	}

	// Test getting developers
	got, err := service.GetDevelopers()
	if err != nil {
		t.Errorf("DeveloperService.GetDevelopers() error = %v", err)
		return
	}

	if len(got) != len(developers) {
		t.Errorf("DeveloperService.GetDevelopers() got = %v developers, want %v developers", len(got), len(developers))
	}

	// Verify the developers' data
	for i, dev := range got {
		if dev.Name != developers[i].Name {
			t.Errorf("DeveloperService.GetDevelopers() got.Name = %v, want %v", dev.Name, developers[i].Name)
		}
		if dev.Productivity != developers[i].Productivity {
			t.Errorf("DeveloperService.GetDevelopers() got.Productivity = %v, want %v", dev.Productivity, developers[i].Productivity)
		}
	}
}
