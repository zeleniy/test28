CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(32) NOT NULL UNIQUE,
    password_hash CHAR(60) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE users IS 'Application users';
COMMENT ON COLUMN users.id IS 'Primary key';
COMMENT ON COLUMN users.login IS 'User login';
COMMENT ON COLUMN users.password_hash IS 'Password bcrypt hash';
COMMENT ON COLUMN users.created_at IS 'Date created';

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE subscriptions IS 'User''s subscriptions';
COMMENT ON COLUMN subscriptions.user_id IS 'Reference to users.id';
COMMENT ON COLUMN subscriptions.service_name IS 'Subscribed service';
COMMENT ON COLUMN subscriptions.price IS 'Subscription price';
COMMENT ON COLUMN subscriptions.start_date IS 'Subscription start date';
COMMENT ON COLUMN subscriptions.end_date IS 'Subscription end date';
