package usecase

import (
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/thread"
	"github.com/artbakulev/techdb/app/vote"
)

type voteUsecase struct {
	voteRepo   vote.Repository
	threadRepo thread.Repository
}

func NewVoteUsecase(voteRepo vote.Repository, threadRepo thread.Repository) vote.Usecase {
	return &voteUsecase{
		voteRepo:   voteRepo,
		threadRepo: threadRepo,
	}
}

func (v voteUsecase) UpsertVote(vote models.Vote) (models.Thread, *models.Error) {
	err := v.voteRepo.Update(vote)
	if err.StatusCode == 404 {
		err = v.voteRepo.Create(vote)
	}
	if err != nil {
		return models.Thread{}, err
	}
	return v.threadRepo.GetByID(vote.Thread)
}
