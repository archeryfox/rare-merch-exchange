package auction

import (
	"time"

	"github.com/google/uuid"
)

// Auction представляет аукцион
type Auction struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	ItemID          uuid.UUID     `json:"item_id" db:"item_id"`
	SellerID        uuid.UUID     `json:"seller_id" db:"seller_id"`
	StartPrice      float64       `json:"start_price" db:"start_price"`
	ReservePrice    *float64      `json:"reserve_price" db:"reserve_price"`
	CurrentPrice    float64       `json:"current_price" db:"current_price"`
	BuyNowPrice     *float64      `json:"buy_now_price" db:"buy_now_price"`
	MinIncrement    float64       `json:"min_increment" db:"min_increment"`
	
	// Временные параметры
	StartTime       time.Time     `json:"start_time" db:"start_time"`
	EndTime         time.Time     `json:"end_time" db:"end_time"`
	SoftCloseMinutes int          `json:"soft_close_minutes" db:"soft_close_minutes"`
	
	// Статус и результаты
	Status          AuctionStatus `json:"status" db:"status"`
	WinnerID        *uuid.UUID    `json:"winner_id" db:"winner_id"`
	FinalPrice      *float64      `json:"final_price" db:"final_price"`
	TotalBids       int           `json:"total_bids" db:"total_bids"`
	
	// Настройки
	AutoExtend      bool          `json:"auto_extend" db:"auto_extend"`
	RequireDeposit  bool          `json:"require_deposit" db:"require_deposit"`
	DepositAmount   float64       `json:"deposit_amount" db:"deposit_amount"`
	
	// Метаданные
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
	EndedAt         *time.Time    `json:"ended_at" db:"ended_at"`
}

// AuctionStatus статус аукциона
type AuctionStatus string

const (
	AuctionStatusDraft     AuctionStatus = "draft"
	AuctionStatusActive    AuctionStatus = "active"
	AuctionStatusEnded     AuctionStatus = "ended"
	AuctionStatusCancelled AuctionStatus = "cancelled"
	AuctionStatusSold      AuctionStatus = "sold"
	AuctionStatusNoBids    AuctionStatus = "no_bids"
	AuctionStatusReserveNotMet AuctionStatus = "reserve_not_met"
)

// Bid представляет ставку
type Bid struct {
	ID        uuid.UUID `json:"id" db:"id"`
	AuctionID uuid.UUID `json:"auction_id" db:"auction_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Amount    float64   `json:"amount" db:"amount"`
	IsAutoBid bool      `json:"is_auto_bid" db:"is_auto_bid"`
	MaxAmount *float64  `json:"max_amount" db:"max_amount"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

// AutoBid представляет автоматическую ставку
type AutoBid struct {
	ID        uuid.UUID `json:"id" db:"id"`
	AuctionID uuid.UUID `json:"auction_id" db:"auction_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	MaxAmount float64   `json:"max_amount" db:"max_amount"`
	Increment float64   `json:"increment" db:"increment"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AuctionDetail детальная информация об аукционе
type AuctionDetail struct {
	Auction
	Item        ItemInfo     `json:"item"`
	Seller      UserInfo     `json:"seller"`
	Winner      *UserInfo    `json:"winner,omitempty"`
	Bids        []BidInfo    `json:"bids"`
	AutoBids    []AutoBidInfo `json:"auto_bids"`
	TimeLeft    int64        `json:"time_left"` // секунды до окончания
	IsActive    bool         `json:"is_active"`
	CanBid      bool         `json:"can_bid"`
	MinBid      float64      `json:"min_bid"`
}

// ItemInfo краткая информация о товаре
type ItemInfo struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	MainPhoto   string    `json:"main_photo"`
	Category    string    `json:"category"`
	Condition   string    `json:"condition"`
}

// UserInfo краткая информация о пользователе
type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Rating   float64   `json:"rating"`
	Verified bool      `json:"verified"`
}

// BidInfo информация о ставке
type BidInfo struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Amount    float64   `json:"amount"`
	IsAutoBid bool      `json:"is_auto_bid"`
	Timestamp time.Time `json:"timestamp"`
}

// AutoBidInfo информация об автоставке
type AutoBidInfo struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	MaxAmount float64   `json:"max_amount"`
	Increment float64   `json:"increment"`
	IsActive  bool      `json:"is_active"`
}

// CreateAuctionRequest запрос на создание аукциона
type CreateAuctionRequest struct {
	ItemID           uuid.UUID  `json:"item_id" validate:"required"`
	StartPrice       float64    `json:"start_price" validate:"required,min=1"`
	ReservePrice     *float64   `json:"reserve_price" validate:"omitempty,min=1"`
	BuyNowPrice      *float64   `json:"buy_now_price" validate:"omitempty,min=1"`
	MinIncrement     float64    `json:"min_increment" validate:"required,min=1"`
	Duration         int        `json:"duration" validate:"required,min=1,max=168"` // часы
	SoftCloseMinutes int        `json:"soft_close_minutes" validate:"omitempty,min=1,max=30"`
	AutoExtend       bool       `json:"auto_extend"`
	RequireDeposit   bool       `json:"require_deposit"`
}

// PlaceBidRequest запрос на размещение ставки
type PlaceBidRequest struct {
	Amount float64 `json:"amount" validate:"required,min=1"`
}

// CreateAutoBidRequest запрос на создание автоставки
type CreateAutoBidRequest struct {
	MaxAmount float64 `json:"max_amount" validate:"required,min=1"`
	Increment float64 `json:"increment" validate:"required,min=1"`
}

// AuctionSearchRequest параметры поиска аукционов
type AuctionSearchRequest struct {
	Query       string         `json:"query"`
	Category    string         `json:"category"`
	Status      AuctionStatus  `json:"status"`
	MinPrice    float64        `json:"min_price"`
	MaxPrice    float64        `json:"max_price"`
	EndingSoon  bool           `json:"ending_soon"`
	SortBy      string         `json:"sort_by"` // price, time_left, bids
	SortOrder   string         `json:"sort_order"` // asc, desc
	Limit       int            `json:"limit" validate:"min=1,max=100"`
	Offset      int            `json:"offset" validate:"min=0"`
}

// AuctionRepository интерфейс для работы с аукционами
type AuctionRepository interface {
	Create(auction *Auction) error
	GetByID(id uuid.UUID) (*Auction, error)
	GetByItemID(itemID uuid.UUID) (*Auction, error)
	GetBySellerID(sellerID uuid.UUID, limit, offset int) ([]*Auction, error)
	Update(auction *Auction) error
	Delete(id uuid.UUID) error
	Search(req *AuctionSearchRequest) ([]*AuctionDetail, int, error)
	GetActive(limit, offset int) ([]*AuctionDetail, error)
	GetEndingSoon(minutes int, limit int) ([]*AuctionDetail, error)
	GetEnded(limit, offset int) ([]*AuctionDetail, error)
	GetByStatus(status AuctionStatus, limit, offset int) ([]*AuctionDetail, error)
	GetDetail(id uuid.UUID) (*AuctionDetail, error)
	
	// Ставки
	CreateBid(bid *Bid) error
	GetBids(auctionID uuid.UUID, limit, offset int) ([]*Bid, error)
	GetHighestBid(auctionID uuid.UUID) (*Bid, error)
	GetUserBids(userID uuid.UUID, limit, offset int) ([]*Bid, error)
	
	// Автоставки
	CreateAutoBid(autoBid *AutoBid) error
	GetAutoBids(auctionID uuid.UUID) ([]*AutoBid, error)
	GetUserAutoBid(auctionID, userID uuid.UUID) (*AutoBid, error)
	UpdateAutoBid(autoBid *AutoBid) error
	DeleteAutoBid(id uuid.UUID) error
}

// AuctionUseCase интерфейс бизнес-логики аукционов
type AuctionUseCase interface {
	CreateAuction(sellerID uuid.UUID, req *CreateAuctionRequest) (*Auction, error)
	GetAuction(id uuid.UUID) (*AuctionDetail, error)
	UpdateAuction(auctionID, userID uuid.UUID, req *CreateAuctionRequest) error
	CancelAuction(auctionID, userID uuid.UUID) error
	SearchAuctions(req *AuctionSearchRequest) ([]*AuctionDetail, int, error)
	GetActiveAuctions(limit, offset int) ([]*AuctionDetail, error)
	GetEndingSoonAuctions(minutes int, limit int) ([]*AuctionDetail, error)
	GetUserAuctions(userID uuid.UUID, limit, offset int) ([]*Auction, error)
	GetUserBids(userID uuid.UUID, limit, offset int) ([]*Bid, error)
	
	// Ставки
	PlaceBid(auctionID, userID uuid.UUID, req *PlaceBidRequest) (*Bid, error)
	GetBids(auctionID uuid.UUID, limit, offset int) ([]*Bid, error)
	
	// Автоставки
	CreateAutoBid(auctionID, userID uuid.UUID, req *CreateAutoBidRequest) (*AutoBid, error)
	UpdateAutoBid(autoBidID, userID uuid.UUID, req *CreateAutoBidRequest) error
	DeleteAutoBid(autoBidID, userID uuid.UUID) error
	GetAutoBids(auctionID uuid.UUID) ([]*AutoBid, error)
	
	// Завершение аукционов
	EndAuction(auctionID uuid.UUID) error
	ProcessEndedAuctions() error
	ExtendAuction(auctionID uuid.UUID) error
}
