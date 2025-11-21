CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE IF NOT EXISTS teams (
    team_name VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    team_name VARCHAR(255) NOT NULL REFERENCES teams(team_name) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id VARCHAR(255) PRIMARY KEY,
    pull_request_name VARCHAR(255) NOT NULL,
    author_id VARCHAR(255) NOT NULL REFERENCES users(user_id),
    status VARCHAR(50) NOT NULL DEFAULT 'OPEN',
    assigned_reviewers TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    merged_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT valid_status CHECK (status IN ('OPEN', 'MERGED'))
);

CREATE INDEX IF NOT EXISTS idx_pr_status ON pull_requests(status);
CREATE INDEX IF NOT EXISTS idx_pr_author ON pull_requests(author_id);
CREATE INDEX IF NOT EXISTS idx_user_team_active ON users(team_name, is_active);
CREATE INDEX IF NOT EXISTS idx_user_active ON users(is_active);

-- данные для тестов
INSERT INTO teams (team_name) VALUES 
('team1'),
('team2'),
('team3')
ON CONFLICT (team_name) DO NOTHING;

INSERT INTO users (user_id, username, team_name, is_active) VALUES 
('u1', 'J', 'team1', true),
('u2', 'Jo', 'team1', true),
('u3', 'Joh', 'team1', true),
('u4', 'John', 'team2', true),
('u5', 'John Doe', 'team2', true)
ON CONFLICT (user_id) DO NOTHING;