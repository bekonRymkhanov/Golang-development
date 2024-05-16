CREATE TABLE IF NOT EXISTS like_comment (
    id bigserial PRIMARY KEY,
    user_id integer,
    episode_id integer,
    comment_text text,
    like_count integer DEFAULT 0,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);