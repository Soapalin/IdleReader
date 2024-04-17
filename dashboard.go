package main

import (
	"game/engine/theme"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DashboardModel struct {
	tabs         []string
	activeTab    int
	ps           PlayerSave
	progress     []progress.Model
	errorMessage string
	cr_cursor    int
	bs_cursor    int
	bookChange   bool
}

func TabBorder(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.BottomRight = right
	border.Bottom = middle

	return border
}

var (
	inactiveTabBorder = TabBorder("*", "─", "*")
	activeTabBorder   = TabBorder("┴", " ", "┴")
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Padding(0, 1, 0, 1).Border(inactiveTabBorder, true).BorderForeground(highlightColor)
	activeTabStyle    = inactiveTabStyle.Copy().Background(highlightColor).Border(activeTabBorder, true)
)

func InitialDashboardModel(ps *PlayerSave) DashboardModel {
	tabs := []string{"My Bookshelf", "Current Reads", "Bookshop", "Library", "Book Club", "Help", "Exit"}
	prog := make([]progress.Model, 3)
	for i := range prog {
		prog[i] = progress.New(progress.WithDefaultGradient())
	}
	return DashboardModel{
		tabs:         tabs,
		activeTab:    1,
		ps:           *ps,
		progress:     prog,
		errorMessage: "",
		cr_cursor:    0,
		bs_cursor:    0,
		bookChange:   false,
	}
}

func (m *DashboardModel) HelpView() string {
	return "HelpView"
}

func (m *DashboardModel) LibraryView() string {
	return "LibraryView"
}

func (m *DashboardModel) BookClubView() string {
	return "BookClubView"
}

func (m *DashboardModel) ExitView() string {
	return "ExitView"
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m DashboardModel) Init() tea.Cmd {
	var cmd []tea.Cmd
	cmd = append(cmd, tickCmd())
	for i := range m.ps.Reader.CurrentReads.Books {
		cr_p := &m.ps.Reader.CurrentReads.Books[i]
		cmd = append(cmd, m.progress[i].IncrPercent(cr_p.Progress))
	}

	return tea.Batch(cmd...)
}
func (m *DashboardModel) ResetBookChangeState(msg string) {
	if m.bookChange && (msg != "r" && msg != "R") {
		m.bookChange = false
	}
}
func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.errorMessage = ""
		switch msg.String() {
		case tea.KeyCtrlC.String():
			m.ps.SavePlayerToFile()
			return m, tea.Quit
		case "q", tea.KeyEsc.String():
			if m.bookChange {
				m.bookChange = false
			} else {
				m.ps.SavePlayerToFile()
				switched := InitialSaveMenuModel()
				return InitialRootModel().SwitchScreen(&switched)
			}
		case tea.KeyTab.String():
			m.NextTab()
		case tea.KeyShiftTab.String():
			m.PreviousTab()
		case tea.KeyUp.String():
			switch m.activeTab {
			case 0:
				m.PreviousItemBookshelf()
			case 1:
				m.PreviousCurrentReads()
			case 2:
				m.ps.Shop.PreviousRow()
			}

		case tea.KeyDown.String():
			switch m.activeTab {
			case 0:
				m.NextItemBookshelf()
			case 1:
				m.NextCurrentReads()
			case 2:
				m.ps.Shop.NextRow()
			}
		case tea.KeyEnter.String():
			switch m.activeTab {
			case 2:
				m.TryBuy()
			case 0:
				if m.bs_cursor > len(m.ps.Library.Books) {
					// Item Details Screen
				} else {
					switched := InitialBookDetailsModel(m.ps.Library.Books[m.bs_cursor], &m.ps)
					return InitialRootModel().SwitchScreen(&switched)
				}

			}
		case "r", "R":
			switch m.activeTab {
			case 0:
				m.TrySwitchBook()
			}
		case "1", "2", "3":
			switch m.activeTab {
			case 0:
				return m.ChooseBook(msg.String())
			}

		}
		defer m.ResetBookChangeState(msg.String())
		// case tea.WindowSizeMsg:
		// 	log.Println("WindowSizeMsg")

		// 	for _, p := range m.progress {
		// 		p.Width = msg.Width - padding*2 - 4
		// if m.progress.Width > maxWidth {
		// 	m.progress.Width = maxWidth
		// }
		// return m, nil
		// 	}

	case tea.WindowSizeMsg:
		// m.width, m.height = msg.Width, msg.Height
		m.ps.Shop.table.Width(msg.Width)
		m.ps.Shop.table.Height(msg.Height)

	case tickMsg:
		var cmd []tea.Cmd
		cmd = append(cmd, tickCmd())

		for i := range m.ps.Reader.CurrentReads.Books {
			cr_p := &m.ps.Reader.CurrentReads.Books[i]
			if cr_p.Progress >= 1.0 {
				r := &m.ps.Reader
				r.FinishedBook(cr_p.ID)
				cmd = append(cmd, m.progress[i].DecrPercent(1))
				// log.Println("Progression | " + strconv.FormatFloat(cr_p.Progress, 'f', -1, 64))
			} else {
				cr_p.Progress += 0.05
				cmd = append(cmd, m.progress[i].IncrPercent(0.05))
				// log.Println("Progression | " + strconv.FormatFloat(cr_p.Progress, 'f', -1, 64))
			}
		}
		return m, tea.Batch(cmd...)
	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress[0].Update(msg)
		m.progress[0] = progressModel.(progress.Model)
		return m, cmd
	}
	return m, nil

}

func (m *DashboardModel) NextTab() {
	m.activeTab = (m.activeTab + 1) % len(m.tabs)
}

func (m *DashboardModel) PreviousTab() {
	m.activeTab -= 1
	if m.activeTab < 0 {
		m.activeTab = len(m.tabs) - 1
	}
}

func (m DashboardModel) TabsView() string {
	var tabRow []string
	for i, t := range m.tabs {
		if i == m.activeTab {
			tabRow = append(tabRow, activeTabStyle.Render(t))
		} else {
			tabRow = append(tabRow, inactiveTabStyle.Render(t))
		}
	}
	return "\n" + lipgloss.JoinHorizontal(lipgloss.Top, tabRow...) + "\n"
}

func (m *DashboardModel) View() string {
	r := m.ps.Reader
	s := theme.Heading1.Render(r.Name)
	s += "\n"
	s += "\n"
	s += "Favourite Book: " + r.FavouriteBook + ", " + r.FavouriteAuthor
	s += "\n"

	k := "Knowledge: " + strconv.Itoa(r.Knowledge)
	iq := "IQ: " + strconv.Itoa(r.IQ) + " (" + r.IQ_Title() + ")"

	p := "Prestige: " + strconv.Itoa(r.Prestige)
	rs := "Reading Speed: " + strconv.Itoa(r.Prestige)

	v1 := lipgloss.NewStyle().Padding(0, 10, 0, 0).Render(lipgloss.JoinVertical(lipgloss.Left, k, iq))
	v2 := lipgloss.JoinVertical(lipgloss.Left, rs, p)
	s += lipgloss.JoinHorizontal(lipgloss.Center, v1, v2)
	s += "\n"

	s += m.TabsView()

	s += "\n"
	switch m.activeTab {
	case 0:
		s += m.BookshelfView()
	case 1:
		s += m.CurrentReadsView()
	case 2:
		m.ps.Shop.Update()
		s += m.BookshopView()
	case 3:
		s += m.LibraryView()
	case 4:
		s += m.BookClubView()
	case 5:
		s += m.HelpView()
	case 6:
		s += m.ExitView()
	}
	s += "\n"

	s += m.errorMessage

	return s
}
