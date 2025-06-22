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
	// bundle.LoadMessageFile("locales/active.en.toml") // TOML files are for message descriptions, not direct translations in this setup

	// Load Brazilian Portuguese translations
	bundle.MustLoadMessageFile("locales/pt-BR.json")
	// bundle.MustLoadMessageFile("locales/active.pt-BR.toml")

	// Add Brazilian Portuguese to supported languages (optional, but good practice)
	// bundle.AddMessages(language.BrazilianPortuguese) // This is not how you add languages, loading message files does it.
}

func GetLocalizer(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(bundle, lang)
}
