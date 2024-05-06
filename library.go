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

type ReorderType int

var (
	Alphabetical        int = 0
	AlphabeticalReverse int = 1
)

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

func (l *Library) ContainsBookByNameAndAuthor(book Book) bool {
	for _, b := range l.Books {
		if b.Author == book.Author && b.Name == book.Name {
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

func (l *Library) GetBookPointerByID(id uuid.UUID) (*Book, error) {
	for i, b := range l.Books {
		if b.ID == id {
			return &l.Books[i], nil
		}
	}
	return &Book{}, errors.New("book does not exist")
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
