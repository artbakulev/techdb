package http

import (
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/post"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"strings"
)

type PostHandler struct {
	usecase post.Usecase
}

func NewPostHandler(router *fasthttprouter.Router, usecase post.Usecase) {
	handler := &PostHandler{
		usecase: usecase,
	}

	router.GET("/api/post/:id/details", handler.GetOne)
	router.POST("/api/post/:id/details", handler.UpdatePost)
	router.POST("/api/thread/:slug_or_id/create", handler.CreatePosts)
	router.GET("/api/thread/:slug_or_id/posts", handler.GetMany)
}

func (p PostHandler) UpdatePost(ctx *fasthttp.RequestCtx) {
	var id int64 = -1
	id = ctx.UserValue("id").(int64)
	if id == -1 {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}

	var buffer models.PostUpdate
	body := ctx.PostBody()
	err := buffer.UnmarshalJSON(body)
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}

	updatedPost, e := p.usecase.UpdatePost(id, buffer)
	if e != nil {
		e.SetToContext(ctx)
		return
	}

	jsonBlob, err := updatedPost.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)
}

func (p PostHandler) CreatePosts(ctx *fasthttp.RequestCtx) {
	//	TODO: доделай завтра
}

func (p PostHandler) GetMany(ctx *fasthttp.RequestCtx) {}

func (p PostHandler) GetOne(ctx *fasthttp.RequestCtx) {
	var id int64 = -1
	id = ctx.UserValue("id").(int64)
	if id == -1 {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	queryParams := strings.Split(ctx.URI().String(), ",")

	var query models.PostsRelatedQuery

	for _, param := range queryParams {
		if param == "user" {
			query.NeedAuthor = true
		} else if param == "forum" {
			query.NeedForum = true
		} else if param == "thread" {
			query.NeedThread = true
		}
	}

	existingPost, err := p.usecase.GetPostDetails(id, query)
	if err != nil {
		err.SetToContext(ctx)
		return
	}

	jsonBlob, e := existingPost.MarshalJSON()
	if e != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)
}
