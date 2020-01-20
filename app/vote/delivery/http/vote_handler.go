package http

import (
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/vote"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"strconv"
)

type VoteHandler struct {
	usecase vote.Usecase
}

func NewVoteHandler(router *fasthttprouter.Router, usecase vote.Usecase) {
	handler := &VoteHandler{
		usecase: usecase,
	}
	router.POST("/thread/:slug_or_id/vote", handler.CreateVote)
}

func (h VoteHandler) CreateVote(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	id, _ := strconv.ParseInt(slugOrID, 10, 64)
	createdVote := models.Vote{}

	if id == 0 {
		id = -1
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

	existingThread, e := h.usecase.UpsertVote(createdVote)

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
