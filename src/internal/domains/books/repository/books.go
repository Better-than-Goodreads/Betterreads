package repository

type Book struct {
	Title           string   `json:"title" binding:"required"`
	Author          string   `json:"author" binding:"required"`
	Description     string   `json:"description" binding:"required"`
	PhotoId         string   `json:"photo_id" binding:"required"`
	AmountOfPages   string   `json:"amount_of_pages" binding:"required"`
	PublicationDate string   `json:"publication_date" binding:"required"`
	Language        string   `json:"language" binding:"required"`
	Genres          []string `json:"genres" binding:"required"`
}

type BooksDatabase interface {
	SaveBook(book Book) error
	GetBookById(id int) (Book, error)
	GetBookByName(name string) (*Book, error)
}
