package main

import tea "github.com/charmbracelet/bubbletea"

type BookDetailsModel struct {
	cursor int
	book   Book
	ps *PlayerSave
}

func InitialBookDetailsModel(book Book, playersave *PlayerSave) BookDetailsModel {
	return BookDetailsModel{cursor: 0, book: book, ps: playersave}
}

func (m BookDetailsModel) Init() tea.Cmd {
	return nil
}

func (m BookDetailsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), tea.KeyEsc.String():
			switched := InitialDashboardModel(m.ps)
			return InitialRootModel().SwitchScreen(&switched)
		}
	}
	return m, nil
}

func (m BookDetailsModel) View() string {
	return "Booooooks"
}
