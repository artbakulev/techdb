package http

import (
	"github.com/artbakulev/techdb/app/user"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	usecase user.Usecase
}

func NewUserHandler(router *fasthttprouter.Router, usecase user.Usecase) {
	handler := &UserHandler{
		usecase: usecase,
	}
	router.GET("/api/user/:nickname/profile", handler.GetUser)
	router.POST("/api/user/:nickname/create", handler.CreateUser)
	router.POST("/api/user/:nickname/profile", handler.UpdateUser)
}

func (u UserHandler) GetUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname")
	foundUser, err := u.usecase.GetUserByEmailOrByNickname(nickname.(string), false)
	if err != nil {
		ctx.SetBody(err.MessageToBytes())
		ctx.SetStatusCode(err.StatusCode)
		return
	}
	jsonBlob, _ := foundUser.MarshalJSON()
	ctx.SetBody(jsonBlob)
}

func (u UserHandler) UpdateUser(ctx *fasthttp.RequestCtx) {}
func (u UserHandler) CreateUser(ctx *fasthttp.RequestCtx) {}
