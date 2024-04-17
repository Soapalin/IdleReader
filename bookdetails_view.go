package main

import tea "github.com/charmbracelet/bubbletea"

type BookDetailsModel struct {
	cursor int
	book   Book
}

func InitialBookDetailsModel(book Book) BookDetailsModel {
	return BookDetailsModel{cursor: 0, book: book}
}

func (m BookDetailsModel) Init() tea.Cmd {
	return nil
}

func (m BookDetailsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m BookDetailsModel) View() string {
	return "Booooooks"
}
