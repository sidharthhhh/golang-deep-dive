-- Complete Database Setup for Auth Service
-- Run this file directly in MySQL Workbench

-- Create database
CREATE DATABASE IF NOT EXISTS auth_service;
USE auth_service;

-- Drop existing tables if they exist (to avoid conflicts)
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS token_blacklist;
DROP TABLE IF EXISTS users;

-- Create users table
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('user', 'admin', 'super_admin') DEFAULT 'user',
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    password_changed_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_email (email),
    INDEX idx_role (role),
    INDEX idx_is_verified (is_verified),
    INDEX idx_password_changed_at (password_changed_at)
);

-- Create token blacklist table
CREATE TABLE token_blacklist (
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

-- Create password reset tokens table
CREATE TABLE password_reset_tokens (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    used BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_token (token),
    INDEX idx_user_id (user_id),
    INDEX idx_expires_at (expires_at),
    INDEX idx_expires_used (expires_at, used),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Verify tables were created
SHOW TABLES;
