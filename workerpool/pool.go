package workerpool

import (
	"sync"
)

// Pool воркера
type Pool struct {
	Tasks   []*Task
	Workers []*Worker

	concurrency   int
	collector     chan *Task
	runBackground chan bool
	wg            sync.WaitGroup
}

// NewPool инициализирует новый пул с заданными задачами и

func NewPool(tasks []*Task, concurrency int) *Pool {
	return &Pool{
		Tasks:       tasks,
		concurrency: concurrency,
		collector:   make(chan *Task, 1000),
	}
}

// Run запускает всю работу в Pool и блокирует ее до тех пор,
// пока она не будет закончена.
func (p *Pool) Run() {
	for i := 1; i <= p.concurrency; i++ {
		worker := NewWorker(p.collector, i)
		worker.Start(&p.wg)
	}

	for i := range p.Tasks {
		p.collector <- p.Tasks[i]
	}
	close(p.collector)

	p.wg.Wait()
}

// AddTask добавляет таски в pool
func (p *Pool) AddTask(task *Task) {
	p.collector <- task
}

// RunBackground запускает pool в фоне
func (p *Pool) RunBackground() {

	for i := 1; i <= p.concurrency; i++ {
		worker := NewWorker(p.collector, i)
		p.Workers = append(p.Workers, worker)
		go worker.StartBackground()
	}

	for i := range p.Tasks {
		p.collector <- p.Tasks[i]
	}

	p.runBackground = make(chan bool)
	<-p.runBackground
}

func (p *Pool) Stop() {
	for i := range p.Workers {
		p.Workers[i].Stop()
	}
	p.runBackground <- true
}
