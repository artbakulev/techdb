package vote

import "github.com/artbakulev/techdb/app/models"

type Usecase interface {
	UpsertVote(vote models.Vote) (models.Thread, *models.Error)
}
