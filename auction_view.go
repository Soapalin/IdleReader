package main

import (
	"fmt"
	"game/engine/theme"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var AUCTION_COST_MODIFIER int = 4

func (m *DashboardModel) AuctionView() string {
	s := ""
	s += theme.Heading1.Render("Auction") + "\n"
	s += theme.Heading1.Render("Search")

	s += fmt.Sprintf(
		`
%s %s
%s %s
		`,
		theme.InputStyle.Width(40).Render("Book Name"),
		theme.InputStyle.Width(40).Render("Author"),
		m.auction_inputs[0].View(),
		m.auction_inputs[1].View(),
	)

	if len(m.auctionLibrary.Books) > 0 {
		start, end := m.auctionPaginator.GetSliceBounds(len(m.auctionLibrary.Books))

		firstLine := []string{
			theme.Heading2.Width(m.width / 2).Render("Book/Author"),
			theme.Heading2.Width(m.width / 8).Render("IQ Req."),
			theme.Heading2.Width(m.width / 4).Render("Knowledge Cost"),
			theme.Heading2.Width(m.width / 8).Render("Owned"),
		}
		lines := []string{lipgloss.JoinHorizontal(lipgloss.Left, firstLine...)}

		for i, b := range m.auctionLibrary.Books[start:end] {
			var textStyle lipgloss.Style
			var currentLine []string
			if m.auc_cursor == (start + i) {
				textStyle = lipgloss.NewStyle().Foreground(theme.Pink).Background(theme.White)
			} else {
				textStyle = lipgloss.NewStyle()
			}
			currentLine = append(currentLine, textStyle.Width(m.width/2).Render(b.Name+", "+b.Author))
			currentLine = append(currentLine, textStyle.Width(m.width/8).Render(strconv.Itoa(b.IntelligenceRequirement)))
			currentLine = append(currentLine, textStyle.Width(m.width/4).Render(strconv.Itoa(b.KnowledgeRequirement*AUCTION_COST_MODIFIER)))
			if m.ps.Reader.Library.ContainsBook(b) {
				currentLine = append(currentLine, textStyle.Width(m.width/8).Render(RenderEmojiOrFallback("✅", []string{"Yes"})))
			} else {
				currentLine = append(currentLine, textStyle.Width(m.width/8).Render(RenderEmojiOrFallback("❌", []string{"No"})))
			}
			lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, currentLine...))

		}

		s += "\n"
		s += lipgloss.JoinVertical(lipgloss.Left, lines...)

		s += "\n"
		paginatorView := lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Position(0.5)).Render(m.auctionPaginator.View())
		pageNumber := "<" + strconv.Itoa(m.auctionPaginator.Page+1) + "/" + strconv.Itoa(m.auctionPaginator.TotalPages) + ">"
		paginatorFull := lipgloss.JoinVertical(0, paginatorView, pageNumber)
		s += lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Position(0.5)).Render(paginatorFull)
		s += "\n" + strconv.Itoa(len(m.auctionLibrary.Books)) + " Results..."
		if len(m.auctionLibrary.Books) > 5 {
			s += "Narrow down your search to see the rest of the results."
		}

	} else {
		s += "\nNo Results..."
	}

	s += "\n\n" + theme.HelpIcon.Render("ctrl+F") + theme.HelpText.Render(" search • ")
	s += theme.HelpIcon.Render("ctrl+X") + theme.HelpText.Render(" clear search • ")
	s += theme.HelpIcon.Render("enter") + theme.HelpText.Render(" submit search • ")
	s += theme.HelpIcon.Render("b") + theme.HelpText.Render(" buy • ")
	s += theme.HelpIcon.Render("i") + theme.HelpText.Render(" book info • ")
	s += theme.HelpIcon.Render("←/→") + theme.HelpText.Render(" switch input • ")
	s += theme.HelpIcon.Render("↑/↓") + theme.HelpText.Render(" switch book • ")
	s += theme.HelpIcon.Render("esc / q") + theme.HelpText.Render(" quit")

	s += "\n"
	return s

}

func (m *DashboardModel) AuctionSwitchInput() {
	if m.auction_inputs[0].Focused() {
		m.auction_inputs[1].Focus()
		m.auction_inputs[0].Blur()
	} else {
		m.auction_inputs[0].Focus()
		m.auction_inputs[1].Blur()
	}
}

func (m *DashboardModel) AuctionInputsFocused() bool {
	for _, input := range m.auction_inputs {
		if input.Focused() {
			return true
		}
	}
	return false
}

func (m *DashboardModel) NextAuctionBook() {
	m.UnfocusAuctionInputs()
	if m.auc_cursor >= len(m.auctionLibrary.Books)-1 {
		// m.auc_cursor = 0
		// m.auctionPaginator.Page = 0

		return
	}
	m.auc_cursor++

	_, end := m.auctionPaginator.GetSliceBounds(len(m.auctionLibrary.Books))
	if m.auc_cursor >= end {
		m.auctionPaginator.Page++
	}
}

func (m *DashboardModel) PreviousAuctionBook() {
	m.auc_cursor--
	m.UnfocusAuctionInputs()
	if m.auc_cursor < 0 {
		m.auc_cursor = 0
		return
	}
	start, _ := m.auctionPaginator.GetSliceBounds(len(m.auctionLibrary.Books))
	if m.auc_cursor < start {
		m.auctionPaginator.Page--
	}
}

func (m *DashboardModel) UnfocusAuctionInputs() {
	m.auction_inputs[1].Blur()
	m.auction_inputs[0].Blur()
}

func (m *DashboardModel) ClearAuctionSearch() {
	m.auction_inputs[0].SetValue("")
	m.auction_inputs[1].SetValue("")
}

func (m *DashboardModel) SubmitAuctionSearch() {
	m.UnfocusAuctionInputs()
	m.auc_cursor = 0
	bookname := strings.TrimSpace(m.auction_inputs[0].Value())
	author := strings.TrimSpace(m.auction_inputs[1].Value())
	lib, err := DB.FindBooksContainsBoth(bookname, author)
	if err != nil {
		m.errorMessage = err.Error()
		return
	}
	log.Println(lib)
	m.auctionLibrary = lib
	m.auctionPaginator.SetTotalPages(len(m.auctionLibrary.Books))
	m.auctionPaginator.Page = 0
	AUCTION_COST_MODIFIER = rand.Intn(5) + 2

}

func (m *DashboardModel) AuctionBuy() TransactionResult {
	log.Println("AuctionBuy | ")
	book := m.auctionLibrary.Books[m.auc_cursor]
	if book.IntelligenceRequirement > m.ps.Reader.IQ {
		return iqMissingTransaction
	}
	if (book.KnowledgeRequirement * AUCTION_COST_MODIFIER) > m.ps.Reader.Knowledge {
		return knowledgeMissingTransaction
	}
	m.ps.Reader.DecreaseKnowledge(book.KnowledgeRequirement * AUCTION_COST_MODIFIER)
	m.ps.Reader.Library.AddBookToLibrary(book)
	return bookTransaction
}

func (m *DashboardModel) TryAuctionBuy() {
	if m.AuctionInputsFocused() {
		return
	}
	if len(m.auctionLibrary.Books) == 0 {
		return
	}
	itemID := m.auctionLibrary.Books[m.auc_cursor].ID
	if m.ps.AlreadyOwned(itemID) {
		log.Println("AlreadyOwned returned True")
		m.errorMessage = "You already owned this item/book!"
	} else {
		typeItem := m.AuctionBuy()
		switch typeItem {
		case unknownTransaction:
			m.errorMessage = "Unable to purchase item/book."
		case knowledgeMissingTransaction:
			m.errorMessage = "Not enough knowledge."
		case iqMissingTransaction:
			m.errorMessage = "Not enough IQ."
		case bookTransaction, itemTransaction:
			m.errorMessage = "Item/Book bought successfully! Check your Bookshelf."
		}
	}
}
