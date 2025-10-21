package user

import (
	"time"

	"github.com/google/uuid"
)

// User представляет пользователя системы
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Username     string    `json:"username" db:"username"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Phone        string    `json:"phone" db:"phone"`
	Avatar       string    `json:"avatar" db:"avatar"`
	
	// Статусы
	Verified     bool      `json:"verified" db:"verified"`
	KYCStatus    KYCStatus `json:"kyc_status" db:"kyc_status"`
	Active       bool      `json:"active" db:"active"`
	
	// Рейтинг и достижения
	Rating       float64   `json:"rating" db:"rating"`
	Rank         UserRank  `json:"rank" db:"rank"`
	Points       int       `json:"points" db:"points"`
	
	// Настройки уведомлений
	EmailNotifications bool `json:"email_notifications" db:"email_notifications"`
	PushNotifications  bool `json:"push_notifications" db:"push_notifications"`
	FCMToken          string `json:"-" db:"fcm_token"`
	
	// Метаданные
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at" db:"last_login_at"`
}

// KYCStatus статус верификации пользователя
type KYCStatus string

const (
	KYCStatusNone     KYCStatus = "none"
	KYCStatusPending  KYCStatus = "pending"
	KYCStatusVerified KYCStatus = "verified"
	KYCStatusRejected KYCStatus = "rejected"
)

// UserRank ранг пользователя
type UserRank string

const (
	UserRankNewbie      UserRank = "newbie"
	UserRankCollector   UserRank = "collector"
	UserRankExpert      UserRank = "expert"
	UserRankLegend      UserRank = "legend"
)

// UserProfile публичный профиль пользователя
type UserProfile struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Avatar    string    `json:"avatar"`
	Verified  bool      `json:"verified"`
	Rating    float64   `json:"rating"`
	Rank      UserRank  `json:"rank"`
	CreatedAt time.Time `json:"created_at"`
}

// UserStats статистика пользователя
type UserStats struct {
	TotalAuctionsWon    int     `json:"total_auctions_won"`
	TotalLotteriesWon   int     `json:"total_lotteries_won"`
	TotalContestsWon    int     `json:"total_contests_won"`
	TotalSpent          float64 `json:"total_spent"`
	TotalEarned         float64 `json:"total_earned"`
	SuccessfulTrades    int     `json:"successful_trades"`
	DisputeRate         float64 `json:"dispute_rate"`
	AverageRating       float64 `json:"average_rating"`
}

// CreateUserRequest запрос на создание пользователя
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Username  string `json:"username" validate:"required,min=3,max=20"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
}

// LoginRequest запрос на вход
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateProfileRequest запрос на обновление профиля
type UpdateProfileRequest struct {
	Username  string `json:"username" validate:"omitempty,min=3,max=20"`
	FirstName string `json:"first_name" validate:"omitempty"`
	LastName  string `json:"last_name" validate:"omitempty"`
	Phone     string `json:"phone" validate:"omitempty"`
	Avatar    string `json:"avatar" validate:"omitempty"`
}

// UpdateNotificationSettingsRequest запрос на обновление настроек уведомлений
type UpdateNotificationSettingsRequest struct {
	EmailNotifications bool   `json:"email_notifications"`
	PushNotifications  bool   `json:"push_notifications"`
	FCMToken          string `json:"fcm_token"`
}

// ChangePasswordRequest запрос на смену пароля
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// AuthResponse ответ с токеном аутентификации
type AuthResponse struct {
	User  UserProfile `json:"user"`
	Token string      `json:"token"`
}

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	UpdateLastLogin(id uuid.UUID) error
	GetProfile(id uuid.UUID) (*UserProfile, error)
	GetStats(id uuid.UUID) (*UserStats, error)
	Search(query string, limit, offset int) ([]*UserProfile, error)
	GetTopRated(limit int) ([]*UserProfile, error)
}

// UserUseCase интерфейс бизнес-логики пользователей
type UserUseCase interface {
	Register(req *CreateUserRequest) (*AuthResponse, error)
	Login(req *LoginRequest) (*AuthResponse, error)
	GetProfile(userID uuid.UUID) (*UserProfile, error)
	UpdateProfile(userID uuid.UUID, req *UpdateProfileRequest) error
	ChangePassword(userID uuid.UUID, req *ChangePasswordRequest) error
	UpdateNotificationSettings(userID uuid.UUID, req *UpdateNotificationSettingsRequest) error
	GetStats(userID uuid.UUID) (*UserStats, error)
	SearchUsers(query string, limit, offset int) ([]*UserProfile, error)
	GetTopRatedUsers(limit int) ([]*UserProfile, error)
	VerifyUser(userID uuid.UUID) error
	UpdateKYCStatus(userID uuid.UUID, status KYCStatus) error
}
