package model

import (
	_ "github.com/go-sql-driver/mysql"
	"go-streaming-media-video-study/logger"
)

// 添加video删除记录
func AddVideoDeletionRecord(vid string) error {
	stmtIns, err := dbConn.Prepare("insert into video_del_rec (video_id) values(?)")
	if err != nil {
		return err
	}

	if _, err = stmtIns.Exec(vid); err != nil {
		logger.Info("add video deletion record error:\t", err)
		return err
	}

	defer stmtIns.Close()
	return nil
}

// 读取video删除记录
func ReadVideoDeletionRecord(count int) ([]string, error) {
	stmtOut, err := dbConn.Prepare("select video_id from video_del_rec limit?")
	var ids []string
	if err != nil {
		return ids, err
	}

	rows, err := stmtOut.Query(count)
	if err != nil {
		logger.Info("read video deletion record error:\t", err)
		return ids, err
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	defer stmtOut.Close()
	return ids, nil
}

// 删除video删除记录
func DelVideoDeletionRecord(vid string) error {
	stmtDel, err := dbConn.Prepare("delete from video_del_rec where video_id = ?")
	if err != nil {
		return err
	}

	if _, err = stmtDel.Exec(vid); err != nil {
		logger.Info("del video deletion record error:\t", err)
		return err
	}
	defer stmtDel.Close()
	return nil
}
