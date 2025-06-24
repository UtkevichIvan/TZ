package all

import (
	"context"
	"sync"
)

type taskQueue struct {
	task JobDataToHandle
	next *taskQueue
}

type Sender struct {
	head *taskQueue
	end  *taskQueue
	ch   chan JobDataToHandle
	ctx  context.Context
	cond *sync.Cond
}

func NewSender(ch chan JobDataToHandle) *Sender {
	return &Sender{ch: ch, ctx: context.Background(), cond: sync.NewCond(&sync.Mutex{})}
}

func (s *Sender) Start() {
	for {
		s.cond.L.Lock()
		if s.head == nil {
			s.cond.Wait()
		}
		s.cond.L.Unlock()
		s.ch <- s.head.task
		s.remove()
	}
}

func (s *Sender) Add(data JobDataToHandle) {
	s.cond.L.Lock()
	if s.head == nil {
		s.head = &taskQueue{task: data}
		s.end = s.head
	} else {
		s.end.next = &taskQueue{task: data}
		s.end = s.end.next
	}
	s.cond.L.Unlock()
	s.cond.Signal()
}

func (s *Sender) remove() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()
	if s.head.next != nil {
		s.head = s.head.next
	} else {
		s.head = nil
	}
}
