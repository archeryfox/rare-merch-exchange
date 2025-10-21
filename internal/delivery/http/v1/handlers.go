package v1

import (
	"rare-merch-exchange/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handlers содержит все HTTP handlers
type Handlers struct {
	User         *UserHandler
	Item         *ItemHandler
	Auction      *AuctionHandler
	Lottery      *LotteryHandler
	Contest      *ContestHandler
	Transaction  *TransactionHandler
	Verification *VerificationHandler
}

// NewHandlers создаёт новый экземпляр handlers
func NewHandlers(useCases *usecase.UseCases, logger *zap.Logger) *Handlers {
	return &Handlers{
		User:         NewUserHandler(useCases.User, logger),
		Item:         NewItemHandler(useCases.Item, logger),
		Auction:      NewAuctionHandler(useCases.Auction, logger),
		Lottery:      NewLotteryHandler(useCases.Lottery, logger),
		Contest:      NewContestHandler(useCases.Contest, logger),
		Transaction:  NewTransactionHandler(useCases.Transaction, logger),
		Verification: NewVerificationHandler(useCases.Verification, logger),
	}
}

// SetupRoutes настраивает маршруты API
func SetupRoutes(router *gin.Engine, handlers *Handlers) {
	v1 := router.Group("/api/v1")
	{
		// Аутентификация
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.User.Register)
			auth.POST("/login", handlers.User.Login)
		}

		// Пользователи
		users := v1.Group("/users")
		users.Use(AuthMiddleware()) // Требует аутентификации
		{
			users.GET("/me", handlers.User.GetProfile)
			users.PUT("/me", handlers.User.UpdateProfile)
			users.PUT("/me/password", handlers.User.ChangePassword)
			users.PUT("/me/notifications", handlers.User.UpdateNotificationSettings)
			users.GET("/me/stats", handlers.User.GetStats)
			users.GET("/search", handlers.User.SearchUsers)
			users.GET("/top-rated", handlers.User.GetTopRatedUsers)
			users.GET("/:id", handlers.User.GetUserProfile)
		}

		// Товары
		items := v1.Group("/items")
		{
			items.GET("/", handlers.Item.SearchItems)
			items.GET("/featured", handlers.Item.GetFeaturedItems)
			items.GET("/category/:category", handlers.Item.GetItemsByCategory)
			items.GET("/:id", handlers.Item.GetItem)
			items.POST("/", AuthMiddleware(), handlers.Item.CreateItem)
			items.PUT("/:id", AuthMiddleware(), handlers.Item.UpdateItem)
			items.DELETE("/:id", AuthMiddleware(), handlers.Item.DeleteItem)
			items.POST("/:id/verify", AuthMiddleware(), handlers.Item.SubmitForVerification)
			items.GET("/:id/authenticity", handlers.Item.GetAuthenticityStatus)
		}

		// Аукционы
		auctions := v1.Group("/auctions")
		{
			auctions.GET("/", handlers.Auction.SearchAuctions)
			auctions.GET("/active", handlers.Auction.GetActiveAuctions)
			auctions.GET("/ending-soon", handlers.Auction.GetEndingSoonAuctions)
			auctions.GET("/:id", handlers.Auction.GetAuction)
			auctions.GET("/:id/bids", handlers.Auction.GetBids)
			auctions.POST("/", AuthMiddleware(), handlers.Auction.CreateAuction)
			auctions.PUT("/:id", AuthMiddleware(), handlers.Auction.UpdateAuction)
			auctions.DELETE("/:id", AuthMiddleware(), handlers.Auction.CancelAuction)
			auctions.POST("/:id/bid", AuthMiddleware(), handlers.Auction.PlaceBid)
			auctions.POST("/:id/autobid", AuthMiddleware(), handlers.Auction.CreateAutoBid)
			auctions.PUT("/:id/autobid/:autobid_id", AuthMiddleware(), handlers.Auction.UpdateAutoBid)
			auctions.DELETE("/:id/autobid/:autobid_id", AuthMiddleware(), handlers.Auction.DeleteAutoBid)
			auctions.GET("/:id/autobids", handlers.Auction.GetAutoBids)
		}

		// Лотереи
		lotteries := v1.Group("/lotteries")
		{
			lotteries.GET("/", handlers.Lottery.SearchLotteries)
			lotteries.GET("/active", handlers.Lottery.GetActiveLotteries)
			lotteries.GET("/ending-soon", handlers.Lottery.GetEndingSoonLotteries)
			lotteries.GET("/:id", handlers.Lottery.GetLottery)
			lotteries.GET("/:id/tickets", AuthMiddleware(), handlers.Lottery.GetUserTickets)
			lotteries.POST("/", AuthMiddleware(), handlers.Lottery.CreateLottery)
			lotteries.PUT("/:id", AuthMiddleware(), handlers.Lottery.UpdateLottery)
			lotteries.DELETE("/:id", AuthMiddleware(), handlers.Lottery.CancelLottery)
			lotteries.POST("/:id/tickets", AuthMiddleware(), handlers.Lottery.BuyTickets)
			lotteries.GET("/:id/verify", handlers.Lottery.VerifyDraw)
		}

		// Конкурсы
		contests := v1.Group("/contests")
		{
			contests.GET("/", handlers.Contest.SearchContests)
			contests.GET("/active", handlers.Contest.GetActiveContests)
			contests.GET("/ending-soon", handlers.Contest.GetEndingSoonContests)
			contests.GET("/:id", handlers.Contest.GetContest)
			contests.GET("/:id/submissions", handlers.Contest.GetSubmissions)
			contests.POST("/", AuthMiddleware(), handlers.Contest.CreateContest)
			contests.PUT("/:id", AuthMiddleware(), handlers.Contest.UpdateContest)
			contests.DELETE("/:id", AuthMiddleware(), handlers.Contest.CancelContest)
			contests.POST("/:id/submit", AuthMiddleware(), handlers.Contest.SubmitContest)
			contests.POST("/:id/vote", AuthMiddleware(), handlers.Contest.VoteContest)
			contests.GET("/:id/votes", handlers.Contest.GetVotes)
		}

		// Транзакции
		transactions := v1.Group("/transactions")
		transactions.Use(AuthMiddleware())
		{
			transactions.GET("/", handlers.Transaction.GetUserTransactions)
			transactions.GET("/as-buyer", handlers.Transaction.GetUserTransactionsAsBuyer)
			transactions.GET("/as-seller", handlers.Transaction.GetUserTransactionsAsSeller)
			transactions.GET("/:id", handlers.Transaction.GetTransaction)
			transactions.POST("/", handlers.Transaction.CreateTransaction)
			transactions.PUT("/:id/shipping", handlers.Transaction.UpdateShipping)
			transactions.POST("/:id/confirm", handlers.Transaction.ConfirmDelivery)
			transactions.POST("/:id/dispute", handlers.Transaction.CreateDispute)
			transactions.GET("/disputes", handlers.Transaction.GetUserDisputes)
			transactions.GET("/disputes/:id", handlers.Transaction.GetDispute)
		}

		// Верификация
		verification := v1.Group("/verification")
		{
			verification.GET("/pending", AuthMiddleware(), handlers.Verification.GetPendingVerifications)
			verification.GET("/stats", handlers.Verification.GetVerificationStats)
			verification.POST("/submit", AuthMiddleware(), handlers.Verification.SubmitVerification)
			verification.GET("/:id", handlers.Verification.GetVerification)
			verification.POST("/:id/verify", AuthMiddleware(), handlers.Verification.VerifyItem)
			verification.GET("/experts", handlers.Verification.GetExperts)
			verification.POST("/experts/register", AuthMiddleware(), handlers.Verification.RegisterExpert)
		}
	}
}

// AuthMiddleware middleware для аутентификации
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: реализовать JWT аутентификацию
		c.Next()
	}
}
