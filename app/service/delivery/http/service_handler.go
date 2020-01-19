package http

import (
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/service"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type ServiceHandler struct {
	usecase service.Usecase
}

func NewServiceHandler(router *fasthttprouter.Router, usecase service.Usecase) {
	handler := &ServiceHandler{
		usecase: usecase,
	}
	router.POST("/api/service/clear", handler.Clear)
	router.GET("/api/service/status", handler.GetStatus)
}

func (s ServiceHandler) Clear(ctx *fasthttp.RequestCtx) {
	err := s.usecase.ClearDB()
	if err != nil {
		err.SetToContext(ctx)
	}
}
func (s ServiceHandler) GetStatus(ctx *fasthttp.RequestCtx) {
	dbStatus, e := s.usecase.GetDBStatus()
	if e != nil {
		e.SetToContext(ctx)
		return
	}

	jsonBlob, err := dbStatus.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
	}

	ctx.SetBody(jsonBlob)
}
