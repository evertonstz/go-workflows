package main

import (
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	addnew "github.com/evertonstz/go-workflows/screens/add_new"
	"github.com/evertonstz/go-workflows/shared"
	"github.com/evertonstz/go-workflows/shared/loc"
	"github.com/jeandeaual/go-locale"
	"golang.org/x/text/language"

	"github.com/charmbracelet/bubbles/help"
	confirmationmodal "github.com/evertonstz/go-workflows/components/confirmation_modal"
	"github.com/evertonstz/go-workflows/components/notification"
	commandlist "github.com/evertonstz/go-workflows/screens/command_list"
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

func newWithContext(ctx context.Context) tea.Model {
	return model{
		confirmationModal: confirmationmodal.NewConfirmationModal("", "", "", nil, nil),
		help:              help.New(),
		addNewScreen:      addnew.NewWithContext(ctx),
		listScreen:        commandlist.New(),
		notification:      notification.New("Workflows"),
		panelsStyle: panelsStyle{
			helpPanelStyle:         helpPanelStyle,
			notificationPanelStyle: notificationPanelStyle,
		},
		currentHelpHeight: 0,
		screenState:       newList,
	}
}

func main() {
	lang := determineLanguage()
	paths := map[string]string{
		"en":     "locales/en.json",
		"pt-BR": "locales/pt-BR.json",
	}
	i18nService, err := shared.NewI18nService(lang, paths)
	if err != nil {
		log.Fatalf("Error initializing i18n service: %v", err)
	}

	ctx := shared.WithI18n(context.Background(), i18nService)

	p := tea.NewProgram(newWithContext(ctx), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
