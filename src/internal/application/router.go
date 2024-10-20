package application

import (
	"fmt"

	booksController "github.com/betterreads/internal/domains/books/controller"
	booksRepository "github.com/betterreads/internal/domains/books/repository"
	booksService "github.com/betterreads/internal/domains/books/service"
	"github.com/jmoiron/sqlx"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"log"

	usersController "github.com/betterreads/internal/domains/users/controller"
	usersRepository "github.com/betterreads/internal/domains/users/repository"
	usersService "github.com/betterreads/internal/domains/users/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Router struct {
	engine *gin.Engine
	address   string
}

func createRouterFromConfig(cfg *Config) *Router {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// gin.DefaultWriter = io.Discard
	// gin.DefaultErrorWriter = io.Discard
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// slog.SetDefault(logger)

	router := &Router{
		engine:  gin.Default(),
		address: cfg.Host + ":" + cfg.Port,
	}

	// router.Engine.Use(middleware.RequestLogger())
	// router.Engine.Use(middleware.ErrorHandler())

	return router
}

func NewRouter(port string) *Router {
	cfg := LoadConfig()
	
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DatabaseUser,
			cfg.DatabasePassword,
			cfg.DatabaseHost,
			cfg.DatabasePort,
			cfg.DatabaseName)

	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("can't connect to db: %v", err)
	}
	

	r := createRouterFromConfig(cfg)
	addUsersHandlers(r, conn)
	addBooksHandlers(r, conn)

    //Adds swagger documentation
    r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func addUsersHandlers(r *Router, conn *sqlx.DB) {
	// userRepo := usersRepository.NewMemoryDatabase()

	userRepo, err := usersRepository.NewPostgresUserRepository(conn)
	if err != nil {
		log.Fatalf("can't create db: %v", err)
	}
	us := usersService.NewUsersService(userRepo)
	uc := usersController.NewUsersController(us)

	r.engine.POST("/users/register-first", uc.RegisterFirstStep)
    r.engine.POST("/users/register-second", uc.RegisterSecondStep)
	r.engine.POST("/users/login", uc.LogIn)
	r.engine.GET("/users", uc.GetUsers)
    r.engine.GET("/users/:id", uc.GetUser)

    // Authenticated routes
    r.engine.GET("/users/welcome", authMiddleware, uc.Welcome)
}

func addBooksHandlers(r *Router, conn *sqlx.DB) {
	booksRepo := booksRepository.NewInmemoryBooksDatabase()
	bs := booksService.NewBooksService(booksRepo)
	bc := booksController.NewBooksController(bs)

	r.engine.POST("/books", bc.PublishBook)
	r.engine.GET("/books/:book-name", bc.GetBook)
	r.engine.POST("/books/rate", bc.RateBook)
	r.engine.DELETE("/books/rate", bc.DeleteRating)
}

func (r *Router) Run() {
	fmt.Println("Server is running on", r.address)
	if err := r.engine.Run(r.address); err != nil {
		log.Fatalln("can't start server: ", err)
	}
}
