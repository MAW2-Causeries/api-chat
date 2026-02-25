# Causerie Messages Service API
## Description

This is the Messages Service API for the Causerie application, built using Go and the Echo framework. It provides endpoints for managing chat messages, including creating, retrieving, updating, and deleting messages.

## Getting Started
### Prerequisites

- Go 1.25.4 or higher
- Git
- [ScyllaDB instance](#setting-up-scylladb)

### Installation

1. Clone the repository
2. Navigate to the project directory
3. Install dependencies using `go mod tidy`
4. Fill in the `.env` file with your ScyllaDB credentials based on the provided `.env.example`

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

Run the tests with coverage report using the command:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Linting

To lint the code, use the command:
```bash
go install golang.org/x/lint/golint@latest # first install golint
golint ./...
```

### Setting Up ScyllaDB
#### Docker Installation (Recommended)

To set up a ScyllaDB instance using Docker, run the following command:
```bash
docker run -d --name scylla \                                           
  -p 9042:9042 \      
  -e SCYLLA_PASSWORD=your_password \
  -e SCYLLA_USER=your_username \
  scylladb/scylla
```
You can create a keyspace and the software will automatically create the necessary tables on startup.
```sql
CREATE KEYSPACE app WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': 1 };
```

#### Local Installation

Follow the instructions on the [ScyllaDB official website](https://www.scylladb.com/download/) to install ScyllaDB locally.

## Directory structure

```shell
.
├── databases       // Database connection and initialization
├── handlers        // HTTP handlers for the API endpoints
├── models          // Data models and database interactions
├── server.go       // Main entry point of the application
├── tests           // Test files for the application
│   ├── handlers    // Tests for the handlers package
│   └── models      // Tests for the models package
└── utils           // Utility functions and helpers
```

## Collaborate
### Commit Guidelines

Use conventional commit messages to describe your changes. 
- https://www.conventionalcommits.org/en/v1.0.0/
- https://gist.github.com/qoomon/5dfcdf8eec66a051ecd85625518cfd13#file-conventional-commits-cheatsheet-md

Current commit types include:
Changes relevant to the API or UI:
- `feat`: Commits that add, adjust or remove a new feature to the API or UI
- `fix`: Commits that fix an API or UI bug of a preceded feat commit
- `refactor`: Commits that rewrite or restructure code without altering API or UI behavior
- `perf`: Commits are special type of refactor commits that specifically improve performance
- `style`: Commits that address code style (e.g., white-space, formatting, missing semi-colons) and do not affect application behavior
- `test`: Commits that add missing tests or correct existing ones
- `docs`: Commits that exclusively affect documentation
- `build`: Commits that affect build-related components such as build tools, dependencies, project version, ...
- `ops`: Commits that affect operational aspects like infrastructure (IaC), deployment scripts, CI/CD pipelines, backups, monitoring, or recovery procedures, ...
- `chore`: Commits that represent tasks like initial commit, modifying .gitignore, ...

### How to propose a new feature (issue, pull request)

1. Create an issue describing the feature you want to propose, including the problem it solves and any relevant details.
2. If you want to implement the feature yourself, create a new branch from the main branch and name it appropriately (e.g., `feature/new-feature-name`) or fork the repository if you don't have write access to the main repository.
3. Implement the feature in your branch, following the commit guidelines for any commits you make.
4. Once your implementation is complete, push your branch to the repository and create a pull request (PR) against the main branch.

### Branching Strategy

We follow a simple branching strategy:
- `main`: The main branch contains the stable codebase that is ready for production.
- `feature/*`: Feature branches are created for developing new features or making significant changes. They are merged back into the main branch once the feature is complete and has been reviewed.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
