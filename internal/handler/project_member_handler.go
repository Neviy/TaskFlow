package handler

import (
	"errors"
	"net/http"
	"strconv"

	"taskflow/internal/model"
	"taskflow/internal/service"

	"github.com/gin-gonic/gin"
)

// ProjectMemberHandler handles project membership HTTP requests.
type ProjectMemberHandler struct {
	memberService *service.ProjectMemberService
}

// NewProjectMemberHandler creates a project member handler backed by memberService.
func NewProjectMemberHandler(
	memberService *service.ProjectMemberService,
) *ProjectMemberHandler {
	return &ProjectMemberHandler{
		memberService: memberService,
	}
}

// =========================
// Request / Response
// =========================

// AddMemberRequest represents the payload for adding a project member.
type AddMemberRequest struct {
	UserID int64             `json:"user_id" binding:"required"`
	Role   model.ProjectRole `json:"role" binding:"required"`
}

// UpdateMemberRoleRequest represents the payload for changing a member's role.
type UpdateMemberRoleRequest struct {
	Role model.ProjectRole `json:"role" binding:"required"`
}

// ProjectMemberResponse represents project membership data returned by the API.
type ProjectMemberResponse struct {
	ProjectID int64             `json:"project_id"`
	UserID    int64             `json:"user_id"`
	Role      model.ProjectRole `json:"role"`
}

func toProjectMemberResponse(
	member *model.ProjectMember,
) ProjectMemberResponse {
	return ProjectMemberResponse{
		ProjectID: member.ProjectID,
		UserID:    member.UserID,
		Role:      member.Role,
	}
}

// =========================
// Add member
// POST /projects/:id/members
// =========================

// AddMember adds a user to a project.
func (h *ProjectMemberHandler) AddMember(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	var req AddMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	if req.UserID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	if req.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project role",
		})
		return
	}

	err = h.memberService.AddMember(
		c.Request.Context(),
		projectID,
		req.UserID,
		req.Role,
	)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "project not found",
			})
			return
		}

		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}

		if errors.Is(err, service.ErrProjectMemberAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "user is already a project member",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to add project member",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "project member added",
	})
}

// =========================
// Get member
// GET /projects/:id/members/:userID
// =========================

// GetMember returns a specific project member.
func (h *ProjectMemberHandler) GetMember(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	member, err := h.memberService.GetMember(
		c.Request.Context(),
		projectID,
		userID,
	)
	if err != nil {
		if errors.Is(err, service.ErrProjectMemberNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "project member not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get project member",
		})
		return
	}

	c.JSON(http.StatusOK, toProjectMemberResponse(member))
}

// =========================
// List members
// GET /projects/:id/members
// =========================

// ListMembers returns all members of a project.
func (h *ProjectMemberHandler) ListMembers(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	members, err := h.memberService.ListMembers(
		c.Request.Context(),
		projectID,
	)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "project not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list project members",
		})
		return
	}

	response := make([]ProjectMemberResponse, 0, len(members))

	for _, member := range members {
		response = append(response, toProjectMemberResponse(member))
	}

	c.JSON(http.StatusOK, response)
}

// =========================
// Update role
// PATCH /projects/:id/members/:userID
// =========================

// UpdateRole changes a project member's role.
func (h *ProjectMemberHandler) UpdateRole(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	var req UpdateMemberRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	if req.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project role",
		})
		return
	}

	err = h.memberService.UpdateRole(
		c.Request.Context(),
		projectID,
		userID,
		req.Role,
	)
	if err != nil {
		if errors.Is(err, service.ErrProjectMemberNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "project member not found",
			})
			return
		}

		if errors.Is(err, service.ErrCannotChangeOwnerRole) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "cannot change owner role",
			})
			return
		}

		if errors.Is(err, service.ErrInvalidProjectRole) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid project role",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update project member role",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "project member role updated",
	})
}

// =========================
// Remove member
// DELETE /projects/:id/members/:userID
// =========================

// RemoveMember removes a user from a project.
func (h *ProjectMemberHandler) RemoveMember(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	err = h.memberService.RemoveMember(
		c.Request.Context(),
		projectID,
		userID,
	)
	if err != nil {
		if errors.Is(err, service.ErrProjectMemberNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "project member not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to remove project member",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// =========================
// Is member
// GET /projects/:id/members/:userID/check
// =========================

// IsMember reports whether a user belongs to a project.
func (h *ProjectMemberHandler) IsMember(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	isMember, err := h.memberService.IsMember(
		c.Request.Context(),
		projectID,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to check project membership",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_member": isMember,
	})
}
