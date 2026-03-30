-- Add additional performance indexes
-- This migration adds indexes for common query patterns
-- Note: Base indexes (idx_email, idx_role, idx_is_verified, idx_password_changed_at) 
-- are already created in 001_create_users_table.sql

-- Index for token blacklist cleanup queries
CREATE INDEX IF NOT EXISTS idx_token_blacklist_expires ON token_blacklist(expires_at, blacklisted_at);

-- Index for password reset token cleanup
CREATE INDEX IF NOT EXISTS idx_password_reset_expires ON password_reset_tokens(expires_at, used);
