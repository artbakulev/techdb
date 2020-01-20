package queryWorker

import (
	"github.com/valyala/fasthttp"
	"strconv"
)

func GetIntParam(ctx *fasthttp.RequestCtx, key string) int {
	param, _ := strconv.Atoi(string(ctx.URI().QueryArgs().Peek(key)))
	return param
}

func GetInt64Param(ctx *fasthttp.RequestCtx, key string) int64 {
	param, _ := strconv.ParseInt(string(ctx.URI().QueryArgs().Peek(key)), 10, 64)
	return param
}

func GetStringParam(ctx *fasthttp.RequestCtx, key string) string {
	return string(ctx.URI().QueryArgs().Peek(key))
}

func GetBoolParam(ctx *fasthttp.RequestCtx, key string) bool {
	param, _ := strconv.ParseBool(string(ctx.URI().QueryArgs().Peek(key)))
	return param
}
