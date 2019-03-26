package session

import (
	"go-streaming-media-video-study/common"
	"net/http"
)

var (
	HEADER_FIELD_SESSION = "X-Session-Id"
	HEADER_FIELD_UNAME   = "X-User-Name"
)

// 通过session id检查权限
func ValidateUserSession(r *http.Request) bool {
	sid := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		return false
	}

	uname, ok := IsSessionExpired(sid)
	if ok {
		return false
	}

	r.Header.Add(HEADER_FIELD_UNAME, uname)
	return true
}

// 验证用户
func ValidateUser(w http.ResponseWriter, r *http.Request) bool {
	uName := r.Header.Get(HEADER_FIELD_UNAME)
	if len(uName) == 0 {
		common.SendErrorResponses(w, common.ErrorNotAuthUser)
		return false
	}
	return true
}
