# Go Workflows v2 Database

This document describes the v2 database implementation with folder support for the Go Workflows application.

## Overview

The v2 database is a significant upgrade from the original flat structure, introducing:

- **Hierarchical folder organization** - Organize workflows in nested folders
- **Enhanced search capabilities** - Search within specific folders and by tags
- **Rich metadata support** - Store custom metadata with each item
- **Backward compatibility** - Automatic migration from v1 to v2
- **Flexible structure** - Designed for future extensions

## Key Features

### 1. Version Management

- Database version field set to "2.0"
- Automatic detection of database version
- Seamless migration from v1 to v2 format
- Backup creation during migration

### 2. Folder Structure

- Hierarchical folder organization with unlimited nesting
- Each folder has a unique path (e.g., `/dev/scripts/build`)
- Folders contain metadata: name, description, creation/modification dates
- Support for folder-specific operations

### 3. Enhanced Items

- Items can be organized within folders
- Rich metadata support with key-value pairs
- Tag system for categorization
- Unique IDs for reliable item management
- Full-text search across title, description, and command

### 4. Search and Filtering

- Search by text query across all fields
- Filter by folder path
- Filter by tags
- Date range filtering
- Combined search criteria

## Database Structure

### DatabaseV2

```go
type DatabaseV2 struct {
    Version string     `json:"version"`        // Always "2.0"
    Folders []FolderV2 `json:"folders"`        // Folder definitions
    Items   []ItemV2   `json:"items"`          // Workflow items
}
```

### FolderV2

```go
type FolderV2 struct {
    ID          string            `json:"id"`           // Unique identifier
    Name        string            `json:"name"`         // Display name
    Description string            `json:"description"`  // Optional description
    Path        string            `json:"path"`         // Full path from root
    ParentPath  string            `json:"parent_path"`  // Parent folder path
    DateAdded   time.Time         `json:"date_added"`   // Creation timestamp
    DateUpdated time.Time         `json:"date_updated"` // Last modification
    Metadata    map[string]string `json:"metadata"`     // Custom metadata
}
```

### ItemV2

```go
type ItemV2 struct {
    ID          string            `json:"id"`           // Unique identifier
    Title       string            `json:"title"`        // Item title
    Desc        string            `json:"description"`  // Description
    Command     string            `json:"command"`      // Command to execute
    DateAdded   time.Time         `json:"date_added"`   // Creation timestamp
    DateUpdated time.Time         `json:"date_updated"` // Last modification
    Tags        []string          `json:"tags"`         // Categorization tags
    Metadata    map[string]string `json:"metadata"`     // Custom metadata
    FolderPath  string            `json:"folder_path"`  // Container folder path
}
```

## Usage Examples

### Basic Operations

#### Creating Folders

```go
// Create root level folder
folder, err := manager.CreateFolder("scripts", "Build scripts", "/")

// Create subfolder
subfolder, err := manager.CreateFolder("deployment", "Deployment scripts", "/scripts")
```

#### Creating Items

```go
item, err := manager.CreateItem(
    "Deploy to Production",
    "Deploy application to production environment",
    "./deploy.sh production",
    "/scripts/deployment",
    []string{"deploy", "production", "critical"},
    map[string]string{
        "author": "ops-team",
        "approval": "required",
    },
)
```

#### Searching

```go
// Search by text
result := manager.Search(models.SearchCriteria{
    Query: "deploy",
})

// Search within folder
result := manager.SearchInFolder("/scripts", "build")

// Search by tags
result := manager.SearchByTags([]string{"production", "critical"})

// Complex search
result := manager.Search(models.SearchCriteria{
    Query:      "test",
    FolderPath: "/scripts",
    Tags:       []string{"automation"},
})
```

### Advanced Features

#### Database Statistics

```go
stats := manager.GetStatistics()
fmt.Printf("Total Folders: %d\n", stats["total_folders"])
fmt.Printf("Total Items: %d\n", stats["total_items"])

// Items by folder
itemsByFolder := stats["items_by_folder"].(map[string]int)

// Tag usage statistics
tagUsage := stats["tag_usage"].(map[string]int)
```

#### Folder Tree Structure

```go
tree := manager.GetFolderTree()
// Returns hierarchical representation of all folders and items
```

#### Database Validation

```go
issues := manager.ValidateDatabase()
if len(issues) == 0 {
    fmt.Println("Database is valid")
} else {
    for _, issue := range issues {
        fmt.Printf("Issue: %s\n", issue)
    }
}
```

## Migration from v1 to v2

The migration process is automatic and safe:

1. **Detection**: System automatically detects v1 format
2. **Backup**: Creates backup file with `.v1.backup` extension
3. **Migration**: Converts all v1 items to v2 format
4. **Placement**: All migrated items are placed in the root folder (`/`)
5. **IDs**: Generates unique IDs for all migrated items

### Manual Migration

```go
// Check current version
version, err := persistence.GetDatabaseVersion()

// Migrate to v2
err = persistence.MigrateToV2()
```

### Loading v2 Data

```go
// Load as v2 (auto-migrates if needed)
database, err := persistence.LoadDataV2()

// Load as v1 (for backward compatibility)
items, err := persistence.LoadData()
```

## CLI Demo Tool

A comprehensive CLI demo is available at `examples/v2_cli_demo.go`:

```bash
# Create folder structure
go run examples/v2_cli_demo.go -action=create-folder -name=scripts -desc="Development scripts" -parent=/

# Create workflow item
go run examples/v2_cli_demo.go -action=create-item -name="Deploy Script" -desc="Deploy to production" -command="deploy.sh" -folder=/scripts -tags="deploy,production"

# List folder contents
go run examples/v2_cli_demo.go -action=list -folder=/scripts

# Search workflows
go run examples/v2_cli_demo.go -action=search -query=deploy

# Get database statistics
go run examples/v2_cli_demo.go -action=stats -pretty

# Validate database
go run examples/v2_cli_demo.go -action=validate

# Export database
go run examples/v2_cli_demo.go -action=export -pretty
```

## Best Practices

### Folder Organization

- Use clear, descriptive folder names
- Organize by purpose: `/development`, `/production`, `/maintenance`
- Keep folder hierarchies reasonable (3-4 levels max)
- Use consistent naming conventions

### Item Management

- Use descriptive titles and descriptions
- Leverage tags for cross-cutting concerns
- Store relevant metadata (author, environment, etc.)
- Keep commands concise but complete

### Search Strategy

- Use specific queries for better results
- Combine search criteria for precision
- Leverage folder-scoped searches
- Use tags for categorical searches

## Performance Considerations

- **In-Memory Operations**: All operations are performed in memory
- **Batch Updates**: Group multiple changes before saving
- **Search Optimization**: Search is performed on all items/folders
- **File I/O**: Only occurs during load/save operations

## Backward Compatibility

The v2 implementation maintains full backward compatibility:

- **Reading**: Can read both v1 and v2 formats
- **Migration**: Automatic and transparent
- **Fallback**: v2 data can be converted back to v1 if needed
- **API**: Original v1 API methods still work

## Future Extensions

The v2 structure is designed to support future enhancements:

- **Permissions**: Folder-level access control
- **Templates**: Workflow templates per folder
- **History**: Track item modification history
- **Sync**: Folder-based synchronization
- **Export**: Selective export by folder

## Testing

Comprehensive test suites are available:

- `models/item_v2_test.go` - Core model tests
- `shared/di/services/persistence_v2_test.go` - Persistence tests
- `shared/di/services/database_manager_v2_test.go` - Manager tests

Run tests with:

```bash
go test ./models -v
go test ./shared/di/services -v
```

## Schema Example

Sample v2 database JSON structure:

```json
{
  "version": "2.0",
  "folders": [
    {
      "id": "folder_1640995200000000000",
      "name": "scripts",
      "description": "Development scripts",
      "path": "/scripts",
      "parent_path": "/",
      "date_added": "2024-01-01T12:00:00Z",
      "date_updated": "2024-01-01T12:00:00Z",
      "metadata": {}
    }
  ],
  "items": [
    {
      "id": "item_1640995200000000001",
      "title": "Deploy Script",
      "description": "Deploy to production",
      "command": "deploy.sh production",
      "date_added": "2024-01-01T12:00:00Z",
      "date_updated": "2024-01-01T12:00:00Z",
      "tags": ["deploy", "production"],
      "metadata": {
        "author": "ops-team",
        "environment": "production"
      },
      "folder_path": "/scripts"
    }
  ]
}
```

This v2 implementation provides a robust, scalable foundation for organizing and managing workflows while maintaining compatibility with existing v1 databases.
