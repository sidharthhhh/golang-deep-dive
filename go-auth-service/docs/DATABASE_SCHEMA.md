# Database Schema Documentation - Go Auth Service

Complete database schema documentation for the Go Auth Service.

## Table of Contents

1. [Overview](#overview)
2. [Database Tables](#database-tables)
3. [Entity Relationships](#entity-relationships)
4. [Indexes](#indexes)
5. [Migrations](#migrations)
6. [Data Types](#data-types)

---

## Overview

The Go Auth Service uses MySQL 8.0+ as its primary database. The schema is designed to support:

- User authentication and authorization
- JWT token management
- Token blacklisting for logout
- Password reset functionality
- Role-based access control (RBAC)

### Database Information

- **Database Name:** `auth_service`
- **Character Set:** `utf8mb4`
- **Collation:** `utf8mb4_unicode_ci`
- **Engine:** InnoDB
- **Total Tables:** 4

---

## Database Tables

### 1. users

Stores user account information.

**Table Name:** `users`

**Columns:**

| Column | Type | Null | Default | Description |
|--------|------|------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key, unique user identifier |
| email | VARCHAR(255) | NO | - | User email address (unique) |
| password_hash | VARCHAR(255) | NO | - | Bcrypt hashed password |
| role | ENUM('user', 'admin', 'super_admin') | NO | 'user' | User role for RBAC |
| is_verified | TINYINT(1) | NO | 0 | Email verification status |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Account creation timestamp |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP ON UPDATE | Last update timestamp |

**Indexes:**

- PRIMARY KEY: `id`
- UNIQUE KEY: `email`
- INDEX: `role` (for role-based queries)

**Constraints:**

- `email` must be unique
- `password_hash` minimum length enforced by application
- `role` restricted to enum values

**SQL Definition:**

```sql
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('user', 'admin', 'super_admin') NOT NULL DEFAULT 'user',
    is_verified TINYINT(1) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_role (role),
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**Sample Data:**

```sql
-- Regular user
INSERT INTO users (email, password_hash, role, is_verified) 
VALUES ('user@example.com', '$2a$12$...', 'user', 0);

-- Admin user
INSERT INTO users (email, password_hash, role, is_verified) 
VALUES ('admin@example.com', '$2a$12$...', 'admin', 1);

-- Super admin
INSERT INTO users (email, password_hash, role, is_verified) 
VALUES ('superadmin@example.com', '$2a$12$...', 'super_admin', 1);
```

---

### 2. token_blacklist

Stores blacklisted JWT tokens (for logout functionality).

**Table Name:** `token_blacklist`

**Columns:**

| Column | Type | Null | Default | Description |
|--------|------|------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| token_jti | VARCHAR(255) | NO | - | JWT ID (jti claim) |
| user_id | BIGINT | NO | - | User who owns the token |
| expires_at | DATETIME | NO | - | Token expiration time |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Blacklist timestamp |

**Indexes:**

- PRIMARY KEY: `id`
- UNIQUE KEY: `token_jti`
- INDEX: `user_id` (foreign key reference)
- INDEX: `expires_at` (for cleanup queries)

**Constraints:**

- `token_jti` must be unique
- `user_id` references `users.id`
- `expires_at` must be in the future when created

**SQL Definition:**

```sql
CREATE TABLE token_blacklist (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    token_jti VARCHAR(255) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_token_jti (token_jti),
    INDEX idx_user_id (user_id),
    INDEX idx_expires_at (expires_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**Sample Data:**

```sql
INSERT INTO token_blacklist (token_jti, user_id, expires_at) 
VALUES ('550e8400-e29b-41d4-a716-446655440000', 1, '2024-01-15 12:00:00');
```

**Cleanup Query:**

```sql
-- Remove expired tokens (run periodically)
DELETE FROM token_blacklist WHERE expires_at <= NOW();
```

---

### 3. password_reset_tokens

Stores password reset tokens.

**Table Name:** `password_reset_tokens`

**Columns:**

| Column | Type | Null | Default | Description |
|--------|------|------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| user_id | BIGINT | NO | - | User requesting reset |
| token | VARCHAR(255) | NO | - | Reset token (hashed) |
| expires_at | DATETIME | NO | - | Token expiration (1 hour) |
| used | TINYINT(1) | NO | 0 | Whether token has been used |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Token creation timestamp |

**Indexes:**

- PRIMARY KEY: `id`
- UNIQUE KEY: `token`
- INDEX: `user_id` (foreign key reference)
- INDEX: `expires_at` (for cleanup queries)

**Constraints:**

- `token` must be unique
- `user_id` references `users.id`
- `expires_at` typically set to 1 hour from creation
- `used` prevents token reuse

**SQL Definition:**

```sql
CREATE TABLE password_reset_tokens (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    used TINYINT(1) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_token (token),
    INDEX idx_user_id (user_id),
    INDEX idx_expires_at (expires_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**Sample Data:**

```sql
INSERT INTO password_reset_tokens (user_id, token, expires_at, used) 
VALUES (1, 'abc123def456...', DATE_ADD(NOW(), INTERVAL 1 HOUR), 0);
```

**Cleanup Query:**

```sql
-- Remove expired or used tokens
DELETE FROM password_reset_tokens 
WHERE expires_at <= NOW() OR used = 1;
```

---

### 4. refresh_tokens (Optional - Future Enhancement)

Stores refresh tokens for token refresh functionality.

**Table Name:** `refresh_tokens`

**Columns:**

| Column | Type | Null | Default | Description |
|--------|------|------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| user_id | BIGINT | NO | - | User who owns the token |
| token | VARCHAR(255) | NO | - | Refresh token (hashed) |
| expires_at | DATETIME | NO | - | Token expiration |
| revoked | TINYINT(1) | NO | 0 | Whether token is revoked |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | Token creation timestamp |

**Indexes:**

- PRIMARY KEY: `id`
- UNIQUE KEY: `token`
- INDEX: `user_id` (foreign key reference)
- INDEX: `expires_at` (for cleanup queries)

**SQL Definition:**

```sql
CREATE TABLE refresh_tokens (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    revoked TINYINT(1) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_token (token),
    INDEX idx_user_id (user_id),
    INDEX idx_expires_at (expires_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

---

## Entity Relationships

### ER Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│                         users                               │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ id (PK)                                              │  │
│  │ email (UNIQUE)                                       │  │
│  │ password_hash                                        │  │
│  │ role (ENUM)                                          │  │
│  │ is_verified                                          │  │
│  │ created_at                                           │  │
│  │ updated_at                                           │  │
│  └──────────────────────────────────────────────────────┘  │
│                          │                                  │
│                          │                                  │
│         ┌────────────────┼────────────────┐                │
│         │                │                │                │
│         ▼                ▼                ▼                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │token_       │  │password_    │  │refresh_     │       │
│  │blacklist    │  │reset_tokens │  │tokens       │       │
│  ├─────────────┤  ├─────────────┤  ├─────────────┤       │
│  │id (PK)      │  │id (PK)      │  │id (PK)      │       │
│  │token_jti    │  │user_id (FK) │  │user_id (FK) │       │
│  │user_id (FK) │  │token        │  │token        │       │
│  │expires_at   │  │expires_at   │  │expires_at   │       │
│  │created_at   │  │used         │  │revoked      │       │
│  └─────────────┘  │created_at   │  │created_at   │       │
│                   └─────────────┘  └─────────────┘       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Relationships

1. **users → token_blacklist** (One-to-Many)
   - One user can have multiple blacklisted tokens
   - Foreign Key: `token_blacklist.user_id` → `users.id`
   - On Delete: CASCADE

2. **users → password_reset_tokens** (One-to-Many)
   - One user can have multiple password reset tokens
   - Foreign Key: `password_reset_tokens.user_id` → `users.id`
   - On Delete: CASCADE

3. **users → refresh_tokens** (One-to-Many)
   - One user can have multiple refresh tokens
   - Foreign Key: `refresh_tokens.user_id` → `users.id`
   - On Delete: CASCADE

---

## Indexes

### Primary Indexes

| Table | Index Name | Columns | Type |
|-------|------------|---------|------|
| users | PRIMARY | id | PRIMARY KEY |
| token_blacklist | PRIMARY | id | PRIMARY KEY |
| password_reset_tokens | PRIMARY | id | PRIMARY KEY |
| refresh_tokens | PRIMARY | id | PRIMARY KEY |

### Unique Indexes

| Table | Index Name | Columns | Purpose |
|-------|------------|---------|---------|
| users | idx_email | email | Prevent duplicate emails |
| token_blacklist | idx_token_jti | token_jti | Prevent duplicate blacklist entries |
| password_reset_tokens | idx_token | token | Prevent duplicate reset tokens |
| refresh_tokens | idx_token | token | Prevent duplicate refresh tokens |

### Foreign Key Indexes

| Table | Index Name | Columns | References |
|-------|------------|---------|------------|
| token_blacklist | idx_user_id | user_id | users.id |
| password_reset_tokens | idx_user_id | user_id | users.id |
| refresh_tokens | idx_user_id | user_id | users.id |

### Performance Indexes

| Table | Index Name | Columns | Purpose |
|-------|------------|---------|---------|
| users | idx_role | role | Fast role-based queries |
| token_blacklist | idx_expires_at | expires_at | Efficient cleanup queries |
| password_reset_tokens | idx_expires_at | expires_at | Efficient cleanup queries |
| refresh_tokens | idx_expires_at | expires_at | Efficient cleanup queries |

---

## Migrations

### Migration Files

Located in `migrations/` directory:

1. **001_create_users_table.sql**
   - Creates `users` table
   - Adds indexes and constraints

2. **002_add_role_column.sql**
   - Adds `role` column to users
   - Sets default role to 'user'

3. **003_create_token_blacklist.sql**
   - Creates `token_blacklist` table
   - Adds foreign key to users

4. **004_create_password_reset_tokens.sql**
   - Creates `password_reset_tokens` table
   - Adds foreign key to users

### Running Migrations

```bash
# Run all migrations
mysql -u root -p auth_service < migrations/001_create_users_table.sql
mysql -u root -p auth_service < migrations/002_add_role_column.sql
mysql -u root -p auth_service < migrations/003_create_token_blacklist.sql
mysql -u root -p auth_service < migrations/004_create_password_reset_tokens.sql

# Or use Makefile
make migrate-up
```

### Rollback Migrations

```bash
# Rollback all migrations
make migrate-down

# Or manually
mysql -u root -p auth_service -e "DROP TABLE IF EXISTS password_reset_tokens;"
mysql -u root -p auth_service -e "DROP TABLE IF EXISTS token_blacklist;"
mysql -u root -p auth_service -e "DROP TABLE IF EXISTS users;"
```

---

## Data Types

### Column Type Reference

| Type | Size | Description | Example |
|------|------|-------------|---------|
| BIGINT | 8 bytes | Large integer for IDs | 1, 1000000 |
| VARCHAR(255) | Variable | String up to 255 chars | "user@example.com" |
| ENUM | 1-2 bytes | Predefined string values | 'user', 'admin' |
| TINYINT(1) | 1 byte | Boolean (0 or 1) | 0, 1 |
| TIMESTAMP | 4 bytes | Date and time | '2024-01-08 12:00:00' |
| DATETIME | 8 bytes | Date and time | '2024-01-08 12:00:00' |

### Why These Types?

**BIGINT for IDs:**
- Supports up to 9,223,372,036,854,775,807 records
- Future-proof for large-scale applications

**VARCHAR(255) for strings:**
- Efficient storage for variable-length strings
- 255 is standard for email addresses and tokens

**ENUM for roles:**
- Enforces valid values at database level
- More efficient than VARCHAR for fixed options

**TINYINT(1) for booleans:**
- Standard MySQL boolean representation
- Minimal storage (1 byte)

**TIMESTAMP vs DATETIME:**
- TIMESTAMP: Auto-updates, timezone-aware
- DATETIME: Fixed time, no timezone

---

## Database Queries

### Common Queries

#### User Queries

```sql
-- Find user by email
SELECT * FROM users WHERE email = 'user@example.com';

-- Get all admin users
SELECT * FROM users WHERE role IN ('admin', 'super_admin');

-- Count users by role
SELECT role, COUNT(*) as count FROM users GROUP BY role;

-- Get recently registered users
SELECT * FROM users ORDER BY created_at DESC LIMIT 10;
```

#### Token Blacklist Queries

```sql
-- Check if token is blacklisted
SELECT COUNT(*) FROM token_blacklist 
WHERE token_jti = '550e8400-e29b-41d4-a716-446655440000' 
AND expires_at > NOW();

-- Get all blacklisted tokens for a user
SELECT * FROM token_blacklist WHERE user_id = 1;

-- Cleanup expired tokens
DELETE FROM token_blacklist WHERE expires_at <= NOW();
```

#### Password Reset Queries

```sql
-- Validate reset token
SELECT user_id FROM password_reset_tokens 
WHERE token = 'abc123...' 
AND expires_at > NOW() 
AND used = 0;

-- Mark token as used
UPDATE password_reset_tokens 
SET used = 1 
WHERE token = 'abc123...';

-- Cleanup old tokens
DELETE FROM password_reset_tokens 
WHERE expires_at <= NOW() OR used = 1;
```

---

## Performance Optimization

### Recommended Indexes

All recommended indexes are already created in migrations.

### Query Optimization Tips

1. **Use indexes for WHERE clauses:**
   ```sql
   -- Good (uses index)
   SELECT * FROM users WHERE email = 'user@example.com';
   
   -- Bad (no index)
   SELECT * FROM users WHERE LOWER(email) = 'user@example.com';
   ```

2. **Limit result sets:**
   ```sql
   -- Good
   SELECT * FROM users ORDER BY created_at DESC LIMIT 10;
   
   -- Bad (returns all rows)
   SELECT * FROM users ORDER BY created_at DESC;
   ```

3. **Use EXPLAIN to analyze queries:**
   ```sql
   EXPLAIN SELECT * FROM users WHERE role = 'admin';
   ```

### Maintenance Tasks

```sql
-- Analyze tables
ANALYZE TABLE users, token_blacklist, password_reset_tokens;

-- Optimize tables
OPTIMIZE TABLE users, token_blacklist, password_reset_tokens;

-- Check table status
SHOW TABLE STATUS LIKE 'users';
```

---

## Backup and Restore

### Backup Database

```bash
# Full backup
mysqldump -u root -p auth_service > backup_$(date +%Y%m%d).sql

# Backup specific tables
mysqldump -u root -p auth_service users token_blacklist > backup_users.sql

# Compressed backup
mysqldump -u root -p auth_service | gzip > backup_$(date +%Y%m%d).sql.gz
```

### Restore Database

```bash
# Restore from backup
mysql -u root -p auth_service < backup_20240108.sql

# Restore from compressed backup
gunzip < backup_20240108.sql.gz | mysql -u root -p auth_service
```

---

## Security Considerations

### Password Storage

- Passwords are hashed using bcrypt with cost factor 12
- Never store plain text passwords
- Password hashes are 60 characters long

### Token Security

- JWT tokens include expiration time
- Blacklisted tokens cannot be reused
- Token JTI (JWT ID) is unique per token

### Data Protection

- Use SSL/TLS for database connections
- Restrict database user permissions
- Regular security audits
- Backup sensitive data

---

## Monitoring

### Key Metrics to Monitor

```sql
-- Total users
SELECT COUNT(*) FROM users;

-- Users by role
SELECT role, COUNT(*) as count FROM users GROUP BY role;

-- Blacklisted tokens count
SELECT COUNT(*) FROM token_blacklist WHERE expires_at > NOW();

-- Active password reset requests
SELECT COUNT(*) FROM password_reset_tokens 
WHERE expires_at > NOW() AND used = 0;

-- Database size
SELECT 
    table_name,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS size_mb
FROM information_schema.TABLES
WHERE table_schema = 'auth_service';
```

---

## Troubleshooting

### Common Issues

#### Issue: Duplicate email error

```sql
-- Check for existing email
SELECT * FROM users WHERE email = 'user@example.com';

-- Solution: Use different email or update existing user
```

#### Issue: Foreign key constraint fails

```sql
-- Check if user exists
SELECT * FROM users WHERE id = 1;

-- Solution: Ensure user exists before creating related records
```

#### Issue: Token not found in blacklist

```sql
-- Check token format
SELECT * FROM token_blacklist WHERE token_jti = 'your-token-jti';

-- Solution: Verify token JTI is correct
```

---

## Schema Version

- **Current Version:** 1.0.0
- **Last Updated:** 2024-01-08
- **Migration Count:** 4
- **Tables:** 4 (3 active, 1 optional)

---

## Support

For schema-related questions:
- Review migration files in `migrations/`
- Check model definitions in `internal/models/`
- Review repository code in `internal/repository/`
- Consult MySQL documentation

---

## Future Enhancements

### Planned Tables

1. **audit_logs** - Track user actions
2. **sessions** - Active user sessions
3. **oauth_providers** - OAuth integration
4. **api_keys** - API key management
5. **permissions** - Granular permissions

### Planned Indexes

1. Composite indexes for common query patterns
2. Full-text search indexes for email
3. Partitioning for large tables

---

**Schema Status:** ✅ PRODUCTION READY

All tables are properly indexed, constrained, and optimized for production use.
