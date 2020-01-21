package server

import (
	http2 "github.com/artbakulev/techdb/app/forum/delivery/http"
	repository2 "github.com/artbakulev/techdb/app/forum/repository"
	usecase2 "github.com/artbakulev/techdb/app/forum/usecase"
	http5 "github.com/artbakulev/techdb/app/post/delivery/http"
	repository5 "github.com/artbakulev/techdb/app/post/repository"
	usecase5 "github.com/artbakulev/techdb/app/post/usecase"
	http3 "github.com/artbakulev/techdb/app/service/delivery/http"
	repository3 "github.com/artbakulev/techdb/app/service/repository"
	usecase3 "github.com/artbakulev/techdb/app/service/usecase"
	http4 "github.com/artbakulev/techdb/app/thread/delivery/http"
	repository4 "github.com/artbakulev/techdb/app/thread/repository"
	usecase4 "github.com/artbakulev/techdb/app/thread/usecase"
	"github.com/artbakulev/techdb/app/user/delivery/http"
	"github.com/artbakulev/techdb/app/user/repository"
	"github.com/artbakulev/techdb/app/user/usecase"
	repository6 "github.com/artbakulev/techdb/app/vote/repository"
	usecase6 "github.com/artbakulev/techdb/app/vote/usecase"
	"github.com/buaazp/fasthttprouter"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
)

type server struct {
	Host   string
	router *fasthttprouter.Router
}

func NewServer(host string, connection *pgx.ConnPool) *server {
	userRepo := repository.NewPostgresUserRepository(connection)
	forumRepo := repository2.NewPostgresForumRepository(connection)
	serviceRepo := repository3.NewPostgresServiceRepository(connection)
	threadRepo := repository4.NewPostgresThreadRepository(connection)
	postRepo := repository5.NewPostgresPostRepository(connection)
	voteRepo := repository6.NewPostgresVoteRepository(connection)

	userUsecase := usecase.NewUserUsecase(userRepo)
	forumUsecase := usecase2.NewForumUsecase(userRepo, forumRepo)
	serviceUsecase := usecase3.NewServiceUsecase(serviceRepo)
	threadUsecase := usecase4.NewThreadUsecase(threadRepo, userRepo, forumRepo)
	postUsecase := usecase5.NewPostUsecase(userRepo, postRepo, threadRepo, forumRepo)
	voteUsecase := usecase6.NewVoteUsecase(voteRepo, threadRepo)

	router := fasthttprouter.New()

	http.NewUserHandler(router, userUsecase)
	http2.NewForumHandler(router, forumUsecase)
	http3.NewServiceHandler(router, serviceUsecase)
	http4.NewThreadHandler(router, threadUsecase, forumUsecase, voteUsecase)
	http5.NewPostHandler(router, postUsecase)

	return &server{
		Host:   host,
		router: router,
	}
}

func (s server) ListenAndServe() error {
	return fasthttp.ListenAndServe(s.Host, DefaultHeaders(s.router.Handler))
}
