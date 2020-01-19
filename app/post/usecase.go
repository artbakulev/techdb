package post

import "github.com/artbakulev/techdb/app/models"

type Usecase interface {
	CreatePosts(slug string, id int64, posts models.Posts) (models.Posts, *models.Error)
	UpdatePost(id int64, newPost models.PostUpdate) (models.Post, *models.Error)
	GetPostDetails(id int64, query models.PostsRelatedQuery) (models.PostFull, *models.Error)
	GetThreadPosts(query models.PostsRequestQuery) (models.Posts, *models.Error)
}
