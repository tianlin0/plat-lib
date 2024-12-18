// Package promise 流程控制
package promise

// Promise represents a promise struct
type Promise struct {
	resolve chan interface{}
	reject  chan error
	done    chan bool
}

// New represents a new instance of Promise struct
func New(fn func(func(interface{}), func(error))) *Promise {
	promise := &Promise{
		resolve: make(chan interface{}, 1),
		reject:  make(chan error, 1),
		done:    make(chan bool, 1),
	}
	go fn(promise.resolveFn, promise.rejectFn)
	return promise
}

// Then represents a next chain of success flow
func (p *Promise) Then(success func(interface{})) *Promise {
	go func() {
		if result, ok := <-p.resolve; ok {
			p.done <- ok
			success(result)
			close(p.reject)
		}
	}()
	return p
}

// Catch represents a next chain of failure flow
func (p *Promise) Catch(failure func(error)) *Promise {
	go func() {
		if result, ok := <-p.reject; ok {
			p.done <- ok
			failure(result)
			close(p.resolve)
		}
	}()
	return p
}

// Wait waits for the promise to be resolved or rejected
func (p *Promise) Wait() {
	<-p.done
}

func (p *Promise) resolveFn(i interface{}) {
	p.resolve <- i
}

func (p *Promise) rejectFn(e error) {
	p.reject <- e
}
