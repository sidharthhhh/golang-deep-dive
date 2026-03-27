-- Create token blacklist table
CREATE TABLE IF NOT EXISTS token_blacklist (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    token_jti CHAR(36) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    blacklisted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    INDEX idx_token_jti (token_jti),
    INDEX idx_user_id (user_id),
    INDEX idx_expires_at (expires_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);