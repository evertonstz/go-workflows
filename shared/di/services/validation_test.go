package services

import (
	"strings"
	"testing"
	"time"

	"github.com/evertonstz/go-workflows/models"
)

func TestValidationService_ValidateItem(t *testing.T) {
	validationService := NewValidationService()

	tests := []struct {
		name          string
		item          models.ItemV2
		expectError   bool
		errorContains string
	}{
		{
			name: "valid item",
			item: models.ItemV2{
				ID:          "test-id",
				Title:       "Valid Item",
				Desc:        "A valid description",
				Command:     "echo hello",
				DateAdded:   time.Now(),
				DateUpdated: time.Now(),
				Tags:        []string{"test", "valid"},
				Metadata:    map[string]string{"author": "test"},
				FolderPath:  "/",
			},
			expectError: false,
		},
		{
			name: "missing title",
			item: models.ItemV2{
				ID:          "test-id",
				Title:       "", // Invalid: empty title
				Desc:        "A valid description",
				Command:     "echo hello",
				DateAdded:   time.Now(),
				DateUpdated: time.Now(),
				FolderPath:  "/",
			},
			expectError:   true,
			errorContains: "Title is required",
		},
		{
			name: "invalid folder path",
			item: models.ItemV2{
				ID:          "test-id",
				Title:       "Valid Item",
				Desc:        "A valid description",
				Command:     "echo hello",
				DateAdded:   time.Now(),
				DateUpdated: time.Now(),
				FolderPath:  "invalid-path", // Invalid: doesn't start with /
			},
			expectError:   true,
			errorContains: "FolderPath must be a valid folder path",
		},
		{
			name: "tag with invalid characters",
			item: models.ItemV2{
				ID:          "test-id",
				Title:       "Valid Item",
				Desc:        "A valid description",
				Command:     "echo hello",
				DateAdded:   time.Now(),
				DateUpdated: time.Now(),
				Tags:        []string{"test!", "valid"}, // Invalid: contains !
				FolderPath:  "/",
			},
			expectError:   true,
			errorContains: "can only contain letters, numbers, spaces, hyphens, and underscores",
		},
		{
			name: "too long title",
			item: models.ItemV2{
				ID:          "test-id",
				Title:       string(make([]byte, 300)), // Invalid: too long
				Desc:        "A valid description",
				Command:     "echo hello",
				DateAdded:   time.Now(),
				DateUpdated: time.Now(),
				FolderPath:  "/",
			},
			expectError:   true,
			errorContains: "Title must be at most 255 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validationService.Validate(tt.item)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected validation error but got none")
					return
				}

				errors := validationService.GetValidationErrors(err)
				found := false
				for _, errMsg := range errors {
					if tt.errorContains != "" && strings.Contains(errMsg, tt.errorContains) {
						found = true
						break
					}
				}

				if !found && tt.errorContains != "" {
					t.Errorf("Expected error containing '%s', but got: %v", tt.errorContains, errors)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error but got: %v", err)
				}
			}
		})
	}
}

func TestValidationService_ValidateFolder(t *testing.T) {
	validationService := NewValidationService()

	tests := []struct {
		name          string
		folder        models.FolderV2
		expectError   bool
		errorContains string
	}{
		{
			name: "valid folder",
			folder: models.FolderV2{
				ID:          "test-id",
				Name:        "Valid Folder",
				Description: "A valid description",
				Path:        "/valid-folder",
				ParentPath:  "/",
				DateAdded:   time.Now(),
				DateUpdated: time.Now(),
				Metadata:    map[string]string{"type": "test"},
			},
			expectError: false,
		},
		{
			name: "missing name",
			folder: models.FolderV2{
				ID:          "test-id",
				Name:        "", // Invalid: empty name
				Description: "A valid description",
				Path:        "/valid-folder",
				DateAdded:   time.Now(),
				DateUpdated: time.Now(),
			},
			expectError:   true,
			errorContains: "Name is required",
		},
		{
			name: "invalid path format",
			folder: models.FolderV2{
				ID:          "test-id",
				Name:        "Valid Folder",
				Description: "A valid description",
				Path:        "invalid//path", // Invalid: double slash
				DateAdded:   time.Now(),
				DateUpdated: time.Now(),
			},
			expectError:   true,
			errorContains: "Path must be a valid folder path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validationService.Validate(tt.folder)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected validation error but got none")
					return
				}

				errors := validationService.GetValidationErrors(err)
				found := false
				for _, errMsg := range errors {
					if tt.errorContains != "" && strings.Contains(errMsg, tt.errorContains) {
						found = true
						break
					}
				}

				if !found && tt.errorContains != "" {
					t.Errorf("Expected error containing '%s', but got: %v", tt.errorContains, errors)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error but got: %v", err)
				}
			}
		})
	}
}

func TestValidationService_ValidateFolderPath(t *testing.T) {
	validationService := NewValidationService()

	tests := []struct {
		path        string
		expectValid bool
	}{
		{"/", true},                   // Root is valid
		{"/folder", true},             // Simple folder
		{"/folder/subfolder", true},   // Nested folder
		{"/folder-name", true},        // Hyphens allowed
		{"/folder_name", true},        // Underscores allowed
		{"/folder name", true},        // Spaces allowed
		{"", true},                    // Empty is valid for optional fields
		{"invalid", false},            // Must start with /
		{"/folder/", false},           // Must not end with / (except root)
		{"/folder//subfolder", false}, // No double slashes
		{"/folder/sub*folder", false}, // Invalid characters
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			err := validationService.ValidateVar(tt.path, "folder_path")

			if tt.expectValid && err != nil {
				t.Errorf("Expected path '%s' to be valid but got error: %v", tt.path, err)
			}

			if !tt.expectValid && err == nil {
				t.Errorf("Expected path '%s' to be invalid but it was accepted", tt.path)
			}
		})
	}
}
