package main

import (
	"fmt"
	"log"

	"github.com/evertonstz/go-workflows/shared/di"
	"github.com/evertonstz/go-workflows/shared/di/services"
)

func main() {
	// Initialize DI container similar to main.go
	localesDir := "../locales"
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

	// Get the persistence service
	persistence := di.GetService[*services.PersistenceService](di.PersistenceServiceKey)

	// Create a new database manager
	dbManager, err := services.NewDatabaseManagerV2(persistence)
	if err != nil {
		log.Fatalf("Failed to create database manager: %v", err)
	}

	// Create some test folders
	fmt.Println("Creating test folders...")

	// Create Development folder
	devFolder, err := dbManager.CreateFolder("Development", "Development tools and scripts", "/")
	if err != nil {
		fmt.Printf("Error creating Development folder: %v\n", err)
	} else {
		fmt.Printf("Created folder: %s\n", devFolder.Path)
	}

	// Create Testing folder
	testFolder, err := dbManager.CreateFolder("Testing", "Testing scripts and tools", "/")
	if err != nil {
		fmt.Printf("Error creating Testing folder: %v\n", err)
	} else {
		fmt.Printf("Created folder: %s\n", testFolder.Path)
	}

	// Create a subfolder in Development
	frontendFolder, err := dbManager.CreateFolder("Frontend", "Frontend development tools", "/Development")
	if err != nil {
		fmt.Printf("Error creating Frontend subfolder: %v\n", err)
	} else {
		fmt.Printf("Created subfolder: %s\n", frontendFolder.Path)
	}

	// Move some existing items to the Development folder
	// First, let's see what items we have
	database := dbManager.GetDatabase()
	rootItems := database.GetItemsByFolder("/")

	fmt.Printf("\nFound %d items in root folder\n", len(rootItems))

	if len(rootItems) > 0 {
		// Move the first few items to Development folder
		for i, item := range rootItems {
			if i >= 3 { // Only move first 3 items
				break
			}
			err := dbManager.MoveItem(item.ID, "/Development")
			if err != nil {
				fmt.Printf("Error moving item %s: %v\n", item.Title, err)
			} else {
				fmt.Printf("Moved item '%s' to Development folder\n", item.Title)
			}
		}
	}

	// Save the database
	err = dbManager.Save()
	if err != nil {
		log.Fatalf("Failed to save database: %v", err)
	}

	fmt.Println("\nTest folders created and database saved!")

	// Show folder structure
	fmt.Println("\nFolder structure:")
	folders, items, _ := dbManager.GetFolderContents("/")
	fmt.Printf("Root folder - %d folders, %d items\n", len(folders), len(items))

	for _, folder := range folders {
		subfolders, subitems, _ := dbManager.GetFolderContents(folder.Path)
		fmt.Printf("  %s - %d folders, %d items\n", folder.Path, len(subfolders), len(subitems))
	}
}
