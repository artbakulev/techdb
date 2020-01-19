package http

import (
	"github.com/artbakulev/techdb/app/forum"
	"github.com/artbakulev/techdb/app/models"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type ForumHandler struct {
	usecase forum.Usecase
}

func NewForumHandler(router *fasthttprouter.Router, usecase forum.Usecase) {
	handler := &ForumHandler{
		usecase: usecase,
	}
	router.POST("/api/forum/create", handler.CreateForum)
	router.GET("/api/forum/:slug/details", handler.GetForumBySlug)
}

func (f ForumHandler) CreateForum(ctx *fasthttp.RequestCtx) {
	var buffer models.Forum
	body := ctx.PostBody()
	err := buffer.UnmarshalJSON(body)
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	createdForum, e := f.usecase.CreateForum(buffer)
	if e != nil {
		e.SetToContext(ctx)
		return
	}
	jsonBlob, err := createdForum.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}
	ctx.SetBody(jsonBlob)
}

func (f ForumHandler) GetForumBySlug(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	if slug == "" {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
	}

	existingForum, e := f.usecase.GetForumBySlug(slug)
	if e != nil {
		e.SetToContext(ctx)
		return
	}

	jsonBlob, err := existingForum.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)
}
