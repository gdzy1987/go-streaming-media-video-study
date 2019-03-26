## 5-1 scheduler介绍
### Scheduler
- 什么是scheduler
- 为什么需要scheduler
- scheduler通常做什么

### scheduler包含什么
- Restful的http server 接收任务
- Timer
- 生产者/消费者模型下的task runner

### 架构概览
```
    Producer/Dispatcher

Timer       Channel

    Consumer/Executor    
```

## 5-2 代码架构搭建
```
├── conf.json
├── handler
│   └── handler.go
├── main.go
├── model
│   ├── api.go
│   └── database.go
└── taskrunner
    ├── runner.go
    ├── runner_test.go
    ├── task.go
    └── trmain.go
```

## 5-3 runner的生产消费者模型实现
```
package taskrunner

import "go-streaming-media-video-study/common"

type controlChan chan string

type dataChan chan interface{}

type fn func(dc dataChan) error

type Runner struct {
	Controller controlChan
	Error      controlChan
	Data       dataChan
	dataSize   int
	longLived  bool
	Dispatcher fn
	Executor   fn
}

func NewRunner(size int, longlived bool, d fn, e fn) *Runner {
	return &Runner{
		Controller: make(chan string, 1),
		Error:      make(chan string, 1),
		Data:       make(chan interface{}, size),
		longLived:  longlived,
		dataSize:   size,
		Dispatcher: d,
		Executor:   e,
	}
}

func (r *Runner) startDispatch() {
	defer func() {
		if !r.longLived {
			close(r.Controller)
			close(r.Data)
			close(r.Error)
		}
	}()
	for {
		select {
		case c := <-r.Controller:
			// 生产者
			if c == common.READY_TO_DISPATCH {
				err := r.Dispatcher(r.Data)
				if err != nil {
					r.Error <- common.CLOSE
				} else {
					r.Controller <- common.READY_TO_EXECUTE
				}
			}
			// 消费者
			if c == common.READY_TO_EXECUTE {
				err := r.Executor(r.Data)
				if err != nil {
					r.Error <- common.CLOSE
				} else {
					r.Controller <- common.READY_TO_DISPATCH
				}
			}
		case e := <-r.Error:
			if e == common.CLOSE {
				return
			}
		default:

		}
	}
}

func (r *Runner) StartAll() {
	r.Controller <- common.READY_TO_DISPATCH
	r.startDispatch()
}
```

## 5-4 runner的使用与测试
- go test
```
package taskrunner

import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	d := func(dc dataChan) error {
		for i := 0; i < 30; i++ {
			dc <- i
			log.Printf("Dispatcher sent: %v", i)
		}
		return nil
	}

	e := func(dc dataChan) error {
		forLoop:
			for {
				select {
				case d := <-dc:
					log.Printf("Executor received:%v", d)
				default:
					break forLoop
				}
			}
			return errors.New("executor")
	}

	runner := NewRunner(30, false, d, e)
	go runner.StartAll()
	time.Sleep(3 * time.Second)
}

```

## 5-5 task示例的实现
```
CREATE TABLE `video_del_rec` (
  `id` varchar(64) NOT NULL ,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```
### 实现机制
```
api -> videoid -> mysql
dispatcher -> mysql -> videoid -> datachannel
executor -> datachannel -> videoid -> delete videos
```

## 5-6 timer的实现
```
timer 
setup
start{trigger->task->runner}

timer, task, runner, longlived.
```

```
type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

func NewWorker(interval time.Duration, r *Runner) *Worker {
	return &Worker{
		ticker: time.NewTicker(interval * time.Second),
		runner: r,
	}
}

// ticker 
func (w *Worker) startWorker() {
	for {
		select {
		case <-w.ticker.C:
			go w.runner.StartAll()
		}
	}
}

func Start() {
	r := NewRunner(3, true, VideoClearDispatcher, VideoClearExecutor)
	w := NewWorker(3, r)
	go w.startWorker()
}

```
## 5-7 api的实现以及scheduler完成
```
func VidDelRecHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	if len(vid) == 0 {
		common.SendResponse(w, 400, "video id should not be empty")
		return
	}
	err := model.AddVideoDeletionRecord(vid)
	if err != nil {
		common.SendResponse(w, 500, "Internal server error")
		return
	}
	common.SendResponse(w, 200, "")
	return
}
```
