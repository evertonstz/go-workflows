package models

import "time"

type Item struct {
	Title, Desc, Command string
	DateAdded   time.Time
	DateUpdated time.Time
}

type Items struct {
	Items []Item
}