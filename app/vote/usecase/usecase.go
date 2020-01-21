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
	if vote.Thread == -1 {
		foundThread, err := v.threadRepo.GetBySlug(vote.ThreadSlug)
		if err != nil {
			return models.Thread{}, err
		}
		vote.Thread = foundThread.ID
	}
	err := v.voteRepo.Update(vote)
	if err != nil && err.StatusCode == 404 {
		err = v.voteRepo.Create(vote)
	}
	if err != nil {
		return models.Thread{}, err
	}
	return v.threadRepo.GetByID(vote.Thread)
}
