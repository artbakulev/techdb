package http

import (
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/user"
	"github.com/artbakulev/techdb/pkg/validation"
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
	foundUser, err := u.usecase.GetUserByNickname(nickname.(string))
	if err != nil {
		err.SetToContext(ctx)
		return
	}
	jsonBlob, _ := foundUser.MarshalJSON()
	ctx.SetBody(jsonBlob)
}

func (u UserHandler) UpdateUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	e := validation.ValidateNickname(nickname)
	if e != nil {
		e.SetToContext(ctx)
		return
	}
	buffer := models.User{}
	body := ctx.PostBody()
	err := buffer.UnmarshalJSON(body)
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	buffer.Nickname = nickname
	returnUser, e := u.usecase.UpdateUser(buffer)
	if e != nil {
		e.SetToContext(ctx)
		return
	}
	jsonBlob, err := returnUser.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody(models.InternalErrorBytes)
		return
	}
	ctx.SetBody(jsonBlob)
}

func (u UserHandler) CreateUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	e := validation.ValidateNickname(nickname)
	if e != nil {
		e.SetToContext(ctx)
		return
	}
	buffer := models.User{}
	body := ctx.PostBody()
	err := buffer.UnmarshalJSON(body)
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(models.BadRequestErrorBytes)
		return
	}
	buffer.Nickname = nickname
	returnUser, e := u.usecase.CreateUser(buffer)
	if e != nil {
		foundUsers := models.Users{}
		ctx.SetStatusCode(409)
		returnUserByNickname, e := u.usecase.GetUserByNickname(nickname)
		if e == nil {
			foundUsers = append(foundUsers, returnUserByNickname)
		}
		returnUserByEmail, e := u.usecase.GetUserByEmail(buffer.Email)
		if e == nil && returnUserByNickname.Email != returnUserByEmail.Email {
			foundUsers = append(foundUsers, returnUserByEmail)
		}
		jsonBlob, err := foundUsers.MarshalJSON()
		if err != nil {
			ctx.SetStatusCode(500)
			ctx.SetBody(models.InternalErrorBytes)
			return
		}
		ctx.SetBody(jsonBlob)
	} else {
		ctx.SetStatusCode(201)
		jsonBlob, err := returnUser.MarshalJSON()
		if err != nil {
			ctx.SetStatusCode(500)
			ctx.SetBody(models.InternalErrorBytes)
			return
		}
		ctx.SetBody(jsonBlob)
	}

}
