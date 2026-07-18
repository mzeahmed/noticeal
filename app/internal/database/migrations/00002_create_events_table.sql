-- +goose Up
CREATE TABLE events
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    source     TEXT     NOT NULL,
    type       TEXT     NOT NULL,
    status     TEXT     NOT NULL,
    title      TEXT     NOT NULL,
    message    TEXT     NOT NULL,
    data       TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE events;