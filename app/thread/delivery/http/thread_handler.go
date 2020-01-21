package http

import (
	"github.com/artbakulev/techdb/app/forum"
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/thread"
	"github.com/artbakulev/techdb/app/vote"
	"github.com/artbakulev/techdb/pkg/queryWorker"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"strconv"
)

type ThreadHandler struct {
	usecase      thread.Usecase
	forumUsecase forum.Usecase
	voteUsecase  vote.Usecase
}

func NewThreadHandler(router *fasthttprouter.Router, usecase thread.Usecase,
	forumUsecase forum.Usecase, voteUsecase vote.Usecase) {
	handler := &ThreadHandler{
		usecase:      usecase,
		forumUsecase: forumUsecase,
		voteUsecase:  voteUsecase,
	}

	router.POST("/api/forum/:slug", handler.CreateForum)
	router.POST("/api/forum/:slug/create", handler.CreateThread)
	router.GET("/api/forum/:slug/threads", handler.GetThreads)
	router.POST("/api/thread/:slug_or_id/vote", handler.CreateVote)
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
	if e != nil && e.StatusCode == 409 {
		ctx.SetStatusCode(409)
	} else if e != nil {
		e.SetToContext(ctx)
		return
	} else {
		ctx.SetStatusCode(201)
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
	createdForum, e := f.forumUsecase.CreateForum(buffer)
	if e != nil && e.StatusCode == 404 {
		e.SetToContext(ctx)
		return
	}

	if e != nil && e.StatusCode == 409 {
		createdForum, e = f.forumUsecase.GetForumBySlug(buffer.Slug)
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

func (h ThreadHandler) CreateVote(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	id, _ := strconv.ParseInt(slugOrID, 10, 64)
	createdVote := models.Vote{}

	if id == 0 {
		createdVote.Thread = -1
		createdVote.ThreadSlug = slugOrID
	} else {
		createdVote.Thread = id
	}

	err := createdVote.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}

	existingThread, e := h.voteUsecase.UpsertVote(createdVote)

	if e != nil {
		e.SetToContext(ctx)
		return
	}
	jsonBlob, err := existingThread.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}

	ctx.SetBody(jsonBlob)
}
