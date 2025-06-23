package i18n

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	localizerInstance *i18n.Localizer
	once              sync.Once
	bundle            *i18n.Bundle
	DefaultLang       = "en" // Default language if none is specified
)

// InitializeLocalizer initializes the localizer singleton.
func InitializeLocalizer(lang string) {
	once.Do(func() {
		localizerInstance = GetLocalizer(lang)
	})
}

// GetLocalizerInstance returns the singleton instance of the localizer.
func GetLocalizerInstance() *i18n.Localizer {
	return localizerInstance
}

// GetBundle ...
func GetBundle() *i18n.Bundle {
	return bundle
}

// GetLocalizer creates a new localizer instance for the given language.
func GetLocalizer(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(bundle, lang)
}

// InitializeBundle ...
func InitializeBundle() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	err := filepath.Walk("locales", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			if _, err := bundle.LoadMessageFile(path); err != nil {
				log.Fatalf("Failed to load translation file %s: %v", path, err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error loading translation files: %v", err)
	}
}

func init() {
	InitializeBundle()
}
