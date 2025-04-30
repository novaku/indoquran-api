# IndoQuran API

Welcome to the IndoQuran API project. This API provides access to Quranic content, including search functionality, surah listings, and detailed ayah information.

## Features

- Search Quranic content with pagination and filtering
- List all surahs or specific surah details
- Get ayahs for a specific surah with pagination
- Retrieve detailed information for a specific ayah
- Rate limiting for API endpoints
- Content Security Policy (CSP) implementation
- CORS support
- Request timeout handling
- Graceful shutdown
- Redis caching for improved performance
- Comprehensive error handling
- Request logging and monitoring

## API Endpoints

### Base URL
`https://indoquran.web.id/api/v1`

### Endpoints

#### Welcome
- `GET /`: Welcome message
  - Response: `{"message": "Welcome to indoquran.web.id API v1.0"}`

#### Search
- `GET /search`: Search Quranic content
  - Query Parameters:
    - `q` (required): Search query
    - `p` (optional): Page number (default: 1)
    - `juz` (optional): Filter by juz number
    - `surat` (optional): Filter by surah number
    - `n` (optional): Number of results per page (default: 10)
  - Example: `/search?q=allah&p=1&n=10`

#### Surah
- `GET /surat`: List all surahs
  - Query Parameters:
    - `surat` (optional): Filter by specific surah number
  - Example: `/surat?surat=1`
- `GET /surat/:id`: List ayahs in a specific surah
  - Path Parameters:
    - `id`: Surah number
  - Query Parameters:
    - `p` (optional): Page number (default: 1)
    - `n` (optional): Number of ayahs per page (default: 10)
  - Example: `/surat/1?p=1&n=10`

#### Ayah
- `GET /ayat/:id`: Get detailed information for a specific ayah
  - Path Parameters:
    - `id`: Ayah ID
  - Example: `/ayat/1`

## Response Format

### Standard Response
All API responses follow this structure:
```json
{
  "version": "1.0",
  "data": {
    // Response data
  },
  "error": "" // Error message if any
}
```

### Search Response
For search results, the response includes:
```json
{
  "version": "1.0",
  "data": {
    "results": [
      {
        "id": 1,
        "juz": 1,
        "surat": 1,
        "ayat": 1,
        "text_indo": "Dengan nama Allah",
        "text_arabic": "بِسْمِ ٱللَّهِ"
      }
    ],
    "pagination": {
      "current_page": 1,
      "rows_per_page": 10,
      "total_pages": 10,
      "total_rows": 100
    },
    "aggregate": [
      {
        "type": "surat",
        "identifier": 1,
        "count": 1
      }
    ]
  },
  "error": ""
}
```

## Error Handling

The API uses standard HTTP status codes and provides detailed error messages:

- `200 OK`: Request successful
- `400 Bad Request`: Invalid parameters or request format
- `404 Not Found`: Resource not found
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server-side error

Error responses include:
```json
{
  "version": "1.0",
  "data": null,
  "error": "Detailed error message here"
}
```

## Rate Limiting

The API implements rate limiting for all endpoints:
- `/api/v1/search`: 10 requests per second
- `/api/v1/surat`: 10 requests per second
- `/api/v1/surat/*`: 10 requests per second
- `/api/v1/ayat/*`: 10 requests per second

When rate limit is exceeded, the API returns:
- Status code: `429 Too Many Requests`
- Headers:
  - `X-RateLimit-Limit`: Maximum requests allowed
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Time until limit reset

## Security Features

- Content Security Policy (CSP) headers
- CORS configuration
- Request timeout protection
- Rate limiting
- Input validation
- Error message sanitization
- Secure headers

## Technologies Used

- Go (Golang)
- Gin Web Framework
- MySQL
- Redis
- CORS middleware
- Rate limiting
- Timeout middleware
- Content Security Policy (CSP)

## Configuration

The API uses configuration files and environment variables for various settings:

### Configuration Files
- Rate limiting: `./internal/config/rate_limit.yml`
- Local settings: `./internal/config/file/local.yml`
- Docker settings: `./internal/config/file/docker.yml`
- Heroku settings: `./internal/config/file/heroku.yml`

### Environment Variables
- `API_PORT`: API server port
- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `REDIS_HOST`: Redis host
- `REDIS_PORT`: Redis port
- `REDIS_PASSWORD`: Redis password

## Development

### Prerequisites
- Go 1.16 or higher
- PostgreSQL 12 or higher
- Redis 6 or higher
- Make

### Running the Project

#### Local Development
```bash
# Install dependencies
go mod download

# Run the application
make run
```

#### Docker
```bash
# Build and run with Docker
docker-compose up --build
```

### Testing
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

## Project Structure

```shell
indoquran/ 
├── api/                    # API layer containing versioned endpoints
│   └── v1/                # Version 1 of the API
│       └── routes.go       # API routes and handlers
|── bin/                     # Compiled binaries
│   └── indoquran-api        # Executable binary
├── cmd/                   # Command-line applications
│   └── api/              # API server executable
│       └── main.go       # Server initialization and configuration
├── internal/              # Private application code
│   ├── config/           # Configuration files and settings
|   │   ├── file/       # Configuration files
|   │   |   └── local.yml # Local configuration file
|   │   |   └── docker.yml # Docker configuration file
|   │   |   └── heroku.yml # Heroku configuration file
│   │   ├── config.go   # Configuration settings
│   │   └── rate_limit.yml # Rate limiting configuration
│   │   └── vars.go   # Shared variables and structs
│   ├── controllers/      # Request handlers and business logic
│   │   ├── search.go     # Search functionality
│   │   ├── defaultOutput.go # Default output handler
│   │   ├── detail.go      # Detail handler
│   │   ├── list.go       # List handler
│   │   ├── search.go     # Search handler
│   │   └── vars.go       # Shared variables and structs
│   └── models/           # Data models and database interactions
│       └── model_guestbook.go # Model for guestbook entries
│       └── model_id_indonesian.go # Model for ID Indonesian
│       └── model_id_muntakhab.go # Model for ID Muntakhab
│       └── model_quran_ayat.go # Model for Quran ayat
│       └── model_quran_translation.go # Model for Quran ayat
│       └── model_traffic.go # Model for traffic
│       └── vars.go # Shared variables and structs
├── └─── services/         # Business logic and services
│       ├── detail/       # Detail service for ayah details
│       |   └── detail.go  # Detail service implementation
│       └── list/         # List service for surah listings
│       |   └── ayat.go   # List service implementation
│       |   └── surat.go  # List service implementation
│       └── search/       # Search service for searching Quran content
│           └── search.go  # Search service implementation
│           └── tools.go  # Search utility functions
├── pkg/                  # Public packages that can be used by external services
│   └── cache/            # Cache implementation
|       └── redis.go    # Redis cache implementation
│   ├── logger/           # Logging and error handling
│   │   └── logger.go      # Logger implementation
│   ├── database/        # Database connection and utilities
│   │   └── postgres.go  # PostgreSQL implementation
│   └── middleware/      # HTTP middleware components
│       ├── content_security_policy.go # Content Security Policy middleware
│       ├── rate_limit.go # Rate limiting middleware
│       └── timeout.go   # Request timeout handling
│       └── tools.go     # Utility functions
│       └── traffic.go   # Traffic handling middleware
├── go.mod              # Go modules definition
├── go.sum              # Go modules checksum
├── Makefile            # Makefile for running and building the project
└── README.md           # Project documentation
```

## Key Directory Structure

### `api/`
- Contains all API endpoint definitions and routing logic
- Organized by versions (v1) for better API lifecycle management
- Houses route handlers and server initialization

### `internal/`
- Core application logic not meant for external use
- Contains:
  - `config/`: Environment configurations and variables
  - `controllers/`: Request handlers and response formatting
  - `models/`: Data structures and database schema definitions
  - `services/`: Business logic implementation for search, list, and detail operations

### `pkg/`
- Reusable packages that can be imported by other projects
- Features:
  - `cache/`: Redis caching implementation
  - `database/`: PostgreSQL connection and queries
  - `logger/`: Logging utilities
  - `middleware/`: HTTP middleware components (rate limiting, CSP, timeout)

### `cmd/`
- Entry point for the application
- Contains the main server initialization and configuration
- Handles startup procedures and dependency injection

This structure enables clean separation of concerns, making the codebase maintainable and scalable. Each directory serves a specific purpose, following Go's standard project layout patterns and best practices.