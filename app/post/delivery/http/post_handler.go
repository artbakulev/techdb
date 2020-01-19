package http

import (
	"github.com/artbakulev/techdb/app/post"
	"github.com/buaazp/fasthttprouter"
)

type PostHandler struct {
	usecase post.Usecase
}

func NewPostHandler(router *fasthttprouter.Router, usecase post.Usecase) {
	//handler := &PostHandler{
	//	usecase: usecase,
	//}
	//router.GET("/api/user/:nickname/profile", handler.GetUser)
}
