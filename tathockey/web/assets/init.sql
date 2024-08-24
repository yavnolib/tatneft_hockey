CREATE TABLE sessions
(
    user_id    INTEGER      NOT NULL,
    id         VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_session_id ON sessions (id);


create table if not exists users
(
    id       serial primary key,
    email    varchar(129) unique,
    nickname varchar(128) unique,
    password bytea NOT NULL
);

CREATE TABLE videos
(
    id   serial primary key ,
    name TEXT NOT NULL
);

CREATE TABLE posts
(
    id          serial primary key ,
    title      TEXT        NOT NULL,
    preview    TEXT        NOT NULL,
    video_id   INTEGER REFERENCES videos (id),
    creator_id INTEGER REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE gifs
(
    id         serial primary key,
    path       TEXT    NOT NULL,
    video_id   TEXT    NOT NULL,
    class_name INTEGER NOT NULL,
    post_id    INTEGER REFERENCES posts (id)
);

