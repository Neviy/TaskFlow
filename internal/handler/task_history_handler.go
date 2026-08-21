package handler

import (
	"errors"
	"net/http"
	"strconv"

	"taskflow/internal/model"
	"taskflow/internal/service"

	"github.com/gin-gonic/gin"
)

// TaskHistoryHandler handles task history HTTP requests.
type TaskHistoryHandler struct {
	historyService *service.TaskHistoryService
}

// NewTaskHistoryHandler creates a task history handler backed by historyService.
func NewTaskHistoryHandler(
	historyService *service.TaskHistoryService,
) *TaskHistoryHandler {
	return &TaskHistoryHandler{
		historyService: historyService,
	}
}

// TaskHistoryResponse represents a task status change returned by the API.
type TaskHistoryResponse struct {
	ID        int64            `json:"id"`
	TaskID    int64            `json:"task_id"`
	ChangedBy int64            `json:"changed_by"`
	OldStatus model.TaskStatus `json:"old_status"`
	NewStatus model.TaskStatus `json:"new_status"`
	CreatedAt string           `json:"created_at"`
}

func toTaskHistoryResponse(
	history *model.TaskHistory,
) TaskHistoryResponse {
	return TaskHistoryResponse{
		ID:        history.ID,
		TaskID:    history.TaskID,
		ChangedBy: history.ChangedBy,
		OldStatus: history.OldStatus,
		NewStatus: history.NewStatus,
		CreatedAt: history.CreatedAt.Format(
			"2006-01-02T15:04:05Z07:00",
		),
	}
}

// GetByID returns a task history record by ID.
func (h *TaskHistoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task history id",
		})
		return
	}
	history, err := h.historyService.GetByID(
		c.Request.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTaskHistoryID) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid task history id",
			})
			return
		}
		if errors.Is(err, service.ErrTaskHistoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "task history not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get task history",
		})
		return
	}
	c.JSON(
		http.StatusOK,
		toTaskHistoryResponse(history),
	)
}

// GetByTaskID returns history of a task.
func (h *TaskHistoryHandler) GetByTaskID(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}
	history, err := h.historyService.GetByTaskID(
		c.Request.Context(),
		taskID,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTaskID) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid task id",
			})
			return
		}
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get task history",
		})
		return
	}
	response := make(
		[]TaskHistoryResponse,
		0,
		len(history),
	)
	for _, item := range history {
		response = append(
			response,
			toTaskHistoryResponse(item),
		)
	}
	c.JSON(http.StatusOK, response)
}

// GetByProjectID returns history of all tasks in a project.
func (h *TaskHistoryHandler) GetByProjectID(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}
	history, err := h.historyService.GetByProjectID(
		c.Request.Context(),
		projectID,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidProjectID) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid project id",
			})
			return
		}
		if errors.Is(err, service.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "project not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get project history",
		})
		return
	}
	response := make(
		[]TaskHistoryResponse,
		0,
		len(history),
	)
	for _, item := range history {
		response = append(
			response,
			toTaskHistoryResponse(item),
		)
	}
	c.JSON(http.StatusOK, response)
}
