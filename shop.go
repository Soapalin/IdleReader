package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"

	_ "embed"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Item struct {
	ID            uuid.UUID
	Name          string
	Description   string
	Cost          int
	IqRequirement int
	Bought        bool
	Effect        string
}

type GameItemDatabase struct {
	Items []Item
}

func (g *GameItemDatabase) ContainsItem(item Item) bool {
	for _, i := range g.Items {
		if i.ID == item.ID {
			return true
		}
	}
	return false
}

func (g *GameItemDatabase) ContainsItemByName(item Item) bool {
	for _, i := range g.Items {
		if i.Name == item.Name {
			return true
		}
	}
	return false
}

func (g *GameItemDatabase) AddItem(item Item) {
	g.Items = append(g.Items, item)
}

type Shop struct {
	Books      Library
	Items      GameItemDatabase
	ShopSize   int
	Modified   time.Time
	table      table.Table
	TableIndex int
	TableLen   int
}

type TransactionResult int

const (
	bookTransaction TransactionResult = iota
	itemTransaction
	knowledgeMissingTransaction
	iqMissingTransaction
	unknownTransaction
)

func (s *Shop) Buy(ps *PlayerSave) TransactionResult {
	log.Println("Buy | ")
	index := s.TableIndex - 1
	if index < len(s.Books.Books) {
		log.Println(s.Books.Books[index].Name)
		if s.Books.Books[index].IntelligenceRequirement > ps.Reader.IQ {
			return iqMissingTransaction
		}
		if s.Books.Books[index].KnowledgeRequirement > ps.Reader.Knowledge {
			return knowledgeMissingTransaction
		}
		ps.Reader.DecreaseKnowledge(s.Books.Books[index].KnowledgeRequirement)
		ps.Reader.Library.AddBookToLibrary(s.Books.Books[index])
		s.LoadShopTable(&ps.Reader.Library)
		return bookTransaction
	} else if (index - len(s.Books.Books)) < len(s.Items.Items) {
		log.Println(s.Items.Items[index-len(s.Books.Books)].Name)
		if s.Items.Items[index-len(s.Books.Books)].IqRequirement > ps.Reader.IQ {
			return iqMissingTransaction
		}
		if s.Items.Items[index-len(s.Books.Books)].Cost > ps.Reader.Knowledge {
			return knowledgeMissingTransaction
		}
		ps.Reader.DecreaseKnowledge(s.Items.Items[index-len(s.Books.Books)].Cost)
		ps.Reader.Inventory.AddItem(s.Items.Items[index-len(s.Books.Books)])
		return itemTransaction
	}
	return unknownTransaction

}

func (s *Shop) GetShopItemByIndex() uuid.UUID {
	log.Println("GetShopItemByIndex | ")
	shopIndex := s.TableIndex - 1
	if shopIndex < len(s.Books.Books) {
		log.Println(s.Books.Books[shopIndex].Name)
		return s.Books.Books[shopIndex].ID
	} else if (shopIndex - len(s.Books.Books)) < len(s.Items.Items) {
		log.Println(s.Items.Items[shopIndex-len(s.Books.Books)].Name)
		return s.Items.Items[shopIndex-len(s.Books.Books)].ID
	}
	return uuid.New()
}

func (s *Shop) NextRow() {
	if s.TableIndex >= s.TableLen {
		return
	}
	s.TableIndex++

}

func (s *Shop) PreviousRow() {
	s.TableIndex--
	if s.TableIndex < 1 {
		s.TableIndex = 1
	}
}

func (s *Shop) Update(reader_p *Reader) {
	if s.Modified.Add(time.Minute * 1).Before(time.Now()) {
		log.Println("Shop | Update()")
		*s = InitShop(reader_p)
		s.Modified = time.Now()
	}
}

func (s *Shop) LoadShopTable(lib_p *Library) {
	columns := []string{"Name", "Details/Author", "IQ Required", "Knowledge Cost"}
	var rows [][]string
	n := 0
	for n < len(s.Books.Books) {
		var owned_string string
		owned := lib_p.ContainsBook(s.Books.Books[n])
		if owned {
			owned_string = "*"
		}
		rows = append(rows, []string{owned_string + s.Books.Books[n].Name, s.Books.Books[n].Author, strconv.Itoa(s.Books.Books[n].IntelligenceRequirement), strconv.Itoa(s.Books.Books[n].KnowledgeRequirement)})
		n += 1
	}
	n = 0
	for n < len(s.Items.Items) {
		rows = append(rows, []string{s.Items.Items[n].Name, s.Items.Items[n].Description, strconv.Itoa(s.Items.Items[n].IqRequirement), strconv.Itoa(s.Items.Items[n].Cost)})
		n += 1
	}

	styleFunc := func(row, col int) lipgloss.Style {
		if row == 0 {
			return lipgloss.NewStyle().Padding(0, 1)
		}
		if col == 0 {
			return lipgloss.NewStyle().Width(18).Padding(1)
		}
		return lipgloss.NewStyle().Width(20).Padding(1, 1)
	}
	s.table = *table.New().Headers(columns...).Rows(rows...).StyleFunc(styleFunc)
	s.TableLen = len(rows)
}

func InitShop(reader_p *Reader) Shop {
	n := 0
	var books Library
	var items GameItemDatabase
	pf := message.NewPrinter(language.English)

	columns := []string{"Name", "Details/Author", "IQ Required", "Knowledge Cost"}
	var rows [][]string
	allBooks, err := DB.GetBooksByFilter("IntelligenceRequirement <= " + strconv.Itoa(reader_p.IQ))

	if err != nil {
		log.Println(err)
		panic(err)
	}
	allItems, err := DB.GetAllItems()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	for n < 4 {
		randomIndex := rand.Intn(len(allBooks.Books) - 1)
		if !books.ContainsBook(allBooks.Books[randomIndex]) {
			books.AddBookToLibrary(allBooks.Books[randomIndex])
			b := allBooks.Books[randomIndex]
			var owned_string string
			owned := reader_p.Library.ContainsBook(b)
			if owned {
				owned_string = "*"
			}
			rows = append(rows, []string{owned_string + b.Name, b.Author, pf.Sprintf("%d", b.IntelligenceRequirement), pf.Sprintf("%d", b.KnowledgeRequirement)})
			n += 1
		}
	}

	n = 0
	for n < 1 {
		randomIndex := rand.Intn(len(allItems.Items) - 1)
		if !items.ContainsItem(allItems.Items[randomIndex]) {
			items.AddItem(allItems.Items[randomIndex])
			i := allItems.Items[randomIndex]
			rows = append(rows, []string{i.Name, i.Description, pf.Sprintf("%d", i.IqRequirement), pf.Sprintf("%d", i.Cost)})
			n += 1
		}

	}

	styleFunc := func(row, col int) lipgloss.Style {
		if row == 0 {
			return lipgloss.NewStyle().Padding(0, 1)
		}
		if col == 0 || col == 1 {
			return lipgloss.NewStyle().Width(25).Padding(0)
		}
		return lipgloss.NewStyle().Width(5).Padding(0)
	}

	return Shop{
		Books:      books,
		Items:      items,
		ShopSize:   len(books.Books) + len(items.Items),
		Modified:   time.Now(),
		table:      *table.New().Headers(columns...).Rows(rows...).StyleFunc(styleFunc),
		TableIndex: 1,
		TableLen:   len(rows),
	}
}
