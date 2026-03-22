-- Add role column to existing users table
ALTER TABLE users 
ADD COLUMN role ENUM('user', 'admin', 'super_admin') DEFAULT 'user' AFTER password_hash;

-- Add index on role column
CREATE INDEX idx_role ON users(role);
