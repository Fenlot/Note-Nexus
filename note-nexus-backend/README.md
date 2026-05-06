# NoteNexus Backend API

The backend service for NoteNexus is a RESTful API built with Go and the Chi router framework. It handles user authentication, workspace management, note storage, and task tracking with a PostgreSQL database backend.

## Overview

This backend service provides the core functionality for the NoteNexus application. It manages user accounts, authentication tokens, workspace collaboration features, and all associated data operations for notes and tasks.

**Live Application:** https://note-nexusgh.onrender.com

**API Base URL (Production):** https://note-nexus-anxm.onrender.com

**API Base URL (Local):** http://localhost:8080

## Technology Stack

- Language: Go 1.25.1
- HTTP Framework: Chi v5.2.3 - lightweight and composable HTTP router
- Database: PostgreSQL
- Authentication: JWT tokens with golang-jwt/jwt
- Password Hashing: golang.org/x/crypto
- Environment Management: godotenv for configuration
- CORS Support: go-chi/cors middleware

## Project Structure

```
cmd/api/
├── main.go          - Server initialization and configuration
├── auth.go          - Authentication handlers (signup, login)
├── middleware.go    - Authentication and authorization middleware
├── notes.go         - Note management endpoints
├── tasks.go         - Task/todo management endpoints
└── workspaces.go    - Workspace management endpoints

internal/data/
├── models.go        - Data model registry and initialization
├── users.go         - User database operations
├── workspaces.go    - Workspace database operations
├── notes.go         - Note database operations
└── tasks.go         - Task database operations

migrations/
└── 001_create_tables.sql - Initial database schema

go.mod              - Go module dependencies
Dockerfile          - Container image definition
LICENSE             - License information
```

## Getting Started

### Prerequisites

- Go 1.25.1 or later
- PostgreSQL 15 or later
- Make (optional, for running build commands)

### Local Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/Fenlot/Note-Nexus.git
   cd note-nexus-backend
   ```

2. Install Go dependencies:
   ```bash
   go mod download
   ```

3. Create a `.env` file with the following configuration:
   ```bash
   DATABASE_URL=postgres://username:password@localhost:5432/notenexus
   PORT=8080
   JWT_SECRET=your-secret-key-here
   ```

4. Set up the PostgreSQL database:
   ```bash
   createdb notenexus
   psql notenexus < migrations/001_create_tables.sql
   ```

5. Run the server:
   ```bash
   go run cmd/api/main.go
   ```

6. Verify the server is running:
   ```bash
   curl http://localhost:8080/health
   ```

### Docker Setup

1. Build the Docker image:
   ```bash
   docker build -t note-nexus-backend .
   ```

2. Run the container with environment variables:
   ```bash
   docker run \
     -p 8080:8080 \
     -e DATABASE_URL=postgres://user:password@host:5432/notenexus \
     -e JWT_SECRET=your-secret-key \
     note-nexus-backend
   ```

3. For complete stack setup, use Docker Compose from the root directory:
   ```bash
   docker-compose up --build
   ```

## API Endpoints

### Authentication Endpoints

All requests and responses use JSON format.

#### Sign Up
- Method: POST
- Path: `/v1/signup`
- Request Body:
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword"
  }
  ```
- Response: 201 Created
  ```json
  {
    "id": 1,
    "email": "user@example.com",
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
  ```

#### Login
- Method: POST
- Path: `/v1/login`
- Request Body:
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword"
  }
  ```
- Response: 200 OK
  ```json
  {
    "id": 1,
    "email": "user@example.com",
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
  ```

### Protected Endpoints

All protected endpoints require the `Authorization` header with a valid JWT token:
```
Authorization: Bearer <token>
```

### Workspace Endpoints

#### List Workspaces
- Method: GET
- Path: `/v1/workspaces`
- Authentication: Required
- Response: 200 OK
  ```json
  [
    {
      "id": 1,
      "name": "My Workspace",
      "owner_id": 1,
      "subscription_tier": "free",
      "created_at": "2026-05-06T10:30:00Z"
    }
  ]
  ```

### Note Endpoints

All note endpoints are scoped to a workspace.

#### Create Note
- Method: POST
- Path: `/v1/workspaces/{workspaceId}/notes`
- Authentication: Required
- Authorization: User must be workspace member
- Request Body:
  ```json
  {
    "title": "Meeting Notes",
    "content": "Discussed project timeline and deliverables"
  }
  ```
- Response: 201 Created

#### List Notes
- Method: GET
- Path: `/v1/workspaces/{workspaceId}/notes`
- Authentication: Required
- Authorization: User must be workspace member
- Response: 200 OK
  ```json
  [
    {
      "id": 1,
      "workspace_id": 1,
      "user_id": 1,
      "title": "Meeting Notes",
      "content": "Discussed project timeline and deliverables",
      "updated_at": "2026-05-06T10:30:00Z"
    }
  ]
  ```

#### Update Note
- Method: PUT
- Path: `/v1/workspaces/{workspaceId}/notes/{id}`
- Authentication: Required
- Request Body:
  ```json
  {
    "title": "Updated Title",
    "content": "Updated content"
  }
  ```
- Response: 200 OK

#### Delete Note
- Method: DELETE
- Path: `/v1/workspaces/{workspaceId}/notes/{id}`
- Authentication: Required
- Response: 204 No Content

### Task/Todo Endpoints

All task endpoints are scoped to a workspace.

#### Create Task
- Method: POST
- Path: `/v1/workspaces/{workspaceId}/todos`
- Authentication: Required
- Request Body:
  ```json
  {
    "title": "Complete project documentation"
  }
  ```
- Response: 201 Created

#### List Tasks
- Method: GET
- Path: `/v1/workspaces/{workspaceId}/todos`
- Authentication: Required
- Response: 200 OK
  ```json
  [
    {
      "id": 1,
      "workspace_id": 1,
      "user_id": 1,
      "title": "Complete project documentation",
      "is_completed": false,
      "created_at": "2026-05-06T10:30:00Z"
    }
  ]
  ```

#### Update Task Status
- Method: PATCH
- Path: `/v1/workspaces/{workspaceId}/todos/{id}`
- Authentication: Required
- Request Body:
  ```json
  {
    "is_completed": true
  }
  ```
- Response: 200 OK

#### Update Task Content
- Method: PUT
- Path: `/v1/workspaces/{workspaceId}/todos/{id}`
- Authentication: Required
- Request Body:
  ```json
  {
    "title": "Updated task title"
  }
  ```
- Response: 200 OK

#### Delete Task
- Method: DELETE
- Path: `/v1/workspaces/{workspaceId}/todos/{id}`
- Authentication: Required
- Response: 204 No Content

## Data Models

### User Model
- id: Integer, primary key
- email: String, unique, required
- password_hash: String, hashed password
- created_at: Timestamp, auto-set

### Workspace Model
- id: Integer, primary key
- name: String, required
- owner_id: Integer, foreign key to users
- subscription_tier: String, defaults to 'free'
- created_at: Timestamp, auto-set

### Workspace Members Model
- workspace_id: Integer, foreign key to workspaces
- user_id: Integer, foreign key to users
- role: String, defaults to 'member' (owner, admin, member)
- joined_at: Timestamp, auto-set
- Primary key: (workspace_id, user_id)

### Note Model
- id: Integer, primary key
- workspace_id: Integer, foreign key to workspaces
- user_id: Integer, foreign key to users (creator)
- title: String, required
- content: String, optional
- updated_at: Timestamp, auto-set

### Task Model
- id: Integer, primary key
- workspace_id: Integer, foreign key to workspaces
- user_id: Integer, foreign key to users (creator/assignee)
- title: String, required
- is_completed: Boolean, defaults to false
- created_at: Timestamp, auto-set

## Architecture Details

### Request/Response Flow

1. Request arrives at the Chi router with middleware chain
2. Logger middleware logs the request
3. CORS middleware handles cross-origin requests
4. If protected route: authentication middleware validates JWT
5. If workspace-scoped: authorization middleware verifies workspace membership
6. Route handler processes the request
7. Data models execute database operations
8. Handler returns JSON response

### Authentication and Authorization

JWT tokens are generated upon successful login or signup. The token contains the user ID and email. Every protected request must include a valid, non-expired JWT token in the Authorization header.

Workspace authorization is handled by checking if the user is a member of the workspace before allowing data access. This prevents users from accessing data outside their workspaces.

### Database Connections

The application uses a single PostgreSQL connection pool initialized at startup. Connection pool settings can be configured for production deployments. Migrations are applied automatically on server startup.

### Error Handling

The API returns appropriate HTTP status codes:
- 200: Success
- 201: Resource created
- 204: No content (successful delete)
- 400: Bad request (validation error)
- 401: Unauthorized (missing or invalid token)
- 403: Forbidden (insufficient permissions)
- 404: Not found
- 500: Internal server error

## Configuration

Configuration is managed through environment variables. Create a `.env` file in the backend root directory:

```bash
# Database connection string
DATABASE_URL=postgres://user:password@localhost:5432/notenexus

# Server port
PORT=8080

# JWT signing secret key
JWT_SECRET=your-secret-key-change-this-in-production

# Allowed CORS origins (comma-separated)
ALLOWED_ORIGINS=http://localhost:3000,https://example.com
```

## Development Guidelines

### Adding a New Endpoint

1. Create a handler function in the appropriate file (auth.go, notes.go, tasks.go, etc.)
2. Define the request/response structs
3. Add validation and error handling
4. Register the route in main.go
5. Write tests for the handler
6. Update this README with endpoint documentation

### Adding a New Model

1. Create a new file in internal/data/ (e.g., comments.go)
2. Define the data struct and database operations
3. Register the model in models.go
4. Create any necessary migrations
5. Import and use in handlers

### Database Migrations

Migrations are SQL files in the migrations/ directory. To add a new migration:

1. Create a new file: `migrations/002_add_feature.sql`
2. Write the SQL DDL statements
3. Update applyMigrations() to include the new migration
4. Test locally before committing

## Testing

### Unit Tests

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run tests for a specific package:
```bash
go test ./internal/data
```

### Integration Tests

For integration tests that require a database, use a test database:

```bash
TEST_DATABASE_URL=postgres://user:password@localhost:5432/notenexus_test go test ./...
```

### Manual Testing

Use curl or Postman to test endpoints:

```bash
# Health check
curl http://localhost:8080/health

# Sign up
curl -X POST http://localhost:8080/v1/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# List workspaces (using token from login)
curl http://localhost:8080/v1/workspaces \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Performance Considerations

- Use database indexes on frequently queried fields (email, workspace_id, user_id)
- Connection pooling is enabled by default
- Consider caching for frequently accessed data like workspace lists
- Monitor slow queries in production

## Security Best Practices

- Always use HTTPS in production
- Change JWT_SECRET to a strong random value in production
- Never commit .env files with secrets
- Validate all user input before database operations
- Regularly update dependencies for security patches
- Use parameterized queries to prevent SQL injection (chi/sql already handles this)

## Debugging

Enable detailed logging by running with debug flags:
```bash
go run -v cmd/api/main.go
```

Check database connection issues:
```bash
psql $DATABASE_URL -c "SELECT 1"
```

View Go module dependencies:
```bash
go list -m all
```

## Troubleshooting

### Database Connection Errors
- Verify DATABASE_URL format: postgres://user:password@host:port/dbname
- Ensure PostgreSQL server is running
- Check firewall rules allowing database access
- Verify user credentials and permissions

### JWT Token Issues
- Ensure JWT_SECRET is set and consistent across restarts
- Check token expiration time
- Verify Authorization header format: "Bearer <token>"

### CORS Errors
- Verify frontend origin is in allowed CORS origins list
- Check browser console for specific origin being rejected
- Test with curl to confirm backend is accessible

## Contributing

When contributing to the backend:

1. Follow Go conventions and idioms
2. Use meaningful variable and function names
3. Keep functions focused and testable
4. Add comments for exported functions
5. Write tests for new functionality
6. Keep the dependency list lean
7. Update this README for significant changes

## Deployment

The backend is containerized and can be deployed to:
- Docker Swarm
- Kubernetes
- AWS ECS or Fargate
- Azure Container Instances
- Google Cloud Run
- Render, Railway, or similar container hosting

For production deployments:
- Use a managed database service (AWS RDS, Azure Database, etc.)
- Set environment variables securely
- Configure proper logging and monitoring
- Use SSL certificates for HTTPS
- Set up database backups
- Monitor application performance

## License

This backend service is part of NoteNexus and is licensed under the MIT License.

## Further Reading

- Chi Framework: https://github.com/go-chi/chi
- PostgreSQL Documentation: https://www.postgresql.org/docs/
- JWT Implementation: https://github.com/golang-jwt/jwt
