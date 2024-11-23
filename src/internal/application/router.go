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

	recommendationsController "github.com/betterreads/internal/domains/recommendations/controller"
	recommendationsRepository "github.com/betterreads/internal/domains/recommendations/repository"
	recommendationsService "github.com/betterreads/internal/domains/recommendations/service"

	friendsController "github.com/betterreads/internal/domains/friends/controller"
	friendsRepository "github.com/betterreads/internal/domains/friends/repository"
	friendsService "github.com/betterreads/internal/domains/friends/service"

	communitiesController "github.com/betterreads/internal/domains/communities/controller"
	communitiesRepository "github.com/betterreads/internal/domains/communities/repository"
	communitiesService "github.com/betterreads/internal/domains/communities/service"

	feedController "github.com/betterreads/internal/domains/feed/controller"
	feedRepository "github.com/betterreads/internal/domains/feed/repository"
	feedService "github.com/betterreads/internal/domains/feed/service"

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
	users := addUsersHandlers(r, conn)
	books, booksRepo := addBooksHandlers(r, conn)
	AddBookshelfHandlers(r, conn, books)
	AddRecommendationsHandlers(r, conn, books, booksRepo)
	addFriendsHandlers(r, users, conn)
	AddCommunitiesHandlers(r, conn)
	addFeedHandlers(r, users, conn)

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

func addUsersHandlers(r *Router, conn *sqlx.DB) usersService.UsersService {

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
		public.GET("/search", uc.SearchUsers)
	}

	private := r.engine.Group("/users")
	private.Use(middlewares.AuthMiddleware)
	{
		private.GET("/", uc.GetUsers)
		private.POST("/picture", uc.PostPicture)
	}
	return us
}

func addBooksHandlers(r *Router, conn *sqlx.DB) (booksService.BooksService, booksRepository.BooksDatabase) {
	booksRepo, err := booksRepository.NewPostgresBookRepository(conn)
	if err != nil {
		fmt.Println("error: %w", err)
	}
	if booksRepo == nil {
		fmt.Println("booksRepo is nil")
	}
	bs := booksService.NewBooksServiceImpl(booksRepo)
	bc := booksController.NewBooksController(bs)

	public := r.engine.Group("/books")
	public.Use(middlewares.AuthPublicMiddleware)
	{
		public.GET("/:id/picture", bc.GetBookPicture)
		public.GET("/:id/info", bc.GetBookInfo)
		public.GET("/info", bc.GetBooksInfo)
		public.GET("/info/search", bc.SearchBooksInfo)
		public.GET("/:id/reviews", bc.GetBookReviews)
		public.GET("/genres", bc.GetGenres)
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

	return bs, booksRepo
}

func AddBookshelfHandlers(r *Router, conn *sqlx.DB, books booksService.BooksService) {
	bookshelfRepo, err := bookshelfRepository.NewPostgresBookShelfRepository(conn)
	if err != nil {
		fmt.Println("error: %w", err)
	}
	bs := bookshelfService.NewBookShelfServiceImpl(bookshelfRepo, books)
	bc := bookshelfController.NewBookshelfController(bs)

	public := r.engine.Group("/users")
	{
		public.GET("/:id/shelf", bc.GetBookShelf)
		public.GET("/:id/shelf/search", bc.SearchBookShelf)
	}

	private := r.engine.Group("users/shelf")
	private.Use(middlewares.AuthMiddleware)
	{
		private.POST("/", bc.AddBookToShelf)
		private.PUT("/", bc.EditBookInShelf)
		private.DELETE("/", bc.DeleteBookFromShelf)
	}
}

func AddRecommendationsHandlers(r *Router, conn *sqlx.DB, books booksService.BooksService, booksRepo booksRepository.BooksDatabase) {
	recommendationsRepo := recommendationsRepository.NewPostgresRecommendationsRepository(conn, booksRepo)
	rs := recommendationsService.NewRecommendationsServiceImpl(recommendationsRepo, books)
	rc := recommendationsController.NewRecommendationsController(rs)
	private := r.engine.Group("users/recommendations")
	private.Use(middlewares.AuthMiddleware)
	{
		private.GET("/more", rc.GetMoreRecommendations)
		private.GET("/", rc.GetRecommendations)
		private.GET("/friends", rc.GetFriendsRecommendations)
	}
}

func addFriendsHandlers(r *Router, users usersService.UsersService, conn *sqlx.DB) {
	friendsRepo, err := friendsRepository.NewPostgresFriendsRepository(conn)
	if err != nil {
		fmt.Println("error: %w", err)
	}
	fs := friendsService.NewFriendsServiceImpl(friendsRepo, users)
	fc := friendsController.NewFriendsController(fs)
	public := r.engine.Group("users")
	{
		public.GET("/:id/friends", fc.GetFriends)
	}

	private := r.engine.Group("users/friends")
	private.Use(middlewares.AuthMiddleware)
	{
		private.POST("/", fc.AddFriend)
		private.DELETE("/", fc.DeleteFriend)
		private.POST("/requests", fc.AcceptFriendRequest)
		private.DELETE("/requests", fc.RejectFriendRequest)
		private.GET("/requests/sent", fc.GetFriendsRequestSent)
		private.GET("/requests/received", fc.GetFriendRequestsReceived)
	}

}

func AddCommunitiesHandlers(r *Router, conn *sqlx.DB) {
	communitiesRepo, err := communitiesRepository.NewPostgresCommunitiesRepository(conn)
	if err != nil {
		fmt.Println("error: %w", err)
	}
	cs := communitiesService.NewCommunitiesServiceImpl(communitiesRepo)
	cc := communitiesController.NewCommunitiesController(cs)

	public := r.engine.Group("communities")
	{
		public.GET("/:id/picture", cc.GetCommunityPicture)
	}

	private := r.engine.Group("communities")
	private.Use(middlewares.AuthMiddleware)
	{
		private.POST("/", cc.CreateCommunity)
		private.GET("/", cc.GetCommunities)
		private.GET("/search", cc.SearchCommunities)
		// private.GET("/:id", cc.GetCommunity)
		private.POST("/:id/join", cc.JoinCommunity)
		private.GET("/:id/users", cc.GetCommunityUsers)
		// private.DELETE("/:id/leave", cc.LeaveCommunity)
		// private.GET("/:id/members", cc.GetCommunityMembers)
		// gracias por tanto copilot perdon por tan poco
	}
}

func addFeedHandlers(r *Router, users usersService.UsersService, conn *sqlx.DB) {
	feedRepo := feedRepository.NewPostgresFeedRepository(conn)
	fs := feedService.NewFeedServiceImpl(feedRepo, users)
	fc := feedController.NewFeedController(fs)
	private := r.engine.Group("feed")
	private.Use(middlewares.AuthMiddleware)
	{
		private.GET("/", fc.GetFeed)

	}
}

func (r *Router) Run() {
	fmt.Println("Server is running on", r.address)
	if err := r.engine.Run(r.address); err != nil {
		log.Fatalln("can't start server: ", err)
	}
}
