package repository

import (
	"strconv"
)

type InmemoryBooksDatabase struct {
	books map[int]Book
	newId int
}

func NewInmemoryBooksDatabase() BooksDatabase {
	inmemoryBooksDatabase := new(InmemoryBooksDatabase)
	inmemoryBooksDatabase.books = make(map[int]Book)
	inmemoryBooksDatabase.newId = 1
	return inmemoryBooksDatabase
}

func (db *InmemoryBooksDatabase) SaveBook(book Book) error {
	db.books[db.GenerateBookId()] = book
	return nil
}

func (db *InmemoryBooksDatabase) GetBookById(id int) (*Book, error) {
	if len(db.books) == 0 {
		return nil, nil
	}

	book, ok := db.books[id]

	if !ok {
		return nil, nil
	}

	return &book, nil
	
}

func (db *InmemoryBooksDatabase) GetBookByName(name string) (*Book, error) {
	if len(db.books) == 0 {
		return nil, nil
	}
	for _, book := range db.books {
		if book.Title == name {
			return &book, nil
		}
	}
	return nil, nil
}

func (db *InmemoryBooksDatabase) GenerateBookId() int {
	id := db.newId
	db.newId++
	return id
}

func (db *InmemoryBooksDatabase) RateBook(bookId int, userId int, rating int) error {
	var book = db.books[bookId]
	
	ratingId := db.createRateId(bookId, userId)

	book.Ratings[ratingId] = rating


	db.books[bookId] = book

	return nil
}

func (db *InmemoryBooksDatabase) DeleteRating(bookId int, userId int) error {
	var book = db.books[bookId]
	
	ratingId := db.createRateId(bookId, userId)

	delete(book.Ratings, ratingId)

	db.books[bookId] = book
	return nil
}

func (db *InmemoryBooksDatabase) createRateId(bookId int, userId int) int {
	strA := strconv.Itoa(bookId)
	strB := strconv.Itoa(userId)

	concatenated := strA + strB

	result, err := strconv.Atoi(concatenated)
	if err != nil {
		return -1
	}

	return result
}