# Documentation - Go Auth Service

Complete documentation for the Go Auth Service.

## 📚 Documentation Index

### Getting Started

1. **[Quick Start Guide](QUICK_START.md)** ⭐
   - Get up and running in 10 minutes
   - Basic setup and configuration
   - First API calls
   - Common use cases

### API Reference

2. **[API Documentation](API_DOCUMENTATION.md)** 📖
   - Complete API reference
   - All endpoints documented
   - Request/response examples
   - Error codes and handling
   - Rate limiting information

### Testing

3. **[Postman Guide](POSTMAN_GUIDE.md)** 🧪
   - Postman collection setup
   - Environment configuration
   - Automated testing scripts
   - Test workflows
   - Tips and tricks

### Integration

4. **[Integration Guide](INTEGRATION_GUIDE.md)** 🔗
   - Microservices integration
   - Code examples (Go, Node.js, Python)
   - Authentication middleware
   - Best practices
   - Troubleshooting

### Database

5. **[Database Schema](DATABASE_SCHEMA.md)** 🗄️
   - Complete schema documentation
   - Table structures and relationships
   - Indexes and constraints
   - Migration guides
   - Query examples

---

## Quick Links

### For Developers

- **New to the project?** Start with [Quick Start Guide](QUICK_START.md)
- **Need API details?** Check [API Documentation](API_DOCUMENTATION.md)
- **Testing with Postman?** See [Postman Guide](POSTMAN_GUIDE.md)
- **Integrating with your service?** Read [Integration Guide](INTEGRATION_GUIDE.md)

### For DevOps

- **Deployment:** See `../DEPLOYMENT_GUIDE.md`
- **Docker:** Check `../docker-compose.yml`
- **Kubernetes:** Review `../k8s/` directory
- **Database Schema:** See [Database Schema](DATABASE_SCHEMA.md)
- **Monitoring:** See Integration Guide monitoring section

### For QA/Testers

- **API Testing:** Use [Postman Guide](POSTMAN_GUIDE.md)
- **Test Suite:** Check `../tests/` directory
- **Test Documentation:** See `../TESTING_GUIDE.md`

---

## Documentation Structure

```
docs/
├── README.md                    # This file
├── QUICK_START.md              # 10-minute quick start
├── API_DOCUMENTATION.md        # Complete API reference
├── POSTMAN_GUIDE.md            # Postman testing guide
├── INTEGRATION_GUIDE.md        # Integration with other services
├── DATABASE_SCHEMA.md          # Database schema documentation
└── postman/                    # Postman collections
    └── Go_Auth_Service.postman_collection.json
```

---

## Common Tasks

### I want to...

#### ...get started quickly
→ Follow [Quick Start Guide](QUICK_START.md)

#### ...understand an API endpoint
→ Check [API Documentation](API_DOCUMENTATION.md)

#### ...test the API
→ Use [Postman Guide](POSTMAN_GUIDE.md)

#### ...integrate with my service
→ Read [Integration Guide](INTEGRATION_GUIDE.md)

#### ...deploy to production
→ See `../DEPLOYMENT_GUIDE.md`

#### ...understand the database schema
→ Check [Database Schema](DATABASE_SCHEMA.md)

#### ...run tests
→ Check `../TESTING_GUIDE.md`

---

## API Overview

### Authentication Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/auth/register` | POST | Register new user |
| `/v1/auth/login` | POST | Login and get JWT token |
| `/v1/auth/logout` | POST | Logout and blacklist token |
| `/v1/auth/validate` | POST | Validate JWT token |
| `/v1/auth/refresh` | POST | Refresh JWT token |
| `/v1/auth/token-info` | GET | Get token information |

### Health Check Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Basic health check |
| `/health/ready` | GET | Readiness probe |
| `/health/live` | GET | Liveness probe |

### Password Management

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/auth/forgot-password` | POST | Request password reset |
| `/v1/auth/reset-password` | POST | Reset password with token |
| `/v1/auth/change-password` | POST | Change password (authenticated) |

### User Management

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/api/profile` | GET | Get user profile |
| `/v1/admin/users` | GET | List all users (admin) |
| `/v1/admin/users/:id` | GET | Get user by ID (admin) |
| `/v1/admin/users/:id` | PUT | Update user (admin) |
| `/v1/admin/users/:id` | DELETE | Delete user (admin) |

### Admin Operations

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/admin/dashboard` | GET | Admin dashboard |
| `/v1/super-admin/promote` | POST | Promote user to admin |

---

## Features

### Security
- ✅ JWT-based authentication
- ✅ Password hashing with bcrypt
- ✅ Token blacklisting on logout
- ✅ Role-based access control (RBAC)
- ✅ Rate limiting per endpoint
- ✅ CORS configuration

### User Management
- ✅ User registration and login
- ✅ Password reset flow
- ✅ User profile management
- ✅ Admin user management
- ✅ Role promotion system

### Microservices
- ✅ Token validation endpoint
- ✅ Token refresh mechanism
- ✅ Request ID tracking
- ✅ Structured logging
- ✅ Health check endpoints

### Infrastructure
- ✅ Docker support
- ✅ Kubernetes manifests
- ✅ Database migrations
- ✅ Environment configuration
- ✅ Horizontal pod autoscaling

---

## Architecture

```
┌─────────────┐
│   Client    │
│  (Browser)  │
└──────┬──────┘
       │
       ▼
┌─────────────────┐
│  Auth Service   │
│   (Port 8080)   │
├─────────────────┤
│  • Register     │
│  • Login        │
│  • Validate     │
│  • Manage Users │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   MySQL DB      │
│  • Users        │
│  • Tokens       │
│  • Reset Tokens │
└─────────────────┘
```

### Microservices Architecture

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│Service A│     │Service B│     │Service C│
└────┬────┘     └────┬────┘     └────┬────┘
     │               │               │
     └───────────────┼───────────────┘
                     │
                     ▼
            ┌────────────────┐
            │  Auth Service  │
            │  (Validation)  │
            └────────────────┘
```

---

## Response Format

All API responses follow this standard format:

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "request_id": "uuid",
  "timestamp": "ISO8601"
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description"
  },
  "request_id": "uuid",
  "timestamp": "ISO8601"
}
```

---

## Authentication Flow

### User Registration and Login

```
1. POST /v1/auth/register
   ↓
2. POST /v1/auth/login
   ↓ (receive JWT token)
3. Use token in Authorization header
   ↓
4. Access protected endpoints
   ↓
5. POST /v1/auth/logout (when done)
```

### Token Validation (Microservices)

```
Client Request
   ↓
Service A (receives request with token)
   ↓
POST /v1/auth/validate (validate token)
   ↓
Check response.data.valid
   ↓
If valid: Process request
If invalid: Return 401
```

---

## Environment Variables

### Required

```env
APP_PORT=8080
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=your_password
MYSQL_DB=auth_service
JWT_SECRET=your-secret-key-min-32-chars
SUPER_ADMIN_CODE=your-super-admin-code
```

### Optional

```env
APP_ENV=development
CORS_ALLOWED_ORIGINS=http://localhost:3000
LOG_LEVEL=info
```

---

## Rate Limits

| Endpoint | Limit |
|----------|-------|
| POST /v1/auth/login | 5 requests / 15 minutes |
| POST /v1/auth/register | 3 requests / hour |
| POST /v1/auth/validate | 1000 requests / minute |
| Other endpoints | 100 requests / minute |

---

## Token Expiry

| Role | Expiry |
|------|--------|
| User | 7 days |
| Admin | 7 days |
| Super Admin | 30 days |

---

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| BAD_REQUEST | 400 | Invalid request |
| UNAUTHORIZED | 401 | Invalid/missing auth |
| FORBIDDEN | 403 | Insufficient permissions |
| NOT_FOUND | 404 | Resource not found |
| CONFLICT | 409 | Resource exists |
| RATE_LIMIT_EXCEEDED | 429 | Too many requests |
| INTERNAL_ERROR | 500 | Server error |

---

## Support

### Documentation Issues
- Found a typo? Submit a PR
- Missing information? Open an issue
- Need clarification? Ask in discussions

### Technical Issues
- Check troubleshooting sections
- Review error codes
- Check service logs
- Verify configuration

### Integration Help
- Review [Integration Guide](INTEGRATION_GUIDE.md)
- Check code examples
- Test with Postman first
- Verify token format

---

## Contributing

To improve documentation:

1. Fork the repository
2. Update documentation files
3. Test examples and code snippets
4. Submit pull request

### Documentation Standards

- Use clear, concise language
- Include code examples
- Add request/response samples
- Update table of contents
- Test all commands and code

---

## Version History

- **v1.0.0** - Initial release
  - Complete API documentation
  - Postman guide
  - Integration guide
  - Quick start guide

---

## License

This project is licensed under the MIT License.

---

## Additional Resources

### External Links
- [JWT.io](https://jwt.io/) - JWT debugger
- [Postman](https://www.postman.com/) - API testing tool
- [Go Documentation](https://golang.org/doc/) - Go language docs

### Related Documentation
- `../README.md` - Project README
- `../DEPLOYMENT_GUIDE.md` - Deployment guide
- `../TESTING_GUIDE.md` - Testing guide
- `../BUILD_SUCCESS.md` - Build information

---

**Happy Coding! 🚀**
