package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validators
	registerCustomValidators()
}

// registerCustomValidators registers all custom validators
func registerCustomValidators() {
	validate.RegisterValidation("after_start_date", validateAfterStartDate)
	validate.RegisterValidation("date_format", validateDateFormat)
	validate.RegisterValidation("after_start_date_str", validateAfterStartDateString)
}

// validateAfterStartDate validates that EndDate is after StartDate (for time.Time fields)
func validateAfterStartDate(fl validator.FieldLevel) bool {
	field := fl.Field()

	// The validator gives us the dereferenced value if it's a pointer
	// So we need to check if it's a time.Time struct or zero value
	if field.Kind() != reflect.Struct {
		return false
	}

	// If it's the zero value for time.Time, it means nil pointer (optional field)
	endDate, ok := field.Interface().(time.Time)
	if !ok || endDate.IsZero() {
		return true // Optional field, validation passes
	}

	// Get parent struct to access StartDate
	parent := fl.Parent()

	// Get the StartDate field from parent struct
	startDateField := parent.FieldByName("StartDate")
	if !startDateField.IsValid() {
		return false
	}

	startDate, ok := startDateField.Interface().(time.Time)
	if !ok {
		return false
	}

	// EndDate must be after StartDate (not equal)
	return endDate.After(startDate)
}

// validateDateFormat validates that a string is a valid date format
func validateDateFormat(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()

	// Empty strings are allowed (will be handled by required tag)
	if dateStr == "" {
		return true
	}

	// Try to parse the date
	_, err := utils.ParseDate(dateStr)
	return err == nil
}

// validateAfterStartDateString validates that end_date string is after start_date string
func validateAfterStartDateString(fl validator.FieldLevel) bool {
	endDateStr := fl.Field().String()

	// If empty, validation passes (optional field)
	if endDateStr == "" {
		return true
	}

	// Parse end date
	endDate, err := utils.ParseDate(endDateStr)
	if err != nil {
		return false
	}

	// Get parent struct to access StartDate
	parent := fl.Parent()
	startDateField := parent.FieldByName("StartDate")
	if !startDateField.IsValid() {
		return false
	}

	startDateStr := startDateField.String()
	if startDateStr == "" {
		return false
	}

	// Parse start date
	startDate, err := utils.ParseDate(startDateStr)
	if err != nil {
		return false
	}

	// EndDate must be after StartDate
	return endDate.After(startDate)
}

// ValidateRequest is a generic middleware that validates request body against validation tags
func ValidateRequest[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T

		// First bind the JSON to validate JSON structure
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err)
			c.Abort()
			return
		}

		// Then validate using validator tags
		if err := validate.Struct(req); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			errorMessages := formatValidationErrors(validationErrors)
			utils.RespondWithError(c, http.StatusBadRequest, strings.Join(errorMessages, "; "), err)
			c.Abort()
			return
		}

		// Store validated request in context for handler to use
		c.Set("validatedRequest", req)
		c.Next()
	}
}

// ValidateQuery is a generic middleware that validates query parameters against validation tags
func ValidateQuery[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T

		// Bind query parameters
		if err := c.ShouldBindQuery(&req); err != nil {
			logger.Error("Failed to bind query parameters: %v", err)
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid query parameters", err)
			c.Abort()
			return
		}

		// Validate using validator tags
		if err := validate.Struct(req); err != nil {
			var validationErrors validator.ValidationErrors
			errors.As(err, &validationErrors)
			errorMessages := formatValidationErrors(validationErrors)
			utils.RespondWithError(c, http.StatusBadRequest, strings.Join(errorMessages, "; "), err)
			c.Abort()
			return
		}

		// Store validated query in context for handler to use
		c.Set("validatedQuery", req)
		c.Next()
	}
}

// formatValidationErrors converts validator errors to human-readable messages
func formatValidationErrors(errors validator.ValidationErrors) []string {
	var messages []string

	for _, err := range errors {
		var message string
		field := err.Field()

		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", field)
		case "min":
			// Check if it's numeric validation or string length
			if err.Type().Kind() == reflect.Int || err.Type().Kind() == reflect.Int64 {
				message = fmt.Sprintf("%s must be at least %s", field, err.Param())
			} else {
				message = fmt.Sprintf("%s must be at least %s characters", field, err.Param())
			}
		case "max":
			// Check if it's numeric validation or string length
			if err.Type().Kind() == reflect.Int || err.Type().Kind() == reflect.Int64 {
				message = fmt.Sprintf("%s must not exceed %s", field, err.Param())
			} else {
				message = fmt.Sprintf("%s must not exceed %s characters", field, err.Param())
			}
		case "url":
			message = fmt.Sprintf("%s must be a valid URL", field)
		case "oneof":
			message = fmt.Sprintf("%s must be one of: %s", field, err.Param())
		case "gtefield":
			message = fmt.Sprintf("%s must be after %s", field, err.Param())
		case "after_start_date":
			message = fmt.Sprintf("%s must be after start date", field)
		case "after_start_date_str":
			message = fmt.Sprintf("%s must be after start date", field)
		case "date_format":
			message = fmt.Sprintf("%s must be a valid date (formats: DD/MM/YYYY, DD/MM/YY, DD-MM-YYYY, DD-MM-YY, YYYY-MM-DD)", field)
		default:
			message = fmt.Sprintf("%s is invalid", field)
		}

		messages = append(messages, message)
	}

	return messages
}
