package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/components/list"
	"github.com/evertonstz/go-workflows/components/persist"
	"github.com/evertonstz/go-workflows/models"
	"github.com/evertonstz/go-workflows/shared"
)

func handleSaveItem(m model, msg shared.SaveCommandMsg) (model, tea.Cmd) {
	r, _ := m.list.Update(msg)
	m.list = r.(list.Model)

	var items []models.Item
	for _, i := range m.list.AllItems() {
		items = append(items, models.Item{
			Title:       i.Title(),
			Desc:        i.Description(),
			Command:     i.Command(),
			DateAdded:   i.DateAdded(),
			DateUpdated: i.DateUpdated()})
	}
	data := models.Items{Items: items}

	return m, persist.SaveConfigFile(m.persistPath, data)
}
