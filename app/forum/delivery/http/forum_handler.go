package http

import (
	"github.com/artbakulev/techdb/app/forum"
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/pkg/queryWorker"
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
	router.GET("/api/forum/:slug/details", handler.GetForumBySlug)
	router.GET("/api/forum/:slug/users", handler.GetForumUsers)
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

func (f ForumHandler) GetForumUsers(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	query := models.PostsRequestQuery{
		Limit: queryWorker.GetIntParam(ctx, "limit"),
		Since: queryWorker.GetStringParam(ctx, "since"),
		Desc:  queryWorker.GetBoolParam(ctx, "desc"),
	}
	users, err := f.usecase.GetForumUsers(slug, query)

	if err != nil {
		err.SetToContext(ctx)
		return
	}

	jsonBlob, e := users.MarshalJSON()
	if e != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)
}
