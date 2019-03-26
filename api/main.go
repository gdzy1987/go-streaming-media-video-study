package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go-streaming-media-video-study/api/handler"
	"go-streaming-media-video-study/api/model"
	"go-streaming-media-video-study/api/session"
	"go-streaming-media-video-study/config"
	"go-streaming-media-video-study/logger"
	"log"
	"net/http"
	"os"
	"strings"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

// 实现http.Handler接口中的ServeHttp接口
func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session.ValidateUserSession(r)
	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	// user
	router.POST("/user", handler.CreateUser)
	router.POST("/user/:username", handler.Login)
	router.GET("/user/:username", handler.GetUserInfo)

	// videos
	router.POST("/user/:username/videos", handler.AddNewVideo)
	router.GET("/user/:username/videos", handler.ListAllVideos)
	router.DELETE("/user/:username/videos/:vid-id", handler.DeleteVideo)

	// comments
	router.POST("/videos/:vid-id/comments", handler.PostComment)
	router.GET("/videos/:vid-id/comments", handler.ShowComments)

	return router
}

func InitSession() {
	session.LoadSessionsFromDB()
}

func main() {
	// 初始化配置
	config.InitConfig("./conf.json")
	fmt.Printf("%+v\n", config.DefaultConfig)

	// 日志配置
	fmt.Println("logger init...")
	path := "logs"
	mode := os.ModePerm
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, mode)
	}
	file, _ := os.Create(strings.Join([]string{path, "api_log.txt"}, "/"))
	defer file.Close()
	loger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetDefault(loger)

	// 初始化DB
	fmt.Println("mysql init...")
	model.InitMysql()

	// 初始化session
	fmt.Println("session init...")
	InitSession()

	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r)
	addr := config.DefaultConfig.Address + ":" + config.DefaultConfig.ApiPort

	fmt.Println("handler init...Port:\t", addr)
	http.ListenAndServe(addr, mh)
}
