package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared/di/services"
)

func main() {
	var (
		action   = flag.String("action", "", "Action to perform: create-folder, create-item, list, search, migrate, stats, validate")
		name     = flag.String("name", "", "Name of the folder or item")
		desc     = flag.String("desc", "", "Description")
		command  = flag.String("command", "", "Command for items")
		folder   = flag.String("folder", "/", "Folder path")
		parent   = flag.String("parent", "/", "Parent folder path")
		query    = flag.String("query", "", "Search query")
		tags     = flag.String("tags", "", "Comma-separated tags")
		metadata = flag.String("metadata", "", "JSON metadata")
		dataFile = flag.String("data", "", "Data file path (optional)")
		force    = flag.Bool("force", false, "Force operation")
		pretty   = flag.Bool("pretty", false, "Pretty print JSON output")
	)
	flag.Parse()

	if *action == "" {
		fmt.Println("Go Workflows v2 Database Demo")
		fmt.Println("Usage examples:")
		fmt.Println("  # Create a folder")
		fmt.Println("  go run examples/v2_cli_demo.go -action=create-folder -name=scripts -desc='Development scripts' -parent=/")
		fmt.Println("")
		fmt.Println("  # Create an item")
		fmt.Println("  go run examples/v2_cli_demo.go -action=create-item -name='Deploy Script' -desc='Deploy to production' -command='deploy.sh' -folder=/scripts -tags='deploy,production'")
		fmt.Println("")
		fmt.Println("  # List folder contents")
		fmt.Println("  go run examples/v2_cli_demo.go -action=list -folder=/scripts")
		fmt.Println("")
		fmt.Println("  # Search")
		fmt.Println("  go run examples/v2_cli_demo.go -action=search -query=deploy")
		fmt.Println("")
		fmt.Println("  # Migrate from v1 to v2")
		fmt.Println("  go run examples/v2_cli_demo.go -action=migrate")
		fmt.Println("")
		fmt.Println("  # Get statistics")
		fmt.Println("  go run examples/v2_cli_demo.go -action=stats")
		fmt.Println("")
		fmt.Println("  # Validate database")
		fmt.Println("  go run examples/v2_cli_demo.go -action=validate")
		os.Exit(1)
	}

	// Initialize persistence service
	var persistence *services.PersistenceService
	var err error

	if *dataFile != "" {
		persistence = &services.PersistenceService{}
		// This would need to be properly initialized with the custom data file
		log.Fatal("Custom data file not implemented in this demo")
	} else {
		persistence, err = services.NewPersistenceService("go-workflows-v2-demo")
		if err != nil {
			log.Fatalf("Failed to create persistence service: %v", err)
		}
	}

	// Initialize database manager
	manager, err := services.NewDatabaseManagerV2(persistence)
	if err != nil {
		log.Fatalf("Failed to create database manager: %v", err)
	}

	switch *action {
	case "create-folder":
		if *name == "" {
			log.Fatal("Folder name is required")
		}

		folder, createErr := manager.CreateFolder(*name, *desc, *parent)
		if createErr != nil {
			log.Fatalf("Failed to create folder: %v", createErr)
		}

		fmt.Printf("Created folder: %s (ID: %s)\n", folder.Path, folder.ID)

	case "create-item":
		if *name == "" || *command == "" {
			log.Fatal("Item name and command are required")
		}

		var tagList []string
		if *tags != "" {
			tagList = strings.Split(*tags, ",")
			for i, tag := range tagList {
				tagList[i] = strings.TrimSpace(tag)
			}
		}

		var metadataMap map[string]string
		if *metadata != "" {
			if unmarshalErr := json.Unmarshal([]byte(*metadata), &metadataMap); unmarshalErr != nil {
				log.Fatalf("Invalid metadata JSON: %v", unmarshalErr)
			}
		}

		item, createErr := manager.CreateItem(*name, *desc, *command, *folder, tagList, metadataMap)
		if createErr != nil {
			log.Fatalf("Failed to create item: %v", createErr)
		}

		fmt.Printf("Created item: %s (ID: %s) in folder %s\n", item.Title, item.ID, item.FolderPath)

	case "list":
		subfolders, items, listErr := manager.GetFolderContents(*folder)
		if listErr != nil {
			log.Fatalf("Failed to list folder contents: %v", listErr)
		}

		fmt.Printf("Contents of folder: %s\n", *folder)
		fmt.Printf("Subfolders (%d):\n", len(subfolders))
		for _, subfolder := range subfolders {
			fmt.Printf("  📁 %s - %s\n", subfolder.Name, subfolder.Description)
		}

		fmt.Printf("Items (%d):\n", len(items))
		for _, item := range items {
			fmt.Printf("  📄 %s - %s\n", item.Title, item.Desc)
			if len(item.Tags) > 0 {
				fmt.Printf("     Tags: %s\n", strings.Join(item.Tags, ", "))
			}
		}

	case "search":
		if *query == "" {
			log.Fatal("Search query is required")
		}

		criteria := models.SearchCriteria{Query: *query}
		if *folder != "/" {
			criteria.FolderPath = *folder
		}
		if *tags != "" {
			criteria.Tags = strings.Split(*tags, ",")
		}

		result := manager.Search(criteria)

		fmt.Printf("Search results for query '%s' (%d total):\n", *query, result.Total)
		fmt.Printf("Folders (%d):\n", len(result.Folders))
		for _, folder := range result.Folders {
			fmt.Printf("  📁 %s - %s\n", folder.Name, folder.Description)
		}

		fmt.Printf("Items (%d):\n", len(result.Items))
		for _, item := range result.Items {
			fmt.Printf("  📄 %s - %s (in %s)\n", item.Title, item.Desc, item.FolderPath)
			fmt.Printf("     Command: %s\n", item.Command)
			if len(item.Tags) > 0 {
				fmt.Printf("     Tags: %s\n", strings.Join(item.Tags, ", "))
			}
		}

	case "migrate":
		version, versionErr := persistence.GetDatabaseVersion()
		if versionErr != nil {
			log.Fatalf("Failed to get database version: %v", versionErr)
		}

		if version == "2.0" {
			fmt.Println("Database is already v2.0")
			return
		}

		fmt.Printf("Migrating from version %s to 2.0...\n", version)
		migrateErr := persistence.MigrateToV2()
		if migrateErr != nil {
			log.Fatalf("Failed to migrate to v2: %v", migrateErr)
		}

		fmt.Println("Migration completed successfully!")
		fmt.Printf("Backup of original data saved as: %s.v1.backup\n", persistence.GetDataFilePath())

	case "stats":
		stats := manager.GetStatistics()

		if *pretty {
			data, marshalErr := json.MarshalIndent(stats, "", "  ")
			if marshalErr != nil {
				log.Fatalf("Failed to marshal statistics: %v", marshalErr)
			}
			fmt.Println(string(data))
		} else {
			fmt.Printf("Database Version: %s\n", stats["version"])
			fmt.Printf("Total Folders: %v\n", stats["total_folders"])
			fmt.Printf("Total Items: %v\n", stats["total_items"])

			fmt.Println("\nItems by folder:")
			if itemsByFolder, ok := stats["items_by_folder"].(map[string]int); ok {
				for folder, count := range itemsByFolder {
					fmt.Printf("  %s: %d items\n", folder, count)
				}
			}

			fmt.Println("\nTag usage:")
			if tagUsage, ok := stats["tag_usage"].(map[string]int); ok {
				for tag, count := range tagUsage {
					fmt.Printf("  %s: %d items\n", tag, count)
				}
			}
		}

	case "validate":
		issues := manager.ValidateDatabase()

		if len(issues) == 0 {
			fmt.Println("✅ Database validation passed - no issues found")
		} else {
			fmt.Printf("❌ Database validation found %d issues:\n", len(issues))
			for i, issue := range issues {
				fmt.Printf("  %d. %s\n", i+1, issue)
			}
		}

	case "tree":
		tree := manager.GetFolderTree()
		printTree(tree, 0)

	case "delete-folder":
		if *folder == "" {
			log.Fatal("Folder path is required")
		}

		deleteErr := manager.DeleteFolder(*folder, *force)
		if deleteErr != nil {
			log.Fatalf("Failed to delete folder: %v", deleteErr)
		}

		fmt.Printf("Deleted folder: %s\n", *folder)

	case "export":
		db := manager.GetDatabase()

		var data []byte
		var marshalErr error
		if *pretty {
			data, marshalErr = json.MarshalIndent(db, "", "  ")
		} else {
			data, marshalErr = json.Marshal(db)
		}

		if marshalErr != nil {
			log.Fatalf("Failed to marshal database: %v", marshalErr)
		}

		fmt.Println(string(data))

	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}

func printTree(tree map[string]interface{}, indent int) {
	prefix := strings.Repeat("  ", indent)

	// Print folders
	if folders, ok := tree["folders"].([]map[string]interface{}); ok {
		for _, folderData := range folders {
			if folder, ok := folderData["folder"].(models.FolderV2); ok {
				fmt.Printf("%s📁 %s - %s\n", prefix, folder.Name, folder.Description)

				// Print items in this folder
				if items, ok := folderData["items"].([]models.ItemV2); ok {
					for _, item := range items {
						fmt.Printf("%s  📄 %s - %s\n", prefix, item.Title, item.Desc)
					}
				}

				// Recursively print subfolders
				if subfolders, ok := folderData["subfolders"].([]map[string]interface{}); ok {
					subtree := map[string]interface{}{"folders": subfolders}
					printTree(subtree, indent+1)
				}
			}
		}
	}

	// Print root items
	if items, ok := tree["items"].([]models.ItemV2); ok && indent == 0 {
		fmt.Println("📄 Root Items:")
		for _, item := range items {
			fmt.Printf("  📄 %s - %s\n", item.Title, item.Desc)
		}
	}
}
