CREATE SCHEMA IF NOT EXISTS logs;

CREATE TABLE logs.audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    action TEXT NOT NULL,
    metadata TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);
