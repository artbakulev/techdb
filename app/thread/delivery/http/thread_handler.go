package http

import (
	"github.com/artbakulev/techdb/app/forum"
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/thread"
	"github.com/artbakulev/techdb/pkg/queryWorker"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
)

type ThreadHandler struct {
	usecase      thread.Usecase
	extraUsecase forum.Usecase //	из-за проблемы в fasthttprouter
}

func NewThreadHandler(router *fasthttprouter.Router, usecase thread.Usecase, extraUsecase forum.Usecase) {
	handler := &ThreadHandler{
		usecase:      usecase,
		extraUsecase: extraUsecase,
	}

	//router.POST("/api/forum/:slug/create", handler.GetHandler)
	router.POST("/api/forum/:slug", handler.CreateForum)
	router.POST("/api/forum/:slug/create", handler.CreateThread)
	router.GET("/api/forum/:slug/threads", handler.GetThreads)
	router.POST("/api/thread/:slug_or_id/details", handler.UpdateThread)
	router.GET("/api/thread/:slug_or_id/details", handler.GetThread)
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

	createdThread, e := h.usecase.CreateThread(slug, createdThread)
	log.Printf("%v", e)

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

	ctx.SetStatusCode(201)
	ctx.SetBody(jsonBlob)
}

func (h ThreadHandler) GetThreads(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	query := models.PostsRequestQuery{
		Limit: queryWorker.GetIntParam(ctx, "limit"),
		Since: queryWorker.GetStringParam(ctx, "since"),
		Desc:  queryWorker.GetBoolParam(ctx, "desc"),
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

func (h ThreadHandler) GetThread(ctx *fasthttp.RequestCtx) {
	slugOrId := ctx.UserValue("slug_or_id").(string)
	id, _ := strconv.ParseInt(slugOrId, 10, 64)
	isSlug := false
	if id == 0 {
		isSlug = true
	}
	existingThread, err := h.usecase.GetThreadBySlugOrID(slugOrId, isSlug)

	if err != nil {
		err.SetToContext(ctx)
		return
	}

	jsonBlob, e := existingThread.MarshalJSON()
	if e != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)

}

func (f ThreadHandler) CreateForum(ctx *fasthttp.RequestCtx) {
	var buffer models.Forum
	body := ctx.PostBody()
	err := buffer.UnmarshalJSON(body)
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	createdForum, e := f.extraUsecase.CreateForum(buffer)
	if e != nil {
		createdForum, e = f.extraUsecase.GetForumBySlug(buffer.Slug)
		if e != nil {
			e.SetToContext(ctx)
			return
		}
		ctx.SetStatusCode(409)

	} else {
		ctx.SetStatusCode(201)
	}
	jsonBlob, err := createdForum.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}
	ctx.SetBody(jsonBlob)
}
