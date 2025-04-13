package server

import (
	"todo-planning/internal/planner"
	"todo-planning/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	*gin.Engine
	planner *planner.Planner

	Port int
}

var serverInstance *Server

func NewServer(port int, database *gorm.DB) *Server {
	if serverInstance == nil {
		taskService := service.NewTaskService(database)
		developerService := service.NewDeveloperService(database)
		assignmentService := service.NewAssignmentService(database)

		serverInstance = &Server{
			Port: port,
			planner: planner.NewPlanner(planner.PlanningOptions{
				TaskService:       taskService,
				DeveloperService:  developerService,
				AssignmentService: assignmentService,
				SaveAssignments:   false,
				ChannelManager:    planner.NewDefaultChannelManager(),
			}),
		}

		gin.SetMode(gin.ReleaseMode)

		r := gin.Default()
		r.Use(gin.Logger(), gin.Recovery())

		// Configure CORS
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{"http://localhost:3000"}
		config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
		r.Use(cors.New(config))

		serverInstance.Engine = r
		serverInstance.RegisterRoutes()
	}

	return serverInstance
}

func (s *Server) RegisterRoutes() {
	api := s.Group("/api")
	api.GET("/weekly-plan", s.GetPlan)
}
