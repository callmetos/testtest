-- Safety Sessions
CREATE TABLE IF NOT EXISTS safety_sessions (
    id SERIAL PRIMARY KEY,
    plan_id INTEGER NOT NULL REFERENCES trip_plans(id) ON DELETE CASCADE,
    share_token VARCHAR(255) UNIQUE NOT NULL,
    interval_minutes INTEGER NOT NULL,
    next_due TIMESTAMP NOT NULL,
    active BOOLEAN DEFAULT TRUE NOT NULL,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Heartbeats
CREATE TABLE IF NOT EXISTS heartbeats (
    id SERIAL PRIMARY KEY,
    session_id INTEGER NOT NULL REFERENCES safety_sessions(id) ON DELETE CASCADE,
    due_at TIMESTAMP NOT NULL,
    acked_at TIMESTAMP NULL,
    status VARCHAR(50) DEFAULT 'due' NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Emergency Contacts
CREATE TABLE IF NOT EXISTS emergency_contacts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    email VARCHAR(255) NULL,
    priority INTEGER DEFAULT 1 NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_safety_sessions_plan_id ON safety_sessions(plan_id);
CREATE INDEX IF NOT EXISTS idx_safety_sessions_share_token ON safety_sessions(share_token);
CREATE INDEX IF NOT EXISTS idx_safety_sessions_active ON safety_sessions(active);
CREATE INDEX IF NOT EXISTS idx_heartbeats_session_id ON heartbeats(session_id);
CREATE INDEX IF NOT EXISTS idx_heartbeats_status ON heartbeats(status);
CREATE INDEX IF NOT EXISTS idx_emergency_contacts_user_id ON emergency_contacts(user_id);
CREATE INDEX IF NOT EXISTS idx_emergency_contacts_priority ON emergency_contacts(priority);