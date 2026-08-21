package handler

import (
	"taskflow/internal/auth"
	"taskflow/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter creates and configures the HTTP router.
func SetupRouter(
	userService *service.UserService,
	projectService *service.ProjectService,
	taskService *service.TaskService,
	memberService *service.ProjectMemberService,
	commentService *service.CommentService,
	historyService *service.TaskHistoryService,
) *gin.Engine {

	router := gin.Default()

	// =========================
	// Handlers
	// =========================

	authHandler := NewAuthHandler(userService)
	projectHandler := NewProjectHandler(projectService)
	taskHandler := NewTaskHandler(taskService)
	memberHandler := NewProjectMemberHandler(memberService)
	commentHandler := NewCommentHandler(commentService)
	historyHandler := NewTaskHistoryHandler(historyService)

	// =========================
	// Public routes
	// =========================

	authRoutes := router.Group("/")

	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// =========================
	// Protected routes
	// =========================

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware())

	// =========================
	// Authentication
	// =========================

	protected.GET("/me", authHandler.Me)

	// =========================
	// Projects
	// =========================

	protected.POST(
		"/projects",
		projectHandler.CreateProject,
	)

	protected.GET(
		"/projects",
		projectHandler.GetProjects,
	)

	protected.GET(
		"/projects/:id",
		projectHandler.GetProject,
	)

	protected.PUT(
		"/projects/:id",
		projectHandler.UpdateProject,
	)

	protected.DELETE(
		"/projects/:id",
		projectHandler.DeleteProject,
	)

	// =========================
	// Project members
	// =========================

	protected.POST(
		"/projects/:id/members",
		memberHandler.AddMember,
	)

	protected.GET(
		"/projects/:id/members",
		memberHandler.ListMembers,
	)

	protected.GET(
		"/projects/:id/members/:userID",
		memberHandler.GetMember,
	)

	protected.GET(
		"/projects/:id/members/:userID/check",
		memberHandler.IsMember,
	)

	protected.PATCH(
		"/projects/:id/members/:userID",
		memberHandler.UpdateRole,
	)

	protected.DELETE(
		"/projects/:id/members/:userID",
		memberHandler.RemoveMember,
	)

	// =========================
	// Tasks
	// =========================

	protected.POST(
		"/projects/:projectID/tasks",
		taskHandler.CreateTask,
	)

	protected.GET(
		"/projects/:projectID/tasks",
		taskHandler.GetProjectTasks,
	)

	protected.GET(
		"/tasks/:id",
		taskHandler.GetTask,
	)

	protected.PUT(
		"/tasks/:id",
		taskHandler.UpdateTask,
	)

	protected.DELETE(
		"/tasks/:id",
		taskHandler.DeleteTask,
	)

	protected.PATCH(
		"/tasks/:id/assign",
		taskHandler.AssignTask,
	)

	protected.DELETE(
		"/tasks/:id/assign",
		taskHandler.UnassignTask,
	)

	// =========================
	// Comments
	// =========================

	// Create comment
	// POST /tasks/:id/comments
	protected.POST(
		"/tasks/:id/comments",
		commentHandler.CreateComment,
	)

	// Get all comments for task
	// GET /tasks/:id/comments
	protected.GET(
		"/tasks/:id/comments",
		commentHandler.GetByTaskID,
	)

	// Get comment by ID
	// GET /comments/:id
	protected.GET(
		"/comments/:id",
		commentHandler.GetByID,
	)

	// Delete comment
	// DELETE /comments/:id
	protected.DELETE(
		"/comments/:id",
		commentHandler.Delete,
	)

	// =========================
	// Task history
	// =========================
	// Get task history by task ID
	// GET /tasks/:id/history
	protected.GET(
		"/tasks/:id/history",
		historyHandler.GetByTaskID,
	)
	// Get history record by ID
	// GET /task-history/:id
	protected.GET(
		"/task-history/:id",
		historyHandler.GetByID,
	)
	// Get history of all tasks in project
	// GET /projects/:id/history
	protected.GET(
		"/projects/:id/history",
		historyHandler.GetByProjectID,
	)
	return router
}
