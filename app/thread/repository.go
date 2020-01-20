package thread

import "github.com/artbakulev/techdb/app/models"

type Repository interface {
	GetByID(id int64) (models.Thread, *models.Error)
	GetBySlug(slug string) (models.Thread, *models.Error)
	Create(forum models.Forum, user models.User, thread models.Thread) (models.Thread, *models.Error)
	Update(thread models.Thread, threadUpdate models.ThreadUpdate) (models.Thread, *models.Error)
	GetMany(forum models.Forum, query models.PostsRequestQuery) (models.Threads, *models.Error)
}
