package repository

import (
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/vote"
	"github.com/jackc/pgx"
)

type postgresVoteRepository struct {
	conn *pgx.ConnPool
}

func NewPostgresVoteRepository(conn *pgx.ConnPool) vote.Repository {
	return &postgresVoteRepository{conn: conn}
}

func (p postgresVoteRepository) Create(vote models.Vote) *models.Error {
	resInsert, err := p.conn.Exec(`INSERT INTO forum_vote (nickname, voice, thread) VALUES ($1, $2, $3)`,
		vote.Nickname, vote.Voice, vote.Thread)

	if err != nil {
		return models.NewError(500, models.InternalError)
	}

	if resInsert.RowsAffected() == 0 {
		return models.NewError(404, models.NotFoundError)
	}

	return nil
}

func (p postgresVoteRepository) Update(vote models.Vote) *models.Error {
	res, err := p.conn.Exec(`UPDATE forum_vote SET voice = $1 WHERE nickname = $2 AND thread = $3`,
		vote.Voice, vote.Nickname, vote.Thread)
	if err != nil {
		return models.NewError(500, models.UpdateError)
	}
	if res.RowsAffected() == 0 {
		return models.NewError(404, models.NotFoundError)
	}
	return nil
}

func (p postgresVoteRepository) GetByNicknameAndThreadID(nickname string, threadID int32) (models.Vote, *models.Error) {
	res, err := p.conn.Query(`SELECT * FROM forum_vote WHERE nickname = $1 AND thread = $2`, nickname, threadID)

	if err != nil {
		return models.Vote{}, models.NewError(404, models.NotFoundError)
	}
	defer res.Close()

	existingVote := models.Vote{}

	if res.Next() {
		err := res.Scan(&existingVote.Nickname, &existingVote.Voice, &existingVote.Thread)
		if err != nil {
			return models.Vote{}, models.NewError(500, models.DBParsingError)
		}

		return existingVote, nil
	}
	return models.Vote{}, models.NewError(404, models.NotFoundError)
}
