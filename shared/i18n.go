package shared

import (
	"context"
	"encoding/json"
	"os"
)

type I18nService struct {
	translations map[string]map[string]string
	defaultLang  string
}

func NewI18nService(defaultLang string, paths map[string]string) (*I18nService, error) {
	translations := make(map[string]map[string]string)
	for lang, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var data map[string]string
		if err := json.NewDecoder(file).Decode(&data); err != nil {
			return nil, err
		}
		translations[lang] = data
	}

	return &I18nService{
		translations: translations,
		defaultLang: defaultLang,
	}, nil
}

func (i *I18nService) Translate(lang, key string) string {
	if langTranslations, ok := i.translations[lang]; ok {
		if value, ok := langTranslations[key]; ok {
			return value
		}
	}
	return key // Fallback to key if translation is missing
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
