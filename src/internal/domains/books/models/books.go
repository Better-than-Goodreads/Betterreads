package models

type NewBookRequest struct {
	Title           string   `json:"title" binding:"required"`
	Author          string   `json:"author" binding:"required"`
	Description     string   `json:"description" binding:"required"`
	PhotoId         string   `json:"photo_id" binding:"required"`
	AmountOfPages   string   `json:"amount_of_pages" binding:"required"`
	PublicationDate string   `json:"publication_date" binding:"required"`
	Language        string   `json:"language" binding:"required"`
	Genres          []string `json:"genres" binding:"required"`
}

type NewBookRating struct {
	BookId int		`json:"book_id" binding:"required"`
	UserId int 		`json:"user_id" binding:"required"`
	Rating int		`json:"rating" binding:"required"`
}