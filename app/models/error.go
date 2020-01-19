package models

import (
	"github.com/valyala/fasthttp"
	"log"
)

const (
	CreateError     = "ошибка создания объекта"
	UpdateError     = "ошибка обновления объекта"
	DBParsingError  = "ошибка парсинга результатов"
	NotFoundError   = "объект не найден"
	InternalError   = "внутренняя ошибка"
	ConflictError   = "военный конфликт"
	BadRequestError = "невалидный запрос"
	DBError         = "конфикт данных в базе"
)

var (
	BadRequestErrorBytes = []byte("{\"message\": \"невалидный запрос\"")
	InternalErrorBytes   = []byte("{\"message\": \"внутренняя ошибка\"")
)

type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func NewError(statusCode int, message string, extraMessages ...string) *Error {
	if len(extraMessages) > 0 {
		for _, item := range extraMessages {
			message += ": " + item
		}
	}
	return &Error{
		StatusCode: statusCode,
		Message:    message,
	}
}

func (e Error) MessageToBytes() []byte {
	return append([]byte("{\"message\": \""), []byte(e.Message+"\"}")...)
}

func (e Error) Log() {
	log.Print("Error: ", e.StatusCode, e.Message)
}

func (e Error) SetToContext(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(e.StatusCode)
	ctx.SetBody(e.MessageToBytes())
	e.Log() //	TODO: убрать на проде
}
