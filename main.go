package main

import (
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	addnew "github.com/evertonstz/go-workflows/screens/add_new"
	"github.com/evertonstz/go-workflows/shared"
	"github.com/jeandeaual/go-locale"
	"golang.org/x/text/language"

	"github.com/charmbracelet/bubbles/help"
	confirmationmodal "github.com/evertonstz/go-workflows/components/confirmation_modal"
	"github.com/evertonstz/go-workflows/components/notification"
	commandlist "github.com/evertonstz/go-workflows/screens/command_list"
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
	localesDir := "locales"
	i18nService, err := shared.NewI18nService(lang, localesDir)
	if err != nil {
		log.Fatalf("Error initializing i18n service: %v", err)
	}

	ctx := shared.WithI18n(context.Background(), i18nService)

	p := tea.NewProgram(newWithContext(ctx), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
