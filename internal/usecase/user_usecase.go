package usecase

import (
	"errors"
	"rare-merch-exchange/internal/domain/user"
	"rare-merch-exchange/internal/pkg/config"
	"rare-merch-exchange/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userUseCase struct {
	repo   user.UserRepository
	config *config.Config
	logger *zap.Logger
}

// NewUserUseCase создаёт новый use case для пользователей
func NewUserUseCase(repo user.UserRepository, cfg *config.Config, logger *zap.Logger) user.UserUseCase {
	return &userUseCase{
		repo:   repo,
		config: cfg,
		logger: logger,
	}
}

func (uc *userUseCase) Register(req *user.CreateUserRequest) (*user.AuthResponse, error) {
	// Проверяем, существует ли пользователь с таким email
	existingUser, err := uc.repo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	// Проверяем, существует ли пользователь с таким username
	existingUser, err = uc.repo.GetByUsername(req.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("пользователь с таким именем уже существует")
	}

	// Хешируем пароль
	hashedPassword, err := repository.HashPassword(req.Password)
	if err != nil {
		uc.logger.Error("Failed to hash password", zap.Error(err))
		return nil, errors.New("ошибка при создании пользователя")
	}

	// Создаём пользователя
	newUser := &user.User{
		ID:                uuid.New(),
		Email:             req.Email,
		PasswordHash:      hashedPassword,
		Username:          req.Username,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Phone:             req.Phone,
		Verified:          false,
		KYCStatus:         user.KYCStatusNone,
		Active:            true,
		Rating:            0.0,
		Rank:              user.UserRankNewbie,
		Points:            0,
		EmailNotifications: true,
		PushNotifications:  true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Сохраняем в базу данных
	if err := uc.repo.Create(newUser); err != nil {
		uc.logger.Error("Failed to create user", zap.Error(err))
		return nil, errors.New("ошибка при создании пользователя")
	}

	// Генерируем JWT токен
	token, err := uc.generateJWT(newUser.ID)
	if err != nil {
		uc.logger.Error("Failed to generate JWT", zap.Error(err))
		return nil, errors.New("ошибка при создании токена")
	}

	// Создаём профиль для ответа
	profile := &user.UserProfile{
		ID:        newUser.ID,
		Username:  newUser.Username,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Avatar:    newUser.Avatar,
		Verified:  newUser.Verified,
		Rating:    newUser.Rating,
		Rank:      newUser.Rank,
		CreatedAt: newUser.CreatedAt,
	}

	return &user.AuthResponse{
		User:  *profile,
		Token: token,
	}, nil
}

func (uc *userUseCase) Login(req *user.LoginRequest) (*user.AuthResponse, error) {
	// Получаем пользователя по email
	u, err := uc.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	// Проверяем пароль
	if !repository.CheckPassword(req.Password, u.PasswordHash) {
		return nil, errors.New("неверный email или пароль")
	}

	// Проверяем, активен ли пользователь
	if !u.Active {
		return nil, errors.New("аккаунт заблокирован")
	}

	// Обновляем время последнего входа
	if err := uc.repo.UpdateLastLogin(u.ID); err != nil {
		uc.logger.Error("Failed to update last login", zap.Error(err))
	}

	// Генерируем JWT токен
	token, err := uc.generateJWT(u.ID)
	if err != nil {
		uc.logger.Error("Failed to generate JWT", zap.Error(err))
		return nil, errors.New("ошибка при создании токена")
	}

	// Создаём профиль для ответа
	profile := &user.UserProfile{
		ID:        u.ID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Avatar:    u.Avatar,
		Verified:  u.Verified,
		Rating:    u.Rating,
		Rank:      u.Rank,
		CreatedAt: u.CreatedAt,
	}

	return &user.AuthResponse{
		User:  *profile,
		Token: token,
	}, nil
}

func (uc *userUseCase) GetProfile(userID uuid.UUID) (*user.UserProfile, error) {
	profile, err := uc.repo.GetProfile(userID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (uc *userUseCase) UpdateProfile(userID uuid.UUID, req *user.UpdateProfileRequest) error {
	// Получаем пользователя
	u, err := uc.repo.GetByID(userID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	// Проверяем, не занят ли новый username
	if req.Username != "" && req.Username != u.Username {
		existingUser, err := uc.repo.GetByUsername(req.Username)
		if err == nil && existingUser != nil {
			return errors.New("пользователь с таким именем уже существует")
		}
		u.Username = req.Username
	}

	// Обновляем поля
	if req.FirstName != "" {
		u.FirstName = req.FirstName
	}
	if req.LastName != "" {
		u.LastName = req.LastName
	}
	if req.Phone != "" {
		u.Phone = req.Phone
	}
	if req.Avatar != "" {
		u.Avatar = req.Avatar
	}

	u.UpdatedAt = time.Now()

	// Сохраняем изменения
	if err := uc.repo.Update(u); err != nil {
		uc.logger.Error("Failed to update user profile", zap.Error(err))
		return errors.New("ошибка при обновлении профиля")
	}

	return nil
}

func (uc *userUseCase) ChangePassword(userID uuid.UUID, req *user.ChangePasswordRequest) error {
	// Получаем пользователя
	u, err := uc.repo.GetByID(userID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	// Проверяем текущий пароль
	if !repository.CheckPassword(req.CurrentPassword, u.PasswordHash) {
		return errors.New("неверный текущий пароль")
	}

	// Хешируем новый пароль
	hashedPassword, err := repository.HashPassword(req.NewPassword)
	if err != nil {
		uc.logger.Error("Failed to hash new password", zap.Error(err))
		return errors.New("ошибка при смене пароля")
	}

	// Обновляем пароль
	u.PasswordHash = hashedPassword
	u.UpdatedAt = time.Now()

	if err := uc.repo.Update(u); err != nil {
		uc.logger.Error("Failed to update password", zap.Error(err))
		return errors.New("ошибка при смене пароля")
	}

	return nil
}

func (uc *userUseCase) UpdateNotificationSettings(userID uuid.UUID, req *user.UpdateNotificationSettingsRequest) error {
	// Получаем пользователя
	u, err := uc.repo.GetByID(userID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	// Обновляем настройки
	u.EmailNotifications = req.EmailNotifications
	u.PushNotifications = req.PushNotifications
	u.FCMToken = req.FCMToken
	u.UpdatedAt = time.Now()

	// Сохраняем изменения
	if err := uc.repo.Update(u); err != nil {
		uc.logger.Error("Failed to update notification settings", zap.Error(err))
		return errors.New("ошибка при обновлении настроек уведомлений")
	}

	return nil
}

func (uc *userUseCase) GetStats(userID uuid.UUID) (*user.UserStats, error) {
	stats, err := uc.repo.GetStats(userID)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (uc *userUseCase) SearchUsers(query string, limit, offset int) ([]*user.UserProfile, error) {
	profiles, err := uc.repo.Search(query, limit, offset)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (uc *userUseCase) GetTopRatedUsers(limit int) ([]*user.UserProfile, error) {
	profiles, err := uc.repo.GetTopRated(limit)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (uc *userUseCase) VerifyUser(userID uuid.UUID) error {
	// Получаем пользователя
	u, err := uc.repo.GetByID(userID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	// Верифицируем пользователя
	u.Verified = true
	u.UpdatedAt = time.Now()

	if err := uc.repo.Update(u); err != nil {
		uc.logger.Error("Failed to verify user", zap.Error(err))
		return errors.New("ошибка при верификации пользователя")
	}

	return nil
}

func (uc *userUseCase) UpdateKYCStatus(userID uuid.UUID, status user.KYCStatus) error {
	// Получаем пользователя
	u, err := uc.repo.GetByID(userID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	// Обновляем статус KYC
	u.KYCStatus = status
	u.UpdatedAt = time.Now()

	if err := uc.repo.Update(u); err != nil {
		uc.logger.Error("Failed to update KYC status", zap.Error(err))
		return errors.New("ошибка при обновлении статуса KYC")
	}

	return nil
}

// generateJWT генерирует JWT токен для пользователя
func (uc *userUseCase) generateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(uc.config.GetJWTExpiration()).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.config.JWT.Secret))
}
