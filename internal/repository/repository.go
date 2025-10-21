package repository

import (
	"database/sql"
	"fmt"
	"time"
	"rare-merch-exchange/internal/domain/auction"
	"rare-merch-exchange/internal/domain/contest"
	"rare-merch-exchange/internal/domain/item"
	"rare-merch-exchange/internal/domain/lottery"
	"rare-merch-exchange/internal/domain/transaction"
	"rare-merch-exchange/internal/domain/user"
	"rare-merch-exchange/internal/domain/verification"

	"github.com/redis/go-redis/v9"
)

// Repositories содержит все репозитории
type Repositories struct {
	User         user.UserRepository
	Item         item.ItemRepository
	Auction      auction.AuctionRepository
	Lottery      lottery.LotteryRepository
	Contest      contest.ContestRepository
	Transaction  transaction.TransactionRepository
	Verification verification.VerificationRepository
}

// NewRepositories создаёт новый экземпляр репозиториев
func NewRepositories(db *sql.DB, rdb *redis.Client) *Repositories {
	return &Repositories{
		User:         NewUserRepository(db, rdb),
		Item:         NewItemRepository(db, rdb),
		Auction:      NewAuctionRepository(db, rdb),
		Lottery:      NewLotteryRepository(db, rdb),
		Contest:      NewContestRepository(db, rdb),
		Transaction:  NewTransactionRepository(db, rdb),
		Verification: NewVerificationRepository(db, rdb),
	}
}

// NewPostgresDB создаёт подключение к PostgreSQL
func NewPostgresDB(cfg DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// DatabaseConfig конфигурация базы данных
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}
