package adapters

import (
	"context"
	"sort"
	"sync"
	"time"
)

type Locker struct {
	mx sync.Map
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

// values passed to Lock must be of the same type [either all strings or all ints]
func (l *Locker) Lock(x any, y ...any) func() {
	locks := make([]*sync.Mutex, len(y)+1)
	// sort the keys to avoid deadlocks
	arr := orderedLocks(append([]any{x}, y...))
	sort.Sort(arr)
	for i, key := range arr {
		lock, _ := l.mx.LoadOrStore(key, &sync.Mutex{})
		lock.(*sync.Mutex).Lock()
		locks[i] = lock.(*sync.Mutex)
	}
	return func() {
		for _, lock := range locks {
			lock.Unlock()
		}
	}
}

type orderedLocks []any

func (o orderedLocks) Len() int {
	return len(o)
}

func (o orderedLocks) Less(i int, j int) bool {
	switch o[i].(type) {
	case string:
		return o[i].(string) < o[j].(string)
	case int:
		return o[i].(int) < o[j].(int)
	default:
		panic("Lock key must either be strings or ints")
	}
}

func (o orderedLocks) Swap(i int, j int) {
	o[i], o[j] = o[j], o[i]
}
