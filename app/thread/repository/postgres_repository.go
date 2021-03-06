package repository

import (
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/thread"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
)

type postgresThreadRepository struct {
	conn *pgx.ConnPool
}

func NewPostgresThreadRepository(conn *pgx.ConnPool) thread.Repository {
	return &postgresThreadRepository{conn: conn}
}

func (p postgresThreadRepository) GetByID(id int64) (models.Thread, *models.Error) {
	t := models.Thread{}

	if id == -1 {
		return models.Thread{}, models.NewError(400, models.BadRequestError)
	}

	res, err := p.conn.Query(`SELECT * FROM threads WHERE id = $1`, id)
	if err != nil {
		return models.Thread{}, models.NewError(404, models.NotFoundError)
	}
	defer res.Close()

	nullSlug := &pgtype.Varchar{}
	if res.Next() {
		err = res.Scan(&t.ID, nullSlug, &t.Author, &t.Forum, &t.Title, &t.Message, &t.Created, &t.Votes)
		if err != nil {
			return models.Thread{}, models.NewError(500, models.InternalError)
		}

		t.Slug = nullSlug.String
		return t, nil
	}

	return models.Thread{}, models.NewError(404, models.NotFoundError)
}

func (p postgresThreadRepository) GetBySlug(slug string) (models.Thread, *models.Error) {
	if slug == "" {
		return models.Thread{}, models.NewError(400, models.BadRequestError)
	}

	t := models.Thread{}
	res, err := p.conn.Query(`SELECT * FROM threads WHERE slug = $1`, slug)

	if err != nil {
		return models.Thread{}, models.NewError(404, models.NotFoundError)
	}
	defer res.Close()

	if res.Next() {
		err = res.Scan(&t.ID, &t.Slug, &t.Author, &t.Forum, &t.Title, &t.Message, &t.Created, &t.Votes)
		if err != nil {
			return models.Thread{}, models.NewError(500, models.DBParsingError)
		}

		return t, nil
	}

	return models.Thread{}, models.NewError(404, models.InternalError)
}

func (p postgresThreadRepository) Create(forum models.Forum, user models.User, thread models.Thread) (models.Thread, *models.Error) {
	thread.Forum = forum.Slug
	thread.Author = user.Nickname

	tx, _ := p.conn.Begin()
	defer tx.Rollback()
	if thread.Slug == "" {

		err := tx.QueryRow(`INSERT INTO threads (author, created, forum, message, title) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			thread.Author, thread.Created, thread.Forum, thread.Message,
			thread.Title).Scan(&thread.ID)

		if err == pgx.ErrNoRows || err != nil {
			return models.Thread{}, models.NewError(409, models.ConflictError)
		}

	} else {

		err := tx.QueryRow(`INSERT INTO threads (author, created, forum, message, slug, title) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			thread.Author, thread.Created, thread.Forum, thread.Message, thread.Slug,
			thread.Title).Scan(&thread.ID)

		if err == pgx.ErrNoRows || err != nil {
			return models.Thread{}, models.NewError(409, models.ConflictError)
		}

	}

	err := tx.Commit()
	if err != nil {
		return models.Thread{}, models.NewError(500, models.InternalError)
	}

	return thread, nil
}

func (p postgresThreadRepository) Update(thread models.Thread, threadUpdate models.ThreadUpdate) (models.Thread, *models.Error) {
	if threadUpdate.Message == "" && threadUpdate.Title == "" {
		return thread, nil
	}

	baseSQL := "UPDATE threads SET"
	if threadUpdate.Message == "" {
		baseSQL += " message = message,"
	} else {
		thread.Message = threadUpdate.Message
		baseSQL += " message = '" + threadUpdate.Message + "',"
	}

	if threadUpdate.Title == "" {
		baseSQL += " title = title"
	} else {
		thread.Title = threadUpdate.Title
		baseSQL += " title = '" + threadUpdate.Title + "'"
	}

	baseSQL += " WHERE slug = '" + thread.Slug + "'"

	res, err := p.conn.Exec(baseSQL)
	if err != nil {
		return models.Thread{}, models.NewError(500, models.UpdateError)
	}

	if res.RowsAffected() == 0 {
		return models.Thread{}, models.NewError(404, models.NotFoundError)
	}

	return thread, nil
}

func (p postgresThreadRepository) GetMany(forum models.Forum, query models.PostsRequestQuery) (models.Threads, *models.Error) {
	baseSQL := "SELECT * FROM threads"

	baseSQL += " WHERE forum = '" + forum.Slug + "'"

	if query.Since != "" {
		if query.Desc {
			baseSQL += " AND created <= '" + query.Since + "'"
		} else {
			baseSQL += " AND created >= '" + query.Since + "'"
		}
	}

	if query.Desc {
		baseSQL += " ORDER BY created DESC"
	} else {
		baseSQL += " ORDER BY created"
	}

	if query.Limit > 0 {
		baseSQL += " LIMIT " + query.GetStringLimit()
	}

	res, err := p.conn.Query(baseSQL)
	if err != nil {
		return models.Threads{}, models.NewError(500, models.DBParsingError)
	}

	buffer := models.Thread{}
	existingThreads := models.Threads{}

	for res.Next() {
		nullSlug := pgtype.Varchar{}

		err = res.Scan(&buffer.ID, &nullSlug, &buffer.Author, &buffer.Forum,
			&buffer.Title, &buffer.Message, &buffer.Created, &buffer.Votes)

		if err != nil {
			return models.Threads{}, models.NewError(500, models.InternalError)
		}
		buffer.Slug = nullSlug.String
		existingThreads = append(existingThreads, buffer)
	}

	return existingThreads, nil
}
