package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adrg/xdg"

	"github.com/evertonstz/go-workflows/shared/di"
	"github.com/evertonstz/go-workflows/shared/di/services"
)

func ParseFlags(i18nService *services.I18nService) (showVersion, showHelp, showConfig bool) {
	versionFlag := flag.Bool("version", false, i18nService.Translate("flags_version"))
	versionShortFlag := flag.Bool("v", false, i18nService.Translate("flags_version"))
	helpFlag := flag.Bool("help", false, i18nService.Translate("flags_help"))
	helpShortFlag := flag.Bool("h", false, i18nService.Translate("flags_help"))
	configFlag := flag.Bool("print-config", false, i18nService.Translate("flags_print_config"))

	flag.Parse()

	showVersion = *versionFlag || *versionShortFlag
	showHelp = *helpFlag || *helpShortFlag
	showConfig = *configFlag

	return
}

func HandleFlags(showVersion, showHelp, showConfig bool) {
	i18nService := di.GetService[*services.I18nService](di.I18nServiceKey)

	if showVersion {
		fmt.Printf("%s: %s\n", i18nService.Translate("flags_version"), Version)
		os.Exit(0)
	}

	if showConfig {
		appName := "go-workflows"
		dataFile, err := xdg.DataFile(fmt.Sprintf("%s/data.json", appName))
		if err != nil {
			fmt.Printf("Error determining the location of the configuration file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s: %s\n", i18nService.Translate("flags_print_config"), dataFile)
		os.Exit(0)
	}

	if showHelp {
		fmt.Printf("%s\n", i18nService.Translate("flags_usage"))
		fmt.Printf("  --version, -v       %s\n", i18nService.Translate("flags_version"))
		fmt.Printf("  --help, -h          %s\n", i18nService.Translate("flags_help"))
		fmt.Printf("  --print-config      %s\n", i18nService.Translate("flags_print_config"))
		os.Exit(0)
	}
}
