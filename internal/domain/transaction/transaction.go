package transaction

import (
	"time"

	"github.com/google/uuid"
)

// Transaction представляет транзакцию
type Transaction struct {
	ID            uuid.UUID         `json:"id" db:"id"`
	Type          TransactionType    `json:"type" db:"type"`
	BuyerID       uuid.UUID         `json:"buyer_id" db:"buyer_id"`
	SellerID      uuid.UUID         `json:"seller_id" db:"seller_id"`
	ItemID        uuid.UUID         `json:"item_id" db:"item_id"`
	AuctionID     *uuid.UUID        `json:"auction_id" db:"auction_id"`
	LotteryID     *uuid.UUID        `json:"lottery_id" db:"lottery_id"`
	ContestID     *uuid.UUID        `json:"contest_id" db:"contest_id"`
	
	// Финансовые данные
	Amount        float64           `json:"amount" db:"amount"`
	Commission    float64           `json:"commission" db:"commission"`
	NetAmount     float64           `json:"net_amount" db:"net_amount"`
	Currency      string            `json:"currency" db:"currency"`
	
	// Escrow данные
	EscrowAmount  float64           `json:"escrow_amount" db:"escrow_amount"`
	EscrowStatus  EscrowStatus      `json:"escrow_status" db:"escrow_status"`
	EscrowReleasedAt *time.Time     `json:"escrow_released_at" db:"escrow_released_at"`
	
	// Доставка
	TrackingNumber string           `json:"tracking_number" db:"tracking_number"`
	ShippingMethod string           `json:"shipping_method" db:"shipping_method"`
	ShippingAddress string          `json:"shipping_address" db:"shipping_address"`
	ShippedAt     *time.Time        `json:"shipped_at" db:"shipped_at"`
	DeliveredAt   *time.Time        `json:"delivered_at" db:"delivered_at"`
	
	// Статус и результаты
	Status        TransactionStatus `json:"status" db:"status"`
	DisputeID     *uuid.UUID       `json:"dispute_id" db:"dispute_id"`
	Notes         string           `json:"notes" db:"notes"`
	
	// Метаданные
	CreatedAt     time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at" db:"updated_at"`
	CompletedAt   *time.Time       `json:"completed_at" db:"completed_at"`
}

// TransactionType тип транзакции
type TransactionType string

const (
	TransactionTypeAuction TransactionType = "auction"
	TransactionTypeLottery TransactionType = "lottery"
	TransactionTypeContest TransactionType = "contest"
	TransactionTypeDirect  TransactionType = "direct"
)

// TransactionStatus статус транзакции
type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusPaid       TransactionStatus = "paid"
	TransactionStatusShipped    TransactionStatus = "shipped"
	TransactionStatusDelivered  TransactionStatus = "delivered"
	TransactionStatusCompleted  TransactionStatus = "completed"
	TransactionStatusDisputed   TransactionStatus = "disputed"
	TransactionStatusCancelled  TransactionStatus = "cancelled"
	TransactionStatusRefunded   TransactionStatus = "refunded"
)

// EscrowStatus статус escrow
type EscrowStatus string

const (
	EscrowStatusPending   EscrowStatus = "pending"
	EscrowStatusHeld      EscrowStatus = "held"
	EscrowStatusReleased  EscrowStatus = "released"
	EscrowStatusRefunded  EscrowStatus = "refunded"
	EscrowStatusDisputed  EscrowStatus = "disputed"
)

// Dispute представляет спор по транзакции
type Dispute struct {
	ID            uuid.UUID     `json:"id" db:"id"`
	TransactionID uuid.UUID      `json:"transaction_id" db:"transaction_id"`
	ComplainantID uuid.UUID      `json:"complainant_id" db:"complainant_id"`
	Reason        DisputeReason  `json:"reason" db:"reason"`
	Description   string         `json:"description" db:"description"`
	Evidence      []string       `json:"evidence" db:"evidence"`
	Status        DisputeStatus  `json:"status" db:"status"`
	Resolution    *string        `json:"resolution" db:"resolution"`
	ResolvedBy    *uuid.UUID     `json:"resolved_by" db:"resolved_by"`
	ResolvedAt    *time.Time     `json:"resolved_at" db:"resolved_at"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
}

// DisputeReason причина спора
type DisputeReason string

const (
	DisputeReasonItemNotReceived DisputeReason = "item_not_received"
	DisputeReasonItemNotAsDescribed DisputeReason = "item_not_as_described"
	DisputeReasonItemDamaged     DisputeReason = "item_damaged"
	DisputeReasonItemFake        DisputeReason = "item_fake"
	DisputeReasonSellerNotResponding DisputeReason = "seller_not_responding"
	DisputeReasonOther           DisputeReason = "other"
)

// DisputeStatus статус спора
type DisputeStatus string

const (
	DisputeStatusOpen      DisputeStatus = "open"
	DisputeStatusUnderReview DisputeStatus = "under_review"
	DisputeStatusResolved  DisputeStatus = "resolved"
	DisputeStatusClosed    DisputeStatus = "closed"
)

// TransactionDetail детальная информация о транзакции
type TransactionDetail struct {
	Transaction
	Item        ItemInfo     `json:"item"`
	Buyer       UserInfo     `json:"buyer"`
	Seller      UserInfo     `json:"seller"`
	Dispute     *DisputeInfo `json:"dispute,omitempty"`
	TimeLeft    int64        `json:"time_left"` // секунды до авторелиза
	CanDispute  bool         `json:"can_dispute"`
	CanConfirm  bool         `json:"can_confirm"`
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

// DisputeInfo информация о споре
type DisputeInfo struct {
	ID          uuid.UUID     `json:"id"`
	Reason      DisputeReason `json:"reason"`
	Description string        `json:"description"`
	Status      DisputeStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
}

// CreateTransactionRequest запрос на создание транзакции
type CreateTransactionRequest struct {
	Type            TransactionType `json:"type" validate:"required"`
	ItemID          uuid.UUID       `json:"item_id" validate:"required"`
	AuctionID       *uuid.UUID      `json:"auction_id" validate:"omitempty"`
	LotteryID       *uuid.UUID      `json:"lottery_id" validate:"omitempty"`
	ContestID       *uuid.UUID      `json:"contest_id" validate:"omitempty"`
	Amount          float64         `json:"amount" validate:"required,min=1"`
	ShippingAddress string          `json:"shipping_address" validate:"required"`
	ShippingMethod  string          `json:"shipping_method" validate:"required"`
}

// UpdateShippingRequest запрос на обновление данных доставки
type UpdateShippingRequest struct {
	TrackingNumber string `json:"tracking_number" validate:"required"`
	ShippingMethod string `json:"shipping_method" validate:"required"`
}

// ConfirmDeliveryRequest запрос на подтверждение получения
type ConfirmDeliveryRequest struct {
	Rating    int    `json:"rating" validate:"required,min=1,max=5"`
	Comment   string `json:"comment" validate:"omitempty,max=500"`
}

// CreateDisputeRequest запрос на создание спора
type CreateDisputeRequest struct {
	Reason      DisputeReason `json:"reason" validate:"required"`
	Description string        `json:"description" validate:"required,min=10,max=1000"`
	Evidence    []string      `json:"evidence" validate:"omitempty,max=10"`
}

// ResolveDisputeRequest запрос на разрешение спора
type ResolveDisputeRequest struct {
	Resolution string `json:"resolution" validate:"required,min=10,max=1000"`
	RefundBuyer bool  `json:"refund_buyer"`
	RefundSeller bool `json:"refund_seller"`
}

// TransactionSearchRequest параметры поиска транзакций
type TransactionSearchRequest struct {
	Type     TransactionType    `json:"type"`
	Status   TransactionStatus `json:"status"`
	MinAmount float64          `json:"min_amount"`
	MaxAmount float64          `json:"max_amount"`
	SortBy   string            `json:"sort_by"` // amount, date, status
	SortOrder string           `json:"sort_order"` // asc, desc
	Limit    int               `json:"limit" validate:"min=1,max=100"`
	Offset   int               `json:"offset" validate:"min=0"`
}

// TransactionRepository интерфейс для работы с транзакциями
type TransactionRepository interface {
	Create(transaction *Transaction) error
	GetByID(id uuid.UUID) (*Transaction, error)
	GetByUserID(userID uuid.UUID, limit, offset int) ([]*Transaction, error)
	GetByBuyerID(buyerID uuid.UUID, limit, offset int) ([]*Transaction, error)
	GetBySellerID(sellerID uuid.UUID, limit, offset int) ([]*Transaction, error)
	Update(transaction *Transaction) error
	Search(req *TransactionSearchRequest, userID uuid.UUID) ([]*TransactionDetail, int, error)
	GetDetail(id uuid.UUID) (*TransactionDetail, error)
	GetPendingEscrow(limit, offset int) ([]*Transaction, error)
	GetDisputed(limit, offset int) ([]*Transaction, error)
	
	// Disputes
	CreateDispute(dispute *Dispute) error
	GetDisputeByID(id uuid.UUID) (*Dispute, error)
	GetDisputeByTransactionID(transactionID uuid.UUID) (*Dispute, error)
	UpdateDispute(dispute *Dispute) error
	GetDisputes(limit, offset int) ([]*Dispute, error)
	GetUserDisputes(userID uuid.UUID, limit, offset int) ([]*Dispute, error)
}

// TransactionUseCase интерфейс бизнес-логики транзакций
type TransactionUseCase interface {
	CreateTransaction(buyerID uuid.UUID, req *CreateTransactionRequest) (*Transaction, error)
	GetTransaction(id uuid.UUID) (*TransactionDetail, error)
	GetUserTransactions(userID uuid.UUID, limit, offset int) ([]*Transaction, error)
	GetUserTransactionsAsBuyer(userID uuid.UUID, limit, offset int) ([]*Transaction, error)
	GetUserTransactionsAsSeller(userID uuid.UUID, limit, offset int) ([]*Transaction, error)
	SearchTransactions(req *TransactionSearchRequest, userID uuid.UUID) ([]*TransactionDetail, int, error)
	
	// Доставка
	UpdateShipping(transactionID, userID uuid.UUID, req *UpdateShippingRequest) error
	ConfirmDelivery(transactionID, userID uuid.UUID, req *ConfirmDeliveryRequest) error
	
	// Escrow
	ReleaseEscrow(transactionID uuid.UUID) error
	RefundEscrow(transactionID uuid.UUID, reason string) error
	ProcessPendingEscrow() error
	
	// Споры
	CreateDispute(transactionID, userID uuid.UUID, req *CreateDisputeRequest) (*Dispute, error)
	ResolveDispute(disputeID uuid.UUID, req *ResolveDisputeRequest, resolvedBy uuid.UUID) error
	GetDispute(id uuid.UUID) (*Dispute, error)
	GetUserDisputes(userID uuid.UUID, limit, offset int) ([]*Dispute, error)
	GetDisputes(limit, offset int) ([]*Dispute, error)
	
	// Автоматизация
	ProcessAutoRelease() error
	ProcessDisputeTimeouts() error
}
