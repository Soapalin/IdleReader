package main

import (
	"game/engine/theme"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type BookDetailsModel struct {
	cursor int
	book   Book
	dm     *DashboardModel
}

func InitialBookDetailsModel(book Book, dm *DashboardModel) BookDetailsModel {
	return BookDetailsModel{cursor: 0, book: book, dm: dm}
}

func (m BookDetailsModel) Init() tea.Cmd {
	return nil
}

func (m BookDetailsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), tea.KeyEsc.String():
			return InitialRootModel().SwitchScreen(m.dm)
		}
	}
	return m, nil
}

func (m BookDetailsModel) View() string {
	s := ""
	s += theme.Heading1.Render(m.book.Name)
	s += "\n"
	s += theme.Heading2.Render(m.book.Author)
	s += "\n\n"
	s += "Knowledge Gain: " + strconv.Itoa(m.book.KnowledgeIncrease) + "\n"
	s += "This book increases your IQ by " + strconv.Itoa(m.book.IntelligenceIncrease) + " on your first read.\n"

	return s
}
