# **Callable API**

A robust API built in Go using the Gin framework, providing endpoints for data management with complete validation.

Go Version: *\[Versão do Go\]* Gin: *\[Versão do Gin\]* License: *\[Tipo de Licença\]* Test Coverage: *93.1%*

## **Table of Contents**

1. [Overview](https://flow.ciandt.com/chat-with-docs/c/6835ca72bfebe5983dca984c#overview)  
2. [Architecture](https://flow.ciandt.com/chat-with-docs/c/6835ca72bfebe5983dca984c#architecture)  
3. [Requirements](https://flow.ciandt.com/chat-with-docs/c/6835ca72bfebe5983dca984c#requirements)  
4. [Configuration](https://flow.ciandt.com/chat-with-docs/c/6835ca72bfebe5983dca984c#configuration)  
5. [Execution](https://flow.ciandt.com/chat-with-docs/c/6835ca72bfebe5983dca984c#execution)  
6. [API Documentation](https://flow.ciandt.com/chat-with-docs/c/6835ca72bfebe5983dca984c#api-documentation)  
7. [Authentication](https://flow.ciandt.com/chat-with-docs/c/6835ca72bfebe5983dca984c#authentication)  
8. [Endpoints](https://flow.ciandt.com/chat-with-docs/c/6835ca72bfebe5983dca984c#endpoints)  
9. [Usage Examples](https://flow.ciandt.com/chat-with-docs/c/6835ca72bfebe5983dca984c#usage-examples)  
10. [Important Notes](https://flow.ciandt.com/chat-with-docs/c/6835ca72bfebe5983dca984c#important-notes)  
11. [Tests](https://flow.ciandt.com/chat-with-docs/c/6835ca72bfebe5983dca984c#tests)

## **Overview \<a name="overview"\>\</a\>**

Callable API is a demonstration application that implements a RESTful API for data management. The project serves as an example of best practices in API development with Go, including organized project structure, data validation, Swagger documentation, structured logging, and consistent error handling.

## **Architecture \<a name="architecture"\>\</a\>**

The application follows a layered architecture with clear separation of concerns:

`callable-api/`

├── `cmd/`

│ └── `api/`

│ └── `main.go # Application entry point`

├── `docs/`

│ └── `swagger.json # Generated Swagger documentation`

├── `internal/`

│ ├── `handlers/ # HTTP request handlers`

│ ├── `middleware/ # Middleware (authentication, logging)`

│ └── `models/ # Data structure definitions`

└── `pkg/`

├── `config/ # Application configuration`

└── `logger/ # Custom logging system`

## **Requirements \<a name="requirements"\>\</a\>**

* Go 1.18 or higher  
* Internet access to download dependencies

## **Configuration \<a name="configuration"\>\</a\>**

The API can be configured through environment variables:

| Variable         | Description                                | Default Value             |
| ---------------- | ------------------------------------------ | ------------------------- |
| API\_PORT        | Port on which the API will run             | 8080                      |
| LOG\_LEVEL       | Logging level (debug, info, warn, error)   | debug                     |
| ALLOWED\_ORIGINS | Origins allowed for CORS (comma-separated) | localhost:\*,127.0.0.1:\* |
| DEMO\_API\_TOKEN | Token for demonstration authentication     | api-token-123             |

## **Execution \<a name="execution"\>\</a\>**

### **Local compilation and execution**

bash

**Copy code**

*\# Clone the repository*

git clone https://github.com/your-username/callable-api.git

cd callable-api

*\# Download dependencies*

go mod download

*\# Run the application*

go run cmd/api/main.go

### **Using Docker**

bash

**Copy code**

*\# Build the image*

docker build \-t callable-api .

*\# Run the container*

docker run \-p 8080:8080 \-e LOG\_LEVEL=info callable-api

## **API Documentation \<a name="api-documentation"\>\</a\>**

The API offers interactive Swagger documentation available at:

`http://localhost:8080/swagger/index.html`

## **Authentication \<a name="authentication"\>\</a\>**

The API uses Bearer token authentication for protected endpoints.

⚠️ Security note: For demonstration purposes, this API uses a simple static token defined in the configuration. In a production environment, it would be recommended to implement a more robust authentication system such as JWT with asymmetric keys or integration with an OAuth2 provider.

To access protected endpoints, add the authorization header:

`Authorization: Bearer api-token-123`

or

`Authorization: api-token-123`

## **Endpoints \<a name="endpoints"\>\</a\>**

### **Health Check**

GET /health

Checks if the API is working correctly.

Response:

json

**Copy code**

{

"status": "success",

"message": "API is running"

}

### **Get Data List**

GET /api/v1/data

Returns a paginated list of available items.

Query Parameters:

* page (optional): page number, default: 1  
* limit (optional): items per page (maximum 100), default: 10

Response:

json

**Copy code**

{

"status": "success",

"message": "Data retrieved successfully",

"data": \[

{

"id": "1",

"name": "Item 1",

"value": "ABC123",

"description": "Description for Item 1",

"email": "user1@example.com",

"created\_at": "2023-05-22T14:56:32Z"

},

{

"id": "2",

"name": "Item 2",

"value": "XYZ456",

"description": "Description for Item 2",

"email": "user2@example.com",

"created\_at": "2023-05-23T10:15:45Z"

}

\],

"page": 1,

"page\_size": 10,

"total\_rows": 42

}

### **Get Item by ID**

GET /api/v1/data/{id}

Returns a specific item based on the provided ID.

Path Parameters:

* id: unique identifier of the item

Response:

json

**Copy code**

{

"status": "success",

"message": "Data retrieved successfully",

"data": {

"id": "1",

"name": "Item 1",

"value": "Value-1",

"description": "Description for item 1",

"email": "user1@example.com",

"created\_at": "2023-06-01T09:30:00Z"

}

}

### **Create New Item (Protected)**

POST /api/v1/data

Adds a new item based on the provided data.

Requires Authentication: Yes

Request Body:

json

**Copy code**

{

"name": "New Item",

"value": "NEW123",

"description": "Description for new item",

"email": "new.user@example.com"

}

Validations:

* name: required, between 3 and 50 characters  
* value: required, at least 1 character  
* description: optional, maximum 200 characters  
* email: optional, must be a valid email address  
* created\_at: optional, must be in RFC3339 format

Response (201 Created):

json

**Copy code**

{

"status": "success",

"message": "Data saved successfully",

"data": {

"id": "new-generated-id",

"name": "New Item",

"value": "NEW123",

"description": "Description for new item",

"email": "new.user@example.com",

"created\_at": "2023-06-05T13:45:22Z"

}

}

## **Usage Examples \<a name="usage-examples"\>\</a\>**

### **Check API status**

bash

**Copy code**

curl http://localhost:8080/health

### **List data with pagination**

bash

**Copy code**

curl http://localhost:8080/api/v1/data?page=1&limit\=5

### **Get specific item**

bash

**Copy code**

curl http://localhost:8080/api/v1/data/1

### **Create new item (with authentication)**

bash

**Copy code**

curl \-X POST http://localhost:8080/api/v1/data \\

\-H "Authorization: Bearer api-token-123" \\

\-H "Content-Type: application/json" \\

\-d '{

"name": "Test Item",

"value": "TEST-123",

"description": "Created via API",

"email": "test@example.com"

}'

## **Important Notes \<a name="important-notes"\>\</a\>**

* Security: This API implements a simplified authentication mechanism with a static token for demonstration. In a production environment, it is recommended to implement JWT with expiration time, key rotation, and secure credential storage.  
* Data Persistence: The current implementation simulates responses and does not use a real database. In a production scenario, it would be necessary to integrate a persistence system such as PostgreSQL, MongoDB, or Redis.  
* Error Handling: The API implements a consistent error handling system with standardized responses and detailed logging to facilitate debugging.  
* CORS: The default CORS configuration only allows local access. For production environments, appropriately configure allowed origins through the ALLOWED\_ORIGINS variable.  
* Rate Limiting: This demonstration implementation does not include request rate limiting. In a production environment, consider adding this protection to prevent overload and abuse.

## **Tests \<a name="tests"\>\</a\>**

The project includes comprehensive unit tests focused on ensuring the reliability of core components.

### **Test Coverage**

Current test coverage is at 93.1%, with a focus on thoroughly testing data models, validation logic, and response structures.

### **Types of Tests Implemented**

* Structure Tests: Verify JSON serialization/deserialization and field mapping  
* Validation Tests: Ensure input data meets all business requirements  
* Method Tests: Validate utility methods like pagination helpers and data handling functions  
* Edge Case Tests: Cover scenarios with empty or invalid data

### **Running Tests**

bash

**Copy code**

*\# Run all tests*

go test \-v ./...

*\# Run tests with coverage report*

go test \-v \-cover ./...

*\# Generate detailed coverage report*

go test \-coverprofile=coverage.out ./...

go tool cover \-html=coverage.out \-o coverage.html

### **Test Design Philosophy**

The tests follow these principles:

* Isolated: Each test is independent and doesn't rely on external systems  
* Deterministic: Tests produce the same results on each run  
* Fast: Tests execute quickly to enable rapid feedback during development  
* Comprehensive: Cover both happy paths and error scenarios  
* Maintainable: Tests are well-structured and easy to update when requirements change

## **Contributions**

Contributions are welcome\! Please feel free to submit a Pull Request.

## **License**

This project is licensed under the MIT License \- see the [LICENSE](https://flow.ciandt.com/chat-with-docs/c/LICENSE) file for details.

