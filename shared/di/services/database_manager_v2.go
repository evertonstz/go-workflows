package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/evertonstz/go-workflows/models"
)

// DatabaseManagerV2 provides high-level operations for the v2 database
type DatabaseManagerV2 struct {
	persistenceService *PersistenceService
	database           models.DatabaseV2
}

func NewDatabaseManagerV2(persistenceService *PersistenceService) (*DatabaseManagerV2, error) {
	database, err := persistenceService.LoadDataV2()
	if err != nil {
		return nil, fmt.Errorf("failed to load database: %w", err)
	}

	return &DatabaseManagerV2{
		persistenceService: persistenceService,
		database:           database,
	}, nil
}

// Folder operations

func (dm *DatabaseManagerV2) CreateFolder(name, description, parentPath string) (*models.FolderV2, error) {
	// Normalize parent path
	if parentPath == "" {
		parentPath = "/"
	}
	if parentPath != "/" && !strings.HasPrefix(parentPath, "/") {
		parentPath = "/" + parentPath
	}

	// Generate folder path
	var folderPath string
	if parentPath == "/" {
		folderPath = "/" + name
	} else {
		folderPath = strings.TrimSuffix(parentPath, "/") + "/" + name
	}

	// Validate parent folder exists (if not root)
	if parentPath != "/" {
		if _, found := dm.database.GetFolderByPath(parentPath); !found {
			return nil, fmt.Errorf("parent folder %s does not exist", parentPath)
		}
	}

	// Normalize parent path for storage to match IsChildOf logic
	normalizedParentPath := parentPath
	if parentPath == "/" {
		normalizedParentPath = ""
	}

	folder := models.FolderV2{
		Name:        name,
		Description: description,
		Path:        folderPath,
		ParentPath:  normalizedParentPath,
		DateAdded:   time.Now(),
		DateUpdated: time.Now(),
		Metadata:    make(map[string]string),
	}

	if err := dm.database.AddFolder(folder); err != nil {
		return nil, err
	}

	if err := dm.Save(); err != nil {
		return nil, fmt.Errorf("failed to save after creating folder: %w", err)
	}

	// Return the folder with generated ID
	for _, f := range dm.database.Folders {
		if f.Path == folderPath {
			return &f, nil
		}
	}

	return &folder, nil
}

func (dm *DatabaseManagerV2) GetFolder(path string) (*models.FolderV2, error) {
	folder, found := dm.database.GetFolderByPath(path)
	if !found {
		return nil, fmt.Errorf("folder %s not found", path)
	}
	return folder, nil
}

func (dm *DatabaseManagerV2) GetFolderContents(folderPath string) ([]models.FolderV2, []models.ItemV2, error) {
	// Verify folder exists (unless it's root)
	if folderPath != "/" && folderPath != "" {
		if _, found := dm.database.GetFolderByPath(folderPath); !found {
			return nil, nil, fmt.Errorf("folder %s does not exist", folderPath)
		}
	}

	subfolders := dm.database.GetSubfolders(folderPath)
	items := dm.database.GetItemsByFolder(folderPath)

	return subfolders, items, nil
}

func (dm *DatabaseManagerV2) DeleteFolder(path string, force bool) error {
	if path == "/" {
		return fmt.Errorf("cannot delete root folder")
	}

	// Check if folder has contents
	subfolders := dm.database.GetSubfolders(path)
	items := dm.database.GetItemsByFolder(path)

	if !force && (len(subfolders) > 0 || len(items) > 0) {
		return fmt.Errorf("folder %s is not empty (contains %d subfolders and %d items). Use force=true to delete",
			path, len(subfolders), len(items))
	}

	if force {
		// Delete all items in the folder
		for _, item := range items {
			if err := dm.database.DeleteItem(item.ID); err != nil {
				return fmt.Errorf("failed to delete item %s: %w", item.ID, err)
			}
		}

		// Recursively delete subfolders
		for _, subfolder := range subfolders {
			if err := dm.DeleteFolder(subfolder.Path, true); err != nil {
				return fmt.Errorf("failed to delete subfolder %s: %w", subfolder.Path, err)
			}
		}
	}

	if err := dm.database.DeleteFolder(path); err != nil {
		return err
	}

	return dm.Save()
}

// Item operations

func (dm *DatabaseManagerV2) CreateItem(title, description, command, folderPath string, tags []string, metadata map[string]string) (*models.ItemV2, error) {
	// Normalize folder path
	if folderPath == "" {
		folderPath = "/"
	}

	// Validate folder exists (unless it's root)
	if folderPath != "/" {
		if _, found := dm.database.GetFolderByPath(folderPath); !found {
			return nil, fmt.Errorf("folder %s does not exist", folderPath)
		}
	}

	if tags == nil {
		tags = []string{}
	}
	if metadata == nil {
		metadata = make(map[string]string)
	}

	item := models.ItemV2{
		Title:       title,
		Desc:        description,
		Command:     command,
		FolderPath:  folderPath,
		DateAdded:   time.Now(),
		DateUpdated: time.Now(),
		Tags:        tags,
		Metadata:    metadata,
	}

	if err := dm.database.AddItem(item); err != nil {
		return nil, err
	}

	if err := dm.Save(); err != nil {
		return nil, fmt.Errorf("failed to save after creating item: %w", err)
	}

	// Return the item with generated ID
	for _, i := range dm.database.Items {
		if i.Title == title && i.FolderPath == folderPath {
			return &i, nil
		}
	}

	return &item, nil
}

func (dm *DatabaseManagerV2) GetItem(id string) (*models.ItemV2, error) {
	item, found := dm.database.GetItemByID(id)
	if !found {
		return nil, fmt.Errorf("item %s not found", id)
	}
	return item, nil
}

func (dm *DatabaseManagerV2) UpdateItem(id string, title, description, command, folderPath string, tags []string, metadata map[string]string) error {
	// Get current item
	currentItem, found := dm.database.GetItemByID(id)
	if !found {
		return fmt.Errorf("item %s not found", id)
	}

	// Validate new folder exists (if changing)
	if folderPath != "" && folderPath != currentItem.FolderPath {
		if folderPath != "/" {
			if _, found := dm.database.GetFolderByPath(folderPath); !found {
				return fmt.Errorf("folder %s does not exist", folderPath)
			}
		}
	}

	// Create updated item
	updatedItem := *currentItem
	if title != "" {
		updatedItem.Title = title
	}
	if description != "" {
		updatedItem.Desc = description
	}
	if command != "" {
		updatedItem.Command = command
	}
	if folderPath != "" {
		updatedItem.FolderPath = folderPath
	}
	if tags != nil {
		updatedItem.Tags = tags
	}
	if metadata != nil {
		updatedItem.Metadata = metadata
	}

	if err := dm.database.UpdateItem(id, updatedItem); err != nil {
		return err
	}

	return dm.Save()
}

func (dm *DatabaseManagerV2) DeleteItem(id string) error {
	if err := dm.database.DeleteItem(id); err != nil {
		return err
	}

	return dm.Save()
}

func (dm *DatabaseManagerV2) MoveItem(id, newFolderPath string) error {
	// Validate new folder exists
	if newFolderPath != "/" {
		if _, found := dm.database.GetFolderByPath(newFolderPath); !found {
			return fmt.Errorf("destination folder %s does not exist", newFolderPath)
		}
	}

	// Get current item
	currentItem, found := dm.database.GetItemByID(id)
	if !found {
		return fmt.Errorf("item %s not found", id)
	}

	// Update item with new folder path
	updatedItem := *currentItem
	updatedItem.FolderPath = newFolderPath
	updatedItem.DateUpdated = time.Now()

	if err := dm.database.UpdateItem(id, updatedItem); err != nil {
		return err
	}

	return dm.Save()
}

// Search operations

func (dm *DatabaseManagerV2) Search(criteria models.SearchCriteria) models.SearchResult {
	return dm.database.Search(criteria)
}

func (dm *DatabaseManagerV2) SearchInFolder(folderPath, query string) models.SearchResult {
	criteria := models.SearchCriteria{
		Query:      query,
		FolderPath: folderPath,
	}
	return dm.database.Search(criteria)
}

func (dm *DatabaseManagerV2) SearchByTags(tags []string) models.SearchResult {
	criteria := models.SearchCriteria{
		Tags: tags,
	}
	return dm.database.Search(criteria)
}

// Database operations

func (dm *DatabaseManagerV2) GetDatabase() models.DatabaseV2 {
	return dm.database
}

func (dm *DatabaseManagerV2) Save() error {
	return dm.persistenceService.SaveDataV2(dm.database)
}

func (dm *DatabaseManagerV2) Reload() error {
	database, err := dm.persistenceService.LoadDataV2()
	if err != nil {
		return fmt.Errorf("failed to reload database: %w", err)
	}

	dm.database = database
	return nil
}

// Statistics and metadata

func (dm *DatabaseManagerV2) GetStatistics() map[string]interface{} {
	stats := make(map[string]interface{})

	stats["version"] = dm.database.Version
	stats["total_folders"] = len(dm.database.Folders)
	stats["total_items"] = len(dm.database.Items)

	// Count items by folder
	folderCounts := make(map[string]int)
	for _, item := range dm.database.Items {
		folderCounts[item.FolderPath]++
	}
	stats["items_by_folder"] = folderCounts

	// Count folders by depth
	depthCounts := make(map[int]int)
	for _, folder := range dm.database.Folders {
		depth := folder.GetDepth()
		depthCounts[depth]++
	}
	stats["folders_by_depth"] = depthCounts

	// Collect all tags
	tagCounts := make(map[string]int)
	for _, item := range dm.database.Items {
		for _, tag := range item.Tags {
			tagCounts[tag]++
		}
	}
	stats["tag_usage"] = tagCounts

	return stats
}

func (dm *DatabaseManagerV2) GetFolderTree() map[string]interface{} {
	tree := make(map[string]interface{})

	// Build folder hierarchy
	rootFolders := dm.database.GetSubfolders("/")
	tree["folders"] = dm.buildFolderTree(rootFolders)

	// Add root items
	rootItems := dm.database.GetItemsByFolder("/")
	tree["items"] = rootItems

	return tree
}

func (dm *DatabaseManagerV2) buildFolderTree(folders []models.FolderV2) []map[string]interface{} {
	var tree []map[string]interface{}

	for _, folder := range folders {
		folderNode := make(map[string]interface{})
		folderNode["folder"] = folder

		// Get subfolders
		subfolders := dm.database.GetSubfolders(folder.Path)
		folderNode["subfolders"] = dm.buildFolderTree(subfolders)

		// Get items in this folder
		items := dm.database.GetItemsByFolder(folder.Path)
		folderNode["items"] = items

		tree = append(tree, folderNode)
	}

	return tree
}

// Validation

func (dm *DatabaseManagerV2) ValidateDatabase() []string {
	var issues []string

	// Check for orphaned items (items in non-existent folders)
	for _, item := range dm.database.Items {
		if item.FolderPath != "/" {
			if _, found := dm.database.GetFolderByPath(item.FolderPath); !found {
				issues = append(issues, fmt.Sprintf("Item '%s' is in non-existent folder '%s'", item.Title, item.FolderPath))
			}
		}
	}

	// Check for invalid folder hierarchies
	for _, folder := range dm.database.Folders {
		if folder.ParentPath != "/" && folder.ParentPath != "" {
			if _, found := dm.database.GetFolderByPath(folder.ParentPath); !found {
				issues = append(issues, fmt.Sprintf("Folder '%s' has non-existent parent '%s'", folder.Path, folder.ParentPath))
			}
		}
	}

	// Check for duplicate folder paths
	pathCounts := make(map[string]int)
	for _, folder := range dm.database.Folders {
		pathCounts[folder.Path]++
	}
	for path, count := range pathCounts {
		if count > 1 {
			issues = append(issues, fmt.Sprintf("Duplicate folder path '%s' found %d times", path, count))
		}
	}

	return issues
}
