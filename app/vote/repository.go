package vote

import "github.com/artbakulev/techdb/app/models"

type Repository interface {
	Create(vote models.Vote) (models.Thread, *models.Error)
	Update(vote models.Vote) (models.Thread, *models.Error)
	GetByNicknameAndThreadID(nickname string, threadID int32) (models.Vote, *models.Error)
}
