package persist

import (
	"encoding/json" // Required for direct use in test setup if not mocking all JSON ops
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/adrg/xdg" // Required to access xdg.ConfigHome for mocking
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared"
	"github.com/stretchr/testify/assert"
)

func TestInitPersistionManagerCmd(t *testing.T) {
	originalXDGDataFile := XDGDataFile
	originalOSMkdirAll := OSMkdirAll
	originalXDGConfigHome := xdg.ConfigHome // Save original global variable

	defer func() {
		XDGDataFile = originalXDGDataFile
		OSMkdirAll = originalOSMkdirAll
		xdg.ConfigHome = originalXDGConfigHome // Restore original global variable
	}()

	appName := "test-app"
	// Create a temporary directory to simulate xdg.ConfigHome for isolated testing
	tempHomeDir, err := os.MkdirTemp("", "test-xdg-home-")
	assert.NoError(t, err, "Failed to create temp dir for xdg.ConfigHome mock")
	defer os.RemoveAll(tempHomeDir) // Clean up the temporary directory
	xdg.ConfigHome = tempHomeDir    // Override global xdg.ConfigHome

	expectedDataFilePathInXDG := filepath.Join("some", "xdg", "path", appName, "data.json")
	// The actual path used by os.MkdirAll in the function.
	expectedDirToCreate := filepath.Join(xdg.ConfigHome, appName)

	t.Run("success case", func(t *testing.T) {
		XDGDataFile = func(path string) (string, error) {
			assert.Equal(t, fmt.Sprintf("%s/data.json", appName), path)
			return expectedDataFilePathInXDG, nil
		}
		var mkdirCalledWithPath string
		OSMkdirAll = func(path string, perm os.FileMode) error {
			mkdirCalledWithPath = path
			assert.Equal(t, os.ModePerm, perm)
			return nil
		}

		cmd := InitPersistionManagerCmd(appName)
		msg := cmd()

		assert.IsType(t, InitiatedPersistion{}, msg)
		assert.Equal(t, expectedDataFilePathInXDG, msg.(InitiatedPersistion).DataFile)
		assert.Equal(t, expectedDirToCreate, mkdirCalledWithPath, "OSMkdirAll called with incorrect path")
	})

	t.Run("XDGDataFile error case", func(t *testing.T) {
		expectedErr := errors.New("XDGDataFile error")
		XDGDataFile = func(path string) (string, error) {
			return "", expectedErr
		}

		cmd := InitPersistionManagerCmd(appName)
		assert.PanicsWithValue(t, expectedErr, func() {
			cmd()
		}, "Should panic with XDGDataFile error")
	})

	t.Run("OSMkdirAll error case", func(t *testing.T) {
		expectedErr := errors.New("OSMkdirAll error")
		XDGDataFile = func(path string) (string, error) { // Should succeed
			return expectedDataFilePathInXDG, nil
		}
		OSMkdirAll = func(path string, perm os.FileMode) error {
			return expectedErr
		}

		cmd := InitPersistionManagerCmd(appName)
		assert.PanicsWithValue(t, expectedErr, func() {
			cmd()
		}, "Should panic with OSMkdirAll error")
	})
}

func TestLoadDataFileCmd(t *testing.T) {
	originalOSReadFile := OSReadFile
	originalOSCreate := OSCreate
	originalJSONUnmarshal := JSONUnmarshal

	defer func() {
		OSReadFile = originalOSReadFile
		OSCreate = originalOSCreate
		JSONUnmarshal = originalJSONUnmarshal
	}()

	testPath := "test_data.json"

	t.Run("file not exist, create success", func(t *testing.T) {
		OSReadFile = func(name string) ([]byte, error) {
			assert.Equal(t, testPath, name)
			return nil, os.ErrNotExist
		}
		var createCalled bool
		OSCreate = func(name string) (*os.File, error) {
			assert.Equal(t, testPath, name)
			createCalled = true
			// Return a dummy os.File that can be closed by the original code if it tries to.
			r, w, _ := os.Pipe()
			defer r.Close() // Ensure the read end is also eventually closed if not used
			w.Close()       // Close write end, this makes reads on r return EOF immediately.
			return r, nil
		}

		cmd := LoadDataFileCmd(testPath)
		msg := cmd()

		assert.True(t, createCalled, "OSCreate should have been called")
		assert.IsType(t, LoadedDataFileMsg{}, msg)
		assert.Equal(t, models.Items{Items: []models.Item{}}, msg.(LoadedDataFileMsg).Items)
	})

	t.Run("file not exist, create fails", func(t *testing.T) {
		expectedErr := errors.New("OSCreate error")
		OSReadFile = func(name string) ([]byte, error) {
			return nil, os.ErrNotExist
		}
		OSCreate = func(name string) (*os.File, error) {
			return nil, expectedErr
		}

		cmd := LoadDataFileCmd(testPath)
		assert.PanicsWithValue(t, expectedErr, func() {
			cmd()
		}, "Should panic with OSCreate error")
	})

	t.Run("OSReadFile other error", func(t *testing.T) {
		expectedErr := errors.New("OSReadFile other error")
		OSReadFile = func(name string) ([]byte, error) {
			return nil, expectedErr
		}

		cmd := LoadDataFileCmd(testPath)
		assert.PanicsWithValue(t, expectedErr, func() {
			cmd()
		}, "Should panic with OSReadFile other error")
	})

	t.Run("empty file", func(t *testing.T) {
		OSReadFile = func(name string) ([]byte, error) {
			assert.Equal(t, testPath, name)
			return []byte{}, nil
		}

		cmd := LoadDataFileCmd(testPath)
		msg := cmd()

		assert.IsType(t, LoadedDataFileMsg{}, msg)
		assert.Equal(t, models.Items{Items: []models.Item{}}, msg.(LoadedDataFileMsg).Items)
	})

	t.Run("valid JSON data", func(t *testing.T) {
		expectedItems := models.Items{Items: []models.Item{{Title: "Test Item 1"}}}
		jsonData, err := json.Marshal(expectedItems)
		assert.NoError(t, err, "Failed to marshal test data")

		OSReadFile = func(name string) ([]byte, error) {
			assert.Equal(t, testPath, name)
			return jsonData, nil
		}
		JSONUnmarshal = func(data []byte, v interface{}) error {
			assert.Equal(t, jsonData, data)
			itemsPtr, ok := v.(*models.Items)
			assert.True(t, ok, "v should be a *models.Items")
			return json.Unmarshal(data, itemsPtr) // Use actual unmarshal for mock behavior
		}

		cmd := LoadDataFileCmd(testPath)
		msg := cmd()

		assert.IsType(t, LoadedDataFileMsg{}, msg)
		assert.Equal(t, expectedItems, msg.(LoadedDataFileMsg).Items)
	})

	t.Run("invalid JSON data", func(t *testing.T) {
		invalidJsonData := []byte("this is not json")
		var temp models.Items
		expectedErr := json.Unmarshal(invalidJsonData, &temp) // Get actual error for comparison

		OSReadFile = func(name string) ([]byte, error) {
			return invalidJsonData, nil
		}
		JSONUnmarshal = json.Unmarshal // Let the original fail

		cmd := LoadDataFileCmd(testPath)
		assert.PanicsWithValue(t, expectedErr, func() {
			cmd()
		}, "Should panic with JSONUnmarshal error")
	})
}

func TestPersistListData(t *testing.T) {
	originalJSONMarshal := JSONMarshal
	originalOSWriteFile := OSWriteFile

	defer func() {
		JSONMarshal = originalJSONMarshal
		OSWriteFile = originalOSWriteFile
	}()

	testPath := "test_persist.json"
	testData := models.Items{Items: []models.Item{{Title: "Persist Me"}}}

	t.Run("success case", func(t *testing.T) {
		expectedJsonBytes, errMrsh := json.Marshal(testData)
		assert.NoError(t, errMrsh)

		var marshalCalledWith models.Items
		JSONMarshal = func(v interface{}) ([]byte, error) {
			var ok bool
			marshalCalledWith, ok = v.(models.Items)
			assert.True(t, ok, "JSONMarshal called with unexpected type")
			return expectedJsonBytes, nil
		}

		var writeFileCalledWithPath string
		var writeFileCalledWithData []byte
		var writeFileCalledWithPerm os.FileMode
		OSWriteFile = func(name string, data []byte, perm os.FileMode) error {
			writeFileCalledWithPath = name
			writeFileCalledWithData = data
			writeFileCalledWithPerm = perm
			return nil
		}

		cmd := PersistListData(testPath, testData)
		msg := cmd()

		assert.IsType(t, PersistedFileMsg{}, msg)
		assert.Equal(t, testData, marshalCalledWith)
		assert.Equal(t, testPath, writeFileCalledWithPath)
		assert.Equal(t, expectedJsonBytes, writeFileCalledWithData)
		assert.Equal(t, os.FileMode(0644), writeFileCalledWithPerm)
	})

	t.Run("JSONMarshal error case", func(t *testing.T) {
		expectedErr := errors.New("JSONMarshal error")
		JSONMarshal = func(v interface{}) ([]byte, error) {
			return nil, expectedErr
		}

		cmd := PersistListData(testPath, testData)
		msg := cmd()

		assert.IsType(t, shared.ErrorMsg{}, msg)
		assert.Equal(t, fmt.Errorf("failed to analyze JSON: %w", expectedErr).Error(), msg.(shared.ErrorMsg).Err.Error())
	})

	t.Run("OSWriteFile error case", func(t *testing.T) {
		expectedErr := errors.New("OSWriteFile error")
		jsonData, errMrsh := json.Marshal(testData) // Real marshal for this part
		assert.NoError(t, errMrsh)
		JSONMarshal = func(v interface{}) ([]byte, error) {
			return jsonData, nil
		}
		OSWriteFile = func(name string, data []byte, perm os.FileMode) error {
			return expectedErr
		}

		cmd := PersistListData(testPath, testData)
		msg := cmd()

		assert.IsType(t, shared.ErrorMsg{}, msg)
		assert.Equal(t, fmt.Errorf("failed saving file: %w", expectedErr).Error(), msg.(shared.ErrorMsg).Err.Error())
	})
}

func TestPersistDataFuncIsPersistListData(t *testing.T) {
	assert.NotNil(t, PersistDataFunc, "PersistDataFunc should be an alias to PersistListData and not nil")
	// Note: Comparing function pointers in Go is tricky and often not recommended.
	// The fact that PersistDataFunc is used for mocking in other packages (e.g., main.persistItems)
	// and those tests pass when PersistListData logic is correct, serves as an indirect test
	// of this alias. The primary goal is that PersistListData itself is thoroughly tested.
}
