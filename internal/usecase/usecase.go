package usecase

import (
	"rare-merch-exchange/internal/domain/auction"
	"rare-merch-exchange/internal/domain/contest"
	"rare-merch-exchange/internal/domain/item"
	"rare-merch-exchange/internal/domain/lottery"
	"rare-merch-exchange/internal/domain/transaction"
	"rare-merch-exchange/internal/domain/user"
	"rare-merch-exchange/internal/domain/verification"
	"rare-merch-exchange/internal/pkg/config"
	"rare-merch-exchange/internal/repository"

	"go.uber.org/zap"
)

// UseCases содержит все use cases
type UseCases struct {
	User         user.UserUseCase
	Item         item.ItemUseCase
	Auction      auction.AuctionUseCase
	Lottery      lottery.LotteryUseCase
	Contest      contest.ContestUseCase
	Transaction  transaction.TransactionUseCase
	Verification verification.VerificationUseCase
}

// NewUseCases создаёт новый экземпляр use cases
func NewUseCases(repos *repository.Repositories, cfg *config.Config, logger *zap.Logger) *UseCases {
	return &UseCases{
		User:         NewUserUseCase(repos.User, cfg, logger),
		Item:         NewItemUseCase(repos.Item, cfg, logger),
		Auction:      NewAuctionUseCase(repos.Auction, cfg, logger),
		Lottery:      NewLotteryUseCase(repos.Lottery, cfg, logger),
		Contest:      NewContestUseCase(repos.Contest, cfg, logger),
		Transaction:  NewTransactionUseCase(repos.Transaction, cfg, logger),
		Verification: NewVerificationUseCase(repos.Verification, cfg, logger),
	}
}
