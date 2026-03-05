# Arena

## Version

This is the version of the Arena application. To get the current version, use the `--version` flag.

## Build Instructions

To build the Arena application with versioning, use the following command:

```bash
go build -ldflags='-X main.version=1.0.0'
```

## Run Instructions

To run the application, use the following command:

```bash
go run cmd/api/main.go
```

## Swagger Instructions

To access the API documentation, navigate to `http://localhost:8080/swagger` after starting the server.

## CRUD Route Examples

- **Create**: `POST /api/resource`
- **Read**: `GET /api/resource/{id}`
- **Update**: `PUT /api/resource/{id}`
- **Delete**: `DELETE /api/resource/{id}`
