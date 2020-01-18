package usecase

import (
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/post"
	"github.com/artbakulev/techdb/app/user"
)

//mapUsers := make(map[string]string) - slice

//if _, ok := mapUsers[item.Author]; !ok {
//mapUsers[item.Author] = item.Author
//}

//go func() {
//	for _, val := range mapUsers {
//		AddUser(val, thread.Forum)
//	}
//}()

type postUsecase struct {
	userRepo user.Repository
	postRepo post.Repository
}

func NewPostUsecase(userRepo user.Repository, postRepo post.Repository) post.Usecase {
	return &postUsecase{
		userRepo: userRepo,
		postRepo: postRepo,
	}
}

func (p postUsecase) CreatePosts(posts models.Posts) (models.Posts, *models.Error) {
	panic("implement me")
}

func (p postUsecase) UpdatePost(id int64) (models.Post, *models.Error) {
	panic("implement me")
}

func (p postUsecase) GetPostDetails(id int64, query models.PostsRelatedQuery) (models.PostFull, *models.Error) {
	panic("implement me")
}

func (p postUsecase) GetThreadPosts(query models.PostsRequestQuery) (models.Posts, *models.Error) {
	panic("implement me")
}
