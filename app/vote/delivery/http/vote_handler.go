package http

import (
	"github.com/artbakulev/techdb/app/vote"
	"github.com/buaazp/fasthttprouter"
)

type VoteHandler struct {
	usecase vote.Usecase
}

func NewVoteHandler(router *fasthttprouter.Router, usecase vote.Usecase) {
	//handler := &VoteHandler{
	//	usecase: usecase,
	//}
	//router.GET("/api/user/:nickname/profile", handler.GetUser)
}
