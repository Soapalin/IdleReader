package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"strconv"
	"strings"

	_ "embed"

	"github.com/google/uuid"
)

// WORD := 1
// LINE := 10 * WORD
// PAGE := 20 * LINE * WORD

type BookType int

const (
	Fiction      BookType = 0
	Fantasy      BookType = 1
	Encyclopedia BookType = 2
	Textbook     BookType = 3
)

type Book struct {
	ID                      uuid.UUID
	Name                    string
	Author                  string
	Progress                float64
	KnowledgeIncrease       int
	KnowledgeRequirement    int
	IntelligenceIncrease    int
	IntelligenceRequirement int
	Pages                   int
	Repeat                  int
}

type Library struct {
	Books []Book
}

var AllBooksLibrary Library = LoadAllBooksLibrary()

func (l *Library) String(preceding string) string {
	s := ""
	if preceding == "DIGITS" {
		for i, b := range l.Books {
			s += strconv.Itoa(i+1) + ". " + b.Name + ", " + b.Author + "\n"
		}
		return s
	}
	for _, b := range l.Books {
		s += preceding + " " + b.Name + ", " + b.Author + "\n"
	}

	return s
}

func (l *Library) ReplaceBook(index int, book Book) {
	l.Books[index] = book
}

func (l *Library) FindBookByNameAuthor(author, bookName string) (Book, error) {
	author = strings.ToLower(strings.TrimSpace(author))
	bookName = strings.ToLower(strings.TrimSpace(bookName))

	for _, b := range l.Books {
		if b.Author == author && b.Name == bookName {
			return b, nil
		}
	}
	return Book{}, errors.New("book does not exist")
}

func (l *Library) ContainsBook(book Book) bool {
	for _, b := range l.Books {
		if b.ID == book.ID {
			return true
		}
	}
	return false
}

func (l *Library) AddBookToLibrary(book Book) {
	l.Books = append(l.Books, book)
}

func (l *Library) GetBookByID(id uuid.UUID) (Book, error) {
	for _, b := range l.Books {
		if b.ID == id {
			return b, nil
		}
	}
	return Book{}, errors.New("book does not exist")
}

//go:embed AllBooksLibrary.bin
var abl []byte

func LoadAllBooksLibrary() Library {

	buff := bytes.NewBuffer(abl)
	dec := gob.NewDecoder(buff)

	lib := Library{}

	dec_err := dec.Decode(&lib)
	if dec_err != nil {
		log.Println("LoadAllBooksLibrary | dec.Decode")
		panic(dec_err)
	}
	return lib
}
