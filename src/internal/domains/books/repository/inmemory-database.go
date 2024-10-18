package repository

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

func (db *InmemoryBooksDatabase) GetBookById(id int) (Book, error) {
	//TODO implement me
	panic("implement me")
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
