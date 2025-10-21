package item

import (
	"time"

	"github.com/google/uuid"
)

// Item представляет товар
type Item struct {
	ID          uuid.UUID         `json:"id" db:"id"`
	SellerID    uuid.UUID         `json:"seller_id" db:"seller_id"`
	Title       string            `json:"title" db:"title"`
	Description string            `json:"description" db:"description"`
	Category    Category          `json:"category" db:"category"`
	Brand       string            `json:"brand" db:"brand"`
	Model       string            `json:"model" db:"model"`
	Year        int               `json:"year" db:"year"`
	Condition   Condition         `json:"condition" db:"condition"`
	Size        string            `json:"size" db:"size"`
	Color       string            `json:"color" db:"color"`
	Material    string            `json:"material" db:"material"`
	
	// Фотографии
	Photos      []string          `json:"photos" db:"photos"`
	MainPhoto   string            `json:"main_photo" db:"main_photo"`
	
	// Верификация
	VerificationStatus VerificationStatus `json:"verification_status" db:"verification_status"`
	AuthenticityGrade   AuthenticityGrade  `json:"authenticity_grade" db:"authenticity_grade"`
	VerifiedAt         *time.Time         `json:"verified_at" db:"verified_at"`
	VerifiedBy         *uuid.UUID         `json:"verified_by" db:"verified_by"`
	
	// Цена и состояние
	EstimatedValue float64 `json:"estimated_value" db:"estimated_value"`
	Status         Status  `json:"status" db:"status"`
	
	// Метаданные
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Category категория товара
type Category string

const (
	CategoryClothing     Category = "clothing"
	CategoryShoes        Category = "shoes"
	CategoryAccessories Category = "accessories"
	CategoryToys        Category = "toys"
	CategoryPosters     Category = "posters"
	CategoryCollectibles Category = "collectibles"
	CategoryElectronics Category = "electronics"
	CategoryBooks       Category = "books"
	CategoryOther       Category = "other"
)

// Condition состояние товара
type Condition string

const (
	ConditionNew        Condition = "new"
	ConditionLikeNew    Condition = "like_new"
	ConditionGood       Condition = "good"
	ConditionFair       Condition = "fair"
	ConditionPoor       Condition = "poor"
)

// VerificationStatus статус верификации
type VerificationStatus string

const (
	VerificationStatusPending   VerificationStatus = "pending"
	VerificationStatusVerified  VerificationStatus = "verified"
	VerificationStatusRejected  VerificationStatus = "rejected"
	VerificationStatusNotRequired VerificationStatus = "not_required"
)

// AuthenticityGrade грейд подлинности
type AuthenticityGrade string

const (
	AuthenticityGradeA AuthenticityGrade = "A" // Гарантированный оригинал
	AuthenticityGradeB AuthenticityGrade = "B" // Скорее оригинал
	AuthenticityGradeC AuthenticityGrade = "C" // Требует проверки
	AuthenticityGradeD AuthenticityGrade = "D" // Подозрительный
	AuthenticityGradeF AuthenticityGrade = "F" // Подделка
)

// Status статус товара
type Status string

const (
	StatusDraft     Status = "draft"
	StatusActive    Status = "active"
	StatusSold      Status = "sold"
	StatusReserved  Status = "reserved"
	StatusArchived  Status = "archived"
)

// ItemDetail детальная информация о товаре
type ItemDetail struct {
	Item
	Seller UserInfo `json:"seller"`
}

// UserInfo краткая информация о пользователе
type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Rating   float64   `json:"rating"`
	Verified bool      `json:"verified"`
}

// CreateItemRequest запрос на создание товара
type CreateItemRequest struct {
	Title         string    `json:"title" validate:"required,min=3,max=200"`
	Description   string    `json:"description" validate:"required,min=10,max=2000"`
	Category      Category  `json:"category" validate:"required"`
	Brand         string    `json:"brand" validate:"required"`
	Model         string    `json:"model" validate:"required"`
	Year          int       `json:"year" validate:"required,min=1900,max=2030"`
	Condition     Condition `json:"condition" validate:"required"`
	Size          string    `json:"size" validate:"omitempty"`
	Color         string    `json:"color" validate:"omitempty"`
	Material      string    `json:"material" validate:"omitempty"`
	Photos        []string  `json:"photos" validate:"required,min=1,max=10"`
	EstimatedValue float64  `json:"estimated_value" validate:"required,min=1"`
}

// UpdateItemRequest запрос на обновление товара
type UpdateItemRequest struct {
	Title         string    `json:"title" validate:"omitempty,min=3,max=200"`
	Description   string    `json:"description" validate:"omitempty,min=10,max=2000"`
	Category      Category  `json:"category" validate:"omitempty"`
	Brand         string    `json:"brand" validate:"omitempty"`
	Model         string    `json:"model" validate:"omitempty"`
	Year          int       `json:"year" validate:"omitempty,min=1900,max=2030"`
	Condition     Condition `json:"condition" validate:"omitempty"`
	Size          string    `json:"size" validate:"omitempty"`
	Color         string    `json:"color" validate:"omitempty"`
	Material      string    `json:"material" validate:"omitempty"`
	Photos        []string  `json:"photos" validate:"omitempty,min=1,max=10"`
	EstimatedValue float64  `json:"estimated_value" validate:"omitempty,min=1"`
}

// VerifyItemRequest запрос на верификацию товара
type VerifyItemRequest struct {
	AuthenticityGrade AuthenticityGrade `json:"authenticity_grade" validate:"required"`
	Notes             string            `json:"notes" validate:"omitempty,max=500"`
}

// ItemSearchRequest параметры поиска товаров
type ItemSearchRequest struct {
	Query           string     `json:"query"`
	Category        Category   `json:"category"`
	Brand           string     `json:"brand"`
	Condition       Condition  `json:"condition"`
	MinPrice        float64    `json:"min_price"`
	MaxPrice        float64    `json:"max_price"`
	Year            int        `json:"year"`
	VerificationStatus VerificationStatus `json:"verification_status"`
	SortBy          string     `json:"sort_by"` // price, date, rating
	SortOrder       string     `json:"sort_order"` // asc, desc
	Limit           int        `json:"limit" validate:"min=1,max=100"`
	Offset          int        `json:"offset" validate:"min=0"`
}

// ItemRepository интерфейс для работы с товарами
type ItemRepository interface {
	Create(item *Item) error
	GetByID(id uuid.UUID) (*Item, error)
	GetBySellerID(sellerID uuid.UUID, limit, offset int) ([]*Item, error)
	Update(item *Item) error
	Delete(id uuid.UUID) error
	Search(req *ItemSearchRequest) ([]*ItemDetail, int, error)
	GetByCategory(category Category, limit, offset int) ([]*ItemDetail, error)
	GetFeatured(limit int) ([]*ItemDetail, error)
	UpdateVerificationStatus(id uuid.UUID, status VerificationStatus, grade AuthenticityGrade, verifiedBy uuid.UUID) error
	GetPendingVerification(limit, offset int) ([]*Item, error)
	GetByStatus(status Status, limit, offset int) ([]*ItemDetail, error)
}

// ItemUseCase интерфейс бизнес-логики товаров
type ItemUseCase interface {
	CreateItem(sellerID uuid.UUID, req *CreateItemRequest) (*Item, error)
	GetItem(id uuid.UUID) (*ItemDetail, error)
	UpdateItem(itemID, userID uuid.UUID, req *UpdateItemRequest) error
	DeleteItem(itemID, userID uuid.UUID) error
	SearchItems(req *ItemSearchRequest) ([]*ItemDetail, int, error)
	GetUserItems(userID uuid.UUID, limit, offset int) ([]*Item, error)
	GetFeaturedItems(limit int) ([]*ItemDetail, error)
	GetItemsByCategory(category Category, limit, offset int) ([]*ItemDetail, error)
	SubmitForVerification(itemID, userID uuid.UUID) error
	VerifyItem(itemID uuid.UUID, req *VerifyItemRequest, verifiedBy uuid.UUID) error
	GetPendingVerification(limit, offset int) ([]*Item, error)
	GetItemsByStatus(status Status, limit, offset int) ([]*ItemDetail, error)
}
