package handler

import (
	"errors"
	"net/http"
	"strconv"
	"taskflow/internal/service"

	"github.com/gin-gonic/gin"
)

// ProjectHandler handles project-related HTTP requests.
type ProjectHandler struct {
	projectService *service.ProjectService
}

// NewProjectHandler creates a new ProjectHandler instance.
func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// CreateProjectRequest represents the payload for creating a project.
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateProjectRequest represents the payload for updating a project.
type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// ProjectResponse represents project data in API responses.
type ProjectResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     int64  `json:"owner_id"`
}

// CreateProject creates a new project.
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userID := c.GetInt64("userID")
	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userID"})
		return
	}
	project, err := h.projectService.Create(c.Request.Context(), req.Name, req.Description, userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidProjectName) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project name"})
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create project"})
		return
	}
	response := ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		OwnerID:     project.OwnerID,
	}
	c.JSON(http.StatusCreated, response)
}

// GetProject retrieves a project by ID.
func (h *ProjectHandler) GetProject(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	project, err := h.projectService.GetByID(
		c.Request.Context(),
		projectID,
	)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}

		if errors.Is(err, service.ErrInvalidProjectID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get project",
		})
		return
	}

	response := ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		OwnerID:     project.OwnerID,
	}

	c.JSON(http.StatusOK, response)
}

// GetProjects retrieves all projects owned by the authenticated user.
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	userID := c.GetInt64("userID")

	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	projects, err := h.projectService.GetByOwner(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get projects",
		})
		return
	}

	response := make([]ProjectResponse, 0, len(projects))

	for _, project := range projects {
		response = append(response, ProjectResponse{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			OwnerID:     project.OwnerID,
		})
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProject updates an existing project.
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	project, err := h.projectService.GetByID(
		c.Request.Context(),
		projectID,
	)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get project",
		})
		return
	}
	project.Name = req.Name
	project.Description = req.Description
	if err := h.projectService.Update(
		c.Request.Context(),
		project,
	); err != nil {
		if errors.Is(err, service.ErrInvalidProjectName) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid project name",
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
			"error": "failed to update project",
		})
		return
	}
	response := ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		OwnerID:     project.OwnerID,
	}
	c.JSON(http.StatusOK, response)
}

// DeleteProject removes a project by ID.
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}
	if err := h.projectService.Delete(
		c.Request.Context(),
		projectID,
	); err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "project not found",
			})
			return
		}
		if errors.Is(err, service.ErrInvalidProjectID) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid project id",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete project",
		})
		return
	}
	c.Status(http.StatusNoContent)
}
