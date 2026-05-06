-- REMINDERS TABLE
CREATE TABLE IF NOT EXISTS reminders (
    id SERIAL PRIMARY KEY,
    workspace_id INTEGER REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    target_id INTEGER,                          -- ID of the task or note
    target_type VARCHAR(50) NOT NULL,           -- 'task', 'note', or 'custom'
    reminder_type VARCHAR(100) NOT NULL,        -- 'task_due', 'note_review', 'custom'
    title TEXT NOT NULL,
    description TEXT,
    due_date TIMESTAMP NOT NULL,
    schedule_type VARCHAR(50) DEFAULT 'once',   -- 'once', 'daily', 'weekly', 'monthly'
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- NOTIFICATIONS TABLE
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    reminder_id INTEGER REFERENCES reminders(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    message TEXT,
    notification_type VARCHAR(50) NOT NULL,    -- 'reminder', 'task_update', 'mention', etc
    is_read BOOLEAN DEFAULT false,
    channels TEXT,                              -- JSON: {"in_app": true, "email": true}
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP
);

-- NOTIFICATION SETTINGS TABLE
CREATE TABLE IF NOT EXISTS notification_settings (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    email_reminders BOOLEAN DEFAULT true,
    in_app_reminders BOOLEAN DEFAULT true,
    email_task_updates BOOLEAN DEFAULT true,
    email_mentions BOOLEAN DEFAULT true,
    quiet_hours_start TIME,                     -- e.g., 22:00
    quiet_hours_end TIME,                       -- e.g., 08:00
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- INDEXES FOR PERFORMANCE
CREATE INDEX IF NOT EXISTS idx_reminders_workspace ON reminders(workspace_id);
CREATE INDEX IF NOT EXISTS idx_reminders_user ON reminders(user_id);
CREATE INDEX IF NOT EXISTS idx_reminders_due_date ON reminders(due_date);
CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);
