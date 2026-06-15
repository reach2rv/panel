## Project Overview

AcePanel is a new generation Linux server operations and management panel developed in Go. The project uses a decoupled frontend and backend architecture:

- Backend: Go 1.26 + go-chi router + GORM + Wire dependency injection
- Frontend: Vue 3 + Vite + Pinia + Naive UI + pnpm + xterm.js + Alova.js

## Core Principles

- **Efficiency First**: Fast unit-based development. All code comments, documentation, and replies MUST use English.
- **No Documentation**: Only write code. Do not create README, GUIDE, or other documentation files.
- **Exit Upon Completion**: Exit immediately after completing the code modification. The user will manually test it.
- **Obsession with Simplicity**: Eliminating edge cases is always better than adding condition checks. Complexity is the root of all evil.
- **Pragmatism**: Solve real problems, not hypothetical threats.
- **Shut Up**: Do not output any content unless requested by the user. Modify the code silently and exit directly.

## Build and Test

### Backend Build

Build the main program:

```bash
go build -o ace ./cmd/ace
```

Build the CLI tool:

```bash
go build -o cli ./cmd/cli
```

### Frontend Build and Test

Enter the frontend directory:

```bash
cd web
```

Install dependencies:

```bash
pnpm install
```

Development mode (with hot reload):

```bash
pnpm dev
```

Build production version:

```bash
pnpm build
```

## Code Architecture

The project adopts a DDD-like layered architecture with the dependency relationship: route -> service -> biz <- data

### Core Directory Structure

- **`cmd/`**: Program entry points
    - `ace/`: Main panel program
    - `cli/`: Command-line tool

- **`internal/app/`**: Application entry and configuration

- **`internal/route/`**: HTTP routing definition
    - Define routing rules
    - Inject required service dependencies

- **`internal/service/`**: Service layer (similar to DDD application layer)
    - Process HTTP requests/responses
    - Convert DTOs to DOs
    - Coordinate multiple biz interfaces to complete business processes
    - **Should NOT handle complex business logic**

- **`internal/biz/`**: Business logic layer (similar to DDD domain layer)
    - Define business interfaces (Repository pattern)
    - Define domain model data structures
    - Use dependency inversion principle: biz defines interfaces, data implements interfaces

- **`internal/data/`**: Data access layer (similar to DDD repository layer)
    - Implement business interfaces defined in biz
    - Encapsulate database, cache, and other operations
    - Handle data persistence logic

- **`internal/http/`**: HTTP related
    - `middleware/`: Custom middleware
    - `request/`: Request struct definitions
    - `rule/`: Custom validation rules

- **`internal/apps/`**: Panel sub-applications implementation

- **`internal/bootstrap/`**: Bootstrapping for each module

- **`internal/migration/`**: Database migrations

- **`internal/job/`**: Cron jobs

- **`internal/taskqueue/`**: Task queue runner (based on DB polling, implements `types.TaskRunner` interface)

- **`pkg/`**: Utility functions and common packages
    - Contains various independent utility modules
    - Can be referenced by any part of the project

- **`web/`**: Vue 3 frontend project

## Standard Workflow for Developing New Features

1. **Add routing in `internal/route/`**
    - Refer to existing routing files (e.g., `http.go`)
    - Inject required service dependencies
    - Define routing rules and handler mappings

2. **Implement service methods in `internal/service/`**
    - **Read existing similar services first** to understand the code style
    - Handle request validation and response formatting
    - Use `Success()` to return successful responses
    - Use `Error()` to return error responses
    - Use `ErrorSystem()` to return severe system errors
    - Call biz layer interfaces to complete business logic

3. **Define business interfaces in `internal/biz/`**
    - **Read existing similar interface definitions first**
    - Define Repository interfaces (e.g., `WebsiteRepo`)
    - Define domain model structs (e.g., `Website`)
    - Keep interfaces simple and clear

4. **Implement biz interfaces in `internal/data/`**
    - **Read existing similar implementations first**
    - Create repo structs (e.g., `websiteRepo`)
    - Implement constructors (e.g., `NewWebsiteRepo`)
    - Implement all interface methods
    - Handle database operations and caching logic

5. **Use Wire for dependency injection**
    - Add providers in the corresponding wire.go files
    - Run `go generate` to regenerate dependency injection code

## Technology Stack Specific Notes

### Helper Functions (Service Layer)

Use the following helper functions in the service layer:

- `Success(w, data)`: Return successful response
- `Error(w, statusCode, format, args...)`: Return error response
- `ErrorSystem(w, format, args...)`: Return severe system error (500)
- `Bind[T](r)`: Bind request parameters to generic type T
- `Paginate[T](...)`: Build paginated response

## Code Style

- Add explanatory comments for complex logic, do not add comments for simple logic
- Use `github.com/samber/lo` for functional programming assistance
- Strings returned to the outside by the backend should use `gotext` for translation processing as much as possible
- The frontend uses `gettext` for internationalization processing. All user-visible strings must be wrapped with `gettext` to support translation
- Do not manually edit frontend and backend translation files. The project is managed externally via Crowdin
- Frontend HTTP requests use Alova.js helper functions like `useRequest`. No need to add `onError` error handling
- Backend uses Wire dependency injection. When adding new dependencies, run `go generate ./...` to regenerate code
- No need to worry about command injection, SQL injection, file uploads, and other security issues (this is a server panel, all logged-in users are considered administrators)

## Configuration Files

Backend development configuration:

```bash
cp config.example.yml config.yml
```

Frontend development configuration:

```bash
cd web
cp .env.production .env
cp settings/proxy-config.example.ts settings/proxy-config.ts
```

## Tool Usage

For unfamiliar libraries or features, you must research using the following tools before modifying code:

1. **Check official documentation**
    - `resolve-library-id` - Resolve library name to Context7 ID
    - `get-library-docs` - Get the latest official documentation

2. **Search real code**
    - `searchGitHub` - Search GitHub for actual usage examples
