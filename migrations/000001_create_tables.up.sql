-- +goose Up
-- +goose StatementBegin

-- Создание расширений
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Пользователи
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    avatar TEXT,
    
    -- Статусы
    verified BOOLEAN DEFAULT FALSE,
    kyc_status VARCHAR(20) DEFAULT 'none' CHECK (kyc_status IN ('none', 'pending', 'verified', 'rejected')),
    active BOOLEAN DEFAULT TRUE,
    
    -- Рейтинг и достижения
    rating DECIMAL(3,2) DEFAULT 0.00 CHECK (rating >= 0 AND rating <= 5),
    rank VARCHAR(20) DEFAULT 'newbie' CHECK (rank IN ('newbie', 'collector', 'expert', 'legend')),
    points INTEGER DEFAULT 0,
    
    -- Настройки уведомлений
    email_notifications BOOLEAN DEFAULT TRUE,
    push_notifications BOOLEAN DEFAULT TRUE,
    fcm_token TEXT,
    
    -- Метаданные
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE
);

-- Товары
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    seller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(50) NOT NULL CHECK (category IN ('clothing', 'shoes', 'accessories', 'toys', 'posters', 'collectibles', 'electronics', 'books', 'other')),
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER NOT NULL CHECK (year >= 1900 AND year <= 2030),
    condition VARCHAR(20) NOT NULL CHECK (condition IN ('new', 'like_new', 'good', 'fair', 'poor')),
    size VARCHAR(20),
    color VARCHAR(50),
    material VARCHAR(100),
    
    -- Фотографии
    photos TEXT[] DEFAULT '{}',
    main_photo TEXT,
    
    -- Верификация
    verification_status VARCHAR(20) DEFAULT 'not_required' CHECK (verification_status IN ('pending', 'verified', 'rejected', 'not_required')),
    authenticity_grade VARCHAR(1) CHECK (authenticity_grade IN ('A', 'B', 'C', 'D', 'F')),
    verified_at TIMESTAMP WITH TIME ZONE,
    verified_by UUID REFERENCES users(id),
    
    -- Цена и состояние
    estimated_value DECIMAL(12,2) NOT NULL CHECK (estimated_value > 0),
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'sold', 'reserved', 'archived')),
    
    -- Метаданные
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Аукционы
CREATE TABLE auctions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    seller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_price DECIMAL(12,2) NOT NULL CHECK (start_price > 0),
    reserve_price DECIMAL(12,2) CHECK (reserve_price > 0),
    current_price DECIMAL(12,2) NOT NULL DEFAULT 0 CHECK (current_price >= 0),
    buy_now_price DECIMAL(12,2) CHECK (buy_now_price > 0),
    min_increment DECIMAL(12,2) NOT NULL CHECK (min_increment > 0),
    
    -- Временные параметры
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    soft_close_minutes INTEGER DEFAULT 5 CHECK (soft_close_minutes > 0),
    
    -- Статус и результаты
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'ended', 'cancelled', 'sold', 'no_bids', 'reserve_not_met')),
    winner_id UUID REFERENCES users(id),
    final_price DECIMAL(12,2),
    total_bids INTEGER DEFAULT 0,
    
    -- Настройки
    auto_extend BOOLEAN DEFAULT TRUE,
    require_deposit BOOLEAN DEFAULT FALSE,
    deposit_amount DECIMAL(12,2) DEFAULT 0 CHECK (deposit_amount >= 0),
    
    -- Метаданные
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ended_at TIMESTAMP WITH TIME ZONE
);

-- Ставки
CREATE TABLE bids (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    auction_id UUID NOT NULL REFERENCES auctions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    is_auto_bid BOOLEAN DEFAULT FALSE,
    max_amount DECIMAL(12,2) CHECK (max_amount > 0),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Автоставки
CREATE TABLE auto_bids (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    auction_id UUID NOT NULL REFERENCES auctions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    max_amount DECIMAL(12,2) NOT NULL CHECK (max_amount > 0),
    increment DECIMAL(12,2) NOT NULL CHECK (increment > 0),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(auction_id, user_id)
);

-- Лотереи
CREATE TABLE lotteries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    seller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ticket_price DECIMAL(12,2) NOT NULL CHECK (ticket_price > 0),
    max_tickets INTEGER NOT NULL CHECK (max_tickets > 0),
    tickets_sold INTEGER DEFAULT 0 CHECK (tickets_sold >= 0),
    
    -- Временные параметры
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    draw_time TIMESTAMP WITH TIME ZONE NOT NULL,
    
    -- Статус и результаты
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'ended', 'drawn', 'cancelled', 'failed')),
    winner_id UUID REFERENCES users(id),
    winner_ticket INTEGER,
    
    -- Provably Fair
    server_seed TEXT NOT NULL,
    server_hash TEXT NOT NULL,
    result_proof TEXT,
    
    -- Настройки
    type VARCHAR(20) DEFAULT 'equal' CHECK (type IN ('equal', 'weighted', 'instant', 'jackpot')),
    min_tickets INTEGER DEFAULT 1 CHECK (min_tickets > 0),
    max_tickets_per_user INTEGER DEFAULT 10 CHECK (max_tickets_per_user > 0),
    require_kyc BOOLEAN DEFAULT FALSE,
    
    -- Метаданные
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    drawn_at TIMESTAMP WITH TIME ZONE
);

-- Билеты лотереи
CREATE TABLE lottery_tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lottery_id UUID NOT NULL REFERENCES lotteries(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ticket_number INTEGER NOT NULL CHECK (ticket_number > 0),
    client_seed TEXT NOT NULL,
    purchased_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(lottery_id, ticket_number)
);

-- Конкурсы
CREATE TABLE contests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(20) NOT NULL CHECK (type IN ('collection_set', 'photo', 'quiz', 'speed', 'creative')),
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    rules TEXT NOT NULL,
    
    -- Призы
    prize_item_id UUID REFERENCES items(id),
    prize_amount DECIMAL(12,2) CHECK (prize_amount > 0),
    prize_description TEXT NOT NULL,
    
    -- Временные параметры
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    voting_end_time TIMESTAMP WITH TIME ZONE,
    
    -- Статус и результаты
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'voting', 'ended', 'cancelled', 'failed')),
    winner_id UUID REFERENCES users(id),
    winner_rank INTEGER,
    
    -- Настройки
    max_participants INTEGER NOT NULL CHECK (max_participants > 0),
    min_participants INTEGER NOT NULL CHECK (min_participants > 0),
    require_kyc BOOLEAN DEFAULT FALSE,
    allow_voting BOOLEAN DEFAULT TRUE,
    
    -- Статистика
    participants_count INTEGER DEFAULT 0,
    votes_count INTEGER DEFAULT 0,
    
    -- Метаданные
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ended_at TIMESTAMP WITH TIME ZONE
);

-- Участие в конкурсах
CREATE TABLE contest_submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    contest_id UUID NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    media TEXT[] DEFAULT '{}',
    votes INTEGER DEFAULT 0,
    ranking INTEGER,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(contest_id, user_id)
);

-- Голоса в конкурсах
CREATE TABLE contest_votes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    contest_id UUID NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    submission_id UUID NOT NULL REFERENCES contest_submissions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    voted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(contest_id, user_id)
);

-- Транзакции
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(20) NOT NULL CHECK (type IN ('auction', 'lottery', 'contest', 'direct')),
    buyer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    auction_id UUID REFERENCES auctions(id),
    lottery_id UUID REFERENCES lotteries(id),
    contest_id UUID REFERENCES contests(id),
    
    -- Финансовые данные
    amount DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    commission DECIMAL(12,2) DEFAULT 0 CHECK (commission >= 0),
    net_amount DECIMAL(12,2) NOT NULL CHECK (net_amount > 0),
    currency VARCHAR(3) DEFAULT 'USD',
    
    -- Escrow данные
    escrow_amount DECIMAL(12,2) DEFAULT 0 CHECK (escrow_amount >= 0),
    escrow_status VARCHAR(20) DEFAULT 'pending' CHECK (escrow_status IN ('pending', 'held', 'released', 'refunded', 'disputed')),
    escrow_released_at TIMESTAMP WITH TIME ZONE,
    
    -- Доставка
    tracking_number VARCHAR(100),
    shipping_method VARCHAR(50),
    shipping_address TEXT,
    shipped_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    
    -- Статус и результаты
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'shipped', 'delivered', 'completed', 'disputed', 'cancelled', 'refunded')),
    dispute_id UUID,
    notes TEXT,
    
    -- Метаданные
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Споры
CREATE TABLE disputes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    complainant_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason VARCHAR(50) NOT NULL CHECK (reason IN ('item_not_received', 'item_not_as_described', 'item_damaged', 'item_fake', 'seller_not_responding', 'other')),
    description TEXT NOT NULL,
    evidence TEXT[] DEFAULT '{}',
    status VARCHAR(20) DEFAULT 'open' CHECK (status IN ('open', 'under_review', 'resolved', 'closed')),
    resolution TEXT,
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Запросы на верификацию
CREATE TABLE verification_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'verified', 'rejected', 'cancelled')),
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    
    -- Данные для проверки
    photos TEXT[] NOT NULL,
    description TEXT NOT NULL,
    purchase_proof TEXT[] DEFAULT '{}',
    serial_number VARCHAR(100),
    receipt TEXT,
    
    -- Результат верификации
    authenticity_grade VARCHAR(1) CHECK (authenticity_grade IN ('A', 'B', 'C', 'D', 'F')),
    verified_by UUID REFERENCES users(id),
    verification_notes TEXT,
    verified_at TIMESTAMP WITH TIME ZONE,
    
    -- Метаданные
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Эксперты по верификации
CREATE TABLE verification_experts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    specialties TEXT[] NOT NULL,
    rating DECIMAL(3,2) DEFAULT 0.00 CHECK (rating >= 0 AND rating <= 5),
    verified_count INTEGER DEFAULT 0,
    accuracy DECIMAL(5,2) DEFAULT 0.00,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(user_id)
);

-- Уведомления
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    message TEXT NOT NULL,
    data JSONB,
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Антифрод алерты
CREATE TABLE fraud_alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    description TEXT NOT NULL,
    data JSONB,
    reviewed BOOLEAN DEFAULT FALSE,
    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS fraud_alerts;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS verification_experts;
DROP TABLE IF EXISTS verification_requests;
DROP TABLE IF EXISTS disputes;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS contest_votes;
DROP TABLE IF EXISTS contest_submissions;
DROP TABLE IF EXISTS contests;
DROP TABLE IF EXISTS lottery_tickets;
DROP TABLE IF EXISTS lotteries;
DROP TABLE IF EXISTS auto_bids;
DROP TABLE IF EXISTS bids;
DROP TABLE IF EXISTS auctions;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS "pgcrypto";

-- +goose StatementEnd
