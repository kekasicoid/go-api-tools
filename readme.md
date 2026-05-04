# Go API Tools

Simple tools API built with Go (Gin) using clean architecture.

## Features
- JSON Formatter API
- [JWT Decode & Validate](docs/jwt.md)
- Clean Architecture
- Rate Limiting
- CORS ready

## Environment Variables
Set the following environment variables (see `.env.example`):
```env
HTTP_PORT=8080
CORS_ORIGIN=https://kekasi.co.id,*
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0
REQUEST_ID_TTL_HOURS=24
JSON_FORMATTER_TTL_HOURS=1
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

## Generate Swagger Docs
1. Install `swag` CLI:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@v1.16.6
   ```
2. Generate Swagger files into `docs/`:
   ```bash
   swag init -g cmd/server/main.go -o docs
   ```
3. If `swag` is not recognized on Windows PowerShell, run this first in the same terminal session:
   ```powershell
   $env:Path += ";$(go env GOPATH)\\bin"
   ```

## Access Swagger UI
1. Run the application with `APP_ENV=development` or `APP_ENV=dev` so the Swagger route is enabled.
   ```bash
   APP_ENV=development go run cmd/server/main.go
   ```
2. Open Swagger UI in the browser:
   ```text
   http://localhost:8080/swagger/index.html
   ```
3. If you use a different port, replace `8080` with the value from `HTTP_PORT`.

## Docker Deployment
1. Build and start using Docker Compose:
   ```bash
   docker-compose up --build -d
   ```
2. Stop services:
   ```bash
   docker-compose down
   ```
