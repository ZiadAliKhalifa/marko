package db

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	Name      string     `json:"name" db:"name"`
	PushToken *string    `json:"push_token,omitempty" db:"push_token"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// Group represents a group that users can join
type Group struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	CreatedBy uuid.UUID  `json:"created_by" db:"created_by"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// GroupMember represents a user's membership in a group
type GroupMember struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	GroupID   uuid.UUID  `json:"group_id" db:"group_id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// UserLocation represents a user's location update
type UserLocation struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	CountryCode string     `json:"country_code" db:"country_code"`
	Status      string     `json:"status" db:"status"` // 'arrived' or 'left'
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// Notification represents a notification sent to users
type Notification struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	GroupID   uuid.UUID  `json:"group_id" db:"group_id"`
	Message   string     `json:"message" db:"message"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}