package services

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/jeandeaual/go-locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type I18nService struct {
	defaultLang string
}

var (
	bundle            *i18n.Bundle
	localizerInstance *i18n.Localizer
	DefaultLang       = "en"
)

func GetSystemLanguage() string {
	userLocale, err := locale.GetLocale()
	if err == nil {
		return userLocale
	}
	return DefaultLang
}

func DetermineLanguage() string {
	userLocaleStr := GetSystemLanguage()

	supportedLangs := []language.Tag{
		language.English,
		language.Portuguese,
	}
	matcher := language.NewMatcher(supportedLangs)

	userLangTag, err := language.Parse(userLocaleStr)
	if err != nil {
		return DefaultLang
	}

	tag, _, _ := matcher.Match(userLangTag)
	return tag.String()
}

func NewI18nServiceWithAutoDetection(localesDir string) (*I18nService, error) {
	lang := DetermineLanguage()
	return NewI18nService(lang, localesDir)
}

func NewI18nService(defaultLang string, localesDir string) (*I18nService, error) {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	err := filepath.Walk(localesDir, func(path string, info os.FileInfo, err error) error {
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

	localizerInstance = i18n.NewLocalizer(bundle, defaultLang)

	return &I18nService{
		defaultLang: defaultLang,
	}, nil
}

func (i *I18nService) Translate(key string) string {
	return localizerInstance.MustLocalize(&i18n.LocalizeConfig{
		MessageID: key,
	})
}

var i18nContextKey = &struct{}{}

func WithI18n(ctx context.Context, service *I18nService) context.Context {
	return context.WithValue(ctx, i18nContextKey, service)
}

func GetI18n(ctx context.Context) *I18nService {
	if service, ok := ctx.Value(i18nContextKey).(*I18nService); ok {
		return service
	}
	return nil
}
