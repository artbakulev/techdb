package repository

import (
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/service"
	"github.com/jackc/pgx"
)

type postgresServiceRepository struct {
	conn *pgx.ConnPool
}

func NewPostgresServiceRepository(conn *pgx.ConnPool) service.Repository {
	return &postgresServiceRepository{conn: conn}
}

func (p postgresServiceRepository) Clear() *models.Error {
	res, err := p.conn.Query("TRUNCATE TABLE users, forums, threads, posts, votes, users_forum CASCADE")
	if err != nil {
		return models.NewError(500, models.InternalError)
	}
	defer res.Close()
	return nil
}

func (p postgresServiceRepository) GetStatus() (models.Status, *models.Error) {
	res, err := p.conn.Query("SELECT * FROM (SELECT count(posts) FROM forums) as f" +
		" CROSS JOIN (SELECT count(id) FROM posts) as p" +
		" CROSS JOIN (SELECT count(id) FROM threads) as t" +
		" CROSS JOIN (SELECT count(nickname) FROM users) as u")

	if err != nil {
		return models.Status{}, models.NewError(500, models.InternalError)
	}

	defer res.Close()

	s := models.Status{}
	for res.Next() {
		err = res.Scan(&s.Forum, &s.Post, &s.Thread, &s.User)

		if err != nil {
			return models.Status{}, models.NewError(500, models.InternalError)
		}

		return s, nil
	}

	return models.Status{}, models.NewError(500, models.InternalError)
}
