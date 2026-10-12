package handler

import (
	"taskflow/internal/auth"
	"taskflow/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter creates and configures the HTTP router.
func SetupRouter(userService *service.UserService, projectService *service.ProjectService, taskService *service.TaskService, memberService *service.ProjectMemberService, commentService *service.CommentService, historyService *service.TaskHistoryService) *gin.Engine {
	router := gin.Default()
	authHandler := NewAuthHandler(userService)
	projectHandler := NewProjectHandler(projectService)
	taskHandler := NewTaskHandler(taskService)
	memberHandler := NewProjectMemberHandler(memberService)
	commentHandler := NewCommentHandler(commentService)
	historyHandler := NewTaskHistoryHandler(historyService)

	// Public routes.
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	// Everything below requires a valid token.
	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware())

	// Current user.
	protected.GET("/me", authHandler.Me)

	// Projects.
	protected.POST("/projects", projectHandler.CreateProject)
	protected.GET("/projects", projectHandler.GetProjects)
	protected.GET("/projects/:id", projectHandler.GetProject)
	protected.PUT("/projects/:id", projectHandler.UpdateProject)
	protected.DELETE("/projects/:id", projectHandler.DeleteProject)

	// Project members and roles.
	protected.POST("/projects/:id/members", memberHandler.AddMember)
	protected.GET("/projects/:id/members", memberHandler.ListMembers)
	protected.GET("/projects/:id/members/:userID", memberHandler.GetMember)
	protected.GET("/projects/:id/members/:userID/check", memberHandler.IsMember)
	protected.PATCH("/projects/:id/members/:userID", memberHandler.UpdateRole)
	protected.DELETE("/projects/:id/members/:userID", memberHandler.RemoveMember)

	// Tasks.
	protected.POST("/projects/:id/tasks", taskHandler.CreateTask)
	protected.GET("/projects/:id/tasks", taskHandler.GetProjectTasks)
	protected.GET("/tasks/:id", taskHandler.GetTask)
	protected.PUT("/tasks/:id", taskHandler.UpdateTask)
	protected.DELETE("/tasks/:id", taskHandler.DeleteTask)
	protected.PATCH("/tasks/:id/assign", taskHandler.AssignTask)
	protected.DELETE("/tasks/:id/assign", taskHandler.UnassignTask)

	// Task comments.
	protected.POST("/tasks/:id/comments", commentHandler.CreateComment)
	protected.GET("/tasks/:id/comments", commentHandler.GetByTaskID)
	protected.GET("/comments/:id", commentHandler.GetByID)
	protected.DELETE("/comments/:id", commentHandler.Delete)

	// Task and project history.
	protected.GET("/tasks/:id/history", historyHandler.GetByTaskID)
	protected.GET("/task-history/:id", historyHandler.GetByID)
	protected.GET("/projects/:id/history", historyHandler.GetByProjectID)

	return router
}
