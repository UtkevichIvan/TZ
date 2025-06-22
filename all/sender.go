package all

type taskQueue struct {
	task    JobDataToHandle
	next    *taskQueue
	isEmpty bool
}

type Sender struct {
	tq *taskQueue
	ch chan JobDataToHandle
}

func NewSender(ch chan JobDataToHandle) *Sender {
	return &Sender{ch: ch, tq: &taskQueue{isEmpty: true}}
}

func (s *Sender) Send() {
	for {
		if !s.tq.isEmpty {
			s.ch <- s.tq.task
			s.remove()
		}
	}
}

func (s *Sender) Add(data JobDataToHandle) {
	if s.tq.isEmpty {
		s.tq = &taskQueue{task: data, isEmpty: false}
	} else {
		s.tq.next = &taskQueue{task: data}
	}
}

func (s *Sender) remove() {
	if s.tq.next != nil {
		s.tq = s.tq.next
	} else {
		s.tq = &taskQueue{isEmpty: true}
	}
}
