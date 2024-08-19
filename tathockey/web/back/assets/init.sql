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
    nickname varchar(128) unique,
    password bytea NOT NULL
);

create table if not exists video
(
    id          serial primary key,
    creator_id int not null references users(id),
    name        text,
    hash        varchar(256) unique not null,
    created_at  TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description text
)

