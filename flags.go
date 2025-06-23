package main

import (
	"flag"
	"os"
)

func ParseFlags() (showVersion, showHelp bool) {
	versionFlag := flag.Bool("version", false, "Exibe a versão do aplicativo")
	versionShortFlag := flag.Bool("v", false, "Exibe a versão do aplicativo (abreviação)")
	helpFlag := flag.Bool("help", false, "Exibe ajuda sobre o aplicativo")
	helpShortFlag := flag.Bool("h", false, "Exibe ajuda sobre o aplicativo (abreviação)")

	flag.Parse()

	showVersion = *versionFlag || *versionShortFlag
	showHelp = *helpFlag || *helpShortFlag

	return
}

func HandleFlags(showVersion, showHelp bool) {
	if showVersion {
		println("Versão:", Version)
		os.Exit(0)
	}

	if showHelp {
		println("Uso:")
		println("  --version, -v    Exibe a versão do aplicativo")
		println("  --help, -h       Exibe ajuda sobre o aplicativo")
		os.Exit(0)
	}
}
