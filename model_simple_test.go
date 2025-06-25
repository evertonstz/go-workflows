package main

import (
	"testing"
)

func TestModelHelpers_IsSmallWidth(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		expected bool
	}{
		{
			name:     "Small width - 50",
			width:    50,
			expected: true,
		},
		{
			name:     "Large width - 120",
			width:    120,
			expected: false,
		},
		{
			name:     "Borderline width - 80",
			width:    80,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testModel := model{
				termDimensions: termDimensions{width: tt.width},
			}

			result := testModel.isSmallWidth()
			if result != tt.expected {
				t.Errorf("Width %d: expected %v, got %v", tt.width, tt.expected, result)
			}
		})
	}
}

func TestTermDimensions(t *testing.T) {
	dimensions := termDimensions{
		width:  100,
		height: 50,
	}

	if dimensions.width != 100 {
		t.Errorf("Expected width 100, got %d", dimensions.width)
	}
	if dimensions.height != 50 {
		t.Errorf("Expected height 50, got %d", dimensions.height)
	}
}

func TestScreenStateConstants(t *testing.T) {
	if addNew != 0 {
		t.Errorf("Expected addNew to be 0, got %d", addNew)
	}
	if newList != 1 {
		t.Errorf("Expected newList to be 1, got %d", newList)
	}
}
