package response

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrInsufficient = errors.New("insufficient balance")
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func Error(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	code := "internal_server_error"

	switch {
	case errors.Is(err, ErrBadRequest):
		status = http.StatusBadRequest
		code = "bad_request"
	case errors.Is(err, ErrUnauthorized):
		status = http.StatusUnauthorized
		code = "unauthorized"
	case errors.Is(err, ErrForbidden):
		status = http.StatusForbidden
		code = "forbidden"
	case errors.Is(err, ErrNotFound):
		status = http.StatusNotFound
		code = "not_found"
	case errors.Is(err, ErrConflict):
		status = http.StatusConflict
		code = "conflict"
	case errors.Is(err, ErrInsufficient):
		status = http.StatusUnprocessableEntity
		code = "insufficient_balance"
	}

	c.JSON(status, ErrorResponse{Error: code, Message: err.Error()})
}

func ValidationError(c *gin.Context, err error) {
	fields := map[string]string{}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			fields[fieldName(fieldError)] = validationMessage(fieldError)
		}
	}

	c.JSON(http.StatusBadRequest, ValidationErrorResponse{
		Error:   "validation_error",
		Message: "invalid request body",
		Fields:  fields,
	})
}

func fieldName(fieldError validator.FieldError) string {
	field := fieldError.Field()
	if field == "" {
		field = fieldError.StructField()
	}
	return toSnakeCase(field)
}

func validationMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fieldError.Param())
	case "gt":
		return fmt.Sprintf("must be greater than %s", fieldError.Param())
	case "oneof":
		return "must be one of: " + fieldError.Param()
	default:
		return "is invalid"
	}
}

func toSnakeCase(value string) string {
	if value == "" {
		return value
	}

	runes := []rune(value)
	var builder strings.Builder
	for i, r := range runes {
		if unicode.IsUpper(r) {
			prevIsLower := i > 0 && unicode.IsLower(runes[i-1])
			nextIsLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if i > 0 && (prevIsLower || nextIsLower) {
				builder.WriteRune('_')
			}
			builder.WriteRune(unicode.ToLower(r))
			continue
		}
		builder.WriteRune(r)
	}
	return builder.String()
}
