package forum

import "github.com/artbakulev/techdb/app/models"

type Usecase interface {
	CreateForum(forumNew models.Forum) (models.Forum, *models.Error)
	GetForumBySlug(slug string) (models.Forum, *models.Error)
	GetForumUsers(slug string, query models.PostsRequestQuery) (models.Users, *models.Error)
}
