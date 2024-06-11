package main

import (
	"game/engine/theme"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type DashboardModel struct {
	tabs             []string
	activeTab        int
	ps               PlayerSave
	progress         []progress.Model
	errorMessage     string
	cr_cursor        int
	bs_cursor        int
	i_cursor         int
	auc_cursor       int
	bookChange       bool
	width            int
	height           int
	spinner          spinner.Model
	helpPaginator    paginator.Model
	bookPaginator    paginator.Model
	auctionPaginator paginator.Model
	exitChoices      []string
	ex_cursor        int
	auction_inputs   []textinput.Model
	bookshelf_input  textinput.Model
	auctionLibrary   Library
	bookshelfLibrary Library
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

func InitialDashboardModel(ps *PlayerSave, activeTab int, bs_cursor int, i_cursor int) DashboardModel {
	tabs := []string{"My Bookshelf", "Current Reads", "Bookshop", "Auction", "Inventory", "Help", "Exit"}
	prog := make([]progress.Model, 3)
	for i := range prog {
		prog[i] = progress.New(progress.WithDefaultGradient())
	}
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = theme.SpinnerStyle
	s.Spinner.FPS = time.Millisecond * 500

	hp := paginator.New()
	hp.Type = paginator.Dots
	hp.PerPage = 1
	hp.ActiveDot = theme.ActiveDotPaginator
	hp.InactiveDot = theme.InactiveDotPaginator
	hp.SetTotalPages(len(HelpSection))

	bp := paginator.New()
	bp.Type = paginator.Dots
	bp.PerPage = 5
	bp.ActiveDot = theme.ActiveDotPaginator
	bp.InactiveDot = theme.InactiveDotPaginator
	bp.SetTotalPages(len(ps.Reader.Library.Books))

	exitChoices := []string{"Save & Go to Main Menu", "Save & Quit", "Quit without Saving"}

	var inputs []textinput.Model = make([]textinput.Model, 2)

	inputs[0] = textinput.New()
	inputs[0].CharLimit = 255
	inputs[0].Width = 40
	inputs[0].Prompt = ""

	inputs[1] = textinput.New()
	inputs[1].CharLimit = 255
	inputs[1].Width = 40
	inputs[1].Prompt = ""

	bsInput := textinput.New()
	bsInput.CharLimit = 255
	bsInput.Width = 40
	bsInput.Prompt = ""

	ap := paginator.New()
	ap.Type = paginator.Dots
	ap.PerPage = 5
	ap.ActiveDot = theme.ActiveDotPaginator
	ap.InactiveDot = theme.InactiveDotPaginator
	ap.SetTotalPages(1)

	return DashboardModel{
		tabs:             tabs,
		activeTab:        activeTab,
		ps:               *ps,
		progress:         prog,
		errorMessage:     "",
		cr_cursor:        0,
		bs_cursor:        bs_cursor,
		i_cursor:         i_cursor,
		bookChange:       false,
		width:            w,
		height:           h,
		spinner:          s,
		helpPaginator:    hp,
		bookPaginator:    bp,
		exitChoices:      exitChoices,
		ex_cursor:        0,
		auction_inputs:   inputs,
		auctionPaginator: ap,
		bookshelf_input:  bsInput,
		auctionLibrary:   Library{},
		bookshelfLibrary: ps.Reader.Library,
	}
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *DashboardModel) updateSize(w, h int) {
	log.Println("updateSize")
}

func (m DashboardModel) Init() tea.Cmd {
	var cmd []tea.Cmd
	cmd = append(cmd, tickCmd())
	for i, id := range m.ps.Reader.CurrentReads.BookIDs {
		cr_p, err := m.ps.Reader.Library.GetBookPointerByID(id)
		if err != nil {
			panic(err)

		}
		cmd = append(cmd, m.progress[i].SetPercent(cr_p.Progress))
	}

	cmd = append(cmd, m.spinner.Tick)

	return tea.Batch(cmd...)
}
func (m *DashboardModel) ResetBookChangeState(msg string) {
	if m.bookChange && (msg != tea.KeyEnter.String()) {
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
			} else if m.bookshelf_input.Focused() {
				m.bookshelf_input.Blur()
			} else if m.auction_inputs[0].Focused() || m.auction_inputs[1].Focused() {
				// do nothing
			} else {
				m.ps.SavePlayerToFile()
				switched := InitialSaveMenuModel()
				return InitialRootModel().SwitchScreen(&switched)
			}
		case "b":
			switch m.activeTab {
			case 3:
				m.TryAuctionBuy()
			}
		case "i":
			switch m.activeTab {
			case 0:
				if !m.bookshelf_input.Focused() {
					if !m.bookChange {
						switched := InitialBookDetailsModel(m.ps.Reader.Library.Books[m.bs_cursor], m)
						return InitialRootModel().SwitchScreen(&switched)
					}
				}

			case 2:
				if m.ps.Shop.TableIndex > len(m.ps.Shop.Books.Books) {
					switched := InitialBookDetailsModel(m.ps.Shop.Items.Items[m.ps.Shop.TableIndex-len(m.ps.Shop.Books.Books)-1], m)
					return InitialRootModel().SwitchScreen(&switched)
				} else {
					switched := InitialBookDetailsModel(m.ps.Shop.Books.Books[m.ps.Shop.TableIndex-1], m)
					return InitialRootModel().SwitchScreen(&switched)
				}
			case 3:
				if !m.AuctionInputsFocused() {
					if len(m.auctionLibrary.Books) > 0 {
						switched := InitialBookDetailsModel(m.auctionLibrary.Books[m.auc_cursor], m)
						return InitialRootModel().SwitchScreen(&switched)
					}
				}

			}
		case "a":
		case tea.KeyCtrlX.String():
			switch m.activeTab {
			case 0:
				m.ClearBookshelfInput()
			case 3:
				m.ClearAuctionSearch()
			}
		case tea.KeyTab.String():
			m.NextTab()
			m.UnfocusAuctionInputs()
			m.UnfocusBookshelfInput()
			m.bookshelfLibrary = m.ps.Reader.Library
		case tea.KeyShiftTab.String():
			m.PreviousTab()
			m.UnfocusAuctionInputs()
			m.UnfocusBookshelfInput()
			m.bookshelfLibrary = m.ps.Reader.Library
		case tea.KeyUp.String():
			switch m.activeTab {
			case 0:
				m.PreviousItemBookshelf()
			case 1:
				m.PreviousCurrentReads()
			case 2:
				m.ps.Shop.PreviousRow()
			case 3:
				m.PreviousAuctionBook()
			case 4:
				m.PreviousItemInventory()
			case 6:
				m.PreviousExitChoice()
			}

		case tea.KeyDown.String():
			switch m.activeTab {
			case 0:
				m.NextItemBookshelf()
			case 1:
				m.NextCurrentReads()
			case 2:
				m.ps.Shop.NextRow()
			case 3:
				m.NextAuctionBook()
			case 4:
				m.NextItemInventory()
			case 6:
				m.NextExitChoice()
			}
		case tea.KeyEnter.String():
			switch m.activeTab {
			case 2:
				m.TryBuy()
			case 0:
				if m.bookshelf_input.Focused() {
					m.SubmitBookshelfSearch()
				} else {
					m.TrySwitchBook()
				}
			case 3:
				m.SubmitAuctionSearch()
			case 4:
				if len(m.ps.Reader.Inventory.Items) > 0 {
					switched := InitialBookDetailsModel(m.ps.Reader.Inventory.Items[m.i_cursor], m)
					return InitialRootModel().SwitchScreen(&switched)
				}
			case 6:
				switch m.exitChoices[m.ex_cursor] {
				case "Save & Go to Main Menu":
					log.Println("Save & Go to Main Menu")
					m.ps.SavePlayerToFile()
					switched := InitialSaveMenuModel()
					return InitialRootModel().SwitchScreen(&switched)
				case "Save & Quit":
					log.Println("Save & Quit")
					m.ps.SavePlayerToFile()
					return m, tea.Quit
				case "Quit without Saving":
					log.Println("Quit without Saving")
					switched := InitialSaveMenuModel()
					return InitialRootModel().SwitchScreen(&switched)

				}
			}
		case tea.KeyLeft.String():
			switch m.activeTab {
			case 0:
				m.PreviousBookPage()
			case 3:
				m.AuctionSwitchInput()
			case 5:
				m.PreviousHelpItem()
			}
		case tea.KeyRight.String():
			switch m.activeTab {
			case 0:
				m.NextBookPage()
			case 3:
				m.AuctionSwitchInput()
			case 5:
				m.NextHelpItem()
			}
		case tea.KeyCtrlF.String():
			switch m.activeTab {
			case 0:
				m.bookshelf_input.Focus()
			case 3:
				m.auction_inputs[0].Focus()
				m.auction_inputs[1].Blur()
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
		// m.ps.Shop.table.Width(msg.Width)
		// m.ps.Shop.table.Height(msg.Height)

	case tickMsg:
		var cmd []tea.Cmd
		cmd = append(cmd, tickCmd())
		m.ps.Shop.Update(&m.ps.Reader)
		m.ps.PeriodicSavePlayerToFile()
		w, h, _ := term.GetSize(int(os.Stdout.Fd()))
		if w != m.width || h != m.height {
			m.updateSize(w, h)
		}
		cmd = append(cmd, func() tea.Msg { return tea.WindowSizeMsg{Width: w, Height: h} })

		for i, id := range m.ps.Reader.CurrentReads.BookIDs {
			r := &m.ps.Reader
			cr_p, err := r.Library.GetBookPointerByID(id)
			if err != nil {
				panic(err)
			}
			log.Println("tickMsg | cr: " + cr_p.Name)
			log.Println(cr_p)
			if cr_p.Progress >= 1.0 {
				r.FinishedBook(cr_p)
				cmd = append(cmd, m.progress[i].SetPercent(0))
			} else {
				r.IncreaseProgress(cr_p)

				cmd = append(cmd, m.progress[i].SetPercent(cr_p.Progress))
				log.Println("Progression | " + strconv.FormatFloat(cr_p.Progress, 'f', -1, 64))
			}
		}
		return m, tea.Batch(cmd...)
	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		var cmds []tea.Cmd
		for i := range m.progress {
			progressModel, cmd := m.progress[i].Update(msg)
			m.progress[i] = progressModel.(progress.Model)
			cmds = append(cmds, cmd)
		}

		return m, tea.Batch(cmds...)
	case spinner.TickMsg:
		var spinner_cmd tea.Cmd
		m.spinner, spinner_cmd = m.spinner.Update(msg)
		return m, spinner_cmd
	}

	switch m.activeTab {
	case 3:
		var keycmds []tea.Cmd = make([]tea.Cmd, len(m.auction_inputs))
		for i := range m.auction_inputs {
			m.auction_inputs[i], keycmds[i] = m.auction_inputs[i].Update(msg)
		}
		return m, tea.Batch(keycmds...)
	case 0:
		var keycmd tea.Cmd
		m.bookshelf_input, keycmd = m.bookshelf_input.Update(msg)
		return m, keycmd

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
	return "\n" + lipgloss.JoinHorizontal(lipgloss.Top, tabRow...)
}

func (m *DashboardModel) View() string {
	pf := message.NewPrinter(language.English)

	r := m.ps.Reader
	s := theme.Heading1.Render(r.Name)
	s += "\n"
	s += "\n"
	s += "Favourite Book: " + r.FavouriteBook + ", " + r.FavouriteAuthor
	s += "\n"

	k := pf.Sprintf("Knowledge: %d", r.Knowledge)
	iq := "IQ: " + strconv.Itoa(r.IQ) + " (" + r.IQ_Title() + ")"

	p := "Prestige: " + strconv.Itoa(r.Prestige)
	rs := "Reading Speed: " + strconv.Itoa(r.Prestige)

	v1 := lipgloss.NewStyle().Padding(0, 10, 0, 0).Render(lipgloss.JoinVertical(lipgloss.Left, k, iq))
	v2 := lipgloss.JoinVertical(lipgloss.Left, rs, p)
	s += lipgloss.JoinHorizontal(lipgloss.Center, v1, v2)
	s += "\n"

	s += lipgloss.NewStyle().Render(m.TabsView())
	m.width = lipgloss.Width(s)

	s += "\n"
	switch m.activeTab {
	case 0:
		s += m.BookshelfView()
	case 1:
		s += m.CurrentReadsView()
	case 2:
		s += m.BookshopView()
	case 3:
		s += m.AuctionView()
	case 4:
		s += m.InventoryView()
	case 5:
		s += m.HelpView()
	case 6:
		s += m.ExitView()
	}
	s += "\n"

	s += m.errorMessage

	return lipgloss.NewStyle().Width(m.width).Height(m.height).Render(s)
}
