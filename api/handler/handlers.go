package handler

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go-streaming-media-video-study/api/model"
	"go-streaming-media-video-study/api/session"
	"go-streaming-media-video-study/api/utils"
	"go-streaming-media-video-study/common"
	"go-streaming-media-video-study/logger"
	"io/ioutil"
	"net/http"
)

// 创建用户
func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &model.UserCredential{}

	if err := json.Unmarshal(res, ubody); err != nil {
		common.SendErrorResponses(w, common.ErrorRequestBodyParseFailed)
		return
	}

	if err := model.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		common.SendErrorResponses(w, common.ErrorDBError)
		return
	}

	id := session.GenerateNewSessionId(ubody.Username)
	su := &model.SignedUp{
		Success:   true,
		SessionId: id,
	}

	if resp, err := json.Marshal(su); err != nil {
		common.SendErrorResponses(w, common.ErrorInternalFaults)
	} else {
		common.SendNormalResponse(w, string(resp), 200)
	}
}

// 用户登陆
func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	logger.Infof("login res:%s\n", res)
	ubody := &model.UserCredential{}

	if err := json.Unmarshal(res, ubody); err != nil {
		logger.Info("login json unmarshal err:\t", err)
		common.SendErrorResponses(w, common.ErrorRequestBodyParseFailed)
		return
	}

	// 验证
	uName := p.ByName("username")
	if uName != ubody.Username {
		common.SendErrorResponses(w, common.ErrorNotAuthUser)
		return
	}

	pwd, err := model.GetUserCredential(ubody.Username)
	if err != nil || len(pwd) == 0 || pwd != ubody.Pwd {
		common.SendErrorResponses(w, common.ErrorNotAuthUser)
		return
	}

	id := session.GenerateNewSessionId(ubody.Username)
	si := &model.SignedIn{
		Success:   true,
		SessionId: id,
	}
	if resp, err := json.Marshal(si); err != nil {
		common.SendErrorResponses(w, common.ErrorInternalFaults)
	} else {
		common.SendNormalResponse(w, string(resp), 200)
	}
}

// 获取用户信息
func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !session.ValidateUser(w, r) {
		logger.Info("get user info unauthorized user")
		return
	}

	uname := p.ByName("username")
	u, err := model.GetUser(uname)
	if err != nil {
		logger.Info("get user info err :\t", err)
		common.SendErrorResponses(w, common.ErrorDBError)
		return
	}

	ui := &model.UserInfo{Id: u.Id}
	if resp, err := json.Marshal(ui); err != nil {
		common.SendErrorResponses(w, common.ErrorInternalFaults)
	} else {
		common.SendNormalResponse(w, string(resp), 200)
	}
}

// 添加新的视频
func AddNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !session.ValidateUser(w, r) {
		logger.Info("add new video unauthorized user")
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	newVBody := &model.NewVideo{}
	if err := json.Unmarshal(res, newVBody); err != nil {
		logger.Info("add new video err:\t", err)
		common.SendErrorResponses(w, common.ErrorRequestBodyParseFailed)
		return
	}

	vi, err := model.AddNewVideo(newVBody.AuthorId, newVBody.Name)
	if err != nil {
		logger.Info("add new video err:\t", err)
		common.SendErrorResponses(w, common.ErrorDBError)
		return
	}

	if resp, err := json.Marshal(vi); err != nil {
		common.SendErrorResponses(w, common.ErrorInternalFaults)
	} else {
		common.SendNormalResponse(w, string(resp), 201)
	}
}

// 列出所有Video
func ListAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !session.ValidateUser(w, r) {
		logger.Info("list all videos unauthorized user")
		return
	}

	uname := p.ByName("username")
	vs, err := model.ListVideoInfo(uname, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		logger.Info("list all videos err:\t", err)
		common.SendErrorResponses(w, common.ErrorDBError)
		return
	}

	vsi := &model.VideosInfo{Videos: vs}
	if resp, err := json.Marshal(vsi); err != nil {
		common.SendErrorResponses(w, common.ErrorInternalFaults)
	} else {
		common.SendNormalResponse(w, string(resp), 200)
	}
}

// 删除video
func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !session.ValidateUser(w, r) {
		logger.Info("delete video unauthorized user")
		return
	}

	vid := p.ByName("vid-id")
	err := model.DeleteVideoInfo(vid)
	if err != nil {
		logger.Info("delete video err:\t", err)
		common.SendErrorResponses(w, common.ErrorDBError)
		return
	}

	go utils.SendDeleteVideoRequest(vid)
	common.SendNormalResponse(w, "删除成功", 204)
}

// 发送评论
func PostComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !session.ValidateUser(w, r) {
		logger.Info("post comment unauthorized user")
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)

	cBody := &model.NewComment{}
	if err := json.Unmarshal(reqBody, cBody); err != nil {
		logger.Info("post comment err:\t", err)
		common.SendErrorResponses(w, common.ErrorRequestBodyParseFailed)
		return
	}

	vid := p.ByName("vid-id")
	if err := model.AddNewComments(vid, cBody.AuthorId, cBody.Content); err != nil {
		logger.Info("post comment err:\t", err)
		common.SendErrorResponses(w, common.ErrorDBError)
	} else {
		common.SendNormalResponse(w, "ok", 201)
	}
}

// 展示评论
func ShowComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !session.ValidateUser(w, r) {
		logger.Info("show comments unauthorized user")
		return
	}

	vid := p.ByName("vid-id")
	cm, err := model.ListComments(vid, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		logger.Info("show comments err:\t", err)
		common.SendErrorResponses(w, common.ErrorDBError)
		return
	}

	cms := &model.Comments{Comments: cm}
	if resp, err := json.Marshal(cms); err != nil {
		common.SendErrorResponses(w, common.ErrorInternalFaults)
	} else {
		common.SendNormalResponse(w, string(resp), 200)
	}
}
