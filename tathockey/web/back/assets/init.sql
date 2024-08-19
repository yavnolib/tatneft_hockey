
CREATE TABLE sessions (
                          user_id INTEGER NOT NULL,
                          id VARCHAR(255) NOT NULL,
                          expires_at TIMESTAMP NOT NULL,
                          data JSONB,
                          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_session_id ON sessions (id);
CREATE INDEX idx_sessions_expires_at ON sessions (expires_at);


create table if not exists users (
    id serial primary key,
    nickname varchar(128) unique,
    password bytea NOT NULL
);

create table if not exists video (
    id serial primary key,
    name text,
    hash varchar(256) unique not null,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description text
)

