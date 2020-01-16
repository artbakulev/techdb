package thread

import "github.com/artbakulev/techdb/app/models"

type Repository interface {
	GetByID(id int64) (models.Thread, *models.Error)
	GetBySlug(slug string) (models.Thread, *models.Error)
	Create(threadNew models.ThreadNew) (models.Thread, *models.Error)
	Update(threadUpdate models.ThreadUpdate) (models.Thread, *models.Error)
}
