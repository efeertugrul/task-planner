# Todo Planning Application

A task planning application that efficiently assigns tasks to developers based on their productivity and weekly capacity. The application uses a greedy heuristic algorithm inspired by the LPT (Longest Processing Time First) scheduling strategy.

## Features

- Task assignment based on developer productivity
- Weekly capacity management


## Tech Stack

### Backend
- **Go** - Main server language
- **Gin** - Web framework
- **GORM** - ORM for database operations
- **SQLite** - Database
    + **(with PostgreSQL support)**

### Frontend
- **React** - UI framework

## Algorithm

The planning algorithm uses a greedy heuristic that:
- Prioritizes tasks based on their estimated duration
- Considers developer productivity when assigning tasks
- Respects weekly capacity limits (45 hours per week)
- Balances workload across developers

This approach resembles the LPT (Longest Processing Time First) scheduling strategy, balancing tasks across developers based on their productivity and remaining weekly capacity. While not optimal, it performs well for bounded scheduling without needing LP solvers.

> **Note**: For future improvements, consider implementing an ILP-based optimal planner (e.g., with Google OR-Tools or SCIP) for more accurate planning when task volume increases.


## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd todo-planning
```

2. Install Go dependencies:
```bash
go mod download
```

3. Install frontend dependencies:
```bash
cd web
npm install
cd ..
```

## Running the Application

First fill the config.yml file then the application can be started with a single command:

```bash
./run.sh
```

This script will:
1. Initialize the database
2. Start the Go server (port 8080)
3. Start the React web app (port 3000)

The web application will be available at `http://localhost:3000`

## Project Structure

```
todo-planning/
├── cmd/
│   ├── api/         # API server
│   └── cli/         # Command line tools
├── internal/
│   ├── db/          # Database connection and migrations
│   ├── model/       # Data models
│   ├── planner/     # Planning algorithm
│   ├── provider/    # Task providers
│   └── service/     # Database operations
├── web/             # React frontend
├── run.sh           # Startup script
└── config.yml       # configurations for backend

```

## API Endpoints

- `GET /api/weekly-plan` - Get the weekly task assignments

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
# Build Go server
go build -o server cmd/api/main.go

# Build React app
cd web
npm run build
```

## CLI Usage

The application provides a command-line interface for managing tasks and developers. Here are the available commands:

### Task Management

```bash
go run cmd/cli/main.go fetch
```

### Database Management

```bash
# Initialize the database
go run cmd/cli/main.go init-db

# Force reinitialize the database (drops existing data)
go run cmd/cli/main.go init-db --force
```
