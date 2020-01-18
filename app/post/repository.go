package post

import "github.com/artbakulev/techdb/app/models"

type Repository interface {
	CreateMany(posts models.Posts, thread models.Thread) (models.Posts, *models.Error)
	Update(post models.Post, postUpdate models.PostUpdate) (models.Post, *models.Error)
	GetMany(thread models.Thread, query models.PostsRequestQuery) (models.Posts, *models.Error)
	GetByID(id int64) (models.Post, *models.Error)
}
