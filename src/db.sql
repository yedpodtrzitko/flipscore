
CREATE TABLE game_key (
    id character varying(32) PRIMARY KEY,
    game_key character varying(36),
    score_ascending boolean default true,
    score_interval  bigint default 0
);

CREATE TABLE score (
    id BIGSERIAL PRIMARY KEY,
    game_id character varying(32) NOT NULL REFERENCES game_key(id),
    player character varying(32) NOT NULL,
    score bigint NOT NULL CHECK (score >= 0),
    content text,
    created_at timestamp with time zone
);
