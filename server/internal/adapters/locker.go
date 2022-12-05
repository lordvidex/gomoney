package adapters

import "sync"

type Locker struct {
	mx sync.Map
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
