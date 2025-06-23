package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/shared"
	"github.com/evertonstz/go-workflows/shared/di"
	"github.com/jeandeaual/go-locale"
	"golang.org/x/text/language"
)

func getSystemLanguage() string {
	return "pt-BR" // TODO: remove
	userLocale, err := locale.GetLocale()
	if err == nil {
		return userLocale
	}
	return shared.DefaultLang
}

func determineLanguage() string {
	userLocaleStr := getSystemLanguage()

	supportedLangs := []language.Tag{
		language.English,
		language.Portuguese,
	}
	matcher := language.NewMatcher(supportedLangs)

	userLangTag, err := language.Parse(userLocaleStr)
	if err != nil {
		return shared.DefaultLang
	}

	tag, _, _ := matcher.Match(userLangTag)
	return tag.String()
}

func main() {
	lang := determineLanguage()
	localesDir := "locales"
	i18nService, err := shared.NewI18nService(lang, localesDir)
	if err != nil {
		log.Fatalf("Error initializing i18n service: %v", err)
	}

	di.RegisterService(di.I18nServiceKey, i18nService)

	p := tea.NewProgram(new(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
