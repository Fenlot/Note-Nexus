# NoteNexus

A modern, collaborative note-taking and task management platform designed for teams and individuals who want to organize their ideas, track tasks, and manage workspaces efficiently.

## Overview

NoteNexus is a full-stack application built with a Go backend and a .NET Blazor WebAssembly frontend. It provides a seamless experience for creating and managing notes, organizing tasks, and collaborating within shared workspaces.

**Try the Live Application:** https://note-nexusgh.onrender.com

## Key Features

- User authentication with secure credential handling
- Multi-workspace support with role-based access control
- Note creation and management with real-time updates
- Task and todo list management with completion tracking
- Workspace-scoped data isolation for security and organization
- RESTful API with comprehensive endpoint coverage
- Cross-origin resource sharing (CORS) support for distributed deployments

## Project Structure

```
NoteNexus/
├── note-nexus-backend/      # Go backend API server
├── note-nexus-client/       # .NET Blazor WebAssembly frontend
├── Docker-compose.yml       # Container orchestration for local development
├── LICENSE
└── README.md
```

## Technology Stack

### Backend
- Language: Go 1.25.1
- Framework: Chi v5 (lightweight HTTP router)
- Database: PostgreSQL
- Authentication: JWT tokens
- Containerization: Docker

### Frontend
- Framework: .NET 9 with Blazor WebAssembly
- Language: C#
- UI Components: Custom Razor components
- Storage: Browser local storage for client state
- Containerization: Docker with Nginx

### Infrastructure
- Container Orchestration: Docker Compose
- Network: Custom bridge network for service communication
- Port Mapping: Backend on 8080, Frontend on 80

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.25.1 (for local backend development)
- .NET 9 SDK (for local frontend development)
- PostgreSQL 15+ (if running without Docker)

### Quick Start with Docker Compose

1. Clone the repository:
   ```bash
   git clone https://github.com/Fenlot/Note-Nexus.git
   cd Note-Nexus
   ```

2. Create the required SQLite database file for the backend:
   ```bash
   touch note-nexus-backend/note.db
   ```

3. Build and start all services:
   ```bash
   docker-compose up --build
   ```

4. Access the application:
   - Frontend: http://localhost
   - Backend Health Check: http://localhost:8080/health

### Production URLs

- Frontend: https://note-nexusgh.onrender.com
- Backend API: https://note-nexus-anxm.onrender.com
- Backend Health Check: https://note-nexus-anxm.onrender.com/health

### Local Development

For detailed setup instructions for each component, refer to:
- [Backend README](note-nexus-backend/README.md) for API development
- [Frontend README](note-nexus-client/README.md) for UI development

## API Overview

The backend exposes a RESTful API with the following main endpoint groups:

### Authentication
- `POST /v1/signup` - Create a new user account
- `POST /v1/login` - Authenticate and receive JWT token

### Workspaces (Protected)
- `GET /v1/workspaces` - List all workspaces for the authenticated user

### Notes (Workspace-scoped)
- `POST /v1/workspaces/{workspaceId}/notes` - Create a new note
- `GET /v1/workspaces/{workspaceId}/notes` - List notes in workspace
- `PUT /v1/workspaces/{workspaceId}/notes/{id}` - Update a note
- `DELETE /v1/workspaces/{workspaceId}/notes/{id}` - Delete a note

### Tasks/Todos (Workspace-scoped)
- `POST /v1/workspaces/{workspaceId}/todos` - Create a new task
- `GET /v1/workspaces/{workspaceId}/todos` - List tasks in workspace
- `PATCH /v1/workspaces/{workspaceId}/todos/{id}` - Update task status
- `PUT /v1/workspaces/{workspaceId}/todos/{id}` - Update task content
- `DELETE /v1/workspaces/{workspaceId}/todos/{id}` - Delete a task

## Architecture

### Backend Architecture

The backend follows a clean architecture pattern with clear separation of concerns:

```
cmd/api/           - HTTP handlers and routing layer
internal/data/     - Database models and data access layer
migrations/        - Database schema and migrations
```

Each handler manages its own HTTP concerns while delegating business logic to the data models. JWT authentication middleware protects routes that require user authentication.

### Frontend Architecture

The frontend is built as a single-page application with component-based UI:

```
Pages/             - Top-level page components (Login, Signup, Notepad, Todo, Home)
Components/        - Reusable UI components (widgets, buttons, toggles)
Services/          - Application services (API calls, authentication, state management)
Layout/            - Layout components and styling
Models/            - Data transfer objects and type definitions
```

## Database Schema

The database consists of five main tables:

- **users**: Stores user account information and authentication credentials
- **workspaces**: Represents collaborative spaces owned by users
- **workspace_members**: Maps users to workspaces with role assignments
- **notes**: User-created notes scoped to workspaces
- **tasks**: Todo items scoped to workspaces with completion tracking

All tables support cascading deletes to maintain referential integrity.

## Security Considerations

- Passwords are hashed before storage (no plain text passwords)
- JWT tokens expire after a configured period
- All workspace-scoped operations verify user membership before granting access
- CORS headers are explicitly configured to allow only trusted origins
- Authorization middleware validates user roles within workspaces

## Development Workflow

### Adding a New Feature

1. Define the database schema changes in a new migration file
2. Implement the data model methods in the backend
3. Create the HTTP handler and route in the backend API
4. Implement the frontend service to call the new endpoint
5. Create Razor components to display the feature in the UI
6. Test the complete flow through Docker Compose

### Testing

Backend tests should cover:
- Database operations for correctness
- Authorization middleware for access control
- API responses for proper HTTP status codes

Frontend tests should cover:
- Component rendering and user interactions
- Service layer for API integration
- Local storage persistence

## Deployment

Both services are containerized and ready for deployment to cloud platforms. The Docker images can be pushed to container registries and deployed to services like:

- Docker Swarm
- Kubernetes
- AWS ECS
- Azure Container Instances
- Render or other container hosting platforms

## Contributing

When reviewing or contributing to this project:

1. Follow the existing code style and patterns
2. Keep commits focused on single features or fixes
3. Write clear commit messages describing the changes
4. Test your changes thoroughly before submitting
5. Update documentation as needed
6. Ensure all CORS origins are reviewed for security

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Support and Documentation

For more detailed information, please see:
- Backend documentation in the backend README
- Frontend documentation in the frontend README
- Database schema details in migrations/001_create_tables.sql
- API endpoint details in each service README

## Contact

For questions or feedback about the project, please open an issue on the GitHub repository.
