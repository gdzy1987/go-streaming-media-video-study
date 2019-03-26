## 4-1 stream server
### Streaming
- 静态视频，非RTMP
- 独立的服务，可独立部署
- 统一的api格式
### Stream Server
- Streaming
- Upload files

## 4-2 stream的架构搭建
- api模块接口短连接
- 长连接使用流控
```
├── conf.json
├── handler
│   └── handlers.go
├── main.go
├── readme.md
├── upload.html
└── utils
    └── limiter.go
```
## 4-3 token bucket算法
- bucket : 20 * token
- channel, shared channel instead of shared memory.

## 4-4 流控模块的实现
- 本质是有长度channel的生成和消费

```
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

```
## 4-5 在http middleware中嵌入流控
```
type middleWareHandler struct {
	r *httprouter.Router
	l *utils.ConnLimiter
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.l = utils.NewConnLimiter(cc)
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.l.GetConn() {
		defines.SendErrorResponse(w, http.StatusTooManyRequests, "Too many requests")
		return
	}
	m.r.ServeHTTP(w, r)
	defer m.l.ReleaseConn()
}

```
## 4-6 streamHandler实现
```
func StreamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	vl := common.VIDEO_DIR + vid

	video, err := os.Open(vl)
	if err != nil {
		logger.Info("stream handler try to open file,error:\t", err)
		common.SendErrorResponse(w, http.StatusInternalServerError, "internal errors")
		return
	}

	w.Header().Set("Content-Type", "/video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)
	defer video.Close()
}
```
## 4-7 验证streamHandler
- http://localhost:9000/testpage
## 4-8 uploadHandler实现
```
func UploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, common.MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(common.MAX_UPLOAD_SIZE); err != nil {
		common.SendErrorResponse(w, http.StatusBadRequest, "File is too big")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		common.SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Info("uploadHandler read file error:\t", err)
		common.SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	fn := p.ByName("vid-id")
	err = ioutil.WriteFile(common.VIDEO_DIR+fn, data, 0666)
	if err != nil {
		logger.Info("uploadHandler write file error:\t", err)
		common.SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	ossfn := "videos/" + fn
	path := "./videos/" + fn
	ret := common.UploadToOss(ossfn, path, config.DefaultConfig.Bucket)
	if !ret {
		common.SendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	os.Remove(path)

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Uploaded successfully")
}
```
## 4-9 验证uploadHandler