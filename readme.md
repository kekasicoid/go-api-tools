# Go API Tools

Simple tools API built with Go (Gin) using clean architecture.

## Features
- JSON Formatter API
- Clean Architecture
- Rate Limiting
- CORS ready

## Environment Variables
Set the following environment variables:
```
PORT=8080
CORS_ORIGIN=https://www.kekasi.co.id
```

## How to Run
1. Install dependencies:
   ```bash
   go mod tidy
   ```
2. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

## How to Build
1. Build for local OS:
   ```bash
   go build -o server cmd/server/main.go
   ```
2. Build for production (static binary):
   ```bash
   CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server cmd/server/main.go
   ```

## Docker Deployment
1. Build and start using Docker Compose:
   ```bash
   docker-compose up --build
   ```
2. Stop services:
   ```bash
   docker-compose down
   ```
