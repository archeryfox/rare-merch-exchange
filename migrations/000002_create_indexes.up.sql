-- +goose Up
-- +goose StatementBegin

-- Индексы для пользователей
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_rating ON users(rating);
CREATE INDEX idx_users_verified ON users(verified);
CREATE INDEX idx_users_kyc_status ON users(kyc_status);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Индексы для товаров
CREATE INDEX idx_items_seller_id ON items(seller_id);
CREATE INDEX idx_items_category ON items(category);
CREATE INDEX idx_items_status ON items(status);
CREATE INDEX idx_items_verification_status ON items(verification_status);
CREATE INDEX idx_items_estimated_value ON items(estimated_value);
CREATE INDEX idx_items_created_at ON items(created_at);
CREATE INDEX idx_items_brand_model ON items(brand, model);

-- Индексы для аукционов
CREATE INDEX idx_auctions_item_id ON auctions(item_id);
CREATE INDEX idx_auctions_seller_id ON auctions(seller_id);
CREATE INDEX idx_auctions_status ON auctions(status);
CREATE INDEX idx_auctions_end_time ON auctions(end_time);
CREATE INDEX idx_auctions_current_price ON auctions(current_price);
CREATE INDEX idx_auctions_status_end_time ON auctions(status, end_time);
CREATE INDEX idx_auctions_active ON auctions(status, end_time) WHERE status = 'active';

-- Индексы для ставок
CREATE INDEX idx_bids_auction_id ON bids(auction_id);
CREATE INDEX idx_bids_user_id ON bids(user_id);
CREATE INDEX idx_bids_timestamp ON bids(timestamp);
CREATE INDEX idx_bids_auction_timestamp ON bids(auction_id, timestamp);
CREATE INDEX idx_bids_amount ON bids(amount);

-- Индексы для автоставок
CREATE INDEX idx_auto_bids_auction_id ON auto_bids(auction_id);
CREATE INDEX idx_auto_bids_user_id ON auto_bids(user_id);
CREATE INDEX idx_auto_bids_active ON auto_bids(auction_id, is_active) WHERE is_active = true;

-- Индексы для лотерей
CREATE INDEX idx_lotteries_item_id ON lotteries(item_id);
CREATE INDEX idx_lotteries_seller_id ON lotteries(seller_id);
CREATE INDEX idx_lotteries_status ON lotteries(status);
CREATE INDEX idx_lotteries_end_time ON lotteries(end_time);
CREATE INDEX idx_lotteries_draw_time ON lotteries(draw_time);
CREATE INDEX idx_lotteries_status_end_time ON lotteries(status, end_time);
CREATE INDEX idx_lotteries_active ON lotteries(status, end_time) WHERE status = 'active';

-- Индексы для билетов лотереи
CREATE INDEX idx_lottery_tickets_lottery_id ON lottery_tickets(lottery_id);
CREATE INDEX idx_lottery_tickets_user_id ON lottery_tickets(user_id);
CREATE INDEX idx_lottery_tickets_ticket_number ON lottery_tickets(ticket_number);
CREATE INDEX idx_lottery_tickets_lottery_user ON lottery_tickets(lottery_id, user_id);

-- Индексы для конкурсов
CREATE INDEX idx_contests_type ON contests(type);
CREATE INDEX idx_contests_status ON contests(status);
CREATE INDEX idx_contests_end_time ON contests(end_time);
CREATE INDEX idx_contests_voting_end_time ON contests(voting_end_time);
CREATE INDEX idx_contests_status_end_time ON contests(status, end_time);
CREATE INDEX idx_contests_active ON contests(status, end_time) WHERE status = 'active';

-- Индексы для участия в конкурсах
CREATE INDEX idx_contest_submissions_contest_id ON contest_submissions(contest_id);
CREATE INDEX idx_contest_submissions_user_id ON contest_submissions(user_id);
CREATE INDEX idx_contest_submissions_votes ON contest_submissions(votes);
CREATE INDEX idx_contest_submissions_ranking ON contest_submissions(ranking);
CREATE INDEX idx_contest_submissions_contest_user ON contest_submissions(contest_id, user_id);

-- Индексы для голосов в конкурсах
CREATE INDEX idx_contest_votes_contest_id ON contest_votes(contest_id);
CREATE INDEX idx_contest_votes_submission_id ON contest_votes(submission_id);
CREATE INDEX idx_contest_votes_user_id ON contest_votes(user_id);
CREATE INDEX idx_contest_votes_contest_user ON contest_votes(contest_id, user_id);

-- Индексы для транзакций
CREATE INDEX idx_transactions_buyer_id ON transactions(buyer_id);
CREATE INDEX idx_transactions_seller_id ON transactions(seller_id);
CREATE INDEX idx_transactions_item_id ON transactions(item_id);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_escrow_status ON transactions(escrow_status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_transactions_type_status ON transactions(type, status);

-- Индексы для споров
CREATE INDEX idx_disputes_transaction_id ON disputes(transaction_id);
CREATE INDEX idx_disputes_complainant_id ON disputes(complainant_id);
CREATE INDEX idx_disputes_status ON disputes(status);
CREATE INDEX idx_disputes_created_at ON disputes(created_at);

-- Индексы для верификации
CREATE INDEX idx_verification_requests_item_id ON verification_requests(item_id);
CREATE INDEX idx_verification_requests_user_id ON verification_requests(user_id);
CREATE INDEX idx_verification_requests_status ON verification_requests(status);
CREATE INDEX idx_verification_requests_priority ON verification_requests(priority);
CREATE INDEX idx_verification_requests_verified_by ON verification_requests(verified_by);
CREATE INDEX idx_verification_requests_created_at ON verification_requests(created_at);

-- Индексы для экспертов
CREATE INDEX idx_verification_experts_user_id ON verification_experts(user_id);
CREATE INDEX idx_verification_experts_rating ON verification_experts(rating);
CREATE INDEX idx_verification_experts_active ON verification_experts(is_active) WHERE is_active = true;

-- Индексы для уведомлений
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_read_at ON notifications(read_at);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, read_at) WHERE read_at IS NULL;

-- Индексы для антифрод алертов
CREATE INDEX idx_fraud_alerts_user_id ON fraud_alerts(user_id);
CREATE INDEX idx_fraud_alerts_type ON fraud_alerts(type);
CREATE INDEX idx_fraud_alerts_severity ON fraud_alerts(severity);
CREATE INDEX idx_fraud_alerts_reviewed ON fraud_alerts(reviewed);
CREATE INDEX idx_fraud_alerts_created_at ON fraud_alerts(created_at);

-- Составные индексы для оптимизации запросов
CREATE INDEX idx_auctions_active_ending ON auctions(status, end_time) WHERE status = 'active';
CREATE INDEX idx_lotteries_active_ending ON lotteries(status, end_time) WHERE status = 'active';
CREATE INDEX idx_contests_active_ending ON contests(status, end_time) WHERE status = 'active';
CREATE INDEX idx_transactions_pending_escrow ON transactions(status, escrow_status) WHERE status = 'pending' AND escrow_status = 'held';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Удаление индексов
DROP INDEX IF EXISTS idx_transactions_pending_escrow;
DROP INDEX IF EXISTS idx_contests_active_ending;
DROP INDEX IF EXISTS idx_lotteries_active_ending;
DROP INDEX IF EXISTS idx_auctions_active_ending;

DROP INDEX IF EXISTS idx_fraud_alerts_created_at;
DROP INDEX IF EXISTS idx_fraud_alerts_reviewed;
DROP INDEX IF EXISTS idx_fraud_alerts_severity;
DROP INDEX IF EXISTS idx_fraud_alerts_type;
DROP INDEX IF EXISTS idx_fraud_alerts_user_id;

DROP INDEX IF EXISTS idx_notifications_user_unread;
DROP INDEX IF EXISTS idx_notifications_created_at;
DROP INDEX IF EXISTS idx_notifications_read_at;
DROP INDEX IF EXISTS idx_notifications_type;
DROP INDEX IF EXISTS idx_notifications_user_id;

DROP INDEX IF EXISTS idx_verification_experts_active;
DROP INDEX IF EXISTS idx_verification_experts_rating;
DROP INDEX IF EXISTS idx_verification_experts_user_id;

DROP INDEX IF EXISTS idx_verification_requests_created_at;
DROP INDEX IF EXISTS idx_verification_requests_verified_by;
DROP INDEX IF EXISTS idx_verification_requests_priority;
DROP INDEX IF EXISTS idx_verification_requests_status;
DROP INDEX IF EXISTS idx_verification_requests_user_id;
DROP INDEX IF EXISTS idx_verification_requests_item_id;

DROP INDEX IF EXISTS idx_disputes_created_at;
DROP INDEX IF EXISTS idx_disputes_status;
DROP INDEX IF EXISTS idx_disputes_complainant_id;
DROP INDEX IF EXISTS idx_disputes_transaction_id;

DROP INDEX IF EXISTS idx_transactions_type_status;
DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_transactions_escrow_status;
DROP INDEX IF EXISTS idx_transactions_status;
DROP INDEX IF EXISTS idx_transactions_item_id;
DROP INDEX IF EXISTS idx_transactions_seller_id;
DROP INDEX IF EXISTS idx_transactions_buyer_id;

DROP INDEX IF EXISTS idx_contest_votes_contest_user;
DROP INDEX IF EXISTS idx_contest_votes_user_id;
DROP INDEX IF EXISTS idx_contest_votes_submission_id;
DROP INDEX IF EXISTS idx_contest_votes_contest_id;

DROP INDEX IF EXISTS idx_contest_submissions_contest_user;
DROP INDEX IF EXISTS idx_contest_submissions_ranking;
DROP INDEX IF EXISTS idx_contest_submissions_votes;
DROP INDEX IF EXISTS idx_contest_submissions_user_id;
DROP INDEX IF EXISTS idx_contest_submissions_contest_id;

DROP INDEX IF EXISTS idx_contests_active;
DROP INDEX IF EXISTS idx_contests_status_end_time;
DROP INDEX IF EXISTS idx_contests_voting_end_time;
DROP INDEX IF EXISTS idx_contests_end_time;
DROP INDEX IF EXISTS idx_contests_status;
DROP INDEX IF EXISTS idx_contests_type;

DROP INDEX IF EXISTS idx_lottery_tickets_lottery_user;
DROP INDEX IF EXISTS idx_lottery_tickets_ticket_number;
DROP INDEX IF EXISTS idx_lottery_tickets_user_id;
DROP INDEX IF EXISTS idx_lottery_tickets_lottery_id;

DROP INDEX IF EXISTS idx_lotteries_active;
DROP INDEX IF EXISTS idx_lotteries_status_end_time;
DROP INDEX IF EXISTS idx_lotteries_draw_time;
DROP INDEX IF EXISTS idx_lotteries_end_time;
DROP INDEX IF EXISTS idx_lotteries_status;
DROP INDEX IF EXISTS idx_lotteries_seller_id;
DROP INDEX IF EXISTS idx_lotteries_item_id;

DROP INDEX IF EXISTS idx_auto_bids_active;
DROP INDEX IF EXISTS idx_auto_bids_user_id;
DROP INDEX IF EXISTS idx_auto_bids_auction_id;

DROP INDEX IF EXISTS idx_bids_amount;
DROP INDEX IF EXISTS idx_bids_auction_timestamp;
DROP INDEX IF EXISTS idx_bids_timestamp;
DROP INDEX IF EXISTS idx_bids_user_id;
DROP INDEX IF EXISTS idx_bids_auction_id;

DROP INDEX IF EXISTS idx_auctions_active;
DROP INDEX IF EXISTS idx_auctions_status_end_time;
DROP INDEX IF EXISTS idx_auctions_current_price;
DROP INDEX IF EXISTS idx_auctions_end_time;
DROP INDEX IF EXISTS idx_auctions_status;
DROP INDEX IF EXISTS idx_auctions_seller_id;
DROP INDEX IF EXISTS idx_auctions_item_id;

DROP INDEX IF EXISTS idx_items_brand_model;
DROP INDEX IF EXISTS idx_items_created_at;
DROP INDEX IF EXISTS idx_items_estimated_value;
DROP INDEX IF EXISTS idx_items_verification_status;
DROP INDEX IF EXISTS idx_items_status;
DROP INDEX IF EXISTS idx_items_category;
DROP INDEX IF EXISTS idx_items_seller_id;

DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_kyc_status;
DROP INDEX IF EXISTS idx_users_verified;
DROP INDEX IF EXISTS idx_users_rating;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;

-- +goose StatementEnd
