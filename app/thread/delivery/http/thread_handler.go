package http

import (
	"github.com/artbakulev/techdb/app/thread"
	"github.com/buaazp/fasthttprouter"
)

type ThreadHandler struct {
	usecase thread.Usecase
}

func NewThreadHandler(router *fasthttprouter.Router, usecase thread.Usecase) {
	//handler := &ThreadHandler{
	//	usecase: usecase,
	//}
	//router.GET("/api/user/:nickname/profile", handler.GetUser)
}
