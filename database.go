package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type Database struct {
	filepath string
	db       *sql.DB
}

var DB Database = createDatabase("IdleReaderDatabase.db")
var DATABASE *sql.DB

func createDatabase(filePath string) Database {
	dir := createDocumentFolder()
	db_filepath := filepath.Join(dir, filePath)
	if _, err := os.Stat(db_filepath); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(db_filepath)
		if err != nil {
			panic(err)
		}
		file.Close()

	}
	DATABASE, _ = sql.Open("sqlite", db_filepath)

	return Database{
		filepath: filePath,
		db:       DATABASE,
	}

}

func (d *Database) CreateAllBooksTable() {

	createAllBooksTable := `CREATE TABLE if not exists Books(
		"ID" TEXT TYPE NOT NULL,
		"Name" TEXT TYPE NOT NULL,
		"Author" TEXT TYPE NOT NULL,
		"Progress" REAL TYPE NOT NULL,
		"KnowledgeIncrease" INTEGER TYPE NOT NULL,
		"KnowledgeRequirement" INTEGER TYPE NOT NULL,
		"IntelligenceIncrease" INTEGER TYPE NOT NULL,
		"IntelligenceRequirement" INTEGER TYPE NOT NULL,
		"Pages" INTEGER NOT NULL,
		"Repeat" INTEGER TYPE NOT NULL,
		PRIMARY KEY ("Name","Author")
	);`

	statement, err := d.db.Prepare(createAllBooksTable)
	if err != nil {
		panic(err)
	}
	statement.Exec()
	log.Println("CreateAllBooksTable | Succcess")

	var abl Library = InitAllBookLibrary()

	for _, book := range abl.Books {
		d.InsertOrIgnoreBook(book)
	}

}

func (d *Database) InsertOrUpdateBook(book Book) {
	insertBook := fmt.Sprintf(`INSERT or REPLACE INTO 
	Books(ID, Name, Author, Progress, KnowledgeIncrease, KnowledgeRequirement, IntelligenceIncrease, IntelligenceRequirement, Pages, Repeat) 
	values ('%s','%s','%s',%f,%d,%d,%d,%d,%d,%d);`,
		book.ID.String(),
		book.Name,
		book.Author,
		book.Progress,
		book.KnowledgeIncrease,
		book.KnowledgeRequirement,
		book.IntelligenceIncrease,
		book.IntelligenceRequirement,
		book.Pages,
		book.Repeat)
	log.Println("InsertOrUpdateBook |" + insertBook)
	statement, err := d.db.Prepare(insertBook)
	if err != nil {
		panic(err)
	}
	statement.Exec()
}

func (d *Database) InsertOrIgnoreBook(book Book) {
	insertBook := fmt.Sprintf(`INSERT or IGNORE INTO 
	Books(ID, Name, Author, Progress, KnowledgeIncrease, KnowledgeRequirement, IntelligenceIncrease, IntelligenceRequirement, Pages, Repeat) 
	values ('%s','%s','%s',%f,%d,%d,%d,%d,%d,%d);`,
		book.ID.String(),
		book.Name,
		book.Author,
		book.Progress,
		book.KnowledgeIncrease,
		book.KnowledgeRequirement,
		book.IntelligenceIncrease,
		book.IntelligenceRequirement,
		book.Pages,
		book.Repeat)
	log.Println("InsertOrUpdateBook |" + insertBook)
	statement, err := d.db.Prepare(insertBook)
	if err != nil {
		panic(err)
	}
	statement.Exec()
}

func (d *Database) CreateAllItemsTable() {

	createAllItemsTable := `CREATE TABLE if not exists Items(
		"ID" TEXT TYPE NOT NULL,
		"Name" TEXT TYPE NOT NULL,
		"Description" TEXT TYPE NOT NULL,
		"Cost" INTEGER TYPE NOT NULL,
		"IqRequirement" INTEGER TYPE NOT NULL,
		"Effect" INTEGER TYPE NOT NULL,
		"Bought" INTEGER,
		PRIMARY KEY ("Name")
	);`

	statement, err := d.db.Prepare(createAllItemsTable)
	if err != nil {
		panic(err)
	}
	statement.Exec()
	log.Println("CreateAllItemsTable | Succcess")

	var agi GameItemDatabase = InitAllGameItemDatabase()

	for _, item := range agi.Items {
		d.InsertOrIgnoreItem(item)
	}

}

func (d *Database) InsertOrUpdateItem(item Item) {
	var boughtInt int
	if item.Bought {
		boughtInt = 1
	}
	insertItem := fmt.Sprintf(`INSERT or REPLACE INTO
	Items(ID, Name, Description, Cost, IqRequirement, Effect, Bought)
	values ('%s','%s','%s',%d,%d,'%s',%d);`,
		item.ID.String(),
		item.Name,
		item.Description,
		item.Cost,
		item.IqRequirement,
		item.Effect,
		boughtInt)
	log.Println("InsertOrUpdateBook |" + insertItem)
	statement, err := d.db.Prepare(insertItem)
	if err != nil {
		panic(err)
	}
	statement.Exec()
}

func (d *Database) InsertOrIgnoreItem(item Item) {
	var boughtInt int
	if item.Bought {
		boughtInt = 1
	}
	insertItem := fmt.Sprintf(`INSERT or IGNORE INTO
	Items(ID, Name, Description, Cost, IqRequirement, Effect, Bought)
	values ('%s','%s','%s',%d,%d,'%s',%d);`,
		item.ID.String(),
		item.Name,
		item.Description,
		item.Cost,
		item.IqRequirement,
		item.Effect,
		boughtInt)
	log.Println("InsertOrUpdateBook |" + insertItem)
	statement, err := d.db.Prepare(insertItem)
	if err != nil {
		panic(err)
	}
	statement.Exec()
}

func (d *Database) GetItemByID(id uuid.UUID) (Item, error) {
	log.Println("GetItemByID | " + id.String())

	selectItemByID := "SELECT * FROM Items WHERE ID = '" + id.String() + "';"
	log.Println(selectItemByID)

	row, err := d.db.Query(selectItemByID)
	if err != nil {
		log.Println(err)
		return Item{}, err
	}
	defer row.Close()

	if row.Next() {
		var i Item
		var uuidstring string

		er := row.Scan(&uuidstring, &i.Name, &i.Description, &i.Cost, &i.IqRequirement, &i.Effect, &i.Bought)
		if er != nil {
			log.Println("GetItemByID | row.Scan")
			log.Println(er)
			return Item{}, er
		}
		i.ID = id
		return i, nil
	}
	return Item{}, errors.New("item not found")
}

func (d *Database) GetAllItems() (GameItemDatabase, error) {
	log.Println("GetAllItems |")
	selectAllItems := "SELECT * FROM Items;"
	row, err := d.db.Query(selectAllItems)

	if err != nil {
		log.Println(err)
		return GameItemDatabase{}, nil
	}
	defer row.Close()
	var allGameItems GameItemDatabase
	for row.Next() {
		var i Item
		var uuidstring string

		er := row.Scan(&uuidstring, &i.Name, &i.Description, &i.Cost, &i.IqRequirement, &i.Effect, &i.Bought)
		if er != nil {
			log.Println("GetItemByID | row.Scan")
			log.Println(er)
			return GameItemDatabase{}, er
		}
		i.ID, _ = uuid.Parse(uuidstring)

		allGameItems.AddItem(i)
	}
	return allGameItems, nil
}

func (d *Database) FindBookByNameAuthor(book, author string) (Book, error) {
	log.Println("FindBookByNameAuthor |" + book + ", " + author)
	selectBookByNameAndAuthor := "SELECT * FROM Books WHERE Name = '" + book + "' AND Author = '" + author + "';"
	log.Println(selectBookByNameAndAuthor)

	row, err := d.db.Query(selectBookByNameAndAuthor)
	if err != nil {
		log.Println(err)
		return Book{}, err
	}
	defer row.Close()

	if row.Next() {
		var b Book
		var uuidstring string

		er := row.Scan(&uuidstring, &b.Name, &b.Author, &b.Progress, &b.KnowledgeIncrease, &b.KnowledgeRequirement, &b.IntelligenceIncrease, &b.IntelligenceRequirement, &b.Pages, &b.Repeat)
		if er != nil {
			log.Println("FindBookByNameAuthor | row.Scan")
			log.Println(er)
		}
		b.ID, _ = uuid.Parse(uuidstring)
		return b, nil
	}
	return Book{}, errors.New("Book not found")
}

func (d *Database) GetBookByID(bookID uuid.UUID) (Book, error) {
	log.Println("GetBookByID | " + bookID.String())
	selectBookByID := "SELECT * FROM Books WHERE ID = '" + bookID.String() + "';"
	// selectBookByID := "SELECT * FROM Books;"
	log.Println(selectBookByID)
	row, err := d.db.Query(selectBookByID)
	if err != nil {
		log.Println(err)
		return Book{}, err
	}
	defer row.Close()
	// var bookFetched Book

	if row.Next() {
		var id string
		var name string
		var author string
		var progress float64
		var knowIncrease int
		var knowCost int
		var iqIncrease int
		var iqCost int
		var pages int
		var repeat int
		er := row.Scan(&id, &name, &author, &progress, &knowIncrease, &knowCost, &iqIncrease, &iqCost, &pages, &repeat)
		log.Println(er)
		if er != nil {
			log.Println("GetBookByID | row.Scan")
			log.Println(er)
		}
		uuidID, _ := uuid.Parse(id)
		return Book{
			ID:                      uuidID,
			Name:                    name,
			Author:                  author,
			Progress:                progress,
			KnowledgeIncrease:       knowIncrease,
			KnowledgeRequirement:    knowCost,
			IntelligenceIncrease:    iqIncrease,
			IntelligenceRequirement: iqCost,
			Pages:                   pages,
			Repeat:                  repeat,
		}, nil
	}
	return Book{}, errors.New("Book not found")

}

func (d *Database) GetAllBooks() (Library, error) {
	log.Println("GetAllBooks")
	selectAll := "SELECT * FROM Books;"

	row, err := d.db.Query(selectAll)
	if err != nil {
		log.Println(err)
		return Library{}, err
	}
	var allBooks Library
	for row.Next() {
		var b Book
		var uuidstring string

		er := row.Scan(&uuidstring, &b.Name, &b.Author, &b.Progress, &b.KnowledgeIncrease, &b.KnowledgeRequirement, &b.IntelligenceIncrease, &b.IntelligenceRequirement, &b.Pages, &b.Repeat)
		if er != nil {
			log.Println("FindBookByNameAuthor | row.Scan")
			log.Println(er)
			return Library{}, er
		}
		b.ID, _ = uuid.Parse(uuidstring)

		allBooks.AddBookToLibrary(b)
	}

	return allBooks, nil

}
