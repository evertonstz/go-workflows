package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/shared/di"
	"github.com/evertonstz/go-workflows/shared/di/services"
	"github.com/jeandeaual/go-locale"
	"golang.org/x/text/language"
)

var (
	Version string
)

func getSystemLanguage() string {
	// return "pt-BR" // TODO: remove
	userLocale, err := locale.GetLocale()
	if err == nil {
		return userLocale
	}
	return services.DefaultLang
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
		return services.DefaultLang
	}

	tag, _, _ := matcher.Match(userLangTag)
	return tag.String()
}

func main() {
	lang := determineLanguage()
	localesDir := "locales"
	i18nService, err := services.NewI18nService(lang, localesDir)
	if err != nil {
		log.Fatalf("Error initializing i18n service: %v", err)
	}

	// Register the i18n service using the updated generic implementation
	di.RegisterService(di.I18nServiceKey, i18nService)

	showVersion, showHelp := ParseFlags()
	HandleFlags(showVersion, showHelp)

	p := tea.NewProgram(new(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
