CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS users_forum CASCADE;

DROP INDEX IF EXISTS idx_users_nickname;
DROP INDEX IF EXISTS idx_users_email;

DROP INDEX IF EXISTS idx_forum_slug;
DROP INDEX IF EXIST idx_forum_

DROP INDEX IF EXISTS idx_users_forum_nickname;
DROP INDEX IF EXISTS idx_users_forum_user;

CREATE TABLE users
(
    nickname CITEXT NOT NULL UNIQUE PRIMARY KEY,
    fullname TEXT,
    email    CITEXT NOT NULL UNIQUE,
    about    TEXT
);


CREATE INDEX idx_users_nickname ON users (nickname);
CREATE INDEX idx_users_email ON users (email);


CREATE TABLE forums
(
    slug    CITEXT  NOT NULL UNIQUE PRIMARY KEY ,
    title   TEXT    NOT NULL,
    "user"  CITEXT  NOT NULL REFERENCES users (nickname),
    threads INTEGER DEFAULT 0,
    posts   INTEGER DEFAULT 0
);

CREATE INDEX idx_forum_user ON forums ("user");


CREATE TABLE threads
(
    id      SERIAL PRIMARY KEY,
    slug    CITEXT DEFAULT NULL UNIQUE,
    author  CITEXT REFERENCES users (nickname) NOT NULL,
    forum   CITEXT REFERENCES forums (slug)    NOT NULL,
    title   TEXT NOT NULL,
    message TEXT   NOT NULL,
    created TIMESTAMPTZ DEFAULT current_timestamp,
    votes   INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE posts
(
    id       SERIAL PRIMARY KEY,
    forum    CITEXT REFERENCES forums (slug),
    parent   INTEGER DEFAULT 0,
    author   CITEXT REFERENCES users (nickname) NOT NULL,
    created  TIMESTAMPTZ DEFAULT current_timestamp,
    isEdited BOOLEAN     DEFAULT FALSE,
    message  TEXT                                     NOT NULL,
    thread   INTEGER REFERENCES threads (id)     NOT NULL,
    path     INTEGER[]   DEFAULT array []::INT[]
);


CREATE TABLE votes
(
    nickname CITEXT REFERENCES  users (nickname) NOT NULL,
    voice    SMALLINT CHECK ( voice IN (-1, 1) ),
    thread   INTEGER REFERENCES threads (id)     NOT NULL,
    UNIQUE (nickname, thread)
);


CREATE TABLE users_forum
(
    nickname CITEXT REFERENCES users (nickname) NOT NULL,
    slug     CITEXT REFERENCES forums (slug) NOT NULL,
    UNIQUE (nickname, slug)
);

CREATE INDEX idx_users_forum_nickname ON users_forum (nickname);
CREATE INDEX idx_users_forum_slug ON users_forum (slug);