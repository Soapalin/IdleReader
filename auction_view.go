package main

import (
	"fmt"
	"game/engine/theme"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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
		bookNames := []string{theme.Heading2.Render("Book/Author"), "-"}
		IQCost := []string{theme.Heading2.Render("IQ Req."), "-"}
		KnowledgeCost := []string{theme.Heading2.Render("Knowledge Cost"), "-"}
		Owned := []string{theme.Heading2.Render("Owned"), "-"}
		start, end := m.auctionPaginator.GetSliceBounds(len(m.auctionLibrary.Books))

		for i, b := range m.auctionLibrary.Books[start:end] {
			var textStyle lipgloss.Style
			if m.auc_cursor == (start + i) {
				textStyle = lipgloss.NewStyle().Foreground(theme.Pink).Background(theme.White)
			} else {
				textStyle = lipgloss.NewStyle()
			}
			bookNames = append(bookNames, textStyle.Render(b.Name))
			IQCost = append(IQCost, textStyle.Render(strconv.Itoa(b.IntelligenceRequirement)))
			KnowledgeCost = append(KnowledgeCost, textStyle.Render(strconv.Itoa(b.KnowledgeRequirement)))
			if m.ps.Reader.Library.ContainsBook(b) {
				Owned = append(Owned, textStyle.Render(RenderEmojiOrFallback("✅", []string{"Yes"})))
			} else {
				Owned = append(Owned, textStyle.Render(RenderEmojiOrFallback("❌", []string{"No"})))
			}

		}

		bookColumn := lipgloss.NewStyle().Width(m.width / 2).Render(lipgloss.JoinVertical(lipgloss.Left, bookNames...))
		IQColumn := lipgloss.NewStyle().Width(m.width / 8).Render(lipgloss.JoinVertical(lipgloss.Left, IQCost...))
		KnowledgeColumn := lipgloss.NewStyle().Width(m.width / 4).Render(lipgloss.JoinVertical(lipgloss.Left, KnowledgeCost...))
		OwnedColumn := lipgloss.NewStyle().Width(m.width / 8).Render(lipgloss.JoinVertical(lipgloss.Left, Owned...))
		s += "\n"
		s += lipgloss.JoinHorizontal(lipgloss.Left, bookColumn, IQColumn, KnowledgeColumn, OwnedColumn)

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
	s += theme.HelpIcon.Render("←/→") + theme.HelpText.Render(" switch input • ")
	s += theme.HelpIcon.Render("↑/↓") + theme.HelpText.Render(" switch book • ")
	s += theme.HelpIcon.Render("esc / q") + theme.HelpText.Render(" quit")

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

func (m *DashboardModel) NextAuctionBook() {
	m.auc_cursor++
	if m.auc_cursor >= len(m.auctionLibrary.Books) {
		m.auc_cursor = 0
		m.auctionPaginator.Page = 0
		return
	}
	_, end := m.auctionPaginator.GetSliceBounds(len(m.auctionLibrary.Books))
	if m.auc_cursor >= end {
		m.auctionPaginator.Page++
	}
}

func (m *DashboardModel) PreviousAuctionBook() {
	m.auc_cursor--
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
	m.auctionLibrary = lib
	m.auctionPaginator.SetTotalPages(len(m.auctionLibrary.Books))
	m.auctionPaginator.Page = 0

}

func (m *DashboardModel) AuctionBuy() TransactionResult {
	return knowledgeMissingTransaction
}

func (m *DashboardModel) TryAuctionBuy() {
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
