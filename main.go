package main

import (

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

import "github.com/jeandeaual/go-locale"
import (
	"log"

	"golang.org/x/text/language"
)

var localizer *i18n.Localizer
var defaultLang = "en"

func getSystemLanguage() string {
	userLocale, err := locale.GetLocale()
	if err == nil {
		return userLocale
	}
	return defaultLang
}

func determineLanguage() string {
	userLocaleStr := getSystemLanguage()

	supportedLangs := bundle.LanguageTags()
	matcher := language.NewMatcher(supportedLangs)

	userLangTag, err := language.Parse(userLocaleStr)
	if err != nil {
		return defaultLang
	}

	tag, _, _ := matcher.Match(userLangTag)
	return tag.String()
}

func main() {
	lang := determineLanguage()
	localizer = GetLocalizer(lang)

	p := tea.NewProgram(new(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "error_starting_app"}), err)
	}
}
