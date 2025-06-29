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

	if err := v.RegisterValidation("folder_path", validateFolderPath); err != nil {
		panic(fmt.Sprintf("failed to register 'folder_path' validation: %v", err))
	}

	if err := v.RegisterValidation("alphanum_space_dash_underscore", validateTagCharacters); err != nil {
		panic(fmt.Sprintf("failed to register 'alphanum_space_dash_underscore' validation: %v", err))
	}

	return &ValidationService{
		validator: v,
	}
}

func (vs *ValidationService) Validate(s interface{}) error {
	return vs.validator.Struct(s)
}

func (vs *ValidationService) ValidateVar(field interface{}, tag string) error {
	return vs.validator.Var(field, tag)
}

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

func validateFolderPath(fl validator.FieldLevel) bool {
	path := fl.Field().String()

	if path == "/" {
		return true
	}

	if path == "" {
		return true
	}

	if !strings.HasPrefix(path, "/") {
		return false
	}

	if len(path) > 1 && strings.HasSuffix(path, "/") {
		return false
	}

	if strings.Contains(path, "//") {
		return false
	}

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

func validateTagCharacters(fl validator.FieldLevel) bool {
	tag := fl.Field().String()
	return isValidTag(tag)
}

func isValidPathSegment(segment string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-_\s\.]+$`, segment)
	return matched
}

func isValidTag(tag string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-_\s]+$`, tag)
	return matched
}
