package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type RootScreenModel struct {
	Model tea.Model
	// choices  []string         // items on the to-do list
	// cursor   int              // which to-do list item our cursor is pointing at
	// selected map[int]struct{} // which to-do items are selected
}

func InitialRootModel() RootScreenModel {
	var rootModel tea.Model
	saveMenu := InitialSaveMenuModel()
	rootModel = &saveMenu

	return RootScreenModel{
		Model: rootModel,
	}
}

func (m RootScreenModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return m.Model.Init()
}

func (m RootScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.Model.Update(msg)
}

func (m RootScreenModel) View() string {
	return m.Model.View()
}

func (m RootScreenModel) SwitchScreen(model tea.Model) (tea.Model, tea.Cmd) {
	m.Model = model
	return m.Model, m.Model.Init()
}
