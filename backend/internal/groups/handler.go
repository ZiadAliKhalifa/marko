package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/marko/backend/internal/auth"
	"github.com/marko/backend/internal/db"
)

// Handler handles group-related HTTP requests
type Handler struct {
	db *db.DB
}

// NewHandler creates a new groups handler
func NewHandler(database *db.DB) *Handler {
	return &Handler{db: database}
}

// RegisterRoutes registers all group-related routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	groups := router.Group("/groups")
	groups.Use(authMiddleware)
	{
		groups.POST("", h.CreateGroup)
		groups.GET("", h.ListUserGroups)
		groups.POST("/:id/join", h.JoinGroup)
		groups.GET("/:id/members", h.GetGroupMembers)
	}
}

// CreateGroupRequest represents the request body for creating a group
type CreateGroupRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

// CreateGroupResponse represents the response for creating a group
type CreateGroupResponse struct {
	Group *db.Group `json:"group"`
}

// CreateGroup creates a new group
func (h *Handler) CreateGroup(c *gin.Context) {
	user, err := auth.GetUserFromGin(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := h.db.CreateGroup(c.Request.Context(), req.Name, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create group")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	// Add the creator as a member
	if err := h.db.AddGroupMember(c.Request.Context(), group.ID, user.ID); err != nil {
		log.Error().Err(err).Msg("Failed to add creator as group member")
		// Don't fail the request, just log the error
	}

	c.JSON(http.StatusCreated, CreateGroupResponse{Group: group})
}

// ListUserGroupsResponse represents the response for listing user groups
type ListUserGroupsResponse struct {
	Groups []*db.Group `json:"groups"`
}

// ListUserGroups lists all groups the user is a member of
func (h *Handler) ListUserGroups(c *gin.Context) {
	user, err := auth.GetUserFromGin(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	groups, err := h.db.ListUserGroups(c.Request.Context(), user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list user groups")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list groups"})
		return
	}

	c.JSON(http.StatusOK, ListUserGroupsResponse{Groups: groups})
}

// JoinGroupRequest represents the request body for joining a group
type JoinGroupRequest struct {
	GroupID string `json:"group_id" binding:"required"`
}

// JoinGroup adds the user to a group
func (h *Handler) JoinGroup(c *gin.Context) {
	user, err := auth.GetUserFromGin(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Check if group exists
	group, err := h.db.GetGroupByID(c.Request.Context(), groupID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get group")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join group"})
		return
	}
	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Add user to group
	if err := h.db.AddGroupMember(c.Request.Context(), groupID, user.ID); err != nil {
		log.Error().Err(err).Msg("Failed to add group member")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined group"})
}

// GetGroupMembersResponse represents the response for getting group members
type GetGroupMembersResponse struct {
	Members []*db.User `json:"members"`
}

// GetGroupMembers gets all members of a group
func (h *Handler) GetGroupMembers(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Check if user is a member of the group
	user, err := auth.GetUserFromGin(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Verify user is a member of the group
	userGroups, err := h.db.ListUserGroups(c.Request.Context(), user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list user groups")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get group members"})
		return
	}

	isMember := false
	for _, ug := range userGroups {
		if ug.ID == groupID {
			isMember = true
			break
		}
	}

	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this group"})
		return
	}

	members, err := h.db.GetGroupMembers(c.Request.Context(), groupID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get group members")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get group members"})
		return
	}

	c.JSON(http.StatusOK, GetGroupMembersResponse{Members: members})
}