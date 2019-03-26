package handler

import (
	"github.com/julienschmidt/httprouter"
	"go-streaming-media-video-study/common"
	"go-streaming-media-video-study/scheduler/model"
	"net/http"
)

func VidDelRecHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	if len(vid) == 0 {
		common.SendResponse(w, 400, "video id should not be empty")
		return
	}
	err := model.AddVideoDeletionRecord(vid)
	if err != nil {
		common.SendResponse(w, 500, "Internal server error")
		return
	}
	common.SendResponse(w, 200, "")
	return
}
