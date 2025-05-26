# Callable API

A robust API built in Go using the Gin framework, providing endpoints for data management with complete validation.

![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8.svg)
![Gin](https://img.shields.io/badge/Framework-Gin-00ADD8.svg)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Requirements](#requirements)
- [Configuration](#configuration)
- [Execution](#execution)
- [API Documentation](#api-documentation)
  - [Authentication](#authentication)
  - [Endpoints](#endpoints)
- [Usage Examples](#usage-examples)
- [Important Notes](#important-notes)
- [Tests](#tests)

## Overview

Callable API is a demonstration application that implements a RESTful API for data management. The project serves as an example of best practices in API development with Go, including organized project structure, data validation, Swagger documentation, structured logging, and consistent error handling.

## Architecture

The application follows a layered architecture with clear separation of concerns:

```
callable-api/
├── cmd/
│   └── api/
│       └── main.go         # Application entry point
├── docs/
│   └── swagger.json        # Generated Swagger documentation
├── internal/
│   ├── handlers/           # HTTP request handlers
│   ├── middleware/         # Middleware (authentication, logging)
│   └── models/             # Data structure definitions
└── pkg/
    ├── config/             # Application configuration
    └── logger/             # Custom logging system
```

## Requirements

- Go 1.18 or higher
- Internet access to download dependencies

## Configuration

The API can be configured through environment variables:

| Variable | Description | Default Value |
|----------|-----------|--------------|
| `API_PORT` | Port on which the API will run | `8080` |
| `LOG_LEVEL` | Logging level (`debug`, `info`, `warn`, `error`) | `debug` |
| `ALLOWED_ORIGINS` | Origins allowed for CORS (comma-separated) | `localhost:*,127.0.0.1:*` |
| `DEMO_API_TOKEN` | Token for demonstration authentication | `api-token-123` |

## Execution

### Local compilation and execution

```bash
# Clone the repository
git clone https://github.com/your-username/callable-api.git
cd callable-api

# Download dependencies
go mod download

# Run the application
go run cmd/api/main.go
```

### Using Docker

```bash
# Build the image
docker build -t callable-api .

# Run the container
docker run -p 8080:8080 -e LOG_LEVEL=info callable-api
```

## API Documentation

The API offers interactive Swagger documentation available at:

```
http://localhost:8080/swagger/index.html
```

### Authentication

The API uses Bearer token authentication for protected endpoints.

**⚠️ Security note**: For demonstration purposes, this API uses a simple static token defined in the configuration. In a production environment, it would be recommended to implement a more robust authentication system such as JWT with asymmetric keys or integration with an OAuth2 provider.

To access protected endpoints, add the authorization header:

```
Authorization: Bearer api-token-123
```

or

```
Authorization: api-token-123
```

### Endpoints

#### Health Check

```
GET /health
```

Checks if the API is working correctly.

**Response**:

```json
{
  "status": "success",
  "message": "API is running"
}
```

#### Get Data List

```
GET /api/v1/data
```

Returns a paginated list of available items.

**Query Parameters**:

- `page` (optional): page number, default: 1
- `limit` (optional): items per page (maximum 100), default: 10

**Response**:

```json
{
  "status": "success",
  "message": "Data retrieved successfully",
  "data": [
    {
      "id": "1",
      "name": "Item 1",
      "value": "ABC123",
      "description": "Description for Item 1",
      "email": "user1@example.com",
      "created_at": "2023-05-22T14:56:32Z"
    },
    {
      "id": "2",
      "name": "Item 2",
      "value": "XYZ456",
      "description": "Description for Item 2",
      "email": "user2@example.com",
      "created_at": "2023-05-23T10:15:45Z"
    }
  ],
  "page": 1,
  "page_size": 10,
  "total_rows": 42
}
```

#### Get Item by ID

```
GET /api/v1/data/{id}
```

Returns a specific item based on the provided ID.

**Path Parameters**:

- `id`: unique identifier of the item

**Response**:

```json
{
  "status": "success",
  "message": "Data retrieved successfully",
  "data": {
    "id": "1",
    "name": "Item 1",
    "value": "Value-1",
    "description": "Description for item 1",
    "email": "user1@example.com",
    "created_at": "2023-06-01T09:30:00Z"
  }
}
```

#### Create New Item (Protected)

```
POST /api/v1/data
```

Adds a new item based on the provided data.

**Requires Authentication**: Yes

**Request Body**:

```json
{
  "name": "New Item",
  "value": "NEW123",
  "description": "Description for new item",
  "email": "new.user@example.com"
}
```

**Validations**:

- `name`: required, between 3 and 50 characters
- `value`: required, at least 1 character
- `description`: optional, maximum 200 characters
- `email`: optional, must be a valid email address
- `created_at`: optional, must be in RFC3339 format

**Response (201 Created)**:

```json
{
  "status": "success",
  "message": "Data saved successfully",
  "data": {
    "id": "new-generated-id",
    "name": "New Item",
    "value": "NEW123",
    "description": "Description for new item",
    "email": "new.user@example.com",
    "created_at": "2023-06-05T13:45:22Z"
  }
}
```

## Usage Examples

### Check API status

```bash
curl http://localhost:8080/health
```

### List data with pagination

```bash
curl http://localhost:8080/api/v1/data?page=1&limit=5
```

### Get specific item

```bash
curl http://localhost:8080/api/v1/data/1
```

### Create new item (with authentication)

```bash
curl -X POST http://localhost:8080/api/v1/data \
  -H "Authorization: Bearer api-token-123" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Item",
    "value": "TEST-123",
    "description": "Created via API",
    "email": "test@example.com"
  }'
```

## Important Notes

1. **Security**: This API implements a simplified authentication mechanism with a static token for demonstration. In a production environment, it is recommended to implement JWT with expiration time, key rotation, and secure credential storage.

2. **Data Persistence**: The current implementation simulates responses and does not use a real database. In a production scenario, it would be necessary to integrate a persistence system such as PostgreSQL, MongoDB, or Redis.

3. **Error Handling**: The API implements a consistent error handling system with standardized responses and detailed logging to facilitate debugging.

4. **CORS**: The default CORS configuration only allows local access. For production environments, appropriately configure allowed origins through the `ALLOWED_ORIGINS` variable.

5. **Rate Limiting**: This demonstration implementation does not include request rate limiting. In a production environment, consider adding this protection to prevent overload and abuse.

## Tests

To run the automated tests:

```bash
go test -v ./...
```

To run tests with coverage:

```bash
go test -cover ./...
```

---

## Contributions

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
