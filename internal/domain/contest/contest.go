package contest

import (
	"time"

	"github.com/google/uuid"
)

// Contest представляет конкурс
type Contest struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	Type        ContestType   `json:"type" db:"type"`
	Title       string        `json:"title" db:"title"`
	Description string        `json:"description" db:"description"`
	Rules       string        `json:"rules" db:"rules"`
	
	// Призы
	PrizeItemID *uuid.UUID    `json:"prize_item_id" db:"prize_item_id"`
	PrizeAmount *float64      `json:"prize_amount" db:"prize_amount"`
	PrizeDescription string   `json:"prize_description" db:"prize_description"`
	
	// Временные параметры
	StartTime   time.Time     `json:"start_time" db:"start_time"`
	EndTime     time.Time     `json:"end_time" db:"end_time"`
	VotingEndTime *time.Time  `json:"voting_end_time" db:"voting_end_time"`
	
	// Статус и результаты
	Status      ContestStatus `json:"status" db:"status"`
	WinnerID    *uuid.UUID   `json:"winner_id" db:"winner_id"`
	WinnerRank  *int         `json:"winner_rank" db:"winner_rank"`
	
	// Настройки
	MaxParticipants int       `json:"max_participants" db:"max_participants"`
	MinParticipants int       `json:"min_participants" db:"min_participants"`
	RequireKYC     bool       `json:"require_kyc" db:"require_kyc"`
	AllowVoting    bool       `json:"allow_voting" db:"allow_voting"`
	
	// Статистика
	ParticipantsCount int      `json:"participants_count" db:"participants_count"`
	VotesCount       int       `json:"votes_count" db:"votes_count"`
	
	// Метаданные
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
	EndedAt     *time.Time    `json:"ended_at" db:"ended_at"`
}

// ContestType тип конкурса
type ContestType string

const (
	ContestTypeCollectionSet ContestType = "collection_set" // коллекционные сеты
	ContestTypePhoto         ContestType = "photo"          // фото-конкурсы
	ContestTypeQuiz          ContestType = "quiz"           // квизы
	ContestTypeSpeed         ContestType = "speed"          // speed-конкурсы
	ContestTypeCreative      ContestType = "creative"       // творческие конкурсы
)

// ContestStatus статус конкурса
type ContestStatus string

const (
	ContestStatusDraft     ContestStatus = "draft"
	ContestStatusActive    ContestStatus = "active"
	ContestStatusVoting    ContestStatus = "voting"
	ContestStatusEnded     ContestStatus = "ended"
	ContestStatusCancelled ContestStatus = "cancelled"
	ContestStatusFailed    ContestStatus = "failed" // недостаточно участников
)

// ContestSubmission представляет участие в конкурсе
type ContestSubmission struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ContestID uuid.UUID `json:"contest_id" db:"contest_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	Media     []string  `json:"media" db:"media"`
	Votes     int       `json:"votes" db:"votes"`
	Ranking   int       `json:"ranking" db:"ranking"`
	SubmittedAt time.Time `json:"submitted_at" db:"submitted_at"`
}

// ContestVote представляет голос в конкурсе
type ContestVote struct {
	ID           uuid.UUID `json:"id" db:"id"`
	ContestID    uuid.UUID `json:"contest_id" db:"contest_id"`
	SubmissionID uuid.UUID `json:"submission_id" db:"submission_id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	VotedAt      time.Time `json:"voted_at" db:"voted_at"`
}

// ContestDetail детальная информация о конкурсе
type ContestDetail struct {
	Contest
	PrizeItem      *ItemInfo        `json:"prize_item,omitempty"`
	Submissions    []SubmissionInfo `json:"submissions"`
	UserSubmission *SubmissionInfo  `json:"user_submission,omitempty"`
	TimeLeft       int64            `json:"time_left"` // секунды до окончания
	IsActive       bool             `json:"is_active"`
	CanParticipate bool             `json:"can_participate"`
	CanVote        bool             `json:"can_vote"`
	UserVoted      bool             `json:"user_voted"`
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

// SubmissionInfo информация об участии
type SubmissionInfo struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
	UserRating  float64   `json:"user_rating"`
	Content     string    `json:"content"`
	Media       []string  `json:"media"`
	Votes       int       `json:"votes"`
	Ranking     int       `json:"ranking"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// CreateContestRequest запрос на создание конкурса
type CreateContestRequest struct {
	Type              ContestType `json:"type" validate:"required"`
	Title             string      `json:"title" validate:"required,min=3,max=200"`
	Description       string      `json:"description" validate:"required,min=10,max=2000"`
	Rules             string      `json:"rules" validate:"required,min=10,max=2000"`
	PrizeItemID       *uuid.UUID  `json:"prize_item_id" validate:"omitempty"`
	PrizeAmount       *float64    `json:"prize_amount" validate:"omitempty,min=1"`
	PrizeDescription  string      `json:"prize_description" validate:"required"`
	Duration          int         `json:"duration" validate:"required,min=1,max=168"` // часы
	VotingDuration    int         `json:"voting_duration" validate:"omitempty,min=1,max=72"` // часы
	MaxParticipants   int         `json:"max_participants" validate:"required,min=2,max=1000"`
	MinParticipants   int         `json:"min_participants" validate:"required,min=2"`
	RequireKYC        bool        `json:"require_kyc"`
	AllowVoting       bool        `json:"allow_voting"`
}

// SubmitContestRequest запрос на участие в конкурсе
type SubmitContestRequest struct {
	Content string   `json:"content" validate:"required,min=10,max=2000"`
	Media   []string `json:"media" validate:"omitempty,max=10"`
}

// VoteContestRequest запрос на голосование
type VoteContestRequest struct {
	SubmissionID uuid.UUID `json:"submission_id" validate:"required"`
}

// ContestSearchRequest параметры поиска конкурсов
type ContestSearchRequest struct {
	Query     string        `json:"query"`
	Type      ContestType   `json:"type"`
	Status    ContestStatus `json:"status"`
	EndingSoon bool         `json:"ending_soon"`
	SortBy    string        `json:"sort_by"` // participants, votes, time_left
	SortOrder string        `json:"sort_order"` // asc, desc
	Limit     int           `json:"limit" validate:"min=1,max=100"`
	Offset    int           `json:"offset" validate:"min=0"`
}

// ContestRepository интерфейс для работы с конкурсами
type ContestRepository interface {
	Create(contest *Contest) error
	GetByID(id uuid.UUID) (*Contest, error)
	Update(contest *Contest) error
	Delete(id uuid.UUID) error
	Search(req *ContestSearchRequest) ([]*ContestDetail, int, error)
	GetActive(limit, offset int) ([]*ContestDetail, error)
	GetEndingSoon(minutes int, limit int) ([]*ContestDetail, error)
	GetEnded(limit, offset int) ([]*ContestDetail, error)
	GetByStatus(status ContestStatus, limit, offset int) ([]*ContestDetail, error)
	GetDetail(id uuid.UUID) (*ContestDetail, error)
	
	// Участие
	CreateSubmission(submission *ContestSubmission) error
	GetSubmissions(contestID uuid.UUID, limit, offset int) ([]*ContestSubmission, error)
	GetUserSubmission(contestID, userID uuid.UUID) (*ContestSubmission, error)
	UpdateSubmission(submission *ContestSubmission) error
	GetSubmissionCount(contestID uuid.UUID) (int, error)
	
	// Голосование
	CreateVote(vote *ContestVote) error
	GetVotes(contestID uuid.UUID, limit, offset int) ([]*ContestVote, error)
	GetUserVote(contestID, userID uuid.UUID) (*ContestVote, error)
	GetVoteCount(contestID uuid.UUID) (int, error)
	GetSubmissionVotes(submissionID uuid.UUID) (int, error)
	
	// Результаты
	UpdateWinner(contestID uuid.UUID, winnerID uuid.UUID, winnerRank int) error
	GetEndCandidates(minutes int) ([]*Contest, error)
}

// ContestUseCase интерфейс бизнес-логики конкурсов
type ContestUseCase interface {
	CreateContest(userID uuid.UUID, req *CreateContestRequest) (*Contest, error)
	GetContest(id uuid.UUID) (*ContestDetail, error)
	UpdateContest(contestID, userID uuid.UUID, req *CreateContestRequest) error
	CancelContest(contestID, userID uuid.UUID) error
	SearchContests(req *ContestSearchRequest) ([]*ContestDetail, int, error)
	GetActiveContests(limit, offset int) ([]*ContestDetail, error)
	GetEndingSoonContests(minutes int, limit int) ([]*ContestDetail, error)
	GetUserContests(userID uuid.UUID, limit, offset int) ([]*Contest, error)
	
	// Участие
	SubmitContest(contestID, userID uuid.UUID, req *SubmitContestRequest) (*ContestSubmission, error)
	GetUserSubmission(contestID, userID uuid.UUID) (*ContestSubmission, error)
	GetSubmissions(contestID uuid.UUID, limit, offset int) ([]*ContestSubmission, error)
	
	// Голосование
	VoteContest(contestID, userID uuid.UUID, req *VoteContestRequest) error
	GetUserVote(contestID, userID uuid.UUID) (*ContestVote, error)
	GetVotes(contestID uuid.UUID, limit, offset int) ([]*ContestVote, error)
	
	// Завершение конкурсов
	EndContest(contestID uuid.UUID) error
	ProcessEndedContests() error
	CalculateRankings(contestID uuid.UUID) error
}
