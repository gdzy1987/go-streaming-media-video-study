package model

// 用户
// request
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

type UserInfo struct {
	Id int `json:"id"`
}

type User struct {
	Id        int
	LoginName string
	Pwd       string
}

// response
type SignedUp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

type SignedIn struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

// session
type UserSession struct {
	Username  string `json:"user_name"`
	SessionId string `json:"session_id"`
}
type SimpleSession struct {
	Username string `json:"username"`
	TTL      int64  `json:"ttl"`
}

// video
type NewVideo struct {
	AuthorId int    `json:"author_id"`
	Name     string `json:"name"`
}

type VideosInfo struct {
	Videos []*VideoInfo `json:"videos"`
}

type VideoInfo struct {
	Id           string `json:"id"`
	AuthorId     int    `json:"author_id"`
	Name         string `json:"name"`
	DisplayCtime string `json:"display_ctime"`
}

// comment
type NewComment struct {
	AuthorId int    `json:"author_id"`
	Content  string `json:"content"`
}

type Comments struct {
	Comments []*Comment `json:"comments"`
}

type Comment struct {
	Id      string `json:"id"`
	VideoId string `json:"video_id"`
	Author  string `json:"author"`
	Content string `json:"content"`
}
