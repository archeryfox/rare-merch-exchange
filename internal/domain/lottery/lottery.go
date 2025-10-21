package lottery

import (
	"time"

	"github.com/google/uuid"
)

// Lottery представляет лотерею
type Lottery struct {
	ID           uuid.UUID     `json:"id" db:"id"`
	ItemID       uuid.UUID     `json:"item_id" db:"item_id"`
	SellerID     uuid.UUID     `json:"seller_id" db:"seller_id"`
	TicketPrice  float64       `json:"ticket_price" db:"ticket_price"`
	MaxTickets   int           `json:"max_tickets" db:"max_tickets"`
	TicketsSold  int           `json:"tickets_sold" db:"tickets_sold"`
	
	// Временные параметры
	StartTime    time.Time     `json:"start_time" db:"start_time"`
	EndTime      time.Time     `json:"end_time" db:"end_time"`
	DrawTime     time.Time     `json:"draw_time" db:"draw_time"`
	
	// Статус и результаты
	Status       LotteryStatus `json:"status" db:"status"`
	WinnerID     *uuid.UUID   `json:"winner_id" db:"winner_id"`
	WinnerTicket *int         `json:"winner_ticket" db:"winner_ticket"`
	
	// Provably Fair
	ServerSeed   string        `json:"server_seed" db:"server_seed"`
	ServerHash   string        `json:"server_hash" db:"server_hash"`
	ResultProof  string        `json:"result_proof" db:"result_proof"`
	
	// Настройки
	Type         LotteryType   `json:"type" db:"type"`
	MinTickets   int           `json:"min_tickets" db:"min_tickets"`
	MaxTicketsPerUser int      `json:"max_tickets_per_user" db:"max_tickets_per_user"`
	RequireKYC   bool          `json:"require_kyc" db:"require_kyc"`
	
	// Метаданные
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" db:"updated_at"`
	DrawnAt      *time.Time    `json:"drawn_at" db:"drawn_at"`
}

// LotteryStatus статус лотереи
type LotteryStatus string

const (
	LotteryStatusDraft     LotteryStatus = "draft"
	LotteryStatusActive    LotteryStatus = "active"
	LotteryStatusEnded     LotteryStatus = "ended"
	LotteryStatusDrawn     LotteryStatus = "drawn"
	LotteryStatusCancelled LotteryStatus = "cancelled"
	LotteryStatusFailed    LotteryStatus = "failed" // недостаточно билетов
)

// LotteryType тип лотереи
type LotteryType string

const (
	LotteryTypeEqual     LotteryType = "equal"      // равные шансы
	LotteryTypeWeighted  LotteryType = "weighted"    // взвешенная по активности
	LotteryTypeInstant   LotteryType = "instant"    // мгновенные выигрыши
	LotteryTypeJackpot   LotteryType = "jackpot"     // джекпот
)

// LotteryTicket представляет билет лотереи
type LotteryTicket struct {
	ID         uuid.UUID `json:"id" db:"id"`
	LotteryID  uuid.UUID `json:"lottery_id" db:"lottery_id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	TicketNumber int     `json:"ticket_number" db:"ticket_number"`
	ClientSeed string    `json:"client_seed" db:"client_seed"`
	PurchasedAt time.Time `json:"purchased_at" db:"purchased_at"`
}

// LotteryDetail детальная информация о лотерее
type LotteryDetail struct {
	Lottery
	Item           ItemInfo     `json:"item"`
	Seller         UserInfo     `json:"seller"`
	Winner         *UserInfo    `json:"winner,omitempty"`
	UserTickets    []int        `json:"user_tickets,omitempty"`
	TimeLeft       int64        `json:"time_left"` // секунды до окончания
	IsActive       bool         `json:"is_active"`
	CanBuyTickets  bool         `json:"can_buy_tickets"`
	TicketsLeft    int          `json:"tickets_left"`
	TotalPrize     float64      `json:"total_prize"`
	UserTicketCount int        `json:"user_ticket_count"`
}

// ItemInfo краткая информация о товаре
type ItemInfo struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	MainPhoto   string    `json:"main_photo"`
	Category    string    `json:"category"`
	Condition   string    `json:"condition"`
	EstimatedValue float64 `json:"estimated_value"`
}

// UserInfo краткая информация о пользователе
type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Rating   float64   `json:"rating"`
	Verified bool      `json:"verified"`
}

// CreateLotteryRequest запрос на создание лотереи
type CreateLotteryRequest struct {
	ItemID            uuid.UUID   `json:"item_id" validate:"required"`
	TicketPrice       float64     `json:"ticket_price" validate:"required,min=1"`
	MaxTickets        int         `json:"max_tickets" validate:"required,min=2,max=10000"`
	Duration          int         `json:"duration" validate:"required,min=1,max=168"` // часы
	Type              LotteryType `json:"type" validate:"required"`
	MinTickets        int         `json:"min_tickets" validate:"required,min=1"`
	MaxTicketsPerUser int         `json:"max_tickets_per_user" validate:"required,min=1"`
	RequireKYC        bool        `json:"require_kyc"`
}

// BuyTicketsRequest запрос на покупку билетов
type BuyTicketsRequest struct {
	TicketCount int `json:"ticket_count" validate:"required,min=1,max=10"`
}

// LotterySearchRequest параметры поиска лотерей
type LotterySearchRequest struct {
	Query     string        `json:"query"`
	Category  string        `json:"category"`
	Status    LotteryStatus `json:"status"`
	Type      LotteryType   `json:"type"`
	MinPrice  float64       `json:"min_price"`
	MaxPrice  float64       `json:"max_price"`
	EndingSoon bool         `json:"ending_soon"`
	SortBy    string        `json:"sort_by"` // price, time_left, tickets_sold
	SortOrder string        `json:"sort_order"` // asc, desc
	Limit     int           `json:"limit" validate:"min=1,max=100"`
	Offset    int           `json:"offset" validate:"min=0"`
}

// DrawResult результат розыгрыша
type DrawResult struct {
	LotteryID     uuid.UUID `json:"lottery_id"`
	WinnerID      uuid.UUID `json:"winner_id"`
	WinnerTicket  int       `json:"winner_ticket"`
	ServerSeed    string    `json:"server_seed"`
	ClientSeeds   []string  `json:"client_seeds"`
	BlockHash     string    `json:"block_hash"`
	ResultHash    string    `json:"result_hash"`
	DrawTime      time.Time `json:"draw_time"`
}

// VerifyDrawRequest запрос на проверку честности розыгрыша
type VerifyDrawRequest struct {
	LotteryID   uuid.UUID `json:"lottery_id" validate:"required"`
	ServerSeed  string    `json:"server_seed" validate:"required"`
	ClientSeeds []string  `json:"client_seeds" validate:"required"`
	BlockHash   string    `json:"block_hash" validate:"required"`
}

// LotteryRepository интерфейс для работы с лотереями
type LotteryRepository interface {
	Create(lottery *Lottery) error
	GetByID(id uuid.UUID) (*Lottery, error)
	GetByItemID(itemID uuid.UUID) (*Lottery, error)
	GetBySellerID(sellerID uuid.UUID, limit, offset int) ([]*Lottery, error)
	Update(lottery *Lottery) error
	Delete(id uuid.UUID) error
	Search(req *LotterySearchRequest) ([]*LotteryDetail, int, error)
	GetActive(limit, offset int) ([]*LotteryDetail, error)
	GetEndingSoon(minutes int, limit int) ([]*LotteryDetail, error)
	GetEnded(limit, offset int) ([]*LotteryDetail, error)
	GetByStatus(status LotteryStatus, limit, offset int) ([]*LotteryDetail, error)
	GetDetail(id uuid.UUID) (*LotteryDetail, error)
	
	// Билеты
	CreateTicket(ticket *LotteryTicket) error
	GetTickets(lotteryID uuid.UUID, limit, offset int) ([]*LotteryTicket, error)
	GetUserTickets(lotteryID, userID uuid.UUID) ([]*LotteryTicket, error)
	GetTicketCount(lotteryID uuid.UUID) (int, error)
	GetUserTicketCount(lotteryID, userID uuid.UUID) (int, error)
	GetAllTickets(lotteryID uuid.UUID) ([]*LotteryTicket, error)
	
	// Розыгрыш
	UpdateDrawResult(lotteryID uuid.UUID, winnerID uuid.UUID, winnerTicket int, resultProof string) error
	GetDrawCandidates(minutes int) ([]*Lottery, error)
}

// LotteryUseCase интерфейс бизнес-логики лотерей
type LotteryUseCase interface {
	CreateLottery(sellerID uuid.UUID, req *CreateLotteryRequest) (*Lottery, error)
	GetLottery(id uuid.UUID) (*LotteryDetail, error)
	UpdateLottery(lotteryID, userID uuid.UUID, req *CreateLotteryRequest) error
	CancelLottery(lotteryID, userID uuid.UUID) error
	SearchLotteries(req *LotterySearchRequest) ([]*LotteryDetail, int, error)
	GetActiveLotteries(limit, offset int) ([]*LotteryDetail, error)
	GetEndingSoonLotteries(minutes int, limit int) ([]*LotteryDetail, error)
	GetUserLotteries(userID uuid.UUID, limit, offset int) ([]*Lottery, error)
	
	// Билеты
	BuyTickets(lotteryID, userID uuid.UUID, req *BuyTicketsRequest) ([]*LotteryTicket, error)
	GetUserTickets(lotteryID, userID uuid.UUID) ([]*LotteryTicket, error)
	GetUserTicketCount(lotteryID, userID uuid.UUID) (int, error)
	
	// Розыгрыш
	DrawLottery(lotteryID uuid.UUID) (*DrawResult, error)
	ProcessEndedLotteries() error
	VerifyDraw(req *VerifyDrawRequest) (bool, error)
	
	// Provably Fair
	GenerateServerSeed() (string, string) // seed, hash
	CalculateWinner(serverSeed string, clientSeeds []string, blockHash string, totalTickets int) int
}
