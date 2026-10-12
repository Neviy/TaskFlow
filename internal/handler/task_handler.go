package handler

import (
	"errors"
	"net/http"
	"strconv"

	"taskflow/internal/model"
	"taskflow/internal/service"

	"github.com/gin-gonic/gin"
)

// TaskHandler handles task requests.
type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTaskRequest holds data for creating a task.
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	AssigneeID  *int64 `json:"assignee_id"`
}

// UpdateTaskRequest holds data for updating a task.
type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	AssigneeID  *int64 `json:"assignee_id"`
}

// AssignTaskRequest holds the ID of the assigned user.
type AssignTaskRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

// TaskResponse holds task data returned by the API.
type TaskResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProjectID   int64  `json:"project_id"`
	AssigneeID  *int64 `json:"assignee_id"`
}

func taskToResponse(task *model.Task) TaskResponse {
	return TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		ProjectID:   task.ProjectID,
		AssigneeID:  task.AssigneeID,
	}
}

// CreateTask creates a task in a project.
func (h *TaskHandler) CreateTask(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	userID := c.GetInt64("userID")
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	task, err := h.taskService.Create(c.Request.Context(), req.Title, req.Description, projectID, userID, req.AssigneeID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidProjectID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if errors.Is(err, service.ErrInvalidTaskTitle) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task title"})
			return
		}
		if errors.Is(err, service.ErrInvalidUserID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		if errors.Is(err, service.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, taskToResponse(task))
}

// GetTask returns a task by ID.
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	userID := c.GetInt64("userID")
	task, err := h.taskService.GetByID(c.Request.Context(), taskID, userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTaskID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get task"})
		return
	}
	c.JSON(http.StatusOK, taskToResponse(task))
}

// GetProjectTasks returns all tasks in a project.
func (h *TaskHandler) GetProjectTasks(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	userID := c.GetInt64("userID")
	tasks, err := h.taskService.GetByProjectID(c.Request.Context(), projectID, userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidProjectID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if errors.Is(err, service.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get project tasks"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	userID := c.GetInt64("userID")
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	task := &model.Task{
		ID:          taskID,
		Title:       req.Title,
		Description: req.Description,
		AssigneeID:  req.AssigneeID,
	}
	if err := h.taskService.Update(c.Request.Context(), task, userID); err != nil {
		if errors.Is(err, service.ErrInvalidTaskID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if errors.Is(err, service.ErrInvalidTaskTitle) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task title"})
			return
		}
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}
	updatedTask, err := h.taskService.GetByID(c.Request.Context(), taskID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get updated task"})
		return
	}
	c.JSON(http.StatusOK, taskToResponse(updatedTask))
}

// DeleteTask removes a task by ID.
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	userID := c.GetInt64("userID")
	if err := h.taskService.Delete(c.Request.Context(), taskID, userID); err != nil {
		if errors.Is(err, service.ErrInvalidTaskID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}
	c.Status(http.StatusNoContent)
}

// AssignTask assigns a task to a user.
func (h *TaskHandler) AssignTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	var req AssignTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userID := c.GetInt64("userID")
	if err := h.taskService.AssignToUser(c.Request.Context(), taskID, userID, req.UserID); err != nil {
		if errors.Is(err, service.ErrInvalidTaskID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
			return
		}
		if errors.Is(err, service.ErrInvalidUserID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign task"})
		return
	}
	c.Status(http.StatusNoContent)
}

// UnassignTask removes a task's assignment.
func (h *TaskHandler) UnassignTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	userID := c.GetInt64("userID")
	if err := h.taskService.UnassignFromUser(c.Request.Context(), taskID, userID); err != nil {
		if errors.Is(err, service.ErrInvalidTaskID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unassign task"})
		return
	}
	c.Status(http.StatusNoContent)
}
