package usecase

import (
	"github.com/artbakulev/techdb/app/forum"
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/user"
)

type forumUsecase struct {
	userRepo  user.Repository
	forumRepo forum.Repository
}

func NewForumUsecase(userRepo user.Repository, forumRepo forum.Repository) forum.Usecase {
	return &forumUsecase{
		userRepo:  userRepo,
		forumRepo: forumRepo,
	}
}

func (f forumUsecase) CreateForum(forumNew models.Forum) (models.Forum, *models.Error) {
	author, err := f.userRepo.GetByNickname(forumNew.User)
	if err != nil {
		return models.Forum{}, err
	}
	return f.forumRepo.Create(author, forumNew)
}

func (f forumUsecase) GetForumBySlug(slug string) (models.Forum, *models.Error) {
	return f.forumRepo.GetBySlug(slug)
}

func (f forumUsecase) GetForumUsers(slug string, query models.PostsRequestQuery) (models.Users, *models.Error) {
	existingForum, err := f.forumRepo.GetBySlug(slug)
	if err != nil {
		return models.Users{}, err
	}
	return f.userRepo.GetByForum(existingForum, query)
}
