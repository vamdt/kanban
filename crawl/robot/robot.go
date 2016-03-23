package robot

import (
	"bytes"
	"container/list"
	"container/ring"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/glog"

	. "../base"
)

const (
	robotIdle int32 = iota
	robotBusy
)

const (
	_ int32 = iota
	TaskDay
	TaskMin1
	TaskMin5
	TaskTick
	TaskRealTick
	TaskRealTicks
)

const DefaultRobotConcurrent int = 3
const maxMultiJobsConcurrent int = 50

type Robot interface {
	Days_download(id string, start time.Time) ([]Tdata, error)
	GetRealtimeTick(ids string) []RealtimeTickRes
	Can(id string, task int32) bool
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
	start  bool
}

func NewWorker(worker Robot) *Worker { return &Worker{worker: worker} }
func NewRobotBox() *RobotBox         { return &RobotBox{} }

var DefaultRobotBox = NewRobotBox()

func Work() {
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

func (p *RobotBox) getJob(robot Robot) *JobItem {
	p.mjob.Lock()
	defer p.mjob.Unlock()
	for e := p.jobs.Front(); e != nil; e = e.Next() {
		job := e.Value.(*JobItem)
		if robot.Can(job.id, job.task) {
			p.jobs.Remove(e)
			return job
		}
	}
	return nil
}

func (p *RobotBox) getJobs(robot Robot, job *JobItem) []*JobItem {
	p.mjob.Lock()
	defer p.mjob.Unlock()
	jobs := []*JobItem{job}
	task := job.task
	for e := p.jobs.Front(); e != nil && len(jobs) < maxMultiJobsConcurrent; {
		job := e.Value.(*JobItem)
		next := e.Next()
		if task == job.task && robot.Can(job.id, task) {
			jobs = append(jobs, job)
			p.jobs.Remove(e)
		}
		e = next
	}
	return jobs
}

func Days_download(id string, start time.Time) ([]Tdata, error) {
	return DefaultRobotBox.Days_download(id, start)
}

type JobItem struct {
	id    string
	start time.Time
	res   chan []Tdata
	task  int32
	cb    func(interface{}, bool) bool
}

func (p *Worker) DoRealTick(jobs []*JobItem) {
	defer atomic.StoreInt32(&p.busy, robotIdle)
	if jobs == nil || len(jobs) < 1 {
		return
	}
	var b bytes.Buffer
	for _, job := range jobs {
		b.WriteString(",")
		b.WriteString(job.id)
	}
	ids := b.String()[1:]
	res := p.worker.GetRealtimeTick(ids)
	for _, job := range jobs {
		ok := false
		var rt *RealtimeTick
		for _, r := range res {
			if job.id == r.Id {
				rt = &r.RealtimeTick
				ok = true
				break
			}
		}
		job.cb(rt, ok)
	}
}

func (p *Worker) Do(job *JobItem) {
	defer atomic.StoreInt32(&p.busy, robotIdle)
	if job == nil {
		return
	}

	switch job.task {
	case TaskDay:
		data, _ := p.worker.Days_download(job.id, job.start)
		job.res <- data
		return
	case TaskRealTick:
		res := p.worker.GetRealtimeTick(job.id)
		var rt *RealtimeTick
		ok := false
		if len(res) > 0 {
			rt = &res[0].RealtimeTick
			ok = true
		}
		job.cb(rt, ok)
		return
	}
}

func (p *RobotBox) Work(once bool) {
	if !once {
		p.mrobot.Lock()
		start := p.start
		if !p.start {
			p.start = true
		}
		p.mrobot.Unlock()
		if start {
			return
		}
	}

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
		count := 0
		busy := 0
		do := 0
		p.robots.Do(func(v interface{}) {
			if v == nil {
				return
			}
			count++
			w := v.(*Worker)
			if !atomic.CompareAndSwapInt32(&w.busy, robotIdle, robotBusy) {
				busy++
				return
			}
			job := p.getJob(w.worker)
			if job == nil {
				atomic.StoreInt32(&w.busy, robotIdle)
				return
			}

			if job.task == TaskRealTick && w.worker.Can(job.id, TaskRealTicks) {
				jobs := p.getJobs(w.worker, job)
				go w.DoRealTick(jobs)
			} else {
				go w.Do(job)
			}
			do++
		})
		p.robots = p.robots.Move(1)
		if do > 0 {
			glog.Infof("%dn %dbusy/robot(%d) %d/jobs(%d)", do, busy, count, l-do, l)
		}
		p.mrobot.Unlock()
	}
}

func (p *RobotBox) Days_download(id string, start time.Time) ([]Tdata, error) {
	job := JobItem{
		id:    id,
		start: start,
		task:  TaskDay,
	}
	job.res = make(chan []Tdata)
	p.mjob.Lock()
	p.jobs.PushBack(&job)
	p.mjob.Unlock()
	res := <-job.res
	return res, nil
}

func GetRealtimeTick(id string, cb func(interface{}, bool) bool) {
	DefaultRobotBox.GetRealtimeTick(id, cb)
}

func (p *RobotBox) GetRealtimeTick(id string, cb func(interface{}, bool) bool) {
	job := JobItem{
		id:   id,
		task: TaskRealTick,
		cb:   cb,
	}
	p.mjob.Lock()
	p.jobs.PushBack(&job)
	p.mjob.Unlock()
}
