package postsSQLGenerator

import (
	"github.com/artbakulev/techdb/app/models"
	"strconv"
)

type PostsSQLGenerator interface {
	FlatSort() string
	TreeSort() string
	ParentTreeSort() string
}

type postsSQLGenerator struct {
	thread models.Thread
	query  models.PostsRequestQuery
}

func (p postsSQLGenerator) FlatSort() string {
	strID := strconv.FormatInt(int64(p.thread.ID), 10)
	baseSQL := ""

	baseSQL = "SELECT author, created, forum, id, isedited, message, parent, thread FROM posts WHERE thread = " + strID

	if p.query.Since != "" {
		if p.query.Desc {
			baseSQL += " AND id < " + p.query.Since
		} else {
			baseSQL += " AND id > " + p.query.Since
		}
	}

	if p.query.Desc {
		baseSQL += " ORDER BY id DESC"
	} else {
		baseSQL += " ORDER BY id"
	}

	baseSQL += " LIMIT " + strconv.Itoa(p.query.Limit)

	return baseSQL
}

func (p postsSQLGenerator) TreeSort() string {
	strID := strconv.FormatInt(int64(p.thread.ID), 10)
	baseSQL := ""

	baseSQL = "SELECT author, created, forum, id, isedited, message, parent, thread FROM posts WHERE thread = " + strID

	if p.query.Since != "" {
		if p.query.Desc {
			baseSQL += " AND path < (SELECT path FROM posts WHERE id = " + p.query.Since + ")"
		} else {
			baseSQL += " AND path > (SELECT path FROM posts WHERE id = " + p.query.Since + ")"
		}
	}

	if p.query.Desc {
		baseSQL += " ORDER BY path DESC, id DESC"
	} else {
		baseSQL += " ORDER BY path, id"
	}

	baseSQL += " LIMIT " + strconv.Itoa(p.query.Limit)

	return baseSQL
}

func (p postsSQLGenerator) ParentTreeSort() string {
	baseSQL := ""

	baseSQL = "SELECT author, created, forum, id, isedited, message, parent, thread FROM posts WHERE path[1]" +
		" IN (SELECT id FROM posts WHERE thread = " + strconv.FormatInt(int64(p.thread.ID), 10) +
		" AND parent = 0"

	if p.query.Since != "" {
		if p.query.Desc {
			baseSQL += " AND path[1] < (SELECT path[1] FROM posts WHERE id = " + p.query.Since + ")"
		} else {
			baseSQL += " AND path[1] > (SELECT path[1] FROM posts WHERE id = " + p.query.Since + ")"
		}
	}

	if p.query.Desc {
		baseSQL += " ORDER BY id DESC"
	} else {
		baseSQL += " ORDER BY id"
	}

	baseSQL += " LIMIT " + strconv.Itoa(p.query.Limit) + ")"

	if p.query.Desc {
		baseSQL += " ORDER BY path[1] DESC, path, id"
	} else {
		baseSQL += " ORDER BY path"
	}

	return baseSQL
}

func NewPostsSQLGenerator(thread models.Thread, query models.PostsRequestQuery) PostsSQLGenerator {
	return &postsSQLGenerator{
		thread: thread,
		query:  query,
	}
}
