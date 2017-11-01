package loadgen

// Bucket represents a limited resource.
type Bucket struct {
	tokens chan struct{}
}

// NewBucket creates a new bucket with size available tokens.
func NewBucket(size int) *Bucket {
	return &Bucket{make(chan struct{}, size)}
}

// Acquire acquires a token from the bucket. If Acquire returns true, it should be matched by a call to Release().
// Acquire is non-blocking and will immediately return true if it was able to acquire a token, false otherwise.
// The nil Bucket represents an infinite sized bucket. For a nil bucket, Acquire will always return true.
func (b *Bucket) Acquire() bool {
	if b == nil || b.tokens == nil {
		return true
	}
	select {
	case b.tokens <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release returns a token to the bucket. Release is non-blocking. Release should always be matched to a previous
// call to Acquire() that has returned true.
func (b *Bucket) Release() {
	if b == nil || b.tokens == nil {
		return
	}
	select {
	case <-b.tokens:
	default:
	}
}
