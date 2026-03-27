-- Add additional performance indexes
-- This migration adds indexes for common query patterns

-- Index for user lookup by verification status
CREATE INDEX IF NOT EXISTS idx_users_is_verified ON users(is_verified);

-- Index for token blacklist cleanup queries
CREATE INDEX IF NOT EXISTS idx_token_blacklist_expires ON token_blacklist(expires_at, blacklisted_at);

-- Index for password reset token cleanup
CREATE INDEX IF NOT EXISTS idx_password_reset_expires ON password_reset_tokens(expires_at, used);
