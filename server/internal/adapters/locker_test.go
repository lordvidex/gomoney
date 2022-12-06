package adapters

import (
	"context"
	"testing"
	"time"
)

func TestLocker_Lock(t *testing.T) {
	type args struct {
		x    any
		y    []any
		f    time.Duration
		wait bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "lock",
			args: args{
				x: "x",
				y: []any{"y", "z"},
				f: time.Minute,
			},
		},
		{
			name: "test cleanup",
			args: args{
				x:    "x",
				y:    []any{"y", "z"},
				f:    time.Second,
				wait: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLocker(context.TODO(), tt.args.f)
			unlock := l.Lock(tt.args.x, tt.args.y...)
			unlock()
			for _, key := range append([]any{tt.args.x}, tt.args.y...) {
				_, ok := l.mx.Load(key)
				if !ok {
					t.Errorf("lock not found")
				}
			}
			if tt.args.wait {
				time.Sleep(tt.args.f * 2)
				for _, key := range append([]any{tt.args.x}, tt.args.y...) {
					_, ok := l.mx.Load(key)
					if ok {
						t.Errorf("lock not cleaned up")
					}
				}
			}
		})
	}
}
