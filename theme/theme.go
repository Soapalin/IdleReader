package theme

import "github.com/charmbracelet/lipgloss"

var Style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(22)

var Heading1 = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("5")).
	PaddingTop(2)

var Heading2 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("5"))

var InputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7"))

var Pink = lipgloss.Color("5")

var CursorArrow = lipgloss.NewStyle().Foreground(Pink).Blink(true).Bold(true).Render(" ←")