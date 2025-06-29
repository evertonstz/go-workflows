package services

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationService struct {
	validator *validator.Validate
}

func NewValidationService() *ValidationService {
	v := validator.New()

	// Register custom validation functions
	v.RegisterValidation("folder_path", validateFolderPath)
	v.RegisterValidation("alphanum_space_dash_underscore", validateTagCharacters)

	return &ValidationService{
		validator: v,
	}
}

// Validate validates a struct using the validation tags
func (vs *ValidationService) Validate(s interface{}) error {
	return vs.validator.Struct(s)
}

// ValidateVar validates a single variable
func (vs *ValidationService) ValidateVar(field interface{}, tag string) error {
	return vs.validator.Var(field, tag)
}

// GetValidationErrors returns user-friendly validation error messages
func (vs *ValidationService) GetValidationErrors(err error) []string {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, validationError := range validationErrors {
			errors = append(errors, vs.formatValidationError(validationError))
		}
	} else {
		errors = append(errors, err.Error())
	}

	return errors
}

// formatValidationError converts a validation error to a user-friendly message
func (vs *ValidationService) formatValidationError(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		if err.Kind().String() == "string" {
			return fmt.Sprintf("%s must be at least %s characters long", field, param)
		}
		return fmt.Sprintf("%s must be at least %s", field, param)
	case "max":
		if err.Kind().String() == "string" {
			return fmt.Sprintf("%s must be at most %s characters long", field, param)
		}
		return fmt.Sprintf("%s must be at most %s", field, param)
	case "eq":
		return fmt.Sprintf("%s must be equal to %s", field, param)
	case "folder_path":
		return fmt.Sprintf("%s must be a valid folder path (e.g., '/', '/folder', '/folder/subfolder')", field)
	case "alphanum_space_dash_underscore":
		return fmt.Sprintf("%s can only contain letters, numbers, spaces, hyphens, and underscores", field)
	default:
		return fmt.Sprintf("%s failed validation for tag '%s'", field, tag)
	}
}

// Custom validation functions

// validateFolderPath validates that a folder path is in the correct format
func validateFolderPath(fl validator.FieldLevel) bool {
	path := fl.Field().String()

	// Root folder is always valid
	if path == "/" {
		return true
	}

	// Empty path is valid only if it's optional
	if path == "" {
		return true
	}

	// Must start with /
	if !strings.HasPrefix(path, "/") {
		return false
	}

	// Must not end with / (except for root)
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		return false
	}

	// Must not contain double slashes
	if strings.Contains(path, "//") {
		return false
	}

	// Check for valid characters: alphanumeric, hyphens, underscores, spaces
	// Split by / and validate each segment
	segments := strings.Split(strings.Trim(path, "/"), "/")
	for _, segment := range segments {
		if segment == "" {
			return false
		}
		if !isValidPathSegment(segment) {
			return false
		}
	}

	return true
}

// validateTagCharacters validates that a tag contains only allowed characters
func validateTagCharacters(fl validator.FieldLevel) bool {
	tag := fl.Field().String()
	return isValidTag(tag)
}

// Helper functions

func isValidPathSegment(segment string) bool {
	// Allow alphanumeric, hyphens, underscores, spaces, and dots
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-_\s\.]+$`, segment)
	return matched
}

func isValidTag(tag string) bool {
	// Allow alphanumeric, hyphens, underscores, and spaces
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-_\s]+$`, tag)
	return matched
}
