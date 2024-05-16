create table if not exists episodes(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    year integer NOT NULL,
    runtime integer NOT NULL,
    characters text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
    );

create table if not exists characters(
    id bigserial PRIMARY KEY,
    name text,
    age integer NOT NULL ,
    version integer NOT NULL DEFAULT 1
);
CREATE TABLE IF NOT EXISTS episode_characters (
    id BIGSERIAL PRIMARY KEY,
    episode_id BIGINT REFERENCES episodes(id),
    character_id BIGINT REFERENCES characters(id)
);