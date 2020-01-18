package forum

import "github.com/artbakulev/techdb/app/models"

type Repository interface {
	GetBySlug(slug string) (models.Forum, *models.Error)
	Create(user models.User, forumNew models.Forum) (models.Forum, *models.Error)
}
