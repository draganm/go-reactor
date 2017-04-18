package reactor

import (
	"sync"
)

func newRemoveDuplicatesScreenUpdater(parent func(*DisplayUpdate)) func(*DisplayUpdate) {

	lock := sync.Mutex{}
	var lastUpdate *DisplayUpdate

	return func(update *DisplayUpdate) {
		lock.Lock()
		defer lock.Unlock()

		if !update.DeepEqual(lastUpdate) {
			parent(update)
			lastUpdate = update
		}
	}

}
