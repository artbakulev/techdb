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

CREATE OR REPLACE FUNCTION new_thread() RETURNS TRIGGER AS
$body$
BEGIN
    UPDATE forums
    SET threads = threads + 1
    WHERE slug = NEW.forum;
    RETURN NEW;
END;
$body$ LANGUAGE plpgsql;

CREATE TRIGGER new_thread_trigger
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE new_thread();


/*CREATE OR REPLACE FUNCTION new_post() RETURNS TRIGGER AS
$body$
BEGIN
    UPDATE forums
    SET posts = posts + 1
    WHERE slug = NEW.forum;
    RETURN NEW;
END;
$body$ LANGUAGE plpgsql;

CREATE TRIGGER new_posts_trigger
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE new_post();*/


CREATE OR REPLACE FUNCTION new_path() RETURNS TRIGGER AS
$body$
BEGIN
    NEW.path = (SELECT path FROM posts WHERE id = NEW.parent) || NEW.id;
    RETURN NEW;
END;
$body$ LANGUAGE plpgsql;

CREATE TRIGGER new_path_trigger
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE new_path();


CREATE OR REPLACE FUNCTION insert_users_forum() RETURNS TRIGGER AS
$body$
BEGIN
    INSERT INTO users_forum(slug, nickname)
    VALUES (NEW.forum, NEW.author)
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$body$ LANGUAGE plpgsql;

CREATE TRIGGER insert_forum_user_trigger
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE insert_users_forum();


-- CREATE TRIGGER insert_forum_user_trigger
--     AFTER INSERT
--     ON posts
--     FOR EACH ROW
-- EXECUTE PROCEDURE insert_users_forum();


CREATE OR REPLACE FUNCTION update_votes() RETURNS TRIGGER AS
$body$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        UPDATE threads
        SET votes = votes + NEW.voice
        WHERE id = NEW.thread;
    ELSE
        UPDATE threads
        SET votes = votes - OLD.voice + NEW.voice
        WHERE id = NEW.thread;
    END IF;
    RETURN NEW;
END;
$body$ LANGUAGE plpgsql;


CREATE TRIGGER update_vote_trigger
    AFTER UPDATE OR INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE update_votes();