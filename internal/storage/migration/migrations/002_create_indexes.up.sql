CREATE INDEX IF NOT EXISTS idx_allow_comms ON posts (allow_comms);
CREATE INDEX IF NOT EXISTS idx_post_id ON comments (post_id);
CREATE INDEX IF NOT EXISTS idx_parent_id ON comments (parent_id);
CREATE INDEX IF NOT EXISTS idx_created_at_posts ON posts (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_created_at_comments ON comments (created_at DESC);
