package main

import (
	"flag"
	"fmt"
	"os"

	"todo-planning/internal/db"
	"todo-planning/internal/logger"
	"todo-planning/internal/model"
	"todo-planning/internal/service"
)

func main() {
	// Check if a subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("expected 'fetch' or 'init-db' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "fetch":
		fetchCmd := flag.NewFlagSet("fetch", flag.ExitOnError)
		fetchCmd.Parse(os.Args[2:])
		fetchTasks()
	case "init-db":
		initDBCmd := flag.NewFlagSet("init-db", flag.ExitOnError)
		force := initDBCmd.Bool("force", false, "Force initialization even if developers exist")
		initDBCmd.Parse(os.Args[2:])
		initializeDatabase(*force)
	default:
		fmt.Printf("unknown subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func fetchTasks() {
	database, err := db.NewConnection()
	if err != nil {
		logger.Error(fmt.Errorf("failed to connect to database: %w", err))
		os.Exit(1)
	}

	// Fetch tasks from providers
	providerService := service.NewProviderService()
	tasks, err := providerService.FetchTasksFromProviders()
	if err != nil {
		logger.Error(fmt.Errorf("failed to fetch tasks from providers: %w", err))
		os.Exit(1)
	}

	// Store tasks in database
	taskService := service.NewTaskService(database)
	if err := taskService.StoreTasks(tasks); err != nil {
		logger.Error(fmt.Errorf("failed to store tasks: %w", err))
		os.Exit(1)
	}

	// Get and display all tasks
	storedTasks, err := taskService.GetTasks()
	if err != nil {
		logger.Error(fmt.Errorf("failed to get tasks: %w", err))
		os.Exit(1)
	}

	// Print tasks in a formatted way
	fmt.Printf("\nFetched %d tasks:\n", len(storedTasks))
	for _, task := range storedTasks {
		fmt.Printf("ID: %d, Name: %s, Difficulty: %.2f, Duration: %.2f, Source: %s\n",
			task.ID, *task.Name, task.Difficulty, task.EstimatedDuration, task.Source)
	}
}

func initializeDatabase(force bool) {
	database, err := db.NewConnection()
	if err != nil {
		logger.Error(fmt.Errorf("failed to connect to database: %w", err))
		os.Exit(1)
	}

	// Run migrations
	if err := db.AutoMigrate(database); err != nil {
		logger.Error(fmt.Errorf("failed to run migrations: %w", err))
		os.Exit(1)
	}

	// Check if developers already exist
	var count int64
	if err := database.Model(&model.Developer{}).Count(&count).Error; err != nil {
		logger.Error(fmt.Errorf("failed to check existing developers: %w", err))
		os.Exit(1)
	}

	if count > 0 && !force {
		fmt.Println("Developers already exist in the database. Use --force to reinitialize.")
		os.Exit(0)
	}

	if force {
		if err := database.Exec("DELETE FROM developers").Error; err != nil {
			logger.Error(fmt.Errorf("failed to delete developers: %w", err))
			os.Exit(1)
		}
	}

	// Create developers
	developers := []model.Developer{
		{Name: "Dev1", Productivity: 1.0},
		{Name: "Dev2", Productivity: 2.0},
		{Name: "Dev3", Productivity: 3.0},
		{Name: "Dev4", Productivity: 4.0},
		{Name: "Dev5", Productivity: 5.0},
	}

	if err := database.Create(&developers).Error; err != nil {
		logger.Error(fmt.Errorf("failed to create developers: %w", err))
		os.Exit(1)
	}

	fmt.Println("Database initialized successfully with 5 developers.")
}
