package main

import (
	"fmt"
	"log"
	"os"

	"github.com/evertonstz/go-workflows/shared/di"
)

func main() {
	// Initialize the database
	container := di.NewContainer()
	dbManager := container.GetDatabaseManager()

	// Get the database path
	dbPath := dbManager.GetDatabasePath()
	fmt.Printf("Database path: %s\n", dbPath)

	// Check if database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("Database file does not exist, creating new one...")
	} else {
		fmt.Println("Database file exists")
	}

	// List all items before deletion
	fmt.Println("\n=== Items before deletion ===")
	root, err := dbManager.GetRootFolder()
	if err != nil {
		log.Fatal("Failed to get root folder:", err)
	}

	_, items, err := dbManager.GetFolderContents(root.Path)
	if err != nil {
		log.Fatal("Failed to get folder contents:", err)
	}

	fmt.Printf("Found %d items:\n", len(items))
	for i, item := range items {
		fmt.Printf("  %d. ID: %s, Title: %s, Description: %s\n", i+1, item.ID, item.Title, item.Desc)
	}

	if len(items) == 0 {
		fmt.Println("No items found to delete")
		return
	}

	// Delete the first item
	firstItem := items[0]
	fmt.Printf("\n=== Deleting item: %s (ID: %s) ===\n", firstItem.Title, firstItem.ID)

	err = dbManager.DeleteItem(firstItem.ID)
	if err != nil {
		log.Fatal("Failed to delete item:", err)
	}
	fmt.Println("Item deleted successfully")

	// List all items after deletion
	fmt.Println("\n=== Items after deletion ===")
	_, itemsAfter, err := dbManager.GetFolderContents(root.Path)
	if err != nil {
		log.Fatal("Failed to get folder contents after deletion:", err)
	}

	fmt.Printf("Found %d items:\n", len(itemsAfter))
	for i, item := range itemsAfter {
		fmt.Printf("  %d. ID: %s, Title: %s, Description: %s\n", i+1, item.ID, item.Title, item.Desc)
	}

	fmt.Printf("\nDeleted successfully! Item count went from %d to %d\n", len(items), len(itemsAfter))
}
