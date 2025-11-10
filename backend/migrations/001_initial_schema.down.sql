-- Drop indexes
DROP INDEX IF EXISTS idx_notifications_created_at;
DROP INDEX IF EXISTS idx_notifications_group_id;
DROP INDEX IF EXISTS idx_notifications_user_id;
DROP INDEX IF EXISTS idx_user_locations_updated_at;
DROP INDEX IF EXISTS idx_user_locations_user_id;
DROP INDEX IF EXISTS idx_group_members_user_id;
DROP INDEX IF EXISTS idx_group_members_group_id;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS user_locations;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS users;