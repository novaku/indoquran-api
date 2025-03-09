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
The organization makes it easy to maintain, test, and scale the application while keeping related code grouped together logically.