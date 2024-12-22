package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertonstz/go-workflows/components/form"
	"github.com/evertonstz/go-workflows/components/list"
	"github.com/evertonstz/go-workflows/components/footer"
	"github.com/evertonstz/go-workflows/shared"
)

type model struct {
	list   tea.Model
	form   form.Model
	footer footer.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdList, cmdForm, cmdFooter tea.Cmd

	switch msg := msg.(type) {
    case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit}
	case shared.SelectedItemMsg:
		m.form, cmdForm = m.form.Update(msg)
	}

	m.list, cmdList = m.list.Update(msg)
	m.footer, cmdFooter = m.footer.Update(msg)

	return m, tea.Batch(cmdList, cmdForm, cmdFooter)
}

func (m model) View() string {
	return m.list.View() + "\n" + m.form.View() + "\n" + m.footer.View()
}

func main() {
	m := model{
		list:   list.New(),
		form:   form.New(),
		footer: footer.New(),
	}

	p := tea.NewProgram(m)

	if err := p.Start(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}
}
