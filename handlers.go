package main

import (
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/components/list"
	"github.com/evertonstz/go-workflows/components/persist"
	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared"
)

func handleClipboardCopy(t string) {
	err := clipboard.WriteAll(t)
	if err != nil {
		panic(err)
	}
}

func handleSaveItem(m model, msg shared.SaveItem) (model, tea.Cmd) {
	// TODO: save to sqlite
	r, _ := m.list.Update(msg)
	m.list = r.(list.Model)

	// iterate over the list and save to file
	var items []models.Item
	for _, i := range m.list.AllItems() {
		items = append(items, models.Item{Title: i.Title(), 
			Desc: i.Description(), 
			Command: i.Command(), 
			DateAdded: i.DateAdded(), 
			DateUpdated: i.DateUpdated()})
	}
	data := models.Items{Items: items}
	//return d inside a list
	return m, persist.SaveConfigFile(m.persistPath, data)
}