// Package handler provides HTTP handlers for task management.
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"taskflow/internal/model"
	"taskflow/internal/service"

	"github.com/gin-gonic/gin"
)

// TaskHandler handles task-related HTTP requests.
type TaskHandler struct {
	taskService *service.TaskService
}

// NewTaskHandler creates a new TaskHandler instance.
func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTaskRequest represents the payload for creating a task.
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	AssigneeID  *int64 `json:"assignee_id"`
}

// UpdateTaskRequest represents the payload for updating a task.
type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	AssigneeID  *int64 `json:"assignee_id"`
}

// AssignTaskRequest represents the payload for assigning a task.
type AssignTaskRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

// TaskResponse represents task data in API responses.
type TaskResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProjectID   int64  `json:"project_id"`
	AssigneeID  *int64 `json:"assignee_id"`
}

// taskToResponse converts a Task model to TaskResponse.
func taskToResponse(task *model.Task) TaskResponse {
	return TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		ProjectID:   task.ProjectID,
		AssigneeID:  task.AssigneeID,
	}
}

// CreateTask creates a new task in a project.
func (h *TaskHandler) CreateTask(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("projectID"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	task, err := h.taskService.Create(
		c.Request.Context(),
		req.Title,
		req.Description,
		projectID,
		req.AssigneeID,
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidProjectID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid project id",
			})

		case errors.Is(err, service.ErrInvalidTaskTitle):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid task title",
			})

		case errors.Is(err, service.ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid user id",
			})

		case errors.Is(err, service.ErrProjectNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "project not found",
			})

		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create task",
			})
		}
		return
	}
	c.JSON(http.StatusCreated, taskToResponse(task))
}

// GetTask retrieves a task by ID.
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}
	task, err := h.taskService.GetByID(
		c.Request.Context(),
		taskID,
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidTaskID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid task id",
			})
		case errors.Is(err, service.ErrTaskNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get task",
			})
		}
		return
	}
	c.JSON(http.StatusOK, taskToResponse(task))
}

// GetProjectTasks retrieves all tasks for a project.
func (h *TaskHandler) GetProjectTasks(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("projectID"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}
	tasks, err := h.taskService.GetByProjectID(
		c.Request.Context(),
		projectID,
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidProjectID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid project id",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get project tasks",
			})
		}
		return
	}
	response := make([]TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		response = append(response, taskToResponse(task))
	}
	c.JSON(http.StatusOK, response)
}

// UpdateTask updates an existing task.
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	task := &model.Task{
		ID:          taskID,
		Title:       req.Title,
		Description: req.Description,
		AssigneeID:  req.AssigneeID,
	}
	if err := h.taskService.Update(
		c.Request.Context(),
		task,
	); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidTaskID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid task id",
			})
		case errors.Is(err, service.ErrInvalidTaskTitle):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid task title",
			})
		case errors.Is(err, service.ErrTaskNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update task",
			})
		}
		return
	}
	updatedTask, err := h.taskService.GetByID(
		c.Request.Context(),
		taskID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get updated task",
		})
		return
	}
	c.JSON(http.StatusOK, taskToResponse(updatedTask))
}

// DeleteTask removes a task by ID.
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}
	if err := h.taskService.Delete(
		c.Request.Context(),
		taskID,
	); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidTaskID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid task id",
			})
		case errors.Is(err, service.ErrTaskNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete task",
			})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

// AssignTask assigns a task to a user.
func (h *TaskHandler) AssignTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}
	var req AssignTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	if err := h.taskService.AssignToUser(
		c.Request.Context(),
		taskID,
		req.UserID,
	); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidTaskID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid task id",
			})
		case errors.Is(err, service.ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid user id",
			})
		case errors.Is(err, service.ErrTaskNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to assign task",
			})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

// UnassignTask removes user assignment from a task.
func (h *TaskHandler) UnassignTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}
	if err := h.taskService.UnassignFromUser(
		c.Request.Context(),
		taskID,
	); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidTaskID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid task id",
			})
		case errors.Is(err, service.ErrTaskNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to unassign task",
			})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
