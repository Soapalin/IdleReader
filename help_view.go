package main

import (
	"game/engine/theme"
	"strconv"

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
	{"How do I purchcase books?", "Purchase your books in the Bookshop or Auction. Bookshop books are refreshed and totally randomised while you can look for your favourtie book in the Auction window. Beware, prices are higher!"},
	{"How do I increase my IQ?", "IQ only increases by a small amount on your first read! You will need to purchase and read different books in order to increase it."},
	{"How do I read the book I have purchased?", "Head to the your Bookshelf and press 'r' on the book you would like to read."},
	{"Where can I find an in-depth game guide?", "Head to our official wiki," + lipgloss.NewStyle().Foreground(theme.Pink).Render("{https://placeholder.com}")},
}

func (m *DashboardModel) HelpView() string {
	s := ""
	start, end := m.helpPaginator.GetSliceBounds(len(HelpSection))
	for _, item := range HelpSection[start:end] {
		s += item.View(m.width)
		s += "\n"
	}
	paginatorView := lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Position(0.5)).Render(m.helpPaginator.View())
	pageNumber := "<" + strconv.Itoa(m.helpPaginator.Page+1) + "/" + strconv.Itoa(m.helpPaginator.TotalPages) + ">"
	paginatorFull := lipgloss.JoinVertical(0, paginatorView, pageNumber)
	s += lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Position(0.5)).Render(paginatorFull)
	s += "\n" + theme.HelpIcon.Render("←") + theme.HelpText.Render(" previous • ")
	s += theme.HelpIcon.Render("→") + theme.HelpText.Render(" next • ")
	s += theme.HelpIcon.Render("tab/shift+tab") + theme.HelpText.Render(" switch tabs • ")
	s += theme.HelpIcon.Render("esc / q") + theme.HelpText.Render(" quit")
	return s
}

func (m *DashboardModel) NextHelpItem() {
	m.helpPaginator.Page++
	if m.helpPaginator.Page >= m.helpPaginator.TotalPages {
		m.helpPaginator.Page = 0
	}
}

func (m *DashboardModel) PreviousHelpItem() {
	m.helpPaginator.Page--
	if m.helpPaginator.Page < 0 {
		m.helpPaginator.Page = 0
	}
}
