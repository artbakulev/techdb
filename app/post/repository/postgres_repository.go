package repository

import (
	"fmt"
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/post"
	"github.com/artbakulev/techdb/pkg/postsSQLGenerator"
	"github.com/jackc/pgx"
	"github.com/lib/pq"
	"strings"
)

type postgresPostRepository struct {
	conn *pgx.ConnPool
}

func (p postgresPostRepository) GetByID(id int64) (models.Post, *models.Error) {

	res, err := p.conn.Query("SELECT author, created, forum, id, isedited, message, parent, thread, path FROM posts WHERE id = $1", id)
	if err != nil {
		return models.Post{}, models.NewError(404, models.NotFoundError, err.Error())
	}
	defer res.Close()

	foundPost := models.Post{}

	for res.Next() {
		err = res.Scan(&foundPost.Author, &foundPost.Created,
			&foundPost.Forum, &foundPost.ID,
			&foundPost.IsEdited, &foundPost.Message,
			&foundPost.Parent, &foundPost.Thread, pq.Array(&foundPost.Path))

		if err != nil {
			return models.Post{}, models.NewError(500, models.DBParsingError, err.Error())
		}

		return foundPost, nil
	}

	return models.Post{}, models.NewError(404, models.NotFoundError)
}

func (p postgresPostRepository) CreateMany(posts models.Posts, thread models.Thread) (models.Posts, *models.Error) {
	if len(posts) == 0 {
		return models.Posts{}, nil
	}

	tx, _ := p.conn.Begin()
	defer tx.Rollback()

	mapParents := make(map[int64]models.Post)

	for _, item := range posts {
		if _, ok := mapParents[item.Parent]; !ok && item.Parent != 0 {
			parentPostQuery, err := p.GetByID(item.Parent)
			if err != nil {
				err.StatusCode = 409
				return models.Posts{}, err
			}

			if parentPostQuery.Thread != thread.ID {
				return models.Posts{}, models.NewError(409, models.BadRequestError)
			}

			mapParents[item.Parent] = parentPostQuery
		}
	}

	postIdsRows, err := tx.Query(fmt.Sprintf(`SELECT nextval(pg_get_serial_sequence('posts', 'id')) FROM generate_series(1, %d);`, len(posts)))
	if err != nil {
		return models.Posts{}, models.NewError(404, models.NotFoundError, err.Error())
	}
	var postIds []int64
	for postIdsRows.Next() {
		var availableId int64
		_ = postIdsRows.Scan(&availableId)
		postIds = append(postIds, availableId)
	}
	postIdsRows.Close()

	if len(postIds) == 0 {
		return models.Posts{}, models.NewError(500, models.DBError)
	}

	posts[0].Path = append(mapParents[posts[0].Parent].Path, postIds[0])

	err = tx.QueryRow(`INSERT INTO posts (id, author, forum, message, parent, thread, path) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING created`,
		postIds[0], posts[0].Author, thread.Forum, posts[0].Message, posts[0].Parent,
		thread.ID,
		"{"+strings.Trim(strings.Replace(fmt.Sprint(posts[0].Path), " ", ",", -1), "[]")+"}").
		Scan(&posts[0].Created)

	if err != nil {
		return models.Posts{}, models.NewError(404, models.CreateError, err.Error())
	}

	now := posts[0].Created

	posts[0].Forum = thread.Forum
	posts[0].Thread = thread.ID
	posts[0].Created = now
	posts[0].ID = postIds[0]

	for i, item := range posts {
		if i == 0 {
			continue
		}

		item.Path = append(mapParents[item.Parent].Path, postIds[i])
		resInsert, err := tx.Exec(`INSERT INTO posts (id, author, created, forum, message, parent, thread, path) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			postIds[i], item.Author, now, thread.Forum, item.Message, item.Parent, thread.ID,
			"{"+strings.Trim(strings.Replace(fmt.Sprint(item.Path), " ", ",", -1), "[]")+"}")

		if err != nil {
			return models.Posts{}, models.NewError(500, models.CreateError, err.Error())
		}

		if resInsert.RowsAffected() == 0 {
			return models.Posts{}, models.NewError(500, models.CreateError)
		}

		posts[i].Forum = thread.Forum
		posts[i].Thread = thread.ID
		posts[i].Created = now
		posts[i].ID = postIds[i]
	}

	_, err = tx.Exec(`UPDATE forums SET posts = posts + $1 WHERE slug = $2`, len(posts), thread.Forum)
	if err != nil {
		return models.Posts{}, models.NewError(500, models.InternalError, err.Error())
	}

	err = tx.Commit()

	if err != nil {
		return models.Posts{}, models.NewError(500, models.DBError, err.Error())
	}

	return posts, nil
}

func (p postgresPostRepository) Update(post models.Post, postUpdate models.PostUpdate) (models.Post, *models.Error) {

	if postUpdate.Message == "" {
		return post, nil
	}

	if post.Message == postUpdate.Message {
		return post, nil
	}

	res, err := p.conn.Exec("UPDATE posts SET message = $1, isedited = true WHERE id = $2", postUpdate.Message, post.ID)
	if err != nil {
		return models.Post{}, models.NewError(409, models.UpdateError, err.Error())
	}

	if res.RowsAffected() == 0 {
		return models.Post{}, models.NewError(500, models.InternalError)
	}

	post.Message = postUpdate.Message
	post.IsEdited = true

	return post, nil
}

func (p postgresPostRepository) GetMany(thread models.Thread, query models.PostsRequestQuery) (models.Posts, *models.Error) {

	generator := postsSQLGenerator.NewPostsSQLGenerator(thread, query)

	if query.Sort == "" {
		query.Sort = models.FLAT
	}

	baseSQL := ""
	sortedPosts := make([]models.Post, 0, 1)

	switch query.Sort {
	case models.FLAT:
		baseSQL = generator.FlatSort()

	case models.TREE:
		baseSQL = generator.TreeSort()

	case models.PARENT_TREE:
		baseSQL = generator.ParentTreeSort()
	}

	res, err := p.conn.Query(baseSQL)
	if err != nil {
		return models.Posts{}, models.NewError(500, models.InternalError, err.Error())
	}
	defer res.Close()

	bufferPost := models.Post{}

	for res.Next() {
		err := res.Scan(&bufferPost.Author, &bufferPost.Created, &bufferPost.Forum, &bufferPost.ID,
			&bufferPost.IsEdited, &bufferPost.Message, &bufferPost.Parent, &bufferPost.Thread)

		if err != nil {
			return models.Posts{}, models.NewError(500, models.DBParsingError, err.Error())
		}
		sortedPosts = append(sortedPosts, bufferPost)
	}

	return sortedPosts, nil
}

func NewPostgresPostRepository(conn *pgx.ConnPool) post.Repository {
	return &postgresPostRepository{conn: conn}
}
