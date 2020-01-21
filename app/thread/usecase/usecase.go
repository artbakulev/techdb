package usecase

import (
	"github.com/artbakulev/techdb/app/forum"
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/thread"
	"github.com/artbakulev/techdb/app/user"
	"strconv"
)

type threadUsecase struct {
	threadRepo thread.Repository
	userRepo   user.Repository
	forumRepo  forum.Repository
}

func NewThreadUsecase(threadRepo thread.Repository, userRepo user.Repository, forumRepo forum.Repository) thread.Usecase {
	return &threadUsecase{
		threadRepo: threadRepo,
		userRepo:   userRepo,
		forumRepo:  forumRepo,
	}
}

func (t threadUsecase) GetThreadBySlugOrID(data string, isSlug bool) (models.Thread, *models.Error) {
	if isSlug {
		return t.threadRepo.GetBySlug(data)
	}
	id, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return models.Thread{}, models.NewError(400, models.BadRequestError)
	}
	return t.threadRepo.GetByID(id)
}

func (t threadUsecase) CreateThread(slug string, thread models.Thread) (models.Thread, *models.Error) {
	foundForum, err := t.forumRepo.GetBySlug(slug)
	if err != nil {
		return models.Thread{}, err
	}
	foundUser, err := t.userRepo.GetByNickname(thread.Author)
	if err != nil {
		return models.Thread{}, err
	}

	thread.Forum = foundForum.Slug
	createdThread, err := t.threadRepo.Create(foundForum, foundUser, thread)

	if err != nil && err.StatusCode == 409 {
		thread, _ = t.threadRepo.GetBySlug(thread.Slug)
		return thread, models.NewError(409, models.ConflictError)
	}

	return createdThread, err
}

func (t threadUsecase) UpdateThread(threadUpdate models.ThreadUpdate) (models.Thread, *models.Error) {
	var foundThread models.Thread
	var err *models.Error
	if threadUpdate.Slug != "" {
		foundThread, err = t.threadRepo.GetBySlug(threadUpdate.Slug)

	} else {
		foundThread, err = t.threadRepo.GetByID(threadUpdate.ID)
	}

	if err != nil {
		return models.Thread{}, err
	}

	return t.threadRepo.Update(foundThread, threadUpdate)
}

func (t threadUsecase) GetThreads(forumSlug string, query models.PostsRequestQuery) (models.Threads, *models.Error) {
	existingForum, err := t.forumRepo.GetBySlug(forumSlug)
	if err != nil {
		return models.Threads{}, err
	}

	return t.threadRepo.GetMany(existingForum, query)
}
