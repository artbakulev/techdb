package thread

import "github.com/artbakulev/techdb/app/models"

type Usecase interface {
	GetThreadBySlugOrID(data string, isSlug bool) (models.Thread, *models.Error)
	CreateThread(threadNew models.ThreadNew) (models.Thread, *models.Error)
	UpdateThread(threadUpdate models.ThreadUpdate) (models.Thread, *models.Error)
}
