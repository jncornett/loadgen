package loadgen

import "sync/atomic"

// Counter represents an atomic counter.
type Counter struct {
	value int64
}

// Inc increments the counter value by one and returns the updated value.
func (c *Counter) Inc() int64 {
	if c == nil {
		return 0
	}
	return atomic.AddInt64(&c.value, 1)
}

// Value returns the counter value.
func (c *Counter) Value() int64 {
	if c == nil {
		return 0
	}
	return atomic.LoadInt64(&c.value)
}

// Reset resets the counter value to zero.
func (c *Counter) Reset() {
	if c == nil {
		return
	}
	atomic.StoreInt64(&c.value, 0)
}
