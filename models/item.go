package models

import "time"

type (
	// Item now implements list.Item interface for bubbletea/list
	Item struct {
		Title, Desc, Command string
		DateAdded            time.Time
		DateUpdated          time.Time
	}

	Items struct {
		Items []Item
	}
)

// FilterValue makes Item implement the list.Item interface.
func (i Item) FilterValue() string { return i.Title }

// Title returns the item's title.
func (i Item) GetTitle() string { return i.Title }

// Description returns the item's description.
func (i Item) GetDescription() string { return i.Desc }
