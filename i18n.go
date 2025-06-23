package main

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/evertonstz/go-workflows/shared/loc"
)

func init() {
	loc.InitializeBundle()
}

func GetLocalizer(lang string) *i18n.Localizer {
	return loc.GetLocalizer(lang)
}