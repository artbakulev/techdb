package http

import (
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/post"
	"github.com/artbakulev/techdb/pkg/queryWorker"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"strconv"
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
	slugOrId := ctx.UserValue("slug_or_id").(string)
	id, _ := strconv.ParseInt(slugOrId, 10, 64)
	if id == 0 {
		id = -1
	}

	posts := models.Posts{}
	err := posts.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	createdPosts, e := p.usecase.CreatePosts(slugOrId, id, posts)

	if e != nil {
		e.SetToContext(ctx)
		return
	}

	jsonBlob, err := createdPosts.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)
}

func (p PostHandler) GetMany(ctx *fasthttp.RequestCtx) {

	slugOrId := ctx.UserValue("slug_or_id").(string)
	id, err := strconv.ParseInt(slugOrId, 10, 64)
	if err != nil {
		id = -1
	}

	query := models.PostsRequestQuery{
		ThreadID:   id,
		ThreadSlug: slugOrId,
	}

	query.Limit = queryWorker.GetIntParam(ctx, "limit")
	query.Since = queryWorker.GetStringParam(ctx, "since")
	query.Sort = queryWorker.GetStringParam(ctx, "sort")
	query.Desc = queryWorker.GetBoolParam(ctx, "desc")

	sortedPosts, e := p.usecase.GetThreadPosts(query)
	if err != nil {
		e.SetToContext(ctx)
		return
	}

	jsonBlob, err := sortedPosts.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)

}

func (p PostHandler) GetOne(ctx *fasthttp.RequestCtx) {
	var id int64 = -1
	id = ctx.UserValue("id").(int64)
	if id == -1 {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	queryParams := strings.Split(string(ctx.URI().QueryArgs().Peek("related")), ",")

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
