package http

import (
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/thread"
	"github.com/artbakulev/techdb/pkg/queryWorker"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"strconv"
)

type ThreadHandler struct {
	usecase thread.Usecase
}

func NewThreadHandler(router *fasthttprouter.Router, usecase thread.Usecase) {
	handler := &ThreadHandler{
		usecase: usecase,
	}

	router.POST("/forum/:slug/create", handler.CreateThread)
	router.GET("/forum/:slug/threads", handler.GetThreads)
	router.POST("/thread/:slug_or_id/details", handler.UpdateThread)
}

func (h ThreadHandler) CreateThread(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	createdThread := models.Thread{}

	err := createdThread.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}

	createdThread.Slug = slug

	createdThread, e := h.usecase.CreateThread(slug, createdThread)

	if e != nil {
		e.SetToContext(ctx)
		return
	}

	jsonBlob, err := createdThread.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)
}

func (h ThreadHandler) GetThreads(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	query := models.PostsRequestQuery{
		ThreadID:   0,
		ThreadSlug: "",
		Limit:      queryWorker.GetIntParam(ctx, "limit"),
		Since:      queryWorker.GetInt64Param(ctx, "since"),
		Sort:       "",
		Desc:       queryWorker.GetBoolParam(ctx, "desc"),
	}

	threads, err := h.usecase.GetThreads(slug, query)

	if err != nil {
		err.SetToContext(ctx)
		return
	}

	jsonBlob, e := threads.MarshalJSON()
	if e != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)
}

func (h ThreadHandler) UpdateThread(ctx *fasthttp.RequestCtx) {
	updateThread := models.ThreadUpdate{}

	slugOrId := ctx.UserValue("slug_or_id").(string)
	id, _ := strconv.ParseInt(slugOrId, 10, 64)
	if id == 0 {
		id = -1
		updateThread.Slug = slugOrId
	}
	updateThread.ID = id

	err := updateThread.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}

	updatedThread, e := h.usecase.UpdateThread(updateThread)

	if e != nil {
		e.SetToContext(ctx)
		return
	}

	jsonBlob, err := updatedThread.MarshalJSON()

	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)
}
