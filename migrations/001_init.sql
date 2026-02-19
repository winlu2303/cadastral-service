CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS queries (
    id VARCHAR(255) PRIMARY KEY,
    cadastral_number VARCHAR(255) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    result BOOLEAN,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE
);

-- indexes for optimization
CREATE INDEX IF NOT EXISTS idx_queries_cadastral ON queries(cadastral_number);
CREATE INDEX IF NOT EXISTS idx_queries_user_id ON queries(user_id);
CREATE INDEX IF NOT EXISTS idx_queries_created_at ON queries(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_queries_status ON queries(status);

CREATE TABLE IF NOT EXISTS query_history (
    id VARCHAR(255) PRIMARY KEY,
    cadastral_number VARCHAR(255) NOT NULL,
    request_data JSONB NOT NULL,
    response_data JSONB,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_query_history_cadastral ON query_history(cadastral_number);

-- creating admin with password admin123
INSERT INTO users (id, username, password_hash, created_at)
VALUES (
    'admin_001',
    'admin',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeS7.2Y5Z1e8Z5c6W5q5k5n5v5c5n5v5c5n',
    CURRENT_TIMESTAMP
) ON CONFLICT (username) DO NOTHING;
