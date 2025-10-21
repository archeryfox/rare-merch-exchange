package repository

import (
	"database/sql"
	"fmt"
	"rare-merch-exchange/internal/domain/user"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	db  *sql.DB
	rdb *redis.Client
}

// NewUserRepository создаёт новый репозиторий пользователей
func NewUserRepository(db *sql.DB, rdb *redis.Client) user.UserRepository {
	return &userRepository{
		db:  db,
		rdb: rdb,
	}
}

func (r *userRepository) Create(u *user.User) error {
	query := `
		INSERT INTO users (
			id, email, password_hash, username, first_name, last_name, phone, avatar,
			verified, kyc_status, active, rating, rank, points,
			email_notifications, push_notifications, fcm_token,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
		)`

	_, err := r.db.Exec(query,
		u.ID, u.Email, u.PasswordHash, u.Username, u.FirstName, u.LastName, u.Phone, u.Avatar,
		u.Verified, u.KYCStatus, u.Active, u.Rating, u.Rank, u.Points,
		u.EmailNotifications, u.PushNotifications, u.FCMToken,
		u.CreatedAt, u.UpdatedAt,
	)

	return err
}

func (r *userRepository) GetByID(id uuid.UUID) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, username, first_name, last_name, phone, avatar,
			   verified, kyc_status, active, rating, rank, points,
			   email_notifications, push_notifications, fcm_token,
			   created_at, updated_at, last_login_at
		FROM users WHERE id = $1`

	var u user.User
	err := r.db.QueryRow(query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Username, &u.FirstName, &u.LastName, &u.Phone, &u.Avatar,
		&u.Verified, &u.KYCStatus, &u.Active, &u.Rating, &u.Rank, &u.Points,
		&u.EmailNotifications, &u.PushNotifications, &u.FCMToken,
		&u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) GetByEmail(email string) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, username, first_name, last_name, phone, avatar,
			   verified, kyc_status, active, rating, rank, points,
			   email_notifications, push_notifications, fcm_token,
			   created_at, updated_at, last_login_at
		FROM users WHERE email = $1`

	var u user.User
	err := r.db.QueryRow(query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Username, &u.FirstName, &u.LastName, &u.Phone, &u.Avatar,
		&u.Verified, &u.KYCStatus, &u.Active, &u.Rating, &u.Rank, &u.Points,
		&u.EmailNotifications, &u.PushNotifications, &u.FCMToken,
		&u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) GetByUsername(username string) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, username, first_name, last_name, phone, avatar,
			   verified, kyc_status, active, rating, rank, points,
			   email_notifications, push_notifications, fcm_token,
			   created_at, updated_at, last_login_at
		FROM users WHERE username = $1`

	var u user.User
	err := r.db.QueryRow(query, username).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Username, &u.FirstName, &u.LastName, &u.Phone, &u.Avatar,
		&u.Verified, &u.KYCStatus, &u.Active, &u.Rating, &u.Rank, &u.Points,
		&u.EmailNotifications, &u.PushNotifications, &u.FCMToken,
		&u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) Update(u *user.User) error {
	query := `
		UPDATE users SET
			email = $2, password_hash = $3, username = $4, first_name = $5, last_name = $6,
			phone = $7, avatar = $8, verified = $9, kyc_status = $10, active = $11,
			rating = $12, rank = $13, points = $14, email_notifications = $15,
			push_notifications = $16, fcm_token = $17, updated_at = $18
		WHERE id = $1`

	_, err := r.db.Exec(query,
		u.ID, u.Email, u.PasswordHash, u.Username, u.FirstName, u.LastName,
		u.Phone, u.Avatar, u.Verified, u.KYCStatus, u.Active,
		u.Rating, u.Rank, u.Points, u.EmailNotifications,
		u.PushNotifications, u.FCMToken, u.UpdatedAt,
	)

	return err
}

func (r *userRepository) UpdateLastLogin(id uuid.UUID) error {
	query := `UPDATE users SET last_login_at = $2 WHERE id = $1`
	_, err := r.db.Exec(query, id, time.Now())
	return err
}

func (r *userRepository) GetProfile(id uuid.UUID) (*user.UserProfile, error) {
	query := `
		SELECT id, username, first_name, last_name, avatar, verified, rating, rank, created_at
		FROM users WHERE id = $1`

	var p user.UserProfile
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Username, &p.FirstName, &p.LastName, &p.Avatar,
		&p.Verified, &p.Rating, &p.Rank, &p.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &p, nil
}

func (r *userRepository) GetStats(id uuid.UUID) (*user.UserStats, error) {
	// Получаем статистику из различных таблиц
	var stats user.UserStats

	// Статистика аукционов
	query := `SELECT COUNT(*) FROM auctions WHERE winner_id = $1`
	err := r.db.QueryRow(query, id).Scan(&stats.TotalAuctionsWon)
	if err != nil {
		return nil, err
	}

	// Статистика лотерей
	query = `SELECT COUNT(*) FROM lotteries WHERE winner_id = $1`
	err = r.db.QueryRow(query, id).Scan(&stats.TotalLotteriesWon)
	if err != nil {
		return nil, err
	}

	// Статистика конкурсов
	query = `SELECT COUNT(*) FROM contests WHERE winner_id = $1`
	err = r.db.QueryRow(query, id).Scan(&stats.TotalContestsWon)
	if err != nil {
		return nil, err
	}

	// Общая потраченная сумма
	query = `SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE buyer_id = $1 AND status = 'completed'`
	err = r.db.QueryRow(query, id).Scan(&stats.TotalSpent)
	if err != nil {
		return nil, err
	}

	// Общая заработанная сумма
	query = `SELECT COALESCE(SUM(net_amount), 0) FROM transactions WHERE seller_id = $1 AND status = 'completed'`
	err = r.db.QueryRow(query, id).Scan(&stats.TotalEarned)
	if err != nil {
		return nil, err
	}

	// Успешные сделки
	query = `SELECT COUNT(*) FROM transactions WHERE (buyer_id = $1 OR seller_id = $1) AND status = 'completed'`
	err = r.db.QueryRow(query, id).Scan(&stats.SuccessfulTrades)
	if err != nil {
		return nil, err
	}

	// Статистика споров
	query = `SELECT COUNT(*) FROM disputes WHERE complainant_id = $1`
	var totalDisputes int
	err = r.db.QueryRow(query, id).Scan(&totalDisputes)
	if err != nil {
		return nil, err
	}

	if stats.SuccessfulTrades > 0 {
		stats.DisputeRate = float64(totalDisputes) / float64(stats.SuccessfulTrades)
	}

	// Средний рейтинг (пока используем рейтинг пользователя)
	query = `SELECT rating FROM users WHERE id = $1`
	err = r.db.QueryRow(query, id).Scan(&stats.AverageRating)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *userRepository) Search(query string, limit, offset int) ([]*user.UserProfile, error) {
	sqlQuery := `
		SELECT id, username, first_name, last_name, avatar, verified, rating, rank, created_at
		FROM users 
		WHERE (username ILIKE $1 OR first_name ILIKE $1 OR last_name ILIKE $1)
		AND active = true
		ORDER BY rating DESC, created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(sqlQuery, "%"+query+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*user.UserProfile
	for rows.Next() {
		var p user.UserProfile
		err := rows.Scan(
			&p.ID, &p.Username, &p.FirstName, &p.LastName, &p.Avatar,
			&p.Verified, &p.Rating, &p.Rank, &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, &p)
	}

	return profiles, nil
}

func (r *userRepository) GetTopRated(limit int) ([]*user.UserProfile, error) {
	query := `
		SELECT id, username, first_name, last_name, avatar, verified, rating, rank, created_at
		FROM users 
		WHERE active = true AND rating > 0
		ORDER BY rating DESC, verified DESC
		LIMIT $1`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*user.UserProfile
	for rows.Next() {
		var p user.UserProfile
		err := rows.Scan(
			&p.ID, &p.Username, &p.FirstName, &p.LastName, &p.Avatar,
			&p.Verified, &p.Rating, &p.Rank, &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, &p)
	}

	return profiles, nil
}

// HashPassword хеширует пароль
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword проверяет пароль
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
