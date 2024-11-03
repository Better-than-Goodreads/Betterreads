package application

import (
	"fmt"

	booksController "github.com/betterreads/internal/domains/books/controller"
	booksRepository "github.com/betterreads/internal/domains/books/repository"
	booksService "github.com/betterreads/internal/domains/books/service"
	"github.com/jmoiron/sqlx"

    bookshelfController "github.com/betterreads/internal/domains/bookshelf/controller"
    bookshelfRepository "github.com/betterreads/internal/domains/bookshelf/repository"
    bookshelfService "github.com/betterreads/internal/domains/bookshelf/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"log"

	usersController "github.com/betterreads/internal/domains/users/controller"
	usersRepository "github.com/betterreads/internal/domains/users/repository"
	usersService "github.com/betterreads/internal/domains/users/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	middlewares "github.com/betterreads/internal/middlewares"
)

type Router struct {
	engine  *gin.Engine
	address string
}

func createRouterFromConfig(cfg *Config) *Router {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// gin.DefaultWriter = io.Discard
	// gin.DefaultErrorWriter = io.Discard
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// slog.SetDefault(logger)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middlewares.ErrorMiddleware)
	engine.Use(middlewares.RequestLogger)

	router := &Router{
		engine:  engine,
		address: cfg.Host + ":" + cfg.Port,
	}

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
	addCorsConfiguration(r)
	addUsersHandlers(r, conn)
    books := addBooksHandlers(r, conn)
    AddBookshelfHandlers(r, conn, books)

	//Adds swagger documentation
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func addCorsConfiguration(r *Router) {
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	r.engine.Use(cors.New(config))
}

func addUsersHandlers(r *Router, conn *sqlx.DB) {

	userRepo, err := usersRepository.NewPostgresUserRepository(conn)
	if err != nil {
		log.Fatalf("can't create db: %v", err)
	}
	us := usersService.NewUsersServiceImpl(userRepo)
	uc := usersController.NewUsersController(us)

	public := r.engine.Group("/users")
	{
		public.POST("/register/basic", uc.RegisterFirstStep)
		public.POST("/register/:id/additional-info", uc.RegisterSecondStep)
		public.POST("/login", uc.LogIn)
		public.GET("/:id", uc.GetUser)
		public.GET("/:id/picture", uc.GetPicture)
	}

	private := r.engine.Group("/users")
	private.Use(middlewares.AuthMiddleware)
	{
		private.GET("/", uc.GetUsers)
		private.POST("/picture", uc.PostPicture)
	}
}

func addBooksHandlers(r *Router, conn *sqlx.DB)  booksService.BooksService{
	booksRepo, err := booksRepository.NewPostgresBookRepository(conn)
	if err != nil {
		fmt.Println("error: %w", err)
	}
	bs := booksService.NewBooksServiceImpl(booksRepo)
	bc := booksController.NewBooksController(bs)

	public := r.engine.Group("/books")
	public.Use(middlewares.AuthPublicMiddleware)
	{
		public.GET("/:id/picture", bc.GetBookPicture)
		public.GET("/:id/info", bc.GetBookInfo)
		public.GET("/info", bc.GetBooksInfo)
		public.GET("/info/search", bc.SearchBooksInfoByName)
        public.GET("/:id/reviews", bc.GetBookReviews)
	}

	private := r.engine.Group("/books")
	private.Use(middlewares.AuthMiddleware)
	{
		private.POST("/", bc.PublishBook)
		private.POST("/:id/reviews", bc.ReviewBook)
		private.POST("/:id/rating", bc.RateBook)
        private.PUT("/:id/rating", bc.UpdateRatingOfBook)
        private.GET("/author/:id", bc.GetBooksOfAuthor)
		private.GET("/user/:id/reviews", bc.GetAllReviewsOfUser)
	}

    return bs
}

func AddBookshelfHandlers(r *Router, conn *sqlx.DB, books booksService.BooksService) {
    bookshelfRepo, err := bookshelfRepository.NewPostgresBookShelfRepository(conn)
    if err != nil {
        fmt.Println("error: %w", err)
    }
    bs := bookshelfService.NewBookShelfServiceImpl(bookshelfRepo, books)
    bc := bookshelfController.NewBookshelfController(bs)

    public := r.engine.Group("/bookshelf")
    {
        public.GET("/:id", bc.GetBookShelf)
    }


    private := r.engine.Group("/bookshelf")
    private.Use(middlewares.AuthMiddleware)
    {
        private.POST("/", bc.AddBookToShelf)
        private.PUT("/", bc.EditBookInShelf)
    }
}

func (r *Router) Run() {
	fmt.Println("Server is running on", r.address)
	if err := r.engine.Run(r.address); err != nil {
		log.Fatalln("can't start server: ", err)
	}
}
