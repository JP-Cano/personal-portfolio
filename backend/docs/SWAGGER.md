# Swagger API Documentation

This project uses Swagger/OpenAPI for API documentation.

## 📖 Viewing Documentation

Once the server is running, visit:
```
http://localhost:8080/swagger/index.html
```

## 🔄 Regenerating Documentation

After making changes to API handlers or adding new endpoints, regenerate the Swagger docs:

```bash
# Using Make (recommended)
make swagger

# Or directly with swag
swag init -g cmd/api/main.go -o docs
```

## 📝 Adding Documentation to New Endpoints

Use godoc-style comments above your handler functions:

```go
// GetExperience godoc
// @Summary Get experience by ID
// @Description Retrieves a single work experience by its ID
// @Tags experiences
// @Accept json
// @Produce json
// @Param id path int true "Experience ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.ExperienceResponse}
// @Failure 404 {object} utils.ErrorResponse
// @Router /experiences/{id} [get]
func (h *ExperienceHandler) GetExperience(c *gin.Context) {
    // handler implementation
}
```

## 🏷️ Annotation Reference

### Common Annotations

- `@Summary` - Short description (one line)
- `@Description` - Detailed description
- `@Tags` - Group endpoints together
- `@Accept` - Request content type (json, xml, etc.)
- `@Produce` - Response content type
- `@Param` - Define parameters
- `@Success` - Success response
- `@Failure` - Error response
- `@Router` - Route path and method

### Parameter Types

```go
// Path parameter
// @Param id path int true "Experience ID"

// Query parameter
// @Param page query int false "Page number"

// Body parameter
// @Param experience body dto.ExperienceRequest true "Experience data"

// Header parameter
// @Param Authorization header string true "Bearer token"
```

## 🔧 Main API Configuration

Main configuration is in `cmd/api/main.go`:

```go
// @title Personal Portfolio API
// @version 1.0
// @description API for managing personal portfolio
// @host localhost:8080
// @BasePath /api/v1
```

## 📚 Available Endpoints

### Experiences

- `GET /api/v1/experiences` - List all experiences
- `GET /api/v1/experiences/{id}` - Get experience by ID
- `POST /api/v1/experiences` - Create new experience
- `PUT /api/v1/experiences/{id}` - Update experience
- `DELETE /api/v1/experiences/{id}` - Delete experience

## 🧪 Testing with Swagger UI

1. Navigate to http://localhost:8080/swagger/index.html
2. Click on an endpoint to expand it
3. Click "Try it out"
4. Fill in required parameters
5. Click "Execute"
6. View the response

## 📦 Generated Files

Running `swag init` generates:
- `docs/docs.go` - Go code for embedding docs
- `docs/swagger.json` - OpenAPI spec in JSON
- `docs/swagger.yaml` - OpenAPI spec in YAML

**Note:** These files are auto-generated. Do not edit them manually!

## 🔗 Resources

- [Swag Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Gin Swagger](https://github.com/swaggo/gin-swagger)
