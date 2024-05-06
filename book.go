package main

import (
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


func (b *Book) String() string {
	return b.Name + ", " + b.Author
}