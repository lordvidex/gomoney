package adapters

import (
	"context"
	"sync"
	"time"
)

type Locker struct {
	// mx is used for locking mutexes for a lockgroup
	mx sync.Map

	// groups is for storing assigning a single item to a lockgroup
	groups sync.Map
}

func NewLocker(c context.Context, cleanupFrequency time.Duration) *Locker {
	l := &Locker{}
	// regular cleanups
	go l.cleanup(c, cleanupFrequency)
	return l
}

func (l *Locker) cleanup(c context.Context, frequency time.Duration) {
	tick := time.NewTicker(frequency)
	for {
		select {
		case <-tick.C:
			l.mx.Range(func(key, value interface{}) bool {
				m, ok := value.(*sync.Mutex)
				if !ok {
					return true // continue
				}
				if m.TryLock() {
					m.Unlock()
					l.mx.Delete(key)
				}
				return true
			})
		case <-c.Done():
			return
		}
	}
}

func (l *Locker) Lock(x any, y ...any) func() {
	items := append([]any{x}, y...)

	for _, item := range items {
		grp, ok := l.groups.Load(item)
		if !ok {
			continue
		}
		if grp == nil {
			l.groups.Delete(item)
			continue
		}
		cast, ok := grp.(*sync.Mutex)
		if !ok {
			continue
		}
		cast.Lock()   // waits for the group to be unlocked
		cast.Unlock() // unlocks the group
		continue
	}
	// locks := make([]*sync.Mutex, len(y)+1)
	// for i, key := range append([]any{x}, y...) {
	// 	lock, _ := l.mx.LoadOrStore(key, &sync.Mutex{})
	// 	lock.(*sync.Mutex).Lock()
	// 	locks[i] = lock.(*sync.Mutex)
	// }
	// return func() {
	// 	for _, lock := range locks {
	// 		lock.Unlock()
	// 	}
	// }
}
