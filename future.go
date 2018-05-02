package future

import (
	"context"
	"sync"
	"time"
)

// Runnable is a function that can later be executed
type Runnable func() interface{}

type Pool struct {
	executor chan Runnable
	enabled  bool
	p        sync.Pool
	ctx      context.Context
}

// Put a callback on the Pool's executor.  If the channel is full, it will block
// until ctx is finished, before returning a nil Future.  It is similar to Java's
// BlockingQueue Put.
func (p *Pool) Put(ctx context.Context, callback Runnable) *Future {
	if p == nil {
		return nil
	}

	select {
	case p.executor <- callback:
		return &Future{
			callback: callback,
			done:     make(chan struct{}, 1),
		}
	case <-ctx.Done():
		return nil
	}
}

func (p *Pool) process() {
	t := time.NewTimer(1)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			p.execute()
			t.Reset(1 * time.Second)
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Pool) execute() {
	for futureChan := range p.executor {
		// TODO: do things
	}
}

// Offer is similar to Java's BlockingQueue offer. Offers a callback to the pool
// and returns a future for the callback. If the channel is full, will return nil.
func (p *Pool) Offer(callback Runnable) *Future {
	return nil
}

// OfferOrDo is a helper around Offer that always returns a non nil Future that
// will eventually be executed, either inline (if the pool is off) or
// asynchronous (if the pool is on and has room).  If you want to make sure your
// callback is executed, this is probably the behavior you want.
func OfferOrDo(p *Pool, callback Runnable) *Future {
	if !p.enabled {

	}
	return nil
}

// A Future is a job sent to the Pool that can either be executed eventually, or
// executed explicitly with Done()
type Future struct {
	// The function to execute
	callback Runnable
	// The result of the calculation
	result interface{}
	// closed when a result is calculated
	done chan struct{}
	// Ensures we calculate a result once
	once sync.Once
	// how long the execute step took
	executeTime time.Duration
}

func (f *Future) run() {
	f.once.Do(func() { f.doRun() })
}

func (f *Future) doRun() {
	f.result = f.callback()
	f.done <- struct{}{}
}
