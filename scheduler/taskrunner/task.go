package taskrunner

import (
	"errors"
	"go-streaming-media-video-study/common"
	"go-streaming-media-video-study/config"
	"go-streaming-media-video-study/logger"
	"go-streaming-media-video-study/scheduler/model"
	"strings"
	"sync"
)

func deleteVideo(vid string) error {
	bn := strings.Trim(config.DefaultConfig.Bucket, "/")
	ossfn := bn + "/" + vid
	ok := common.DeleteObject(ossfn, bn)

	if !ok {
		logger.Info("delete video error,oss operation failed")
		return errors.New("deleting video error")
	}

	return nil
}

func VideoClearDispatcher(dc dataChan) error {
	res, err := model.ReadVideoDeletionRecord(3)
	if err != nil {
		logger.Info("video clear dispatcher error:\t", err)
		return err
	}

	if len(res) == 0 {
		return errors.New("all tasks finished")
	}

	for _, id := range res {
		dc <- id
	}
	return nil
}

func VideoClearExecutor(dc dataChan) error {
	errMap := &sync.Map{}
	var err error
forLoop:
	for {
		select {
		case vid := <-dc:
			go func(id interface{}) {
				if err := deleteVideo(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
				if err := model.DelVideoDeletionRecord(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
			}(vid)
		default:
			break forLoop
		}
	}
	errMap.Range(func(key, value interface{}) bool {
		err = value.(error)
		if err != nil {
			return false
		}
		return true
	})
	return err
}
