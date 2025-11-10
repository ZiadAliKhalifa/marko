package notifications

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/marko/backend/internal/auth"
	"github.com/marko/backend/internal/db"
)

// Handler handles notification-related HTTP requests
type Handler struct {
	db *db.DB
}

// NewHandler creates a new notifications handler
func NewHandler(database *db.DB) *Handler {
	return &Handler{db: database}
}

// RegisterRoutes registers all notification-related routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	notifications := router.Group("/notifications")
	notifications.Use(authMiddleware)
	{
		notifications.GET("", h.ListNotifications)
	}
}

// ListNotificationsResponse represents the response for listing notifications
type ListNotificationsResponse struct {
	Notifications []*db.Notification `json:"notifications"`
}

// ListNotifications lists notifications for the authenticated user
func (h *Handler) ListNotifications(c *gin.Context) {
	user, err := auth.GetUserFromGin(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Get limit parameter (default to 50, max 100)
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	notifications, err := h.db.ListUserNotifications(c.Request.Context(), user.ID, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list user notifications")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list notifications"})
		return
	}

	c.JSON(http.StatusOK, ListNotificationsResponse{Notifications: notifications})
}