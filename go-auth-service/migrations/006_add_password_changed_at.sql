-- Add password_changed_at column to track password changes
-- Tokens issued before this timestamp will be considered invalid
-- NOTE: This column is already included in 001_create_users_table.sql
-- This migration is kept for backward compatibility with existing databases

-- Add column only if it doesn't exist (idempotent)
SET @dbname = DATABASE();
SET @tablename = 'users';
SET @columnname = 'password_changed_at';

SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = @dbname
      AND TABLE_NAME = @tablename
      AND COLUMN_NAME = @columnname
  ) > 0,
  'SELECT 1',
  CONCAT('ALTER TABLE ', @tablename, ' ADD COLUMN ', @columnname, ' TIMESTAMP NULL DEFAULT NULL AFTER updated_at')
));

PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

-- Add index only if it doesn't exist
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = @dbname
      AND TABLE_NAME = @tablename
      AND INDEX_NAME = 'idx_password_changed_at'
  ) > 0,
  'SELECT 1',
  CONCAT('CREATE INDEX idx_password_changed_at ON ', @tablename, '(password_changed_at)')
));

PREPARE indexIfNotExists FROM @preparedStatement;
EXECUTE indexIfNotExists;
DEALLOCATE PREPARE indexIfNotExists;
