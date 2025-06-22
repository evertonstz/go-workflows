package main

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load English translations
	bundle.LoadMessageFile("locales/en.json")

	// Load Brazilian Portuguese translations
	bundle.MustLoadMessageFile("locales/pt-BR.json")
}

func GetLocalizer(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(bundle, lang)
}
