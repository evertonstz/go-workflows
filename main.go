package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/shared/loc"
	"github.com/jeandeaual/go-locale"
	"golang.org/x/text/language"
)

func getSystemLanguage() string {
	userLocale, err := locale.GetLocale()
	if err == nil {
		return userLocale
	}
	return loc.DefaultLang
}

// Fixed reference to bundle by using loc.GetBundle()
func determineLanguage() string {
	userLocaleStr := getSystemLanguage()

	supportedLangs := loc.GetBundle().LanguageTags()
	matcher := language.NewMatcher(supportedLangs)

	userLangTag, err := language.Parse(userLocaleStr)
	if err != nil {
		return loc.DefaultLang
	}

	tag, _, _ := matcher.Match(userLangTag)
	return tag.String()
}

func main() {
	lang := determineLanguage()
	loc.InitializeLocalizer(lang)

	p := tea.NewProgram(new(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
