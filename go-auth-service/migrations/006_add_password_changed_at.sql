-- Add password_changed_at column to track password changes
-- Tokens issued before this timestamp will be considered invalid

ALTER TABLE users 
ADD COLUMN password_changed_at TIMESTAMP NULL DEFAULT NULL AFTER updated_at;

-- Add index for efficient lookups
CREATE INDEX idx_password_changed_at ON users(password_changed_at);
