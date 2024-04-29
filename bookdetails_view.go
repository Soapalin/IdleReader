package main

import (
	"game/engine/theme"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type ObjectDetailsModel struct {
	cursor int
	object any
	dm     *DashboardModel
}

func InitialBookDetailsModel(object any, dm *DashboardModel) ObjectDetailsModel {
	return ObjectDetailsModel{cursor: 0, object: object, dm: dm}
}

func (m ObjectDetailsModel) Init() tea.Cmd {
	return nil
}

func (m ObjectDetailsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), tea.KeyEsc.String():
			return InitialRootModel().SwitchScreen(m.dm)
		}
	}
	return m, nil
}

func (m ObjectDetailsModel) View() string {
	s := ""
	switch obj := m.object.(type) {
	case Book:
		s += theme.Heading1.Render(obj.Name)
		s += "\n"
		s += theme.Heading2.Render(obj.Author)
		s += "\n\n"
		s += "Knowledge Gain: " + strconv.Itoa(obj.KnowledgeIncrease) + "\n"
		s += "This book increases your IQ by " + strconv.Itoa(obj.IntelligenceIncrease) + " on your first read.\n"

	case Item:
		s += theme.Heading1.Render(obj.Name)
		s += "\n\n"
		s += obj.Description
		s += obj.Effect
	}
	return s + "\n\n" + theme.HelpIcon.Render("esc") + theme.HelpText.Render(" back")

}
