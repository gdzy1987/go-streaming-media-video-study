package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go-streaming-media-video-study/config"
	"go-streaming-media-video-study/logger"
	"go-streaming-media-video-study/web/handler"
	"log"
	"net/http"
	"os"
	"strings"
)

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", handler.HomeHandler)
	router.POST("/", handler.HomeHandler)

	router.GET("/userhome", handler.UserHomeHandler)
	router.POST("/userhome", handler.UserHomeHandler)

	router.POST("/api", handler.ApiHandler)

	router.GET("/videos/:vid-id", handler.ProxyVideoHandler)
	router.POST("/upload/:vid-id", handler.ProxyUploadHandler)

	router.ServeFiles("/statics/*filepath", http.Dir("./templates"))

	return router
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
	file, _ := os.Create(strings.Join([]string{path, "web_log.txt"}, "/"))
	defer file.Close()
	loger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetDefault(loger)

	r := RegisterHandler()

	addr := config.DefaultConfig.Address + ":" + config.DefaultConfig.WebPort
	fmt.Println("handler init...Port:\t", addr)
	http.ListenAndServe(addr, r)
}
