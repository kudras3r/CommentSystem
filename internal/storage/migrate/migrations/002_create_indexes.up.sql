CREATE INDEX idx_allow_comms ON posts (allow_comms);
CREATE INDEX idx_post_id ON comments (post_id);
CREATE INDEX idx_parent_id ON comments (parent_id);
CREATE INDEX idx_created_at_posts ON posts (created_at DESC);
CREATE INDEX idx_created_at_comments ON comments (created_at DESC);
