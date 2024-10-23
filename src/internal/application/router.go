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
	addCorsConfiguration(r)
	addUsersHandlers(r, conn)
	addBooksHandlers(r, conn)

	//Adds swagger documentation
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func addCorsConfiguration(r *Router) {
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	r.engine.Use(cors.New(config))
}

func addUsersHandlers(r *Router, conn *sqlx.DB) {
	// userRepo := usersRepository.NewMemoryDatabase()

	userRepo, err := usersRepository.NewPostgresUserRepository(conn)
	if err != nil {
		log.Fatalf("can't create db: %v", err)
	}
	us := usersService.NewUsersService(userRepo)
	uc := usersController.NewUsersController(us)

	public := r.engine.Group("/users")
	{
		public.POST("/register/basic", uc.RegisterFirstStep)
		public.POST("/register/:id/additional-info", uc.RegisterSecondStep)
		public.POST("/login", uc.LogIn)
	}

	private := r.engine.Group("/users")
	private.Use(middlewares.AuthMiddleware)
	{
		private.GET("/", uc.GetUsers)
		private.GET("/:id", uc.GetUser)
	}
}

func addBooksHandlers(r *Router, conn *sqlx.DB) {
	booksRepo, err := booksRepository.NewPostgresBookRepository(conn)
	if err != nil {
		fmt.Println("error: %w", err)
	}
	bs := booksService.NewBooksService(booksRepo)
	bc := booksController.NewBooksController(bs)

	public := r.engine.Group("/books")
	{
		public.GET("/", bc.GetBooks)
		public.GET("/:id", bc.GetBook)
	}

	private := r.engine.Group("/books")
	private.Use(middlewares.AuthMiddleware)
	{
		private.POST("/", bc.PublishBook)
		private.GET("/:id/rating/", bc.GetRatingUser)
		private.POST("/:id/rating", bc.RateBook)
		private.DELETE("/:id/rating/", bc.DeleteRating)
	}

}

func (r *Router) Run() {
	fmt.Println("Server is running on", r.address)
	if err := r.engine.Run(r.address); err != nil {
		log.Fatalln("can't start server: ", err)
	}
}
