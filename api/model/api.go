package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"go-streaming-media-video-study/api/utils"
	"go-streaming-media-video-study/logger"
	"time"
)

// 插入用户loginName+password
func AddUserCredential(loginName string, pwd string) error {
	stmtIns, err := DbConn.Prepare("insert into users (login_name,pwd) values (?,?)")
	if err != nil {
		logger.Info("add user credential err:\t", err)
		return err
	}

	if _, err = stmtIns.Exec(loginName, pwd); err != nil {
		logger.Info("add user credential err:\t", err)
		return err
	}
	defer stmtIns.Close()
	return nil
}

// 根据用户名查询密码
func GetUserCredential(loginName string) (string, error) {
	stmtOut, err := DbConn.Prepare("select pwd from users where login_name = ?")
	if err != nil {
		logger.Info("get user credential err:\t", loginName, err)
		return "", err
	}

	var pwd string
	if err = stmtOut.QueryRow(loginName).Scan(&pwd); err != nil && err != sql.ErrNoRows {
		logger.Info("get user credential err:\t", loginName, err)
		return "", err
	}
	defer stmtOut.Close()
	return pwd, nil
}

// 删除用户
func DeleteUser(loginName string, pwd string) error {
	stmtDel, err := DbConn.Prepare("delete from users where login_name=? and pwd=?")
	if err != nil {
		logger.Info("delete user err:\t", loginName, err)
		return err
	}

	if _, err = stmtDel.Exec(loginName, pwd); err != nil {
		logger.Info("delete user err:\t", loginName, err)
		return err
	}

	defer stmtDel.Close()
	return nil
}

// 根据用户名查询id和pwd
func GetUser(loginName string) (*User, error) {
	stmtOut, err := DbConn.Prepare("select id,pwd from users where login_name =?")
	if err != nil {
		logger.Info("get user err:\t", loginName, err)
		return nil, err
	}

	var id int
	var pwd string

	if err = stmtOut.QueryRow(loginName).Scan(&id, &pwd); err != nil && err != sql.ErrNoRows {
		logger.Info("get user err:\t", loginName, err)
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	res := &User{
		Id:        id,
		LoginName: loginName,
		Pwd:       pwd,
	}
	defer stmtOut.Close()
	return res, nil
}

// 添加新的Video
func AddNewVideo(aid int, name string) (*VideoInfo, error) {
	// create uuid
	vid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}

	ctime := time.Now().Format("Jan 02 2006, 15:04:05")
	stmtIns, err := DbConn.Prepare(`insert into video_info(id, author_id, name, display_ctime) values (?,?,?,?)`)
	if err != nil {
		logger.Info("add new video err:\t", err)
		return nil, err
	}

	if _, err = stmtIns.Exec(vid, aid, name, ctime); err != nil {
		logger.Info("add new video err:\t", err)
		return nil, err
	}

	res := &VideoInfo{
		Id:           vid,
		AuthorId:     aid,
		Name:         name,
		DisplayCtime: ctime,
	}
	defer stmtIns.Close()
	return res, nil
}

// 获取video
func GetVideoInfo(vid string) (*VideoInfo, error) {
	stmtOut, err := DbConn.Prepare("select author_id, name, display_ctime from video_info where id=?")

	var aid int
	var dct string
	var name string

	if err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct); err != nil && err != sql.ErrNoRows {
		logger.Info("get video info err:\t", vid, err)
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	defer stmtOut.Close()

	res := &VideoInfo{
		Id:           vid,
		AuthorId:     aid,
		Name:         name,
		DisplayCtime: dct,
	}
	return res, nil
}

// 查询from->to时间之内的video信息
func ListVideoInfo(uname string, from, to int) ([]*VideoInfo, error) {
	stmtOut, err := DbConn.Prepare(`select
		video_info.id,video_info.author_id,video_info.name,video_info.display_ctime
		from video_info	inner join users on video_info.author_id = users.id
	  	where users.login_name = ? and video_info.create_time > FROM_UNIXTIME(?)
	  	and video_info.create_time <= FROM_UNIXTIME(?)
		order by video_info.create_time desc `)
	var res []*VideoInfo
	if err != nil {
		return res, err
	}

	rows, err := stmtOut.Query(uname, from, to)
	if err != nil {
		logger.Info("list video info err:\t", uname, from, to, err)
		return res, err
	}

	for rows.Next() {
		var id, name, ctime string
		var aid int
		if err := rows.Scan(&id, &aid, &name, &ctime); err != nil {
			return res, err
		}

		vi := &VideoInfo{
			Id:           id,
			AuthorId:     aid,
			Name:         name,
			DisplayCtime: ctime,
		}
		res = append(res, vi)
	}

	defer stmtOut.Close()
	return res, nil
}

// 删除指定id的video信息
func DeleteVideoInfo(vid string) error {
	stmtDel, err := DbConn.Prepare("delete from video_info where id=?")
	if err != nil {
		logger.Info("delete video info err:\t", vid, err)
		return err
	}
	if _, err = stmtDel.Exec(vid); err != nil {
		logger.Info("delete video info err:\t", vid, err)
		return err
	}
	defer stmtDel.Close()
	return nil
}

// 添加评论
func AddNewComments(vid string, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		return err
	}

	stmtIns, err := DbConn.Prepare("insert into comments (id, video_id,author_id,content) values(?,?,?,?)")
	if err != nil {
		logger.Info("add new comments err:\t", err)
		return err
	}

	if _, err = stmtIns.Exec(id, vid, aid, content); err != nil {
		logger.Info("add new comments err:\t", err)
		return err
	}
	defer stmtIns.Close()
	return nil
}

// 查找评论
func ListComments(vid string, from, to int) ([]*Comment, error) {
	stmtOut, err := DbConn.Prepare(`select comments.id, users.login_name,comments.content
		from comments inner join users on comments.author_id = users.id
		where comments.video_id = ? and comments.time > FROM_UNIXTIME(?) and comments.time <= FROM_UNIXTIME(?)
		order by comments.time desc`)

	var res []*Comment
	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		logger.Info("list comments err:\t", vid, from, to, err)
		return res, err
	}

	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return res, err
		}
		c := &Comment{
			Id:      id,
			VideoId: vid,
			Author:  name,
			Content: content,
		}
		res = append(res, c)
	}
	defer stmtOut.Close()
	return res, nil
}
