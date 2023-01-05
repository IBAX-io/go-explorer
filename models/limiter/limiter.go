package limiter

import (
	"sync"
)

type RequestLimitObject struct {
	maxCount int
	lock     sync.Mutex
	reqCount int
	signal   chan int
}

func NewRequestLimitService(maxCnt int) *RequestLimitObject {
	p := &RequestLimitObject{
		maxCount: maxCnt,
	}
	if p.signal == nil {
		p.signal = make(chan int, 1)
	}
	go func() {
		for {
			select {
			case <-p.signal:
				return
			}
		}
	}()

	return p
}

func (p *RequestLimitObject) Increase() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.reqCount += 1
}

func (p *RequestLimitObject) Reduce() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.reqCount -= 1
}

func (p *RequestLimitObject) IsAvailable() bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.reqCount < p.maxCount
}

func (p *RequestLimitObject) GetCount() int {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.reqCount
}

func (p *RequestLimitObject) GetMax() int {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.maxCount
}

func (p *RequestLimitObject) Delete() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.signal <- 1
}
