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
		Controller: make(chan string, 1), // 非阻塞
		Error:      make(chan string, 1),
		Data:       make(chan interface{}, size),
		longLived:  longlived,
		dataSize:   size,
		Dispatcher: d,
		Executor:   e,
	}
}

func (r *Runner) StartAll() {
	r.Controller <- common.READY_TO_DISPATCH
	r.startDispatch()
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
