package application

import (
	booksController "github.com/betterreads/internal/domains/books/controller"
	booksRepository "github.com/betterreads/internal/domains/books/repository"
	booksService "github.com/betterreads/internal/domains/books/service"

	"github.com/betterreads/internal/domains/users/controller"
	"github.com/betterreads/internal/domains/users/repository"
	"github.com/betterreads/internal/domains/users/service"
	"github.com/gin-gonic/gin"
	"log"
)

type Router struct {
	engine *gin.Engine
	port   string
}

func NewRouter(port string) *Router {
	r := gin.Default()
	addUsersHandlers(r)
	addBooksHandlers(r)

	return &Router{
		engine: r,
		port:   port,
	}
}

func addUsersHandlers(r *gin.Engine) {
	userRepo := repository.NewMemoryDatabase()
	us := service.NewUsersService(userRepo)
	uc := controller.NewUsersController(us)

	r.POST("/users/register", uc.Register)
	r.POST("/users/login", uc.LogIn)
	r.GET("/users", uc.GetUsers)
}

func addBooksHandlers(r *gin.Engine) {
	booksRepo := booksRepository.NewInmemoryBooksDatabase()
	bs := booksService.NewBooksService(booksRepo)
	bc := booksController.NewBooksController(bs)
	r.POST("/books", bc.PublishBook)
	r.GET("/books/:book-name", bc.GetBook)
}

func (r *Router) Run() {
	if err := r.engine.Run(r.port); err != nil {
		log.Fatalln("can't start server: ", err)
	}
}
