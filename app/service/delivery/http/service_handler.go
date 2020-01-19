package http

import (
	"github.com/artbakulev/techdb/app/service"
	"github.com/buaazp/fasthttprouter"
)

type ServiceHandler struct {
	usecase service.Usecase
}

func NewServiceHandler(router *fasthttprouter.Router, usecase service.Usecase) {
	//handler := &ServiceHandler{
	//	usecase: usecase,
	//}
	//router.GET("/api/user/:nickname/profile", handler.GetUser)
}
