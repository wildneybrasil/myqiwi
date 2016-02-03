package main

import (
	"github.com/streadway/amqp"
)

// Dispatcher takes care of starting/stoping workers.
type Dispatcher struct {
	concurrency int
	jobChannel  <-chan (amqp.Delivery)
	doneChannel chan (bool)
	killChannel chan (bool)
}

// WorkFn is a function executed by a worker.
type WorkFn func(amqp.Delivery, int)

// NewDispatcher returns a new dispatcher. The dispatcher starts N workers
// specified with concurrency which listen on an amqp delivery channel.
func NewDispatcher(jobChannel <-chan (amqp.Delivery), concurrency int) *Dispatcher {
	return &Dispatcher{
		concurrency: concurrency,
		jobChannel:  jobChannel,
		doneChannel: make(chan (bool)),
		killChannel: make(chan (bool)),
	}
}

// Dispatch starts workers and each worker executes fn.
func (dispatcher *Dispatcher) Dispatch(fn WorkFn) {
	for i := 0; i < dispatcher.concurrency; i++ {
		go wrapWorkFn(dispatcher, fn, i)
	}
}

// Wait watkes for all workers to finish their jobs.
func (dispatcher *Dispatcher) Wait() {
	for i := 0; i < dispatcher.concurrency; i++ {
		<-dispatcher.doneChannel
	}
}

// Kill sends die messages to all workers.
func (dispatcher *Dispatcher) Kill() {
	for i := 0; i < dispatcher.concurrency; i++ {
		dispatcher.killChannel <- true
	}
}

// wrapWorkFn wraps fn and takes care of receiving a die message.
func wrapWorkFn(dispatcher *Dispatcher, fn WorkFn, id int) {
	for {
		select {
		case delivery := <-dispatcher.jobChannel:
			fn(delivery, id)
		case <-dispatcher.killChannel:
			dispatcher.doneChannel <- true
		}
	}
}
