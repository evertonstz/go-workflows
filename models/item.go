package models

import (
	"fmt"
	"strings"
	"time"
)

type (
	Item struct {
		Title, Desc, Command string
		DateAdded            time.Time
		DateUpdated          time.Time
	}

	Items struct {
		Items []Item
	}

	// V2 Database structures with folder support
	ItemV2 struct {
		ID          string            `json:"id"`
		Title       string            `json:"title"`
		Desc        string            `json:"description"`
		Command     string            `json:"command"`
		DateAdded   time.Time         `json:"date_added"`
		DateUpdated time.Time         `json:"date_updated"`
		Tags        []string          `json:"tags,omitempty"`
		Metadata    map[string]string `json:"metadata,omitempty"`
		FolderPath  string            `json:"folder_path"` // Path to the folder containing this item
	}

	FolderV2 struct {
		ID          string            `json:"id"`
		Name        string            `json:"name"`
		Description string            `json:"description,omitempty"`
		Path        string            `json:"path"`                  // Full path from root (e.g., "/folder1/subfolder")
		ParentPath  string            `json:"parent_path,omitempty"` // Path to parent folder
		DateAdded   time.Time         `json:"date_added"`
		DateUpdated time.Time         `json:"date_updated"`
		Metadata    map[string]string `json:"metadata,omitempty"`
	}

	DatabaseV2 struct {
		Version string     `json:"version"`
		Folders []FolderV2 `json:"folders"`
		Items   []ItemV2   `json:"items"`
	}

	// Search and filter structures
	SearchCriteria struct {
		Query      string     `json:"query,omitempty"`
		FolderPath string     `json:"folder_path,omitempty"`
		Tags       []string   `json:"tags,omitempty"`
		DateFrom   *time.Time `json:"date_from,omitempty"`
		DateTo     *time.Time `json:"date_to,omitempty"`
	}

	SearchResult struct {
		Items   []ItemV2   `json:"items"`
		Folders []FolderV2 `json:"folders"`
		Total   int        `json:"total"`
	}
)

// ItemV2 methods
func (i *ItemV2) GenerateID() {
	if i.ID == "" {
		i.ID = fmt.Sprintf("item_%d", time.Now().UnixNano())
	}
}

func (i ItemV2) GetFullPath() string {
	if i.FolderPath == "" || i.FolderPath == "/" {
		return i.Title
	}
	return strings.TrimSuffix(i.FolderPath, "/") + "/" + i.Title
}

func (i ItemV2) MatchesSearch(criteria SearchCriteria) bool {
	// Search in title, description, and command
	if criteria.Query != "" {
		query := strings.ToLower(criteria.Query)
		if !strings.Contains(strings.ToLower(i.Title), query) &&
			!strings.Contains(strings.ToLower(i.Desc), query) &&
			!strings.Contains(strings.ToLower(i.Command), query) {
			return false
		}
	}

	// Filter by folder path
	if criteria.FolderPath != "" && i.FolderPath != criteria.FolderPath {
		return false
	}

	// Filter by tags
	if len(criteria.Tags) > 0 {
		hasTag := false
		for _, tag := range criteria.Tags {
			for _, itemTag := range i.Tags {
				if strings.EqualFold(tag, itemTag) {
					hasTag = true
					break
				}
			}
			if hasTag {
				break
			}
		}
		if !hasTag {
			return false
		}
	}

	// Filter by date range
	if criteria.DateFrom != nil && i.DateAdded.Before(*criteria.DateFrom) {
		return false
	}
	if criteria.DateTo != nil && i.DateAdded.After(*criteria.DateTo) {
		return false
	}

	return true
}

// FolderV2 methods
func (f *FolderV2) GenerateID() {
	if f.ID == "" {
		f.ID = fmt.Sprintf("folder_%d", time.Now().UnixNano())
	}
}

func (f FolderV2) IsRoot() bool {
	return f.Path == "/" || f.Path == ""
}

func (f FolderV2) GetDepth() int {
	if f.IsRoot() {
		return 0
	}
	return len(strings.Split(strings.Trim(f.Path, "/"), "/"))
}

func (f FolderV2) IsChildOf(parentPath string) bool {
	if parentPath == "/" {
		parentPath = ""
	}
	return f.ParentPath == parentPath
}

func (f FolderV2) MatchesSearch(criteria SearchCriteria) bool {
	if criteria.Query != "" {
		query := strings.ToLower(criteria.Query)
		if !strings.Contains(strings.ToLower(f.Name), query) &&
			!strings.Contains(strings.ToLower(f.Description), query) {
			return false
		}
	}
	return true
}

// DatabaseV2 methods
func NewDatabaseV2() DatabaseV2 {
	return DatabaseV2{
		Version: "2.0",
		Folders: []FolderV2{},
		Items:   []ItemV2{},
	}
}

func (db *DatabaseV2) AddFolder(folder FolderV2) error {
	// Ensure ID is set
	folder.GenerateID()

	// Validate folder path
	if folder.Path == "" {
		return fmt.Errorf("folder path cannot be empty")
	}

	// Check for duplicate paths
	for _, existingFolder := range db.Folders {
		if existingFolder.Path == folder.Path {
			return fmt.Errorf("folder with path %s already exists", folder.Path)
		}
	}

	db.Folders = append(db.Folders, folder)
	return nil
}

func (db *DatabaseV2) AddItem(item ItemV2) error {
	// Ensure ID is set
	item.GenerateID()

	// Validate that folder exists if specified
	if item.FolderPath != "" && item.FolderPath != "/" {
		found := false
		for _, folder := range db.Folders {
			if folder.Path == item.FolderPath {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("folder %s does not exist", item.FolderPath)
		}
	}

	db.Items = append(db.Items, item)
	return nil
}

func (db DatabaseV2) GetItemsByFolder(folderPath string) []ItemV2 {
	var items []ItemV2
	for _, item := range db.Items {
		if item.FolderPath == folderPath {
			items = append(items, item)
		}
	}
	return items
}

func (db DatabaseV2) GetSubfolders(parentPath string) []FolderV2 {
	var subfolders []FolderV2
	for _, folder := range db.Folders {
		if folder.IsChildOf(parentPath) {
			subfolders = append(subfolders, folder)
		}
	}
	return subfolders
}

func (db DatabaseV2) Search(criteria SearchCriteria) SearchResult {
	var matchingItems []ItemV2
	var matchingFolders []FolderV2

	// Search items
	for _, item := range db.Items {
		if item.MatchesSearch(criteria) {
			matchingItems = append(matchingItems, item)
		}
	}

	// Search folders
	for _, folder := range db.Folders {
		if folder.MatchesSearch(criteria) {
			matchingFolders = append(matchingFolders, folder)
		}
	}

	return SearchResult{
		Items:   matchingItems,
		Folders: matchingFolders,
		Total:   len(matchingItems) + len(matchingFolders),
	}
}

func (db DatabaseV2) GetFolderByPath(path string) (*FolderV2, bool) {
	for _, folder := range db.Folders {
		if folder.Path == path {
			return &folder, true
		}
	}
	return nil, false
}

func (db DatabaseV2) GetItemByID(id string) (*ItemV2, bool) {
	for _, item := range db.Items {
		if item.ID == id {
			return &item, true
		}
	}
	return nil, false
}

func (db *DatabaseV2) UpdateItem(id string, updatedItem ItemV2) error {
	for i, item := range db.Items {
		if item.ID == id {
			updatedItem.ID = id // Preserve original ID
			updatedItem.DateUpdated = time.Now()
			db.Items[i] = updatedItem
			return nil
		}
	}
	return fmt.Errorf("item with ID %s not found", id)
}

func (db *DatabaseV2) DeleteItem(id string) error {
	for i, item := range db.Items {
		if item.ID == id {
			db.Items = append(db.Items[:i], db.Items[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("item with ID %s not found", id)
}

func (db *DatabaseV2) DeleteFolder(path string) error {
	// Check if folder has items
	items := db.GetItemsByFolder(path)
	if len(items) > 0 {
		return fmt.Errorf("cannot delete folder %s: contains %d items", path, len(items))
	}

	// Check if folder has subfolders
	subfolders := db.GetSubfolders(path)
	if len(subfolders) > 0 {
		return fmt.Errorf("cannot delete folder %s: contains %d subfolders", path, len(subfolders))
	}

	// Delete the folder
	for i, folder := range db.Folders {
		if folder.Path == path {
			db.Folders = append(db.Folders[:i], db.Folders[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("folder with path %s not found", path)
}

// Migration functions for backward compatibility
func MigrateV1ToV2(v1Data Items) DatabaseV2 {
	db := NewDatabaseV2()

	for _, v1Item := range v1Data.Items {
		v2Item := ItemV2{
			Title:       v1Item.Title,
			Desc:        v1Item.Desc,
			Command:     v1Item.Command,
			DateAdded:   v1Item.DateAdded,
			DateUpdated: v1Item.DateUpdated,
			FolderPath:  "/", // All v1 items go to root
			Tags:        []string{},
			Metadata:    make(map[string]string),
		}
		v2Item.GenerateID()
		db.Items = append(db.Items, v2Item)
	}

	return db
}

func (db DatabaseV2) ToV1() Items {
	var v1Items []Item
	for _, v2Item := range db.Items {
		v1Items = append(v1Items, Item{
			Title:       v2Item.Title,
			Desc:        v2Item.Desc,
			Command:     v2Item.Command,
			DateAdded:   v2Item.DateAdded,
			DateUpdated: v2Item.DateUpdated,
		})
	}
	return Items{Items: v1Items}
}
