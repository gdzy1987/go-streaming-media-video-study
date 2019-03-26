package utils

import (
	"go-streaming-media-video-study/logger"
)

type ConnLimiter struct {
	concurrentConn int
	bucket         chan int
}

// 构造函数
func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: cc,
		bucket:         make(chan int, cc),
	}
}

func (cl *ConnLimiter) GetConn() bool {
	if len(cl.bucket) >= cl.concurrentConn {
		logger.Info("Reached the rate limitation.")
		return false
	}

	cl.bucket <- 1
	return true
}

func (cl *ConnLimiter) ReleaseConn() {
	c := <-cl.bucket
	logger.Info("New connection coming:\t", c)
}
