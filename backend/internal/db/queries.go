package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// User queries

// CreateUser creates a new user
func (db *DB) CreateUser(ctx context.Context, email, name string) (*User, error) {
	user := &User{}
	err := db.QueryRowContext(ctx, `
		INSERT INTO users (email, name) 
		VALUES ($1, $2) 
		RETURNING id, email, name, push_token, created_at
	`, email, name).Scan(&user.ID, &user.Email, &user.Name, &user.PushToken, &user.CreatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// GetUserByID gets a user by ID
func (db *DB) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	user := &User{}
	err := db.QueryRowContext(ctx, `
		SELECT id, email, name, push_token, created_at 
		FROM users 
		WHERE id = $1
	`, userID).Scan(&user.ID, &user.Email, &user.Name, &user.PushToken, &user.CreatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpdateUserPushToken updates a user's push token
func (db *DB) UpdateUserPushToken(ctx context.Context, userID uuid.UUID, pushToken string) error {
	_, err := db.ExecContext(ctx, `
		UPDATE users 
		SET push_token = $1 
		WHERE id = $2
	`, pushToken, userID)
	
	if err != nil {
		return fmt.Errorf("failed to update push token: %w", err)
	}
	return nil
}

// Group queries

// CreateGroup creates a new group
func (db *DB) CreateGroup(ctx context.Context, name string, createdBy uuid.UUID) (*Group, error) {
	group := &Group{}
	err := db.QueryRowContext(ctx, `
		INSERT INTO groups (name, created_by) 
		VALUES ($1, $2) 
		RETURNING id, name, created_by, created_at
	`, name, createdBy).Scan(&group.ID, &group.Name, &group.CreatedBy, &group.CreatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}
	return group, nil
}

// GetGroupByID gets a group by ID
func (db *DB) GetGroupByID(ctx context.Context, groupID uuid.UUID) (*Group, error) {
	group := &Group{}
	err := db.QueryRowContext(ctx, `
		SELECT id, name, created_by, created_at 
		FROM groups 
		WHERE id = $1
	`, groupID).Scan(&group.ID, &group.Name, &group.CreatedBy, &group.CreatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	return group, nil
}

// ListUserGroups gets all groups a user is a member of
func (db *DB) ListUserGroups(ctx context.Context, userID uuid.UUID) ([]*Group, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT g.id, g.name, g.created_by, g.created_at 
		FROM groups g 
		INNER JOIN group_members gm ON g.id = gm.group_id 
		WHERE gm.user_id = $1 
		ORDER BY g.created_at DESC
	`, userID)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list user groups: %w", err)
	}
	defer rows.Close()
	
	var groups []*Group
	for rows.Next() {
		group := &Group{}
		if err := rows.Scan(&group.ID, &group.Name, &group.CreatedBy, &group.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, group)
	}
	
	return groups, nil
}

// GroupMember queries

// AddGroupMember adds a user to a group
func (db *DB) AddGroupMember(ctx context.Context, groupID, userID uuid.UUID) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO group_members (group_id, user_id) 
		VALUES ($1, $2) 
		ON CONFLICT (group_id, user_id) DO NOTHING
	`, groupID, userID)
	
	if err != nil {
		return fmt.Errorf("failed to add group member: %w", err)
	}
	return nil
}

// GetGroupMembers gets all members of a group
func (db *DB) GetGroupMembers(ctx context.Context, groupID uuid.UUID) ([]*User, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT u.id, u.email, u.name, u.push_token, u.created_at 
		FROM users u 
		INNER JOIN group_members gm ON u.id = gm.user_id 
		WHERE gm.group_id = $1
	`, groupID)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}
	defer rows.Close()
	
	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.PushToken, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	
	return users, nil
}

// UserLocation queries

// CreateUserLocation creates a new location update
func (db *DB) CreateUserLocation(ctx context.Context, userID uuid.UUID, countryCode, status string) (*UserLocation, error) {
	location := &UserLocation{}
	err := db.QueryRowContext(ctx, `
		INSERT INTO user_locations (user_id, country_code, status) 
		VALUES ($1, $2, $3) 
		RETURNING id, user_id, country_code, status, updated_at
	`, userID, countryCode, status).Scan(&location.ID, &location.UserID, &location.CountryCode, &location.Status, &location.UpdatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create user location: %w", err)
	}
	return location, nil
}

// GetLatestUserLocation gets the latest location for a user
func (db *DB) GetLatestUserLocation(ctx context.Context, userID uuid.UUID) (*UserLocation, error) {
	location := &UserLocation{}
	err := db.QueryRowContext(ctx, `
		SELECT id, user_id, country_code, status, updated_at 
		FROM user_locations 
		WHERE user_id = $1 
		ORDER BY updated_at DESC 
		LIMIT 1
	`, userID).Scan(&location.ID, &location.UserID, &location.CountryCode, &location.Status, &location.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest user location: %w", err)
	}
	return location, nil
}

// Notification queries

// CreateNotification creates a new notification
func (db *DB) CreateNotification(ctx context.Context, userID, groupID uuid.UUID, message string) (*Notification, error) {
	notification := &Notification{}
	err := db.QueryRowContext(ctx, `
		INSERT INTO notifications (user_id, group_id, message) 
		VALUES ($1, $2, $3) 
		RETURNING id, user_id, group_id, message, created_at
	`, userID, groupID, message).Scan(&notification.ID, &notification.UserID, &notification.GroupID, &notification.Message, &notification.CreatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}
	return notification, nil
}

// ListUserNotifications gets notifications for a user
func (db *DB) ListUserNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]*Notification, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT n.id, n.user_id, n.group_id, n.message, n.created_at, g.name as group_name
		FROM notifications n
		INNER JOIN groups g ON n.group_id = g.id
		WHERE n.user_id = $1 
		ORDER BY n.created_at DESC 
		LIMIT $2
	`, userID, limit)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list user notifications: %w", err)
	}
	defer rows.Close()
	
	type NotificationWithGroup struct {
		Notification
		GroupName string `json:"group_name"`
	}
	
	var notifications []*NotificationWithGroup
	for rows.Next() {
		notif := &NotificationWithGroup{}
		if err := rows.Scan(&notif.ID, &notif.UserID, &notif.GroupID, &notif.Message, &notif.CreatedAt, &notif.GroupName); err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notif)
	}
	
	// Convert to regular notifications for now (can extend later)
	result := make([]*Notification, len(notifications))
	for i, n := range notifications {
		result[i] = &n.Notification
		log.Debug().Str("group_name", n.GroupName).Msg("Notification with group info")
	}
	
	return result, nil
}