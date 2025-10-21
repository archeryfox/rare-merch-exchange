package v1

import (
	"net/http"
	"rare-merch-exchange/internal/domain/user"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserHandler struct {
	useCase user.UserUseCase
	logger  *zap.Logger
}

// NewUserHandler создаёт новый handler для пользователей
func NewUserHandler(useCase user.UserUseCase, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// Register регистрация пользователя
// @Summary Регистрация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body user.CreateUserRequest true "Данные для регистрации"
// @Success 201 {object} user.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Неверный формат запроса",
		})
		return
	}

	response, err := h.useCase.Register(&req)
	if err != nil {
		h.logger.Error("Failed to register user", zap.Error(err))
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "registration_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login вход пользователя
// @Summary Вход пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body user.LoginRequest true "Данные для входа"
// @Success 200 {object} user.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Неверный формат запроса",
		})
		return
	}

	response, err := h.useCase.Login(&req)
	if err != nil {
		h.logger.Error("Failed to login user", zap.Error(err))
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "login_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile получение профиля текущего пользователя
// @Summary Получить профиль текущего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} user.UserProfile
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "Пользователь не аутентифицирован",
		})
		return
	}

	profile, err := h.useCase.GetProfile(userID.(uuid.UUID))
	if err != nil {
		h.logger.Error("Failed to get user profile", zap.Error(err))
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "user_not_found",
			Message: "Пользователь не найден",
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile обновление профиля пользователя
// @Summary Обновить профиль пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param request body user.UpdateProfileRequest true "Данные для обновления"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "Пользователь не аутентифицирован",
		})
		return
	}

	var req user.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Неверный формат запроса",
		})
		return
	}

	err := h.useCase.UpdateProfile(userID.(uuid.UUID), &req)
	if err != nil {
		h.logger.Error("Failed to update user profile", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Профиль успешно обновлён",
	})
}

// ChangePassword смена пароля
// @Summary Сменить пароль
// @Tags users
// @Accept json
// @Produce json
// @Param request body user.ChangePasswordRequest true "Данные для смены пароля"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/me/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "Пользователь не аутентифицирован",
		})
		return
	}

	var req user.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Неверный формат запроса",
		})
		return
	}

	err := h.useCase.ChangePassword(userID.(uuid.UUID), &req)
	if err != nil {
		h.logger.Error("Failed to change password", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "password_change_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Пароль успешно изменён",
	})
}

// UpdateNotificationSettings обновление настроек уведомлений
// @Summary Обновить настройки уведомлений
// @Tags users
// @Accept json
// @Produce json
// @Param request body user.UpdateNotificationSettingsRequest true "Настройки уведомлений"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/me/notifications [put]
func (h *UserHandler) UpdateNotificationSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "Пользователь не аутентифицирован",
		})
		return
	}

	var req user.UpdateNotificationSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Неверный формат запроса",
		})
		return
	}

	err := h.useCase.UpdateNotificationSettings(userID.(uuid.UUID), &req)
	if err != nil {
		h.logger.Error("Failed to update notification settings", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Настройки уведомлений обновлены",
	})
}

// GetStats получение статистики пользователя
// @Summary Получить статистику пользователя
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} user.UserStats
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/me/stats [get]
func (h *UserHandler) GetStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "Пользователь не аутентифицирован",
		})
		return
	}

	stats, err := h.useCase.GetStats(userID.(uuid.UUID))
	if err != nil {
		h.logger.Error("Failed to get user stats", zap.Error(err))
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "user_not_found",
			Message: "Пользователь не найден",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// SearchUsers поиск пользователей
// @Summary Поиск пользователей
// @Tags users
// @Accept json
// @Produce json
// @Param query query string true "Поисковый запрос"
// @Param limit query int false "Лимит результатов" default(20)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {array} user.UserProfile
// @Failure 400 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/search [get]
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_query",
			Message: "Поисковый запрос не может быть пустым",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	profiles, err := h.useCase.SearchUsers(query, limit, offset)
	if err != nil {
		h.logger.Error("Failed to search users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "search_failed",
			Message: "Ошибка при поиске пользователей",
		})
		return
	}

	c.JSON(http.StatusOK, profiles)
}

// GetTopRatedUsers получение топ пользователей по рейтингу
// @Summary Получить топ пользователей по рейтингу
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Лимит результатов" default(10)
// @Success 200 {array} user.UserProfile
// @Failure 400 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/top-rated [get]
func (h *UserHandler) GetTopRatedUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	profiles, err := h.useCase.GetTopRatedUsers(limit)
	if err != nil {
		h.logger.Error("Failed to get top rated users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "request_failed",
			Message: "Ошибка при получении списка пользователей",
		})
		return
	}

	c.JSON(http.StatusOK, profiles)
}

// GetUserProfile получение публичного профиля пользователя
// @Summary Получить публичный профиль пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} user.UserProfile
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Неверный ID пользователя",
		})
		return
	}

	profile, err := h.useCase.GetProfile(userID)
	if err != nil {
		h.logger.Error("Failed to get user profile", zap.Error(err))
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "user_not_found",
			Message: "Пользователь не найден",
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// ErrorResponse структура для ошибок
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse структура для успешных ответов
type SuccessResponse struct {
	Message string `json:"message"`
}
