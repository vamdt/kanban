package crawl

import (
	"container/list"
	"container/ring"
	"sync"
	"sync/atomic"
	"time"

	. "./base"
)

const (
	robotIdle int32 = iota
	robotBusy
)

const defaultRobotConcurrent int = 4

type Robot interface {
	Days_download(id string, start time.Time) ([]Tdata, error)
}

type Worker struct {
	worker Robot
	busy   int32
}

type RobotBox struct {
	robots *ring.Ring
	jobs   list.List
	mrobot sync.Mutex
	mjob   sync.Mutex
}

func NewWorker(worker Robot) *Worker { return &Worker{worker: worker} }
func NewRobotBox() *RobotBox         { return &RobotBox{} }

var DefaultRobotBox = NewRobotBox()

func init() {
	go DefaultRobotBox.Work(false)
}

func Registry(robot Robot) {
	DefaultRobotBox.Registry(robot)
}

func (p *RobotBox) Registry(robot Robot) {
	p.mrobot.Lock()
	defer p.mrobot.Unlock()
	s := ring.New(1)
	s.Value = NewWorker(robot)
	if p.robots == nil {
		p.robots = s
	} else {
		p.robots.Link(s)
	}
}

func (p *RobotBox) GetJob() *jobItem {
	p.mjob.Lock()
	defer p.mjob.Unlock()
	e := p.jobs.Front()
	if e == nil {
		return nil
	}
	return p.jobs.Remove(e).(*jobItem)
}

func Days_download(id string, start time.Time) ([]Tdata, error) {
	return DefaultRobotBox.Days_download(id, start)
}

type jobItem struct {
	id    string
	start time.Time
	res   chan []Tdata
}

func (p *Worker) Do(job *jobItem) {
	defer atomic.StoreInt32(&p.busy, robotIdle)
	if job == nil {
		return
	}
	data, _ := p.worker.Days_download(job.id, job.start)
	job.res <- data
}

func (p *RobotBox) Work(once bool) {
	for {
		p.mjob.Lock()
		l := p.jobs.Len()
		p.mjob.Unlock()
		if l < 1 {
			if once {
				break
			} else {
				time.Sleep(time.Second)
				continue
			}
		}

		p.mrobot.Lock()
		p.robots.Do(func(v interface{}) {
			if v == nil {
				return
			}
			w := v.(*Worker)
			if !atomic.CompareAndSwapInt32(&w.busy, robotIdle, robotBusy) {
				return
			}
			go w.Do(p.GetJob())
		})
		p.robots = p.robots.Move(defaultRobotConcurrent)
		p.mrobot.Unlock()
	}
}

func (p *RobotBox) Days_download(id string, start time.Time) ([]Tdata, error) {
	job := jobItem{id: id, start: start}
	job.res = make(chan []Tdata)
	p.mjob.Lock()
	p.jobs.PushBack(&job)
	p.mjob.Unlock()
	res := <-job.res
	return res, nil
}
