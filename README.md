# IndoQuran API

Welcome to the IndoQuran API project. This API provides access to Quranic content, including search functionality, surah listings, and detailed ayah information.

## Features

- Search Quranic content
- List all surahs
- Get ayahs for a specific surah
- Retrieve detailed information for a specific ayah

## API Endpoints

- `GET /api/v1/`: Welcome message
- `GET /api/v1/search`: Search functionality
- `GET /api/v1/surat`: List all surahs
- `GET /api/v1/surat/:id`: List ayahs in a specific surah
- `GET /api/v1/ayat/:id`: Get detailed information for a specific ayah

## Technologies Used

- Go (Golang)
- Gin Web Framework
- CORS middleware
- Rate limiting
- Timeout middleware
- Content Security Policy (CSP)

## Configuration

The API uses configuration files for various settings:

- Rate limiting: `./internal/config/rate_limit.yml`
- API port: Set via the `API_PORT` environment variable

## Running the Project

To run the project in debug mode locally:

`make run`


## Project Structure

```shell
indoquran/ 
├── api/                    # API layer containing versioned endpoints
│   └── v1/                # Version 1 of the API
│       └── routes/        # Route definitions and handlers mapping
│           └── routes.go  # Main routes configuration
├── internal/              # Private application code
│   ├── config/           # Configuration files and settings
│   │   └── rate_limit.yml # Rate limiting configuration
│   └── controllers/      # Request handlers and business logic
│       ├── search_handler.go      # Search functionality
│       ├── list_surat.go         # Surah listing handler
│       ├── list_ayat_in_surat.go # Ayah listing handler
│       └── detail_ayat.go        # Detailed ayah information
├── pkg/                  # Public packages that can be used by external services
│   ├── logger/          # Logging functionality
│   │   └── logger.go    # Logger implementation
│   └── middleware/      # HTTP middleware components
│       ├── logging_middleware.go     # Request logging
│       ├── rate_limiter.go          # Rate limiting implementation
│       ├── timeout_middleware.go     # Request timeout handling
│       └── content_security_policy.go # CSP implementation
├── main.go              # Application entry point
├── go.mod              # Go modules definition
├── go.sum              # Go modules checksum
└── README.md           # Project documentation
```

This structure follows Go best practices with clear separation of concerns:
- `api/` handles routing and API versioning
- `internal/` contains private application code
- `pkg/` houses reusable packages
- Root level contains configuration files and the main entry point

The organization makes it easy to maintain, test, and scale the application while keeping related code grouped together logically.