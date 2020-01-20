package thread

import "github.com/artbakulev/techdb/app/models"

type Usecase interface {
	GetThreadBySlugOrID(data string, isSlug bool) (models.Thread, *models.Error)
	CreateThread(slug string, thread models.Thread) (models.Thread, *models.Error)
	UpdateThread(threadUpdate models.ThreadUpdate) (models.Thread, *models.Error)
	GetThreads(forumSlug string, query models.PostsRequestQuery) (models.Threads, *models.Error)
}
