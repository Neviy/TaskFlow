package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"taskflow/internal/model"
	"taskflow/internal/service"

	"github.com/gin-gonic/gin"
)

// CommentHandler handles task comment HTTP requests.
type CommentHandler struct {
	commentService *service.CommentService
}

// NewCommentHandler creates a comment handler backed by commentService.
func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateCommentRequest represents the payload for creating a comment.
type CreateCommentRequest struct {
	Text string `json:"text" binding:"required"`
}

// UpdateCommentRequest represents the payload for updating a comment.
type UpdateCommentRequest struct {
	Text string `json:"text" binding:"required"`
}

// CommentResponse represents comment data returned by the API.
type CommentResponse struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"task_id"`
	AuthorID  int64     `json:"author_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func toCommentResponse(comment *model.Comment) CommentResponse {
	return CommentResponse{
		ID:        comment.ID,
		TaskID:    comment.TaskID,
		AuthorID:  comment.AuthorID,
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt,
	}
}

// CreateComment creates a new comment for a task.
func (h *CommentHandler) CreateComment(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}
	userID := c.GetInt64("userID")
	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	comment, err := h.commentService.Create(
		c.Request.Context(),
		taskID,
		userID,
		req.Text,
	)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}
		if errors.Is(err, service.ErrInvalidCommentText) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid comment text",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create comment",
		})
		return
	}

	c.JSON(http.StatusCreated, toCommentResponse(comment))
}

// GetByID returns a comment by ID.
func (h *CommentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid comment id",
		})
		return
	}

	comment, err := h.commentService.GetByID(
		c.Request.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "comment not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get comment",
		})
		return
	}

	c.JSON(http.StatusOK, toCommentResponse(comment))
}

// GetByTaskID returns all comments for a task.
func (h *CommentHandler) GetByTaskID(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}

	comments, err := h.commentService.GetByTaskID(
		c.Request.Context(),
		taskID,
	)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get task comments",
		})
		return
	}

	response := make([]CommentResponse, 0, len(comments))

	for _, comment := range comments {
		response = append(response, toCommentResponse(comment))
	}

	c.JSON(http.StatusOK, response)
}

// Delete deletes a comment.
func (h *CommentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid comment id",
		})
		return
	}
	if err := h.commentService.Delete(
		c.Request.Context(),
		id,
	); err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "comment not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete comment",
		})
		return
	}
	c.Status(http.StatusNoContent)
}
