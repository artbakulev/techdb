package forum

import "github.com/artbakulev/techdb/app/models"

type Repository interface {
	GetBySlug(slug string) (models.Forum, *models.Error)
	Create(forumNew models.ForumNew) (models.Forum, *models.Error)
}
