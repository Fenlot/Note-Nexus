# NoteNexus Client

The frontend for NoteNexus is a modern single-page application built with .NET 9 and Blazor WebAssembly. It provides an interactive user interface for creating notes, managing tasks, and collaborating within workspaces.

## Overview

This frontend client delivers a responsive, real-time user experience for the NoteNexus platform. Built with Blazor WebAssembly, it runs entirely in the browser while communicating with the backend API for data persistence and collaboration features.

**Try the Live Application:** https://note-nexusgh.onrender.com

## Technology Stack

- Framework: .NET 9 with Blazor WebAssembly
- Language: C# with Razor markup
- UI Styling: Custom CSS with Bootstrap components
- State Management: Local AppState service with Cascading Parameters
- Authentication: JWT-based with CustomAuthStateProvider
- HTTP Client: Custom HttpClientHandler for request configuration
- Storage: Browser LocalStorage for client-side persistence
- Server: Nginx for production builds

## Project Structure

```
Components/
├── NexusCheckbox.razor          - Reusable checkbox component
├── NexusDeleteButton.razor      - Consistent delete button component
├── NoteWidget.razor             - Display note card component
├── TodoWidget.razor             - Display todo item component
├── RedirectToLogin.razor        - Authentication guard component
└── ThemeToggle.razor            - Dark/light theme switcher

Layout/
├── MainLayout.razor             - Primary layout with navigation
├── AuthLayout.razor             - Layout for login/signup pages
├── NavMenu.razor                - Navigation menu component
└── MainLayout.razor.css         - Layout styling

Pages/
├── Home.razor                   - Dashboard landing page
├── Login.razor                  - User login page
├── Signup.razor                 - User registration page
├── Notepad.razor                - Note creation and editing
└── Todo.razor                   - Task/todo management

Services/
├── AppState.cs                  - Global application state management
├── AuthService.cs               - Authentication API operations
├── CustomAuthStateProvider.cs   - Blazor authentication provider
├── CustomHttpHandler.cs         - HTTP client configuration
├── LocalStorageService.cs       - Browser storage operations
├── NoteServices.cs              - Note API operations
├── TodoServices.cs              - Task API operations
└── WorkspaceService.cs          - Workspace API operations

Models/
├── AuthModels.cs                - Auth request/response objects
├── Note.cs                      - Note data model
├── TodoTasks.cs                 - Task data model
└── Workspace.cs                 - Workspace data model

wwwroot/
├── index.html                   - Entry point
├── css/
│   ├── app.css                  - Global application styles
│   └── dashboard-redesign.css   - Dashboard-specific styles
├── js/
│   ├── theme.js                 - Theme switching logic
│   └── app.js                   - Client-side utilities
├── lib/                         - Bootstrap and dependencies
└── images/                      - Application assets

note-nexus-client.csproj         - Project file with dependencies
Program.cs                       - Application startup configuration
App.razor                        - Root component
_Imports.razor                   - Global imports
Dockerfile                       - Container image definition
nginx.conf                       - Nginx server configuration
build.sh                         - Build script for Docker
```

## Getting Started

### Prerequisites

- .NET 9 SDK or later
- Node.js 18+ (for development tooling)
- Docker (optional, for containerized development)
- Modern web browser (Chrome, Firefox, Edge, Safari)

### Local Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/Fenlot/Note-Nexus.git
   cd note-nexus-client
   ```

2. Restore NuGet packages:
   ```bash
   dotnet restore
   ```

3. Configure the backend API endpoint. Update `Program.cs` if needed:
   ```csharp
   builder.Services.AddScoped(sp => new HttpClient 
   { 
       BaseAddress = new Uri("http://localhost:8080") 
   });
   ```
   
   For production, use: `https://note-nexus-anxm.onrender.com`

4. Run the development server:
   ```bash
   dotnet watch run
   ```

5. Open your browser and navigate to `https://localhost:5173` (or the port shown in terminal)

### Docker Setup

1. Build the Docker image:
   ```bash
   docker build -t note-nexus-client .
   ```

2. Run the container:
   ```bash
   docker run -p 80:80 note-nexus-client
   ```

3. For complete stack setup, use Docker Compose from the root directory:
   ```bash
   docker-compose up --build
   ```

## Project Configuration

### Program.cs

The startup configuration registers services and configures the Blazor WebAssembly application:

```csharp
// HTTP client with base address for API calls
builder.Services.AddScoped(sp => new HttpClient { BaseAddress = new Uri("http://localhost:8080") });

// Authentication state provider
builder.Services.AddScoped<AuthenticationStateProvider, CustomAuthStateProvider>();

// Application services
builder.Services.AddScoped<AppState>();
builder.Services.AddScoped<AuthService>();
builder.Services.AddScoped<NoteServices>();
builder.Services.AddScoped<TodoServices>();
builder.Services.AddScoped<WorkspaceService>();
```

## Component Architecture

### Page Components

#### Login.razor
Handles user authentication with email and password. On successful login, stores the JWT token and redirects to the dashboard.

#### Signup.razor
Manages user registration. Creates a new account and automatically logs the user in.

#### Home.razor
Dashboard landing page displaying workspaces, recent notes, and pending tasks. Provides quick navigation to other sections.

#### Notepad.razor
Rich interface for creating and editing notes. Displays list of notes in current workspace with inline editing capabilities.

#### Todo.razor
Task management interface with create, read, update, and delete functionality. Shows task completion status with visual indicators.

### Reusable Components

#### NoteWidget.razor
Displays individual note cards with title, preview content, and last updated timestamp. Provides quick actions for editing and deletion.

#### TodoWidget.razor
Shows individual todo items with checkbox for marking complete and delete button. Supports inline editing of task titles.

#### NexusCheckbox.razor
Custom checkbox component with consistent styling and event binding for form usage.

#### NexusDeleteButton.razor
Standardized delete button with confirmation dialog to prevent accidental deletion.

#### ThemeToggle.razor
Dark/light mode switcher that persists preference to LocalStorage.

#### RedirectToLogin.razor
Route guard component that redirects unauthenticated users to login page.

## Service Layer

### AppState.cs
Centralized state management for the entire application. Tracks current workspace, user information, and application theme.

```csharp
public class AppState
{
    public string CurrentWorkspace { get; set; }
    public string CurrentUser { get; set; }
    public string Theme { get; set; } // "light" or "dark"
    
    public event Action OnStateChange;
    public void NotifyStateChanged() => OnStateChange?.Invoke();
}
```

### AuthService.cs
Manages authentication API calls:
- `SignupAsync(email, password)` - Register new user
- `LoginAsync(email, password)` - Authenticate user
- `LogoutAsync()` - Clear authentication state

### CustomAuthStateProvider.cs
Implements Blazor's AuthenticationStateProvider. Reads JWT token from LocalStorage and maintains authentication state across page reloads.

### NoteServices.cs
CRUD operations for notes:
- `GetNotesAsync(workspaceId)` - Fetch all notes
- `CreateNoteAsync(note)` - Create new note
- `UpdateNoteAsync(note)` - Update existing note
- `DeleteNoteAsync(noteId)` - Delete note

### TodoServices.cs
Task management operations:
- `GetTasksAsync(workspaceId)` - Fetch all tasks
- `CreateTaskAsync(task)` - Create new task
- `UpdateTaskAsync(task)` - Update task
- `DeleteTaskAsync(taskId)` - Delete task

### WorkspaceService.cs
Workspace operations:
- `GetWorkspacesAsync()` - Fetch user workspaces
- `GetWorkspaceDetailsAsync(workspaceId)` - Get workspace info

### LocalStorageService.cs
Browser LocalStorage wrapper:
- `GetItemAsync(key)` - Retrieve value
- `SetItemAsync(key, value)` - Store value
- `RemoveItemAsync(key)` - Delete value

### CustomHttpHandler.cs
Configures HTTP client with:
- JWT token injection in Authorization headers
- Error handling and retry logic
- Request/response logging for debugging

## Styling

### Global Styles (app.css)
Base styling including typography, colors, spacing, and responsive breakpoints. Defines CSS custom properties for theming.

### Component-Specific Styles
Each component has an associated `.razor.css` file for scoped styling that doesn't affect other components.

### Theme System
Implements light and dark mode through CSS variables. Theme preference is persisted to LocalStorage and applied on app load.

### Bootstrap Integration
Utilizes Bootstrap 5 for responsive layout and pre-built components. Located in `wwwroot/lib/bootstrap/`.

## Authentication Flow

1. User navigates to login page
2. Enters email and password
3. AuthService calls backend `/v1/login` endpoint
4. Backend returns JWT token
5. Token is stored in LocalStorage
6. CustomAuthStateProvider reads token and updates auth state
7. AuthenticationStateProvider notifies Blazor of authentication change
8. Application redirects to dashboard
9. Future API requests include JWT in Authorization header
10. On logout, token is cleared from LocalStorage and auth state

## API Integration

All API communication goes through the services layer. Each service has methods that:
1. Call the backend API with appropriate HTTP method
2. Pass JWT token in Authorization header (handled by CustomHttpHandler)
3. Parse JSON responses into C# objects
4. Handle errors and return appropriate status

Example from NoteServices:
```csharp
public async Task<List<Note>> GetNotesAsync(int workspaceId)
{
    var response = await Http.GetAsync($"/v1/workspaces/{workspaceId}/notes");
    if (!response.IsSuccessStatusCode)
        throw new Exception($"Failed to fetch notes: {response.StatusCode}");
    
    var json = await response.Content.ReadAsStringAsync();
    return JsonSerializer.Deserialize<List<Note>>(json);
}
```

## State Management

Application state is managed through:
1. Cascading Parameters for passing data down component tree
2. AppState service for global state
3. Component parameters for local state
4. LocalStorage for persisting user preferences

State changes trigger component re-rendering automatically in Blazor.

## Responsive Design

The application is fully responsive:
- Mobile: 320px and up
- Tablet: 768px and up
- Desktop: 1024px and up

Breakpoints are defined in CSS and media queries adjust layout accordingly.

## Dark Mode Implementation

Dark mode is implemented through CSS custom properties:
```css
:root {
    --bg-color: #ffffff;
    --text-color: #000000;
    /* more colors */
}

@media (prefers-color-scheme: dark) {
    :root {
        --bg-color: #1a1a1a;
        --text-color: #ffffff;
    }
}
```

User can toggle theme manually, preference is saved to LocalStorage.

## Building for Production

1. Publish the application:
   ```bash
   dotnet publish -c Release
   ```

2. Output is in `bin/Release/net9.0/publish/wwwroot`

3. The Dockerfile handles the complete build:
   ```bash
   docker build -t note-nexus-client:prod .
   ```

## Performance Optimization

- Lazy loading of components for large pages
- Code splitting through Razor Components
- Compression configured in Nginx
- Browser caching for static assets with versioned filenames
- Minimal JavaScript payload (Blazor WebAssembly handles UI)

## Development Guidelines

### Adding a New Page

1. Create `Pages/FeatureName.razor` component
2. Add route directive: `@page "/feature"`
3. Implement the UI with Razor markup
4. Inject required services
5. Add navigation link in NavMenu.razor

### Adding a New Component

1. Create `Components/ComponentName.razor`
2. Define parameters with `[Parameter]` attributes
3. Implement event callbacks for parent communication
4. Add associated `.razor.css` file for styling
5. Export from _Imports.razor for global availability

### Adding a New Service

1. Create `Services/FeatureService.cs`
2. Inject `HttpClient` in constructor
3. Implement methods for API operations
4. Register in Program.cs: `builder.Services.AddScoped<FeatureService>()`
5. Inject in components where needed

### Form Validation

Use Blazor's built-in EditForm and DataAnnotationsValidator:
```razor
<EditForm Model="@model" OnValidSubmit="@HandleSubmit">
    <DataAnnotationsValidator />
    <ValidationSummary />
    <InputText @bind-Value="model.Name" />
    <button type="submit">Submit</button>
</EditForm>
```

## Error Handling

Global error handling through:
1. Try-catch blocks in service methods
2. Error state displayed in UI
3. Logging to browser console
4. User-friendly error messages

## Testing

### Unit Tests
Create `*.Tests.cs` files alongside components:
```bash
dotnet new xunit -n NoteNexus.Client.Tests
```

### End-to-End Tests
Use Playwright or similar for browser automation testing.

### Manual Testing Checklist
- Authentication flow (signup, login, logout)
- CRUD operations (create, read, update, delete)
- Theme switching
- Responsive layout on different screen sizes
- Browser compatibility
- Error handling for failed API calls

## Debugging

### Browser Developer Tools
- Inspect elements
- View console logs
- Network requests
- LocalStorage contents
- Application state in Performance tab

### Visual Studio Debugging
1. Set breakpoints in C# code
2. Run with `dotnet watch run`
3. Breakpoints trigger in browser
4. Use Debug > Windows > Debug Output

### Logging
Enable detailed logging in Program.cs:
```csharp
builder.Logging.SetMinimumLevel(LogLevel.Debug);
```

## Troubleshooting

### API Connection Issues
- Verify backend is running on correct port
- Check CORS configuration on backend
- Inspect Network tab in browser DevTools
- Check browser console for CORS errors

### Authentication Problems
- Verify token is stored in LocalStorage
- Check token expiration time
- Clear LocalStorage and re-login
- Verify Authorization header format

### Component Not Rendering
- Check console for JavaScript errors
- Verify component is properly imported
- Check if authentication state prevents rendering
- Ensure parameters are passed correctly

### Styling Issues
- Check CSS scoped to component
- Verify Bootstrap is loaded
- Check theme variables are defined
- Clear browser cache

## Browser Support

Tested and supported on:
- Chrome 90+
- Firefox 88+
- Edge 90+
- Safari 14+

Requires WebAssembly support in browser.

## Performance Monitoring

Monitor application performance through:
- Browser DevTools Performance tab
- Network request timing
- Component render time
- Memory usage

## Contributing

When contributing to the frontend:

1. Follow C# and Razor conventions
2. Use meaningful component and variable names
3. Keep components focused and testable
4. Add comments for complex logic
5. Maintain consistent styling with existing code
6. Test on multiple browsers and screen sizes
7. Update this README for significant changes

## Deployment

The frontend is containerized and can be deployed to:
- Docker Swarm
- Kubernetes
- AWS S3 with CloudFront
- Azure Static Web Apps
- Netlify, Vercel, or similar
- Traditional web servers with Nginx

For production deployments:
- Use Release build configuration
- Enable gzip compression
- Set up CDN for static assets
- Configure security headers in Nginx
- Enable HTTPS with certificates
- Set up monitoring and error tracking

## License

This frontend is part of NoteNexus and is licensed under the MIT License.

## Further Reading

- Blazor Documentation: https://docs.microsoft.com/en-us/aspnet/core/blazor/
- .NET 9 Release Notes: https://docs.microsoft.com/en-us/dotnet/core/whats-new/dotnet-9
- Bootstrap Documentation: https://getbootstrap.com/docs/
- MDN Web Docs: https://developer.mozilla.org/
