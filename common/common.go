package common

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go-streaming-media-video-study/config"
	"go-streaming-media-video-study/logger"
)

func UploadToOss(filename string, path string, bn string) bool {
	client, err := oss.New(config.DefaultConfig.OssAddr, config.DefaultConfig.OssID, config.DefaultConfig.OssSecret)
	if err != nil {
		logger.Info("upload to oss service, new client error:\t", err)
		return false
	}

	bucket, err := client.Bucket(bn)
	if err != nil {
		logger.Info("upload to oss service, get bucket error:\t", err)
		return false
	}

	if err = bucket.UploadFile(filename, path, 500*1024, oss.Routines(3)); err != nil {
		logger.Info("upload to oss service,uploading object error:\t", err)
		return false
	}
	return true
}

func DeleteObject(filename string, bn string) bool {
	client, err := oss.New(config.DefaultConfig.OssAddr, config.DefaultConfig.OssID, config.DefaultConfig.OssSecret)
	if err != nil {
		logger.Info("delete object, new client error:\t", err)
		return false
	}

	bucket, err := client.Bucket(bn)
	if err != nil {
		logger.Info("delete object, get bucket error:\t", err)
		return false
	}

	err = bucket.DeleteObject(filename)
	if err != nil {
		logger.Info("delete object, delete object error:\t", err)
		return false
	}
	return true
}
