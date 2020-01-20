package repository

import (
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/user"
	"github.com/jackc/pgx"
)

type postgresUserRepository struct {
	conn *pgx.ConnPool
}

func NewPostgresUserRepository(db *pgx.ConnPool) user.Repository {
	return &postgresUserRepository{conn: db}
}

func (p postgresUserRepository) GetByNickname(nickname string) (models.User, *models.Error) {
	u := models.User{}
	res, err := p.conn.Query(`SELECT about, email, fullname, nickname FROM users WHERE nickname = $1`, nickname)
	if err != nil {
		return models.User{}, models.NewError(500, models.InternalError, err.Error())
	}
	defer res.Close()

	if res.Next() {
		err = res.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
		if err != nil {
			return models.User{}, models.NewError(404, models.DBParsingError, err.Error())
		}
		return u, nil
	}
	return models.User{}, models.NewError(404, models.NotFoundError)
}

func (p postgresUserRepository) GetByEmail(email string) (models.User, *models.Error) {
	u := models.User{}
	res, err := p.conn.Query(`SELECT about, email, fullname, nickname FROM users WHERE email = $1`, email)
	if err != nil {
		return models.User{}, models.NewError(500, models.NotFoundError, err.Error())
	}
	defer res.Close()

	if res.Next() {
		err := res.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
		if err != nil {
			return models.User{}, models.NewError(500, models.DBParsingError, err.Error())
		}
		return u, nil
	}
	return models.User{}, models.NewError(404, models.NotFoundError)
}

func (p postgresUserRepository) Create(userNew models.User) (models.User, *models.Error) {
	res, err := p.conn.Exec(`INSERT INTO users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4)`,
		userNew.Nickname, userNew.Fullname, userNew.Email, userNew.About)
	if err != nil {
		return models.User{}, models.NewError(500, models.CreateError, err.Error())
	}

	if res.RowsAffected() == 0 {
		return models.User{}, models.NewError(500, models.CreateError)
	}

	return userNew, nil
}

func (p postgresUserRepository) Update(userUpdate models.User) (models.User, *models.Error) {

	if userUpdate.About == "" && userUpdate.Email == "" && userUpdate.Fullname == "" {
		updatedUser, _ := p.GetByNickname(userUpdate.Nickname)
		return updatedUser, nil
	}

	baseSQL := "UPDATE users SET"
	if userUpdate.Fullname == "" {
		baseSQL += " fullname = fullname,"
	} else {
		baseSQL += " fullname = '" + userUpdate.Fullname + "',"
	}

	if userUpdate.Email == "" {
		baseSQL += " email = email,"
	} else {
		baseSQL += " email = '" + userUpdate.Email + "',"
	}

	if userUpdate.About == "" {
		baseSQL += " about = about"
	} else {
		baseSQL += " about = '" + userUpdate.About + "'"
	}

	baseSQL += " WHERE nickname = '" + userUpdate.Nickname + "'"

	res, err := p.conn.Exec(baseSQL)
	if err != nil {
		return models.User{}, models.NewError(500, models.UpdateError, err.Error())
	}

	if res.RowsAffected() == 0 {
		return models.User{}, models.NewError(404, models.NotFoundError)
	}

	updatedUser, _ := p.GetByNickname(userUpdate.Nickname)

	return updatedUser, nil
}

func (p postgresUserRepository) GetByForum(forum models.Forum, query models.PostsRequestQuery) (models.Users, *models.Error) {
	baseSQL := `SELECT about, email, fullname, fu.nickname FROM users_forum JOIN users u ON u.nickname = users_forum.nickname`

	baseSQL += ` where slug = '` + forum.Slug + `'`
	if query.Since > 0 {
		if query.Desc {
			baseSQL += ` AND u.nickname < '` + query.GetStringSince() + `'`
		} else {
			baseSQL += ` AND u.nickname > '` + query.GetStringSince() + `'`
		}
	}

	if query.Desc {
		baseSQL += " ORDER BY nickname DESC"
	} else {
		baseSQL += " ORDER BY nickname ASC"
	}

	if query.Limit > 0 {
		baseSQL += " LIMIT " + query.GetStringLimit()
	}

	res, err := p.conn.Query(baseSQL)
	if err != nil {
		return models.Users{}, models.NewError(500, models.DBParsingError, err.Error())
	}
	defer res.Close()

	foundUsers := models.Users{}
	buffer := models.User{}

	for res.Next() {
		err = res.Scan(&buffer.About, &buffer.Email, &buffer.Fullname, &buffer.Nickname)

		if err != nil {
			return models.Users{}, models.NewError(500, models.DBParsingError, err.Error())
		}
		foundUsers = append(foundUsers, buffer)
	}

	return foundUsers, nil
}
