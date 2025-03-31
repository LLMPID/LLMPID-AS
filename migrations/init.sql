CREATE TABLE IF NOT EXISTS classification_logs (
    id BIGSERIAL PRIMARY KEY,
    request_text TEXT NOT NULL,
    result VARCHAR(32),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE, 
    password_hash TEXT NOT NULL,
    role VARCHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
    id BIGSERIAL PRIMARY KEY,
    sub TEXT NOT NULL,
    session_id VARCHAR(64) NOT NULL UNIQUE,
    expires_at BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL
);

INSERT INTO users (username, password_hash, role) VALUES (
    'admin',
    '$argon2id$v=19$m=65536,t=1,p=10$ff+Is1j1GoKrkiiYvLLyGQ$xKmunDT6s3/xoa2+ajvex9tFDNdDLN5aSOFgVzqNMWo',
    'user'
);