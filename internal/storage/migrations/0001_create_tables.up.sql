CREATE TABLE IF NOT EXISTS users
(
    id        uuid PRIMARY KEY,
    firstname varchar(255) NOT NULL,
    lastname  varchar(255) NOT NULL,
    is_admin  boolean      NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions
(
    id           uuid PRIMARY KEY,
    user_id      uuid         NOT NULL,
    current_step integer DEFAULT 1,
    max_progress integer DEFAULT 1,
    repo         varchar(255) NOT NULL,
    started      timestamp,
    submitted    timestamp,
    time_limit   integer      NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);