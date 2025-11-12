# Causerie Messages Service API
## Description
This is the Messages Service API for the Causerie application, built using Go and the Echo framework. It provides endpoints for managing chat messages, including creating, retrieving, updating, and deleting messages.

## Getting Started
### Prerequisites
- Go 1.25.4 or higher
- Git

### Installation
1. Clone the repository
2. Navigate to the project directory
3. Install dependencies using `go mod tidy`

### Running the Server
Run the server using the command:
```bash
go run server.go
```

The server will start on `http://localhost:1323`.

### Running Tests
Run the tests using the command:
```bash
go test ./...
```

### Linting
To lint the code, use the command:
```bash
go install golang.org/x/lint/golint@latest # first install golint
golint ./...
```