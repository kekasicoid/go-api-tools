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
