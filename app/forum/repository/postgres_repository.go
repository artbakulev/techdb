package repository

import (
	"github.com/artbakulev/techdb/app/forum"
	"github.com/artbakulev/techdb/app/models"
	"github.com/jackc/pgx"
)

type postgresForumRepository struct {
	conn *pgx.ConnPool
}

func NewPostgresForumRepository(db *pgx.ConnPool) forum.Repository {
	return &postgresForumRepository{conn: db}
}

func (p postgresForumRepository) GetBySlug(slug string) (models.Forum, *models.Error) {
	res, err := p.conn.Query(`SELECT * FROM forums WHERE slug = $1`, slug)
	if err != nil {
		return models.Forum{}, models.NewError(500, models.InternalError)
	}
	defer res.Close()

	f := models.Forum{}

	if res.Next() {
		err := res.Scan(&f.Slug, &f.Title, &f.User, &f.Threads, &f.Posts)
		if err != nil {
			return models.Forum{}, models.NewError(500, models.DBParsingError)
		}
		return f, nil
	}

	return models.Forum{}, models.NewError(404, models.NotFoundError)
}

func (p postgresForumRepository) Create(user models.User, forumNew models.Forum) (models.Forum, *models.Error) {
	forumNew.User = user.Nickname

	_, err := p.conn.Exec(`INSERT INTO forums (slug, title, "user", posts, threads) VALUES ($1, $2, $3, $4, $5)`,
		forumNew.Slug, forumNew.Title, forumNew.User, forumNew.Posts, forumNew.Threads)
	if err != nil {
		return models.Forum{}, models.NewError(409, models.CreateError)
	}

	return forumNew, nil
}
