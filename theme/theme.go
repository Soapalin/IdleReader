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
	PaddingTop(1)

var Heading2 = lipgloss.NewStyle().
	Foreground(HotPink)

var InputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7"))

var HotPink = lipgloss.Color("#FF06B7")
var Pink = lipgloss.Color("5")
var White = lipgloss.Color("#F8F8F8")

var CursorArrow = lipgloss.NewStyle().Foreground(Pink).Blink(true).Bold(true).Render(" ←")

var HelpIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#7B7B7B"))
var HelpText = lipgloss.NewStyle().Foreground(lipgloss.Color("#424242"))

var SpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

var ActiveDotPaginator = lipgloss.NewStyle().Foreground(Pink).Render("•")
var InactiveDotPaginator = lipgloss.NewStyle().Render("•")
