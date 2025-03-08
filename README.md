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


export ENV=local
go run main.go


## Project Structure

```shell
indoquran/ 
├── api/ 
│ └── v1/ 
│ └── routes/ 
│ └── routes.go 
├── internal/ 
│ ├── config/ 
│ │ └── rate_limit.yml 
│ └── controllers/ 
│ ├── search_handler.go 
│ ├── list_surat.go 
│ ├── list_ayat_in_surat.go 
│ └── detail_ayat.go 
├── pkg/ 
│ ├── logger/ 
│ │ └── logger.go 
│ └── middleware/ 
│ ├── logging_middleware.go 
│ ├── rate_limiter.go 
│ ├── timeout_middleware.go 
│ └── content_security_policy.go 
├── main.go 
├── go.mod 
├── go.sum 
└── README.md
```