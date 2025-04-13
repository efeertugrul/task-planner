package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMockOneClient_FetchTasks(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tasks := []*MockOneTask{
			{
				ID:                1,
				Value:             3.5,
				EstimatedDuration: 2.0,
			},
			{
				ID:                2,
				Value:             5.0,
				EstimatedDuration: 3.0,
			},
		}
		json.NewEncoder(w).Encode(tasks)
	}))
	defer server.Close()

	// Create client
	client := NewMockOneClient(server.URL)

	// Test successful case
	tasks, err := client.FetchTasks()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(tasks))
	}

	// Verify first task
	if tasks[0].ExternalID != "1" {
		t.Errorf("Expected ExternalID '1', got '%s'", tasks[0].ExternalID)
	}
	if tasks[0].Difficulty != 3.5 {
		t.Errorf("Expected Difficulty 3.5, got %f", tasks[0].Difficulty)
	}
	if tasks[0].EstimatedDuration != 2.0 {
		t.Errorf("Expected EstimatedDuration 2.0, got %f", tasks[0].EstimatedDuration)
	}
	if *tasks[0].Name != "Mock One Task 1" {
		t.Errorf("Expected Name 'Mock One Task 1', got '%s'", *tasks[0].Name)
	}
	if tasks[0].Source != "mock-one" {
		t.Errorf("Expected Source 'mock-one', got '%s'", tasks[0].Source)
	}

	// Test error case - invalid response
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client = NewMockOneClient(server.URL)
	_, err = client.FetchTasks()
	if err == nil {
		t.Error("Expected error for invalid response, got nil")
	}
}

func TestMockTwoClient_FetchTasks(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tasks := []*MockTwoTask{
			{
				ID:     1,
				Zorluk: 3.5,
				Sure:   2.0,
			},
			{
				ID:     2,
				Zorluk: 5.0,
				Sure:   3.0,
			},
		}
		json.NewEncoder(w).Encode(tasks)
	}))
	defer server.Close()

	// Create client
	client := NewMockTwoClient(server.URL)

	// Test successful case
	tasks, err := client.FetchTasks()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(tasks))
	}

	// Verify first task
	if tasks[0].ExternalID != "1" {
		t.Errorf("Expected ExternalID '1', got '%s'", tasks[0].ExternalID)
	}
	if tasks[0].Difficulty != 3.5 {
		t.Errorf("Expected Difficulty 3.5, got %f", tasks[0].Difficulty)
	}
	if tasks[0].EstimatedDuration != 2.0 {
		t.Errorf("Expected EstimatedDuration 2.0, got %f", tasks[0].EstimatedDuration)
	}
	if *tasks[0].Name != "Mock Two Task 1" {
		t.Errorf("Expected Name 'Mock Two Task 1', got '%s'", *tasks[0].Name)
	}
	if tasks[0].Source != "mock-two" {
		t.Errorf("Expected Source 'mock-two', got '%s'", tasks[0].Source)
	}

	// Test error case - invalid response
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client = NewMockTwoClient(server.URL)
	_, err = client.FetchTasks()
	if err == nil {
		t.Error("Expected error for invalid response, got nil")
	}
}
