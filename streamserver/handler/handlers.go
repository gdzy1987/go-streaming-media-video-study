package handler

import (
	"github.com/julienschmidt/httprouter"
	"go-streaming-media-video-study/common"
	"go-streaming-media-video-study/config"
	"go-streaming-media-video-study/logger"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func TestPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles("./upload.html")
	t.Execute(w, nil)
}

func StreamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	targetUrl := "http://" + strings.Trim(config.DefaultConfig.Bucket, "/") + "." + config.DefaultConfig.OssAddr + "/" + p.ByName("vid-id")
	logger.Info("Enter the stream Handler:\t", targetUrl)
	http.Redirect(w, r, targetUrl, 301)
	return
}

func UploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	logger.Info("stream server uploadHandler")
	r.Body = http.MaxBytesReader(w, r.Body, common.MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(common.MAX_UPLOAD_SIZE); err != nil {
		common.SendErrorResponse(w, http.StatusBadRequest, "File is too big")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		common.SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Info("uploadHandler read file error:\t", err)
		common.SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	fn := p.ByName("vid-id")
	err = ioutil.WriteFile(common.VIDEO_DIR+fn, data, 0666)
	if err != nil {
		logger.Info("uploadHandler write file error:\t", err)
		common.SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	// ossfn := config.DefaultConfig.Bucket + fn
	ossfn := fn
	path := common.VIDEO_DIR + fn
	ret := common.UploadToOss(ossfn, path, strings.Trim(config.DefaultConfig.Bucket, "/"))
	if !ret {
		common.SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	os.Remove(path)

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Uploaded successfully")
}
