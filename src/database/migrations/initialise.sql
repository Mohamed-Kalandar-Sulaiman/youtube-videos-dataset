
CREATE TABLE IF NOT EXISTS videos (
    id VARCHAR(255) PRIMARY KEY, 
    title VARCHAR(255) NOT NULL,
    description TEXT,
    published_date TIMESTAMPTZ NOT NULL,
    thumbnail_url VARCHAR(255),
    channel_name VARCHAR(255) NOT NULL,
    channel_id VARCHAR(255) NOT NULL
);


CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Create an index for published_date 
CREATE INDEX IF NOT EXISTS idx_videos_published_date ON videos(published_date DESC);




-- Create a GIN index for title 
CREATE INDEX IF NOT EXISTS idx_videos_title_trgm
ON videos USING GIN (title gin_trgm_ops);

-- Create a GIN index for channel_name 
CREATE INDEX IF NOT EXISTS idx_videos_channel_name_trgm
ON videos USING GIN (channel_name gin_trgm_ops);