-- Add updated_at column to messages table if not exists
ALTER TABLE messages ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();

-- Create index for message search
CREATE INDEX IF NOT EXISTS idx_messages_content ON messages USING gin(to_tsvector('english', content));

-- Create index for read status
CREATE INDEX IF NOT EXISTS idx_message_read_status_user_id ON message_read_status(user_id);
CREATE INDEX IF NOT EXISTS idx_message_read_status_message_id ON message_read_status(message_id);
