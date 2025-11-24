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

### Linting
To lint the code, use the command:
```bash
go install golang.org/x/lint/golint@latest # first install golint
golint ./...
```

## Setting Up ScyllaDB
### Docker Installation (Recommended)
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

### Local Installation
Follow the instructions on the [ScyllaDB official website](https://www.scylladb.com/download/) to install ScyllaDB locally.