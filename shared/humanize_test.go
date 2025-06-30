package shared

import (
	"strings"
	"testing"
	"time"
)

func TestHumanizeTimeWithService_Past(t *testing.T) {
	// Create mock i18n service for English
	i18nEN := createMockI18nService("en")

	now := time.Now()

	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "just now",
			time:     now.Add(-1 * time.Second),
			expected: "just now",
		},
		{
			name:     "30 seconds ago",
			time:     now.Add(-30 * time.Second),
			expected: "30 seconds ago",
		},
		{
			name:     "1 minute ago",
			time:     now.Add(-1 * time.Minute),
			expected: "1 minute ago",
		},
		{
			name:     "5 minutes ago",
			time:     now.Add(-5 * time.Minute),
			expected: "5 minutes ago",
		},
		{
			name:     "1 hour ago",
			time:     now.Add(-1 * time.Hour),
			expected: "1 hour ago",
		},
		{
			name:     "3 hours ago",
			time:     now.Add(-3 * time.Hour),
			expected: "3 hours ago",
		},
		{
			name:     "1 day ago",
			time:     now.Add(-24 * time.Hour),
			expected: "1 day ago",
		},
		{
			name:     "3 days ago",
			time:     now.Add(-3 * 24 * time.Hour),
			expected: "3 days ago",
		},
		{
			name:     "1 week ago",
			time:     now.Add(-7 * 24 * time.Hour),
			expected: "1 week ago",
		},
		{
			name:     "2 weeks ago",
			time:     now.Add(-14 * 24 * time.Hour),
			expected: "2 weeks ago",
		},
		{
			name:     "1 month ago",
			time:     now.Add(-30 * 24 * time.Hour),
			expected: "1 month ago",
		},
		{
			name:     "3 months ago",
			time:     now.Add(-90 * 24 * time.Hour),
			expected: "3 months ago",
		},
		{
			name:     "1 year ago",
			time:     now.Add(-365 * 24 * time.Hour),
			expected: "1 year ago",
		},
		{
			name:     "2 years ago",
			time:     now.Add(-2 * 365 * 24 * time.Hour),
			expected: "2 years ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HumanizeTimeWithService(tt.time, i18nEN)
			if result != tt.expected {
				t.Errorf("HumanizeTimeWithService() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHumanizeTimeWithService_Future(t *testing.T) {
	// Create mock i18n service for English
	i18nEN := createMockI18nService("en")

	now := time.Now()

	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "in a few seconds",
			time:     now.Add(30 * time.Second),
			expected: "in a few seconds",
		},
		{
			name:     "in 1 minute",
			time:     now.Add(61 * time.Second), // Use 61 seconds to ensure we get "in 1 minute"
			expected: "in 1 minute",
		},
		{
			name:     "in 5 minutes",
			time:     now.Add(5*time.Minute + 30*time.Second), // Add buffer
			expected: "in 5 minutes",
		},
		{
			name:     "in 1 hour",
			time:     now.Add(1*time.Hour + 30*time.Second), // Add buffer
			expected: "in 1 hour",
		},
		{
			name:     "in 3 hours",
			time:     now.Add(3*time.Hour + 30*time.Second), // Add buffer
			expected: "in 3 hours",
		},
		{
			name:     "in 1 day",
			time:     now.Add(24*time.Hour + 30*time.Second), // Add buffer
			expected: "in 1 day",
		},
		{
			name:     "in 3 days",
			time:     now.Add(3*24*time.Hour + 30*time.Second), // Add buffer
			expected: "in 3 days",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HumanizeTimeWithService(tt.time, i18nEN)
			if result != tt.expected {
				t.Errorf("HumanizeTimeWithService() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHumanizeTimeWithService_Portuguese(t *testing.T) {
	// Create mock i18n service for Portuguese
	i18nPT := createMockI18nService("pt-BR")

	now := time.Now()

	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "agora há pouco",
			time:     now.Add(-1 * time.Second),
			expected: "agora há pouco",
		},
		{
			name:     "há 30 segundos",
			time:     now.Add(-30 * time.Second),
			expected: "há 30 segundos",
		},
		{
			name:     "há 1 minuto",
			time:     now.Add(-1 * time.Minute),
			expected: "há 1 minuto",
		},
		{
			name:     "há 5 minutos",
			time:     now.Add(-5 * time.Minute),
			expected: "há 5 minutos",
		},
		{
			name:     "há 1 hora",
			time:     now.Add(-1 * time.Hour),
			expected: "há 1 hora",
		},
		{
			name:     "há 1 dia",
			time:     now.Add(-24 * time.Hour),
			expected: "há 1 dia",
		},
		{
			name:     "há 1 semana",
			time:     now.Add(-7 * 24 * time.Hour),
			expected: "há 1 semana",
		},
		{
			name:     "há 1 mês",
			time:     now.Add(-30 * 24 * time.Hour),
			expected: "há 1 mês",
		},
		{
			name:     "há 1 ano",
			time:     now.Add(-365 * 24 * time.Hour),
			expected: "há 1 ano",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HumanizeTimeWithService(tt.time, i18nPT)
			if result != tt.expected {
				t.Errorf("HumanizeTimeWithService() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHumanizeTimeWithService_PortugueseFuture(t *testing.T) {
	// Create mock i18n service for Portuguese
	i18nPT := createMockI18nService("pt-BR")

	now := time.Now()

	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "em alguns segundos",
			time:     now.Add(30 * time.Second),
			expected: "em alguns segundos",
		},
		{
			name:     "em 1 minuto",
			time:     now.Add(61 * time.Second), // Use 61 seconds to ensure we get "em 1 minuto"
			expected: "em 1 minuto",
		},
		{
			name:     "em 5 minutos",
			time:     now.Add(5*time.Minute + 30*time.Second), // Add buffer
			expected: "em 5 minutos",
		},
		{
			name:     "em 1 hora",
			time:     now.Add(1*time.Hour + 30*time.Second), // Add buffer
			expected: "em 1 hora",
		},
		{
			name:     "em 1 dia",
			time:     now.Add(24*time.Hour + 30*time.Second), // Add buffer
			expected: "em 1 dia",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HumanizeTimeWithService(tt.time, i18nPT)
			if result != tt.expected {
				t.Errorf("HumanizeTimeWithService() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHumanizeTimePastWithService_EdgeCases(t *testing.T) {
	i18nEN := createMockI18nService("en")

	tests := []struct {
		name     string
		diff     time.Duration
		expected string
	}{
		{
			name:     "exactly 1 second",
			diff:     1 * time.Second,
			expected: "just now",
		},
		{
			name:     "exactly 60 seconds",
			diff:     60 * time.Second,
			expected: "1 minute ago",
		},
		{
			name:     "exactly 60 minutes",
			diff:     60 * time.Minute,
			expected: "1 hour ago",
		},
		{
			name:     "exactly 24 hours",
			diff:     24 * time.Hour,
			expected: "1 day ago",
		},
		{
			name:     "exactly 7 days",
			diff:     7 * 24 * time.Hour,
			expected: "1 week ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := humanizeTimePastWithService(tt.diff, i18nEN)
			if result != tt.expected {
				t.Errorf("humanizeTimePastWithService() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHumanizeTimeFutureWithService_EdgeCases(t *testing.T) {
	i18nEN := createMockI18nService("en")

	tests := []struct {
		name     string
		diff     time.Duration
		expected string
	}{
		{
			name:     "exactly 1 second",
			diff:     1 * time.Second,
			expected: "in a few seconds",
		},
		{
			name:     "exactly 60 seconds",
			diff:     60 * time.Second,
			expected: "in 1 minute",
		},
		{
			name:     "exactly 60 minutes",
			diff:     60 * time.Minute,
			expected: "in 1 hour",
		},
		{
			name:     "exactly 24 hours",
			diff:     24 * time.Hour,
			expected: "in 1 day",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := humanizeTimeFutureWithService(tt.diff, i18nEN)
			if result != tt.expected {
				t.Errorf("humanizeTimeFutureWithService() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Additional comprehensive test cases

func TestHumanizeTimeWithService_ZeroTime(t *testing.T) {
	i18nEN := createMockI18nService("en")

	// Test with zero time (should show many years ago)
	zeroTime := time.Time{}
	result := HumanizeTimeWithService(zeroTime, i18nEN)

	// Zero time should result in a very old date
	if !strings.Contains(result, "year") {
		t.Errorf("Expected zero time to show years, got: %v", result)
	}
}

func TestHumanizeTimeWithService_SameTime(t *testing.T) {
	i18nEN := createMockI18nService("en")

	// Test with exactly the same time
	now := time.Now()
	result := HumanizeTimeWithService(now, i18nEN)

	// Should show "just now" for same time
	expected := "just now"
	if result != expected {
		t.Errorf("HumanizeTimeWithService() = %v, want %v", result, expected)
	}
}

func TestHumanizeTimeWithService_MultiplierEdgeCases(t *testing.T) {
	i18nEN := createMockI18nService("en")

	now := time.Now()

	tests := []struct {
		name           string
		time           time.Time
		expectContains string
	}{
		{
			name:           "59 seconds ago",
			time:           now.Add(-59 * time.Second),
			expectContains: "seconds ago",
		},
		{
			name:           "59 minutes ago",
			time:           now.Add(-59 * time.Minute),
			expectContains: "minutes ago",
		},
		{
			name:           "23 hours ago",
			time:           now.Add(-23 * time.Hour),
			expectContains: "hours ago",
		},
		{
			name:           "6 days ago",
			time:           now.Add(-6 * 24 * time.Hour),
			expectContains: "days ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HumanizeTimeWithService(tt.time, i18nEN)
			if !strings.Contains(result, tt.expectContains) {
				t.Errorf("HumanizeTimeWithService() = %v, expected to contain %v", result, tt.expectContains)
			}
		})
	}
}

func TestHumanizeTimeWithService_Fallback(t *testing.T) {
	// Test with mock service that returns keys as-is (fallback behavior)
	fallbackService := &MockI18nService{language: "unknown"}

	now := time.Now()
	pastTime := now.Add(-5 * time.Minute)

	result := HumanizeTimeWithService(pastTime, fallbackService)

	// Should still work even with unknown language (returns key)
	if result == "" {
		t.Error("HumanizeTimeWithService() should not return empty string even with fallback")
	}
}

func TestMockI18nService_Coverage(t *testing.T) {
	// Test to ensure our mock service covers all expected keys
	enService := createMockI18nService("en")
	ptService := createMockI18nService("pt-BR")

	testKeys := []string{
		"time_just_now",
		"time_seconds_ago",
		"time_minute_ago",
		"time_minutes_ago",
		"time_hour_ago",
		"time_hours_ago",
		"time_day_ago",
		"time_days_ago",
		"time_week_ago",
		"time_weeks_ago",
		"time_month_ago",
		"time_months_ago",
		"time_year_ago",
		"time_years_ago",
		"time_in_seconds",
		"time_in_minute",
		"time_in_minutes",
		"time_in_hour",
		"time_in_hours",
		"time_in_day",
		"time_in_days",
	}

	for _, key := range testKeys {
		enResult := enService.Translate(key)
		ptResult := ptService.Translate(key)

		if enResult == key {
			t.Errorf("English mock service missing translation for key: %s", key)
		}
		if ptResult == key {
			t.Errorf("Portuguese mock service missing translation for key: %s", key)
		}
		if enResult == ptResult {
			t.Errorf("English and Portuguese translations are identical for key: %s", key)
		}
	}
}

// createMockI18nService creates a mock i18n service with predefined translations for testing
func createMockI18nService(lang string) I18nTranslator {
	// Create a mock service that returns translations based on language
	return &MockI18nService{language: lang}
}

// MockI18nService is a mock implementation of the I18nService for testing
type MockI18nService struct {
	language string
}

func (m *MockI18nService) Translate(key string) string {
	if m.language == "pt-BR" {
		return getPortugueseTranslation(key)
	}
	return getEnglishTranslation(key)
}

func getEnglishTranslation(key string) string {
	translations := map[string]string{
		"time_just_now":    "just now",
		"time_seconds_ago": "%d seconds ago",
		"time_minute_ago":  "1 minute ago",
		"time_minutes_ago": "%d minutes ago",
		"time_hour_ago":    "1 hour ago",
		"time_hours_ago":   "%d hours ago",
		"time_day_ago":     "1 day ago",
		"time_days_ago":    "%d days ago",
		"time_week_ago":    "1 week ago",
		"time_weeks_ago":   "%d weeks ago",
		"time_month_ago":   "1 month ago",
		"time_months_ago":  "%d months ago",
		"time_year_ago":    "1 year ago",
		"time_years_ago":   "%d years ago",
		"time_in_seconds":  "in a few seconds",
		"time_in_minute":   "in 1 minute",
		"time_in_minutes":  "in %d minutes",
		"time_in_hour":     "in 1 hour",
		"time_in_hours":    "in %d hours",
		"time_in_day":      "in 1 day",
		"time_in_days":     "in %d days",
	}

	if translation, exists := translations[key]; exists {
		return translation
	}
	return key // Return key if translation not found
}

func getPortugueseTranslation(key string) string {
	translations := map[string]string{
		"time_just_now":    "agora há pouco",
		"time_seconds_ago": "há %d segundos",
		"time_minute_ago":  "há 1 minuto",
		"time_minutes_ago": "há %d minutos",
		"time_hour_ago":    "há 1 hora",
		"time_hours_ago":   "há %d horas",
		"time_day_ago":     "há 1 dia",
		"time_days_ago":    "há %d dias",
		"time_week_ago":    "há 1 semana",
		"time_weeks_ago":   "há %d semanas",
		"time_month_ago":   "há 1 mês",
		"time_months_ago":  "há %d meses",
		"time_year_ago":    "há 1 ano",
		"time_years_ago":   "há %d anos",
		"time_in_seconds":  "em alguns segundos",
		"time_in_minute":   "em 1 minuto",
		"time_in_minutes":  "em %d minutos",
		"time_in_hour":     "em 1 hora",
		"time_in_hours":    "em %d horas",
		"time_in_day":      "em 1 dia",
		"time_in_days":     "em %d dias",
	}

	if translation, exists := translations[key]; exists {
		return translation
	}
	return key // Return key if translation not found
}

// Benchmark tests for performance
func BenchmarkHumanizeTimeWithService(b *testing.B) {
	i18nEN := createMockI18nService("en")
	now := time.Now()
	testTime := now.Add(-5 * time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HumanizeTimeWithService(testTime, i18nEN)
	}
}

func BenchmarkHumanizeTimeWithService_Portuguese(b *testing.B) {
	i18nPT := createMockI18nService("pt-BR")
	now := time.Now()
	testTime := now.Add(-5 * time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HumanizeTimeWithService(testTime, i18nPT)
	}
}
