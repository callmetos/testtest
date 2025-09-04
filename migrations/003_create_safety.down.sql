DROP INDEX IF EXISTS idx_emergency_contacts_priority;
DROP INDEX IF EXISTS idx_emergency_contacts_user_id;
DROP INDEX IF EXISTS idx_heartbeats_status;
DROP INDEX IF EXISTS idx_heartbeats_session_id;
DROP INDEX IF EXISTS idx_safety_sessions_active;
DROP INDEX IF EXISTS idx_safety_sessions_share_token;
DROP INDEX IF EXISTS idx_safety_sessions_plan_id;

DROP TABLE IF EXISTS emergency_contacts;
DROP TABLE IF EXISTS heartbeats;
DROP TABLE IF EXISTS safety_sessions;