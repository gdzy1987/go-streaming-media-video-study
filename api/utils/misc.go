package utils

import (
	"crypto/rand"
	"fmt"
	"go-streaming-media-video-study/config"
	"go-streaming-media-video-study/logger"
	"io"
	"net/http"
	"strconv"
	"time"
)

func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}

	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func GetCurrentTimestampSec() int {
	ts, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000, 10))
	return ts
}

// 删除视频
func SendDeleteVideoRequest(id string) {
	addr := config.DefaultConfig.Address + ":" + config.DefaultConfig.StreamServerPort
	url := "http://" + addr + "/video-delete-record/" + id
	logger.Info("删除视频:\t", url)
	if _, err := http.Get(url); err != nil {
		logger.Info("Sending deleting video request error:\t", id, err)
	}
}
