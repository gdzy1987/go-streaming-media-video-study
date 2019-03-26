package session

import (
	"go-streaming-media-video-study/api/model"
	"go-streaming-media-video-study/api/utils"
	"go-streaming-media-video-study/logger"
	"sync"
	"time"
)

// 使用*sync.Map作为cache
var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

func nowInMilli() int64 {
	return time.Now().UnixNano() / 1000000
}

// 删除session
func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	DeleteSession(sid)
}

// 加载session
func LoadSessionsFromDB() {
	r, err := RetrieveAllSessions()
	if err != nil {
		logger.Info("load sessions from db err:\t", err)
		return
	}

	r.Range(func(key, value interface{}) bool {
		ss := value.(*model.SimpleSession)
		sessionMap.Store(key, ss)
		return true
	})
}

// 生成一个新的session id,并存储
func GenerateNewSessionId(un string) string {
	id, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 30*60*1000

	ss := &model.SimpleSession{
		Username: un,
		TTL:      ttl,
	}
	sessionMap.Store(id, ss)
	InsertSession(id, ttl, un)
	return id
}

//  判断session是否过期
func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	ct := nowInMilli()
	if ok {
		// 从sessionMap获取到session
		if ss.(*model.SimpleSession).TTL < ct {
			// 已经过期
			deleteExpiredSession(sid)
			return "", true
		}
		return ss.(*model.SimpleSession).Username, false
	} else {
		// 未从sessionMap获取到session, 从数据库获取session
		ss, err := RetrieveSession(sid)
		if err != nil || ss == nil {
			return "", true
		}

		if ss.TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}
		// 未过期，存储到sessionMap里
		sessionMap.Store(sid, ss)
		return ss.Username, false
	}
}
