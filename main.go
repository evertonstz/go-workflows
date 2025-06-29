package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	helpkeys "github.com/evertonstz/go-workflows/components/keys"
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

	validationService := services.NewValidationService()
	di.RegisterService(di.ValidationServiceKey, validationService)

	helpkeys.InitializeGlobalKeys(i18nService)

	appName := "go-workflows"
	persistenceService, err := services.NewPersistenceService(appName)
	if err != nil {
		log.Fatalf("Error initializing persistence service: %v", err)
	}
	di.RegisterService(di.PersistenceServiceKey, persistenceService)

	showVersion, showHelp, showConfig := ParseFlags(i18nService)
	HandleFlags(showVersion, showHelp, showConfig)

	p := tea.NewProgram(new(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
