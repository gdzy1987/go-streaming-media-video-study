package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go-streaming-media-video-study/common"
	"go-streaming-media-video-study/config"
	"go-streaming-media-video-study/logger"
	"go-streaming-media-video-study/streamserver/handler"
	"go-streaming-media-video-study/streamserver/utils"
	"log"
	"net/http"
	"os"
	"strings"
)

type middleWareHandler struct {
	r *httprouter.Router
	l *utils.ConnLimiter
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.l = utils.NewConnLimiter(cc)
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.l.GetConn() {
		common.SendErrorResponse(w, http.StatusTooManyRequests, "Too many requests")
		return
	}
	m.r.ServeHTTP(w, r)
	defer m.l.ReleaseConn()
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/videos/:vid-id", handler.StreamHandler)
	router.POST("/upload/:vid-id", handler.UploadHandler)
	router.GET("/testpage", handler.TestPageHandler)

	return router
}

func main() {
	// 初始化配置
	fmt.Println("config init...")
	config.InitConfig("./conf.json")
	fmt.Printf("%+v\n", config.DefaultConfig)

	// 日志配置
	fmt.Println("logger init...")
	path := "logs"
	mode := os.ModePerm
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, mode)
	}
	file, _ := os.Create(strings.Join([]string{path, "streamserver_log.txt"}, "/"))
	defer file.Close()
	loger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetDefault(loger)

	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, 2)
	addr := config.DefaultConfig.Address + ":" + config.DefaultConfig.StreamServerPort
	fmt.Println("streamServer start...\t", addr)
	http.ListenAndServe(addr, mh)
}
