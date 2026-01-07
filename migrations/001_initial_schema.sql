-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create services table
CREATE TABLE IF NOT EXISTS services (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    owner VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    version VARCHAR(50),
    repository VARCHAR(500),
    docs_url VARCHAR(500),
    tags TEXT[],
    created_by VARCHAR(50) NOT NULL,
    updated_by VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    lead VARCHAR(100) NOT NULL,
    members TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create metrics table
CREATE TABLE IF NOT EXISTS metrics (
    id BIGSERIAL PRIMARY KEY,
    service_id BIGINT NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    unit VARCHAR(20),
    labels JSONB,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_services_name ON services(name);
CREATE INDEX IF NOT EXISTS idx_services_type ON services(type);
CREATE INDEX IF NOT EXISTS idx_services_status ON services(status);
CREATE INDEX IF NOT EXISTS idx_services_owner ON services(owner);
CREATE INDEX IF NOT EXISTS idx_services_tags ON services USING GIN(tags);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_metrics_service_id ON metrics(service_id);
CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON metrics(timestamp);

-- Insert demo user (password: admin123)
INSERT INTO users (username, email, password_hash, first_name, last_name, role)
VALUES ('admin', 'admin@devportal.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Admin', 'User', 'admin')
ON CONFLICT (username) DO NOTHING;

-- Insert demo services
INSERT INTO services (name, description, owner, type, status, version, repository, tags, created_by, updated_by)
VALUES 
    ('User Service', 'Manages user authentication and profiles', 'Platform Team', 'microservice', 'active', '2.1.0', 'https://github.com/example/user-service', ARRAY['authentication', 'users'], 'admin', 'admin'),
    ('Payment API', 'Handles payment processing and transactions', 'Finance Team', 'api', 'active', '1.5.0', 'https://github.com/example/payment-api', ARRAY['payments', 'transactions'], 'admin', 'admin'),
    ('UI Component Library', 'Reusable React components for the platform', 'Frontend Team', 'library', 'active', '3.0.2', 'https://github.com/example/ui-library', ARRAY['frontend', 'react'], 'admin', 'admin'),
    ('Legacy Auth Service', 'Old authentication service', 'Platform Team', 'microservice', 'deprecated', '1.0.0', 'https://github.com/example/legacy-auth', ARRAY['authentication', 'legacy'], 'admin', 'admin')
ON CONFLICT DO NOTHING;
