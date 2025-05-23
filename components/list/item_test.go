package list

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMyItem(t *testing.T) {
	title := "Test Title"
	desc := "Test Description"
	command := "echo 'test command'"
	now := time.Now()
	dateAdded := now.Add(-24 * time.Hour)
	dateUpdated := now.Add(-1 * time.Hour)

	item := NewMyItem(title, desc, command, dateAdded, dateUpdated)

	assert.Equal(t, title, item.title, "Title should match input")
	assert.Equal(t, desc, item.desc, "Description should match input")
	assert.Equal(t, command, item.command, "Command should match input")
	assert.True(t, dateAdded.Equal(item.dateAdded), "DateAdded should match input")
	assert.True(t, dateUpdated.Equal(item.dateUpdated), "DateUpdated should match input")
}

func TestMyItem_Getters(t *testing.T) {
	title := "Getter Title"
	desc := "Getter Description"
	command := "getter_command"
	now := time.Now()
	dateAdded := now.Round(time.Second) // Round to avoid potential sub-microsecond differences from some time sources
	dateUpdated := now.Add(time.Hour).Round(time.Second)

	item := NewMyItem(title, desc, command, dateAdded, dateUpdated)

	assert.Equal(t, title, item.Title(), "Title() getter failed")
	assert.Equal(t, desc, item.Description(), "Description() getter failed")
	assert.Equal(t, command, item.Command(), "Command() getter failed")
	assert.True(t, dateAdded.Equal(item.DateAdded()), "DateAdded() getter failed")
	assert.True(t, dateUpdated.Equal(item.DateUpdated()), "DateUpdated() getter failed")
}

func TestMyItem_FilterValue(t *testing.T) {
	title := "Filter Title"
	desc := "Some Desc"
	command := "filter_cmd"
	now := time.Now()

	item := NewMyItem(title, desc, command, now, now)

	assert.Equal(t, title, item.FilterValue(), "FilterValue() should return the title")
}
