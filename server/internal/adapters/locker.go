package adapters

import (
	"context"
	"sync"
	"time"
)

type Locker struct {
	mx sync.Map
}

func NewLocker(c context.Context, cleanupFrequency time.Duration) *Locker {
	l := &Locker{}

	// regular cleanups
	go func() {
		tick := time.NewTicker(cleanupFrequency)
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
	}()
	return l
}

func (l *Locker) Lock(x any, y ...any) func() {
	locks := make([]*sync.Mutex, len(y)+1)
	lock, _ := l.mx.LoadOrStore(x, &sync.Mutex{})
	lock.(*sync.Mutex).Lock()
	locks[0] = lock.(*sync.Mutex)
	for i, key := range y {
		lock, _ = l.mx.LoadOrStore(key, &sync.Mutex{})
		lock.(*sync.Mutex).Lock()
		locks[i+1] = lock.(*sync.Mutex)
	}
	return func() {
		for _, lock = range locks {
			lock.(*sync.Mutex).Unlock()
		}
	}
}
