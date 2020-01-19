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
		return models.User{}, models.NewError(500, "cannot get user by nickname", err.Error())
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
	res, err := p.conn.Exec(`INSERT INTO forum_users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4)`,
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

	baseSQL := "Update forum_users SET"
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
		return models.User{}, models.NewError(409, models.UpdateError, err.Error())
	}

	if res.RowsAffected() == 0 {
		return models.User{}, models.NewError(404, models.NotFoundError)
	}

	updatedUser, _ := p.GetByNickname(userUpdate.Nickname)

	return updatedUser, nil
}
