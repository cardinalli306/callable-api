A robust RESTful API built in Go using the Gin framework, with integrated Swagger documentation and a full suite of automated tests.

Overview
The Callable API provides data management endpoints with token authentication, input validation, and interactive documentation capabilities.
This API demonstrates best practices for Go development, including:

Modular and decoupled architecture
Efficient routing system with Gin
Robust data validation
Token authentication
Integrated Swagger documentation
Comprehensive unit and integration tests

Project Structure:

callable-api/
├── main.go # Application entry point
├── handlers.go # Endpoint handlers
├── middleware.go # Middlewares (authentication, logging)
├── models.go # Data model definition
├── *_test.go # Test files
└── docs/ # Generated Swagger documentation

API Endpoints:

Health and Monitoring
GET /health - Checks the availability of the API

Data Management
GET /api/v1/data - Retrieves a list of all items
GET /api/v1/data/{id} - Retrieves a specific item by ID
POST /api/v1/data - Creates a new item (requires authentication)

Requirements:

Go 1.20 or higher
Git (to clone the repository)

Installation and Execution:

1. Clone the repository
git clone https://github.com/cardinalli306/callable-api.git
cd callable-api

2. Install the dependencies
go mod download

3. Run the application
go run .

The API will be available at http://localhost:8080. Access the Swagger documentation at http://localhost:8080/swagger/index.html.

Testing the Endpoints:

Health Check
curl -X GET http://localhost:8080/health
Response:
{
"status": "success",
"message": "API working normally",
"data": null
}

Get List of Data
curl -X GET http://localhost:8080/api/v1/data

Get Item by ID
curl -X GET http://localhost:8080/api/v1/data/123

Create New Item (requires authentication)
curl -X POST http://localhost:8080/api/v1/data \
-H "Content-Type: application/json" \
-H "Authorization: Bearer api-token-123" \
-d '{
"name": "New Item",
"value": "456DEF",
"description": "Item added via curl",
"email": "user@example.com"
}'

Note: For testing purposes, use the token "api-token-123".

Swagger Documentation:

The API comes with built-in Swagger documentation accessible at http://localhost:8080/swagger/index.html. The Swagger interface allows you to:

Explore all available endpoints
View data models and expected requests
Test endpoints directly from the browser
Examine different response codes

Data Models:

InputData
{
"name": "Item Name", // required, 3-50 characters
"value": "123ABC", // required, min 1 character
"description": "Description", // optional, max 200 characters
"email": "user@example.com", // optional, email format
"created_at": "2023-05-22T14:56:32Z" // optional, ISO8601 format
}
Response (Standard Response)
{
"status": "success|error", // operation status
"message": "Descriptive message",
"data": {} // returned data or null
}

Automated Tests:

The project includes unit and integration tests for all components. Run them with:

go test -v

The test suite includes:

Model validation
Individual endpoint testing
Authentication rule verification
Integration testing of complete flows
Input error validation

Authentication:

The API uses token-based authentication to protect sensitive endpoints. To access these endpoints:

Include the Authorization: Bearer api-token-123 header in requests
Unauthenticated endpoints return a 401 code with an error message
The token is validated by a dedicated middleware

Implemented Best Practices:

Standardized Responses: Consistent format across all responses
Robust Error Handling: Informative error messages
Input Validation: Checks on required fields and formats
Modular Middleware: Logging, authentication, and error handling
Comprehensive Documentation: Swagger for real-time documentation
Comprehensive Testing: High coverage of automated tests

Contributions
Contributions are welcome! Feel free to:

Report bugs
Suggest new features
Submit pull requests
MIT License

Developed as a RESTful API demo project in Go.
