package config

import (
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	JWT        JWTConfig
	Upload     UploadConfig
	Auction    AuctionConfig
	Lottery    LotteryConfig
	Commission CommissionConfig
	Escrow     EscrowConfig
	Antifraud  AntifraudConfig
	Logging    LoggingConfig
}

type ServerConfig struct {
	Port string
	Host string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type UploadConfig struct {
	MaxSize int64
	Path    string
}

type AuctionConfig struct {
	SoftCloseMinutes     int
	MinIncrementPercent  int
}

type LotteryConfig struct {
	MaxTicketsPerUser int
	MinTickets        int
}

type CommissionConfig struct {
	RateVerified float64
	RateRegular  float64
}

type EscrowConfig struct {
	AutoReleaseDays     int
	DisputePeriodDays   int
}

type AntifraudConfig struct {
	MaxBidsPerMinute              int
	SuspiciousActivityThreshold   int
}

type LoggingConfig struct {
	Level  string
	Format string
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Устанавливаем значения по умолчанию
	setDefaults()

	// Читаем переменные окружения
	viper.AutomaticEnv()

	// Читаем конфигурационный файл (если есть)
	if err := viper.ReadInConfig(); err != nil {
		// Игнорируем ошибку если файл не найден
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "rare_merch_exchange"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_HOST", "localhost") + ":" + getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
			ExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 24),
		},
		Upload: UploadConfig{
			MaxSize: int64(getEnvAsInt("UPLOAD_MAX_SIZE", 10485760)), // 10MB
			Path:    getEnv("UPLOAD_PATH", "./uploads"),
		},
		Auction: AuctionConfig{
			SoftCloseMinutes:    getEnvAsInt("AUCTION_SOFT_CLOSE_MINUTES", 5),
			MinIncrementPercent: getEnvAsInt("AUCTION_MIN_INCREMENT_PERCENT", 5),
		},
		Lottery: LotteryConfig{
			MaxTicketsPerUser: getEnvAsInt("LOTTERY_MAX_TICKETS_PER_USER", 10),
			MinTickets:        getEnvAsInt("LOTTERY_MIN_TICKETS", 1),
		},
		Commission: CommissionConfig{
			RateVerified: getEnvAsFloat("COMMISSION_RATE_VERIFIED", 0.05),
			RateRegular:  getEnvAsFloat("COMMISSION_RATE_REGULAR", 0.08),
		},
		Escrow: EscrowConfig{
			AutoReleaseDays:   getEnvAsInt("ESCROW_AUTO_RELEASE_DAYS", 7),
			DisputePeriodDays: getEnvAsInt("ESCROW_DISPUTE_PERIOD_DAYS", 3),
		},
		Antifraud: AntifraudConfig{
			MaxBidsPerMinute:            getEnvAsInt("ANTIFRAUD_MAX_BIDS_PER_MINUTE", 10),
			SuspiciousActivityThreshold:  getEnvAsInt("ANTIFRAUD_SUSPICIOUS_ACTIVITY_THRESHOLD", 5),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
}

func setDefaults() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("gin.mode", "debug")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.name", "rare_merch_exchange")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("jwt.secret", "your-super-secret-jwt-key-change-in-production")
	viper.SetDefault("jwt.expire_hours", 24)
	viper.SetDefault("upload.max_size", 10485760)
	viper.SetDefault("upload.path", "./uploads")
	viper.SetDefault("auction.soft_close_minutes", 5)
	viper.SetDefault("auction.min_increment_percent", 5)
	viper.SetDefault("lottery.max_tickets_per_user", 10)
	viper.SetDefault("lottery.min_tickets", 1)
	viper.SetDefault("commission.rate_verified", 0.05)
	viper.SetDefault("commission.rate_regular", 0.08)
	viper.SetDefault("escrow.auto_release_days", 7)
	viper.SetDefault("escrow.dispute_period_days", 3)
	viper.SetDefault("antifraud.max_bids_per_minute", 10)
	viper.SetDefault("antifraud.suspicious_activity_threshold", 5)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func (c *Config) GetJWTExpiration() time.Duration {
	return time.Duration(c.JWT.ExpireHours) * time.Hour
}
