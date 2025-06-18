-- Migration: Create audit_logs table with indexes
-- Date: 2024-11-17
-- Description: Add audit log functionality for tracking system changes

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add foreign key constraint to users table
ALTER TABLE audit_logs 
ADD CONSTRAINT fk_audit_logs_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_action ON audit_logs(user_id, action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_date_range ON audit_logs(created_at DESC);

-- Create partial index for recent activities (last 30 days)
CREATE INDEX IF NOT EXISTS idx_audit_logs_recent 
ON audit_logs(created_at DESC) 
WHERE created_at >= NOW() - INTERVAL '30 days';

-- Add comments for documentation
COMMENT ON TABLE audit_logs IS 'System audit log for tracking user actions and changes';
COMMENT ON COLUMN audit_logs.user_id IS 'ID of the user who performed the action';
COMMENT ON COLUMN audit_logs.action IS 'Type of action performed (create, update, delete, etc.)';
COMMENT ON COLUMN audit_logs.entity_type IS 'Type of entity affected (claim, user, etc.)';
COMMENT ON COLUMN audit_logs.entity_id IS 'ID of the specific entity affected';
COMMENT ON COLUMN audit_logs.old_values IS 'JSON representation of old values before change';
COMMENT ON COLUMN audit_logs.new_values IS 'JSON representation of new values after change';
COMMENT ON COLUMN audit_logs.ip_address IS 'IP address of the user who performed the action';
COMMENT ON COLUMN audit_logs.user_agent IS 'Browser/client user agent string';
COMMENT ON COLUMN audit_logs.created_at IS 'Timestamp when the action was performed';