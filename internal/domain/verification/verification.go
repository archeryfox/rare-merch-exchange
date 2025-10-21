package verification

import (
	"time"

	"github.com/google/uuid"
)

// VerificationRequest представляет запрос на верификацию
type VerificationRequest struct {
	ID            uuid.UUID           `json:"id" db:"id"`
	ItemID        uuid.UUID           `json:"item_id" db:"item_id"`
	UserID        uuid.UUID           `json:"user_id" db:"user_id"`
	Status        VerificationStatus  `json:"status" db:"status"`
	Priority      VerificationPriority `json:"priority" db:"priority"`
	
	// Данные для проверки
	Photos        []string            `json:"photos" db:"photos"`
	Description   string              `json:"description" db:"description"`
	PurchaseProof []string            `json:"purchase_proof" db:"purchase_proof"`
	SerialNumber  string              `json:"serial_number" db:"serial_number"`
	Receipt       string              `json:"receipt" db:"receipt"`
	
	// Результат верификации
	AuthenticityGrade AuthenticityGrade `json:"authenticity_grade" db:"authenticity_grade"`
	VerifiedBy        *uuid.UUID        `json:"verified_by" db:"verified_by"`
	VerificationNotes  string            `json:"verification_notes" db:"verification_notes"`
	VerifiedAt        *time.Time        `json:"verified_at" db:"verified_at"`
	
	// Метаданные
	CreatedAt    time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at" db:"updated_at"`
}

// VerificationStatus статус верификации
type VerificationStatus string

const (
	VerificationStatusPending   VerificationStatus = "pending"
	VerificationStatusInProgress VerificationStatus = "in_progress"
	VerificationStatusVerified   VerificationStatus = "verified"
	VerificationStatusRejected   VerificationStatus = "rejected"
	VerificationStatusCancelled  VerificationStatus = "cancelled"
)

// VerificationPriority приоритет верификации
type VerificationPriority string

const (
	VerificationPriorityLow    VerificationPriority = "low"
	VerificationPriorityNormal VerificationPriority = "normal"
	VerificationPriorityHigh   VerificationPriority = "high"
	VerificationPriorityUrgent VerificationPriority = "urgent"
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

// VerificationExpert представляет эксперта по верификации
type VerificationExpert struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Specialties  []string  `json:"specialties" db:"specialties"`
	Rating       float64   `json:"rating" db:"rating"`
	VerifiedCount int      `json:"verified_count" db:"verified_count"`
	Accuracy     float64   `json:"accuracy" db:"accuracy"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// VerificationDetail детальная информация о верификации
type VerificationDetail struct {
	VerificationRequest
	Item        ItemInfo     `json:"item"`
	User        UserInfo     `json:"user"`
	Expert      *ExpertInfo  `json:"expert,omitempty"`
	TimeLeft    int64        `json:"time_left"` // секунды до дедлайна
	CanVerify   bool         `json:"can_verify"`
	CanCancel   bool         `json:"can_cancel"`
}

// ItemInfo краткая информация о товаре
type ItemInfo struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	MainPhoto   string    `json:"main_photo"`
	Category    string    `json:"category"`
	Brand       string    `json:"brand"`
	Model       string    `json:"model"`
	EstimatedValue float64 `json:"estimated_value"`
}

// UserInfo краткая информация о пользователе
type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Rating   float64   `json:"rating"`
	Verified bool      `json:"verified"`
}

// ExpertInfo информация об эксперте
type ExpertInfo struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Rating       float64   `json:"rating"`
	VerifiedCount int      `json:"verified_count"`
	Accuracy     float64   `json:"accuracy"`
	Specialties  []string  `json:"specialties"`
}

// SubmitVerificationRequest запрос на отправку на верификацию
type SubmitVerificationRequest struct {
	Photos        []string `json:"photos" validate:"required,min=3,max=20"`
	Description   string   `json:"description" validate:"required,min=10,max=1000"`
	PurchaseProof []string `json:"purchase_proof" validate:"omitempty,max=10"`
	SerialNumber  string   `json:"serial_number" validate:"omitempty"`
	Receipt       string   `json:"receipt" validate:"omitempty"`
	Priority      VerificationPriority `json:"priority" validate:"omitempty"`
}

// VerifyItemRequest запрос на верификацию товара
type VerifyItemRequest struct {
	AuthenticityGrade AuthenticityGrade `json:"authenticity_grade" validate:"required"`
	VerificationNotes  string            `json:"verification_notes" validate:"required,min=10,max=1000"`
}

// VerificationSearchRequest параметры поиска верификаций
type VerificationSearchRequest struct {
	Status   VerificationStatus  `json:"status"`
	Priority VerificationPriority `json:"priority"`
	Grade    AuthenticityGrade   `json:"grade"`
	ExpertID *uuid.UUID         `json:"expert_id"`
	SortBy   string             `json:"sort_by"` // date, priority, grade
	SortOrder string            `json:"sort_order"` // asc, desc
	Limit    int                `json:"limit" validate:"min=1,max=100"`
	Offset   int                `json:"offset" validate:"min=0"`
}

// VerificationStats статистика верификации
type VerificationStats struct {
	TotalRequests     int     `json:"total_requests"`
	PendingRequests   int     `json:"pending_requests"`
	VerifiedItems     int     `json:"verified_items"`
	RejectedItems     int     `json:"rejected_items"`
	AverageTime       float64 `json:"average_time"` // часы
	AccuracyRate      float64 `json:"accuracy_rate"`
	GradeDistribution map[AuthenticityGrade]int `json:"grade_distribution"`
}

// VerificationRepository интерфейс для работы с верификацией
type VerificationRepository interface {
	Create(request *VerificationRequest) error
	GetByID(id uuid.UUID) (*VerificationRequest, error)
	GetByItemID(itemID uuid.UUID) (*VerificationRequest, error)
	Update(request *VerificationRequest) error
	Delete(id uuid.UUID) error
	Search(req *VerificationSearchRequest) ([]*VerificationDetail, int, error)
	GetDetail(id uuid.UUID) (*VerificationDetail, error)
	GetPending(limit, offset int) ([]*VerificationDetail, error)
	GetByStatus(status VerificationStatus, limit, offset int) ([]*VerificationDetail, error)
	GetByExpert(expertID uuid.UUID, limit, offset int) ([]*VerificationDetail, error)
	GetStats() (*VerificationStats, error)
	
	// Эксперты
	CreateExpert(expert *VerificationExpert) error
	GetExpertByUserID(userID uuid.UUID) (*VerificationExpert, error)
	GetExperts(limit, offset int) ([]*VerificationExpert, error)
	UpdateExpert(expert *VerificationExpert) error
	GetExpertStats(expertID uuid.UUID) (*VerificationStats, error)
}

// VerificationUseCase интерфейс бизнес-логики верификации
type VerificationUseCase interface {
	SubmitVerification(itemID, userID uuid.UUID, req *SubmitVerificationRequest) (*VerificationRequest, error)
	GetVerification(id uuid.UUID) (*VerificationDetail, error)
	GetUserVerifications(userID uuid.UUID, limit, offset int) ([]*VerificationRequest, error)
	CancelVerification(id, userID uuid.UUID) error
	SearchVerifications(req *VerificationSearchRequest) ([]*VerificationDetail, int, error)
	GetPendingVerifications(limit, offset int) ([]*VerificationDetail, error)
	GetVerificationStats() (*VerificationStats, error)
	
	// Экспертная верификация
	AssignExpert(verificationID, expertID uuid.UUID) error
	VerifyItem(verificationID, expertID uuid.UUID, req *VerifyItemRequest) error
	GetExpertVerifications(expertID uuid.UUID, limit, offset int) ([]*VerificationDetail, error)
	GetExpertStats(expertID uuid.UUID) (*VerificationStats, error)
	
	// Управление экспертами
	RegisterExpert(userID uuid.UUID, specialties []string) (*VerificationExpert, error)
	UpdateExpertSpecialties(expertID uuid.UUID, specialties []string) error
	DeactivateExpert(expertID uuid.UUID) error
	GetExperts(limit, offset int) ([]*VerificationExpert, error)
	
	// Автоматизация
	ProcessExpiredVerifications() error
	AssignPendingVerifications() error
	UpdateExpertRatings() error
}
