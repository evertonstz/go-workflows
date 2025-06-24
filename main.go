package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/evertonstz/go-workflows/shared/di"
	"github.com/evertonstz/go-workflows/shared/di/services"
)

var Version string

func main() {
	localesDir := "locales"
	i18nService, err := services.NewI18nServiceWithAutoDetection(localesDir)
	if err != nil {
		log.Fatalf("Error initializing i18n service: %v", err)
	}

	di.RegisterService(di.I18nServiceKey, i18nService)

	appName := "go-workflows"
	persistenceService, err := services.NewPersistenceService(appName)
	if err != nil {
		log.Fatalf("Error initializing persistence service: %v", err)
	}
	di.RegisterService(di.PersistenceServiceKey, persistenceService)

	showVersion, showHelp := ParseFlags()
	HandleFlags(showVersion, showHelp)

	p := tea.NewProgram(new(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
