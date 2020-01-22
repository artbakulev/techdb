package usecase

import (
	"github.com/artbakulev/techdb/app/forum"
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/post"
	"github.com/artbakulev/techdb/app/thread"
	"github.com/artbakulev/techdb/app/user"
)

type postUsecase struct {
	userRepo   user.Repository
	postRepo   post.Repository
	threadRepo thread.Repository
	forumRepo  forum.Repository
}

func NewPostUsecase(userRepo user.Repository, postRepo post.Repository,
	threadRepo thread.Repository, forumRepo forum.Repository) post.Usecase {
	return &postUsecase{
		userRepo:   userRepo,
		postRepo:   postRepo,
		threadRepo: threadRepo,
		forumRepo:  forumRepo,
	}
}

func (p postUsecase) CreatePosts(slug string, id int64, posts models.Posts) (models.Posts, *models.Error) {
	var foundThread models.Thread
	var err *models.Error

	if id == -1 {
		foundThread, err = p.threadRepo.GetBySlug(slug)
	} else {
		foundThread, err = p.threadRepo.GetByID(id)
	}
	if err != nil {
		return models.Posts{}, err
	}
	posts, err = p.postRepo.CreateMany(posts, foundThread)
	if err != nil {
		return models.Posts{}, err
	}

	for _, item := range posts {
		p.userRepo.AddUserToForum(item.Author, foundThread.Forum)
	}

	return posts, nil
}

func (p postUsecase) UpdatePost(id int64, newPost models.PostUpdate) (models.Post, *models.Error) {
	existingPost, err := p.postRepo.GetByID(id)
	if err != nil {
		return models.Post{}, err
	}
	return p.postRepo.Update(existingPost, newPost)
}

func (p postUsecase) GetPostDetails(id int64, query models.PostsRelatedQuery) (models.PostFull, *models.Error) {
	var postFull models.PostFull
	existingPost, err := p.postRepo.GetByID(id)
	if err != nil {
		return models.PostFull{}, err
	}
	postFull.Post = &existingPost

	if query.NeedAuthor {
		author, err := p.userRepo.GetByNickname(existingPost.Author)
		if err != nil {
			return models.PostFull{}, err
		}
		postFull.Author = &author
	}

	if query.NeedForum {
		existingForum, err := p.forumRepo.GetBySlug(existingPost.Forum)
		if err != nil {
			return models.PostFull{}, err
		}
		postFull.Forum = &existingForum
	}

	if query.NeedThread {
		existingThread, err := p.threadRepo.GetByID(existingPost.Thread)
		if err != nil {
			return models.PostFull{}, err
		}
		postFull.Thread = &existingThread
	}

	return postFull, nil
}

func (p postUsecase) GetThreadPosts(query models.PostsRequestQuery) (models.Posts, *models.Error) {
	var existingThread models.Thread
	var err *models.Error
	if query.ThreadID == -1 {
		existingThread, err = p.threadRepo.GetBySlug(query.ThreadSlug)
	} else {
		existingThread, err = p.threadRepo.GetByID(query.ThreadID)
	}
	if err != nil {
		return models.Posts{}, err
	} else if existingThread.Author == "" {
		return models.Posts{}, models.NewError(400, models.BadRequestError)
	}
	return p.postRepo.GetMany(existingThread, query)
}
