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

	res, err := p.conn.Query(`SELECT * FROM forum_thread WHERE id = $1`, id)
	if err != nil {
		return models.Thread{}, models.NewError(404, models.NotFoundError)
	}
	defer res.Close()

	if res.Next() {
		nullString := pgtype.Text{}
		err = res.Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.Message, &nullString, &t.Title, &t.Votes)
		if err != nil {
			return models.Thread{}, models.NewError(500, models.InternalError)
		}

		t.Slug = nullString.String

		return t, nil
	}

	return models.Thread{}, models.NewError(404, models.NotFoundError)
}

func (p postgresThreadRepository) GetBySlug(slug string) (models.Thread, *models.Error) {
	if slug == "" {
		return models.Thread{}, models.NewError(400, models.BadRequestError)
	}

	t := models.Thread{}
	res, err := p.conn.Query(`SELECT * FROM forum_thread WHERE slug = $1`, slug)

	if err != nil {
		return models.Thread{}, models.NewError(404, models.NotFoundError)
	}
	defer res.Close()

	if res.Next() {
		nullString := pgtype.Text{}
		err = res.Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.Message, &nullString, &t.Title, &t.Votes)
		if err != nil {
			return models.Thread{}, models.NewError(500, models.DBParsingError)
		}

		t.Slug = nullString.String

		return t, nil
	}

	return models.Thread{}, models.NewError(500, models.InternalError)
}

func (p postgresThreadRepository) Create(forum models.Forum, user models.User, thread models.Thread) (models.Thread, *models.Error) {
	thread.Forum = forum.Slug
	thread.Author = user.Nickname

	tx, _ := p.conn.Begin()
	defer tx.Rollback()

	if thread.Slug == "" {
		err := tx.QueryRow(`INSERT INTO forum_thread (author, forum, message, slug, title) VALUES ($1, $2, $3, $4, NULL, $5) RETURNING id`,
			thread.Author, thread.Created, thread.Forum, thread.Message,
			thread.Title).Scan(&thread.ID)

		if err == pgx.ErrNoRows || err != nil {
			return models.Thread{}, models.NewError(409, models.ConflictError)
		}
	} else {
		err := tx.QueryRow(`INSERT INTO forum_thread (author, created, forum, message, slug, title) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			thread.Author, thread.Created, thread.Forum, thread.Message, thread.Slug,
			thread.Title).Scan(&thread.ID)

		if err == pgx.ErrNoRows || err != nil {
			return models.Thread{}, models.NewError(409, models.ConflictError)
		}
	}

	//TODO: update forum stat (AddUser)
	err := tx.Commit()
	if err != nil {
		return models.Thread{}, models.NewError(500, models.InternalError, err.Error())
	}

	return thread, nil
}

func (p postgresThreadRepository) Update(thread models.Thread, threadUpdate models.ThreadUpdate) (models.Thread, *models.Error) {
	if threadUpdate.Message == "" && threadUpdate.Title == "" {
		return thread, nil
	}

	baseSQL := "UPDATE forum_thread SET"
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
