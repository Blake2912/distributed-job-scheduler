# Swagger Documentation Setup for Go

This project uses **Swagger (OpenAPI)** to generate interactive API documentation for the Go backend.

We use:
- `swaggo/swag` → to generate Swagger spec from Go comments  
- `swaggo/gin-swagger` → to serve Swagger UI  
- `swaggo/files` → static Swagger UI assets  

---

## 1. Install Dependencies

Install Swagger generator:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Verify Installation
```bash
swag --version
```
## 2. Adding new APIs into swagger docs
Refer existing handlers for the swagger document, add godocs for your handler
Run inside the `scheduler-service/`
```bash
swag init -g cmd/server/main.go
```
The above command will add the newly added routes and apis into the swagger doc

Swagger UI will be available in
```
http://localhost:8081/swagger/index.html
```
