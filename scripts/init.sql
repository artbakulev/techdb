CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS users_forum CASCADE;

CREATE UNLOGGED TABLE users
(
    nickname CITEXT NOT NULL UNIQUE PRIMARY KEY,
    fullname TEXT,
    email    CITEXT NOT NULL UNIQUE,
    about    TEXT
);


CREATE INDEX IF NOT EXISTS idx_users_nickname ON users (nickname);
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);


CREATE UNLOGGED TABLE forums
(
    slug    CITEXT  NOT NULL UNIQUE PRIMARY KEY ,
    title   TEXT    NOT NULL,
    "user"  CITEXT  NOT NULL REFERENCES users (nickname),
    threads INTEGER DEFAULT 0,
    posts   INTEGER DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_forum_user ON forums ("user");


CREATE UNLOGGED TABLE threads
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

CREATE INDEX IF NOT EXISTS idx_threads_slug ON threads (slug);
CREATE INDEX IF NOT EXISTS idx_threads_forum_created ON threads (forum, created);
CREATE INDEX IF NOT EXISTS idx_threads_author_forum ON threads (author, forum);



CREATE UNLOGGED TABLE posts
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

CREATE INDEX IF NOT EXISTS idx_posts_path_id ON posts (id, (path [1]));
CREATE INDEX IF NOT EXISTS idx_posts_path ON posts (path);
CREATE INDEX IF NOT EXISTS idx_posts_path_1 ON posts ((path [1]));
CREATE INDEX IF NOT EXISTS idx_posts_thread_id ON posts (thread, id);
CREATE INDEX IF NOT EXISTS idx_posts_thread ON posts (thread);
CREATE INDEX IF NOT EXISTS idx_posts_thread_path_id ON posts (thread, path, id);
CREATE INDEX IF NOT EXISTS idx_posts_thread_id_path_parent ON posts (thread, id, (path[1]), parent);
CREATE INDEX IF NOT EXISTS idx_posts_author_forum ON posts (author, forum);


CREATE UNLOGGED TABLE votes
(
    nickname CITEXT REFERENCES  users (nickname) NOT NULL,
    voice    SMALLINT CHECK ( voice IN (-1, 1) ),
    thread   INTEGER REFERENCES threads (id)     NOT NULL,
    UNIQUE (nickname, thread)
);

CREATE INDEX IF NOT EXISTS idx_votes_nickname_thread ON votes (nickname, thread);

CREATE UNLOGGED TABLE users_forum
(
    nickname CITEXT REFERENCES users (nickname) NOT NULL,
    slug     CITEXT REFERENCES forums (slug) NOT NULL,
    UNIQUE (nickname, slug)
);

CREATE INDEX IF NOT EXISTS idx_users_forum_nickname_slug ON users_forum (nickname, slug);

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