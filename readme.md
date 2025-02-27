# Link In Bio API

A high-performance REST API for a Link In Bio management system built with Golang, Gin, and MongoDB. The system is designed to handle a high traffic load of 10M+ visits per day with peak traffic reaching 200k concurrent visits per minute.

## Features

- 🔗 Create, read, update, and delete bio links
- 📊 Track link visits and analytics
- ⏰ Automatic expired link cleanup
- 🔒 Authentication support
- 🚀 Optimized for high concurrency and performance

## Tech Stack

- **Backend**: Go (Golang)
- **Web Framework**: Gin Gonic
- **Database**: MongoDB
- **Containerization**: Docker & Docker Compose

## System Requirements

- Docker and Docker Compose
- Go 1.19+ (for local development)

## Getting Started

### Using Docker (Recommended)

1. Clone the repository:
   ```
   git clone https://github.com/Ammar022/take-home-assignment.git
   cd take-home-assignment
   ```

2. Start the application:
   ```
   docker-compose up -d
   ```

3. The API will be available at `http://localhost:8080`

### Local Development

1. Clone the repository:
   ```
   git clone https://github.com/Ammar022/take-home-assignment.git
   cd take-home-assignment
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Run MongoDB (using Docker or a local installation)
   ```
   docker run -d -p 27017:27017 mongo:5
   ```

4. Set environment variables:
   ```
   export LINKBIO_MONGODB_URI=mongodb://localhost:27017
   export LINKBIO_MONGODB_DATABASE=linkbio
   export LINKBIO_SERVER_ADDRESS=:8080
   export LINKBIO_CLEANUP_INTERVAL=15m
   ```

5. Run the application:
   ```
   go run cmd/api/main.go
   ```

## API Endpoints

### Link Management

| Method | Endpoint           | Description                            |
|--------|-------------------|----------------------------------------|
| GET    | /api/links         | Get all links                          |
| GET    | /api/links/:id     | Get a specific link                    |
| POST   | /api/links         | Create a new link                      |
| PUT    | /api/links/:id     | Update a link                          |
| DELETE | /api/links/:id     | Delete a link                          |

### Click Tracking

| Method | Endpoint           | Description                            |
|--------|-------------------|----------------------------------------|
| GET    | /visit/:id         | Visit a link (increment click count)   |
| GET    | /api/links/:id/visits | Get visit analytics for a link      |

## Authentication

For demonstration purposes, the application uses mock authentication. In a production environment, you would implement JWT-based authentication or similar.

To authenticate API requests:
1. Include an `Authorization` header with `Bearer <token>`
2. For testing, any non-empty token will be accepted

## Performance Optimizations

The application includes several optimizations for high-traffic scenarios:

- **Goroutines & Channels**: Used for concurrent processing of visits
- **Connection Pooling**: MongoDB connection pool is configured for high throughput
- **Context Propagation**: All operations support proper context handling
- **Rate Limiting**: Configurable rate limiting for public endpoints
- **Database Indexing**: Optimized indexes for common query patterns

## Running Tests

```
go test -v ./tests/...
```
