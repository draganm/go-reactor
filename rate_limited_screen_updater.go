package reactor

import (
	"sync"
	"time"
)

type rateLimitedScreenUpdater struct {
	sync.Mutex
	ticker          *time.Ticker
	done            chan interface{}
	pendingUpdate   *DisplayUpdate
	parent          func(*DisplayUpdate)
	didUpdateParent bool
	closed          bool
}

func NewrateLimitedScreenUpdater(minDuration time.Duration, parent func(*DisplayUpdate)) *rateLimitedScreenUpdater {
	ticker := time.NewTicker(minDuration)
	done := make(chan interface{})

	updater := &rateLimitedScreenUpdater{
		ticker: ticker,
		parent: parent,
		done:   done,
	}

	go func() {

		for {
			select {
			case <-ticker.C:
				updater.tick()
			case _, ok := <-done:
				if !ok {
					return
				}
			}
		}
	}()

	return updater
}

func (r *rateLimitedScreenUpdater) update(update *DisplayUpdate) {
	r.Lock()
	defer r.Unlock()

	if r.closed {
		return
	}

	if !r.didUpdateParent {
		r.parent(update)
		r.didUpdateParent = true
		return
	}

	// do sent the update if there is something to eval, new title or new location
	if update.Eval != "" || update.Title != "" || update.Location != "" {
		r.parent(update)
	}

	r.pendingUpdate = update
}

func (r *rateLimitedScreenUpdater) tick() {
	r.Lock()
	defer r.Unlock()

	if r.pendingUpdate != nil {
		r.parent(r.pendingUpdate)
	}

	r.didUpdateParent = false
	r.pendingUpdate = nil
}

func (r *rateLimitedScreenUpdater) close() {
	r.Lock()
	defer r.Unlock()
	r.closed = true
	r.ticker.Stop()
	close(r.done)

}
