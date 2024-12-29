package models

import "time"

type (
	Item struct {
		Title, Desc, Command string
		DateAdded            time.Time
		DateUpdated          time.Time
	}

	Items struct {
		Items []Item
	}
)
