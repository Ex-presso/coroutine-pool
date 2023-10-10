package main

import (
	"fmt"
	"time"
)

type Task struct {
	todo func() error
}

// CreateTask create a new task
func CreateTask(todo func() error) *Task {
	t := Task{
		todo: todo,
	}
	return &t
}

func (t *Task) Execute() {
	err := t.todo()
	if err != nil {
		return
	}
}

type Pool struct {
	EntryChannel chan *Task
	workerNum    int
	TaskChannel  chan *Task
}

func CreatePool(cap int) *Pool {
	p := Pool{
		EntryChannel: make(chan *Task),
		workerNum:    cap,
		TaskChannel:  make(chan *Task),
	}
	return &p
}

func (p *Pool) createWorker(workId int) {
	for task := range p.TaskChannel {
		task.Execute()
		fmt.Printf("worker %d finish task", workId)
	}

}

func (p *Pool) Run() {
	for i := 0; i < p.workerNum; i++ {
		go p.createWorker(i)
	}

	for task := range p.EntryChannel {
		p.TaskChannel <- task
	}

	close(p.TaskChannel)
	close(p.EntryChannel)
}

func main() {
	t := CreateTask(func() error {
		fmt.Println(time.Now())
		return nil
	})

	p := CreatePool(3)

	go func() {
		for {
			p.EntryChannel <- t
		}
	}()

	p.Run()
}
