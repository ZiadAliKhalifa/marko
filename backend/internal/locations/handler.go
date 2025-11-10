package locations

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/marko/backend/internal/auth"
	"github.com/marko/backend/internal/db"
	"github.com/marko/backend/internal/notifications"
)

// Handler handles location-related HTTP requests
type Handler struct {
	db                    *db.DB
	notificationService   *notifications.Service
}

// NewHandler creates a new locations handler
func NewHandler(database *db.DB, notificationService *notifications.Service) *Handler {
	return &Handler{
		db:                    database,
		notificationService:   notificationService,
	}
}

// RegisterRoutes registers all location-related routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	locations := router.Group("/locations")
	locations.Use(authMiddleware)
	{
		locations.POST("", h.UpdateLocation)
	}
}

// UpdateLocationRequest represents the request body for updating location
type UpdateLocationRequest struct {
	CountryCode string `json:"countryCode" binding:"required,len=2"`
	Status      string `json:"status" binding:"required,oneof=arrived left"`
}

// UpdateLocationResponse represents the response for updating location
type UpdateLocationResponse struct {
	Location *db.UserLocation `json:"location"`
	Message  string           `json:"message"`
}

// UpdateLocation updates a user's location and notifies group members
func (h *Handler) UpdateLocation(c *gin.Context) {
	user, err := auth.GetUserFromGin(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var req UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the location update
	location, err := h.db.CreateUserLocation(c.Request.Context(), user.ID, req.CountryCode, req.Status)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user location")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location"})
		return
	}

	// Get user's groups to notify members
	userGroups, err := h.db.ListUserGroups(c.Request.Context(), user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list user groups")
		// Don't fail the request, just log the error
		c.JSON(http.StatusCreated, UpdateLocationResponse{
			Location: location,
			Message:  "Location updated successfully",
		})
		return
	}

	// Notify group members
	for _, group := range userGroups {
		// Get group members (excluding the user who triggered the update)
		members, err := h.db.GetGroupMembers(c.Request.Context(), group.ID)
		if err != nil {
			log.Error().Err(err).Str("group_id", group.ID.String()).Msg("Failed to get group members")
			continue
		}

		// Create notification message
		var message string
		if req.Status == "arrived" {
			message = fmt.Sprintf("%s has arrived in %s", user.Name, req.CountryCode)
		} else {
			message = fmt.Sprintf("%s has left %s", user.Name, req.CountryCode)
		}

		// Notify each member (excluding the user who triggered the update)
		for _, member := range members {
			if member.ID == user.ID {
				continue // Don't notify the user who triggered the update
			}

			// Create notification in database
			_, err := h.db.CreateNotification(c.Request.Context(), member.ID, group.ID, message)
			if err != nil {
				log.Error().Err(err).Str("member_id", member.ID.String()).Msg("Failed to create notification")
				continue
			}

			// Send push notification (stub for now)
			if member.PushToken != nil && *member.PushToken != "" {
				h.notificationService.SendPushNotification(*member.PushToken, message)
			}
		}
	}

	c.JSON(http.StatusCreated, UpdateLocationResponse{
		Location: location,
		Message:  "Location updated and notifications sent",
	})
}