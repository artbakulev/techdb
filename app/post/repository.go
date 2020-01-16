package post

import "github.com/artbakulev/techdb/app/models"

type Repository interface {
	Create(post models.Post) (models.Post, *models.Error)
	Update(id int64) (models.Post, *models.Error)
	GetFull(id int64, query models.PostsRelatedQuery) (models.PostFull, *models.Error)
	GetMany(query models.PostsRequestQuery) (models.Posts, *models.Error)
}
