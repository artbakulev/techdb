package server

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type server struct {
	Host   string
	router *fasthttprouter.Router
}

func NewServer(host string, router *fasthttprouter.Router) *server {
	return &server{
		Host:   host,
		router: router,
	}
}

func (s server) ListenAndServe() error {
	return fasthttp.ListenAndServe(s.Host, DefaultHeaders(s.router.Handler))
}
