package main

import (
	"github.com/atotto/clipboard"
	"github.com/evertonstz/go-workflows/components/list"
	"github.com/evertonstz/go-workflows/shared"
)

func handleClipboardCopy(t string) {
	err := clipboard.WriteAll(t)
	if err != nil {
		panic(err)
	}
}

func handleSaveItem(m model, msg shared.SaveItem) {
	// TODO: save to sqlite
	r, _ := m.list.Update(msg)
	m.list = r.(list.Model)
}