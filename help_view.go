package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type HelpItem struct {
	question string
	answer   string
}

func (h *HelpItem) View(maxWidth int) string {
	s := ""
	s += wordwrap.String(lipgloss.NewStyle().Padding(1).Render("Q: "+h.question), maxWidth-10)
	s += "\n"
	s += wordwrap.String(lipgloss.NewStyle().Padding(1).Render("A: "+h.answer), maxWidth-10)

	return lipgloss.NewStyle().Width(maxWidth - 2).Border(lipgloss.NormalBorder()).Render(s)
}

var HelpSection []HelpItem = []HelpItem{
	{"How do I play this game?", "Purchase, Read, Gain Knowlege and IQ, and finally collect your favourite books!"},
	{"Where can I find an in-depth game guide?", "Head to our official website, {https://placeholder.com}"},
}

func (m *DashboardModel) HelpView() string {
	s := ""
	start, end := m.paginator.GetSliceBounds(len(HelpSection))
	for _, item := range HelpSection[start:end] {
		s += item.View(m.width)
		s += "\n"
	}
	s += lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Position(0.5)).Render(m.paginator.View())
	return s
}

func (m *DashboardModel) NextHelpItem() {
	m.paginator.Page++
	if m.paginator.Page >= m.paginator.TotalPages {
		m.paginator.Page = 0
	}
}

func (m *DashboardModel) PreviousHelpItem() {
	m.paginator.Page--
	if m.paginator.Page < 0 {
		m.paginator.Page = 0
	}
}
