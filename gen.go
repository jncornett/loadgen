package loadgen

import (
	"time"
)

// LoadGen maintains state for a load generator run.
type LoadGen struct {
	Passed          Counter
	Failed          Counter
	callback        func() error
	tickInterval    time.Duration
	rateLimit       *Bucket
	maxTransactions int
	maxDuration     time.Duration
	cancel          chan struct{}
	done            chan struct{}
}

// New creates a new load generator with the given config.
func New(c Config) *LoadGen {
	var rateLimit *Bucket
	if c.MaxConcurrency > 0 {
		rateLimit = NewBucket(c.MaxConcurrency)
	}
	if c.Tps == 0.0 {
		c.Tps = DefaultTps
	}
	interval := time.Duration(float64(time.Second) / c.Tps)
	transactionFunc := c.TransactionFunc
	if transactionFunc == nil {
		transactionFunc = DefaultTransactionFunc
	}
	return &LoadGen{
		callback:        transactionFunc,
		tickInterval:    interval,
		rateLimit:       rateLimit,
		maxTransactions: c.MaxTransactions,
		maxDuration:     c.MaxDuration,
		cancel:          make(chan struct{}, 1),
		done:            make(chan struct{}),
	}
}

// Run executes the load generator. Run is blocking. Run should only be called once.
func (g *LoadGen) Run() {
	tick := time.Tick(g.tickInterval)
	if g.maxDuration > time.Duration(0) {
		time.AfterFunc(g.maxDuration, g.Cancel)
	}
	for {
		select {
		case <-g.cancel:
			g.done <- struct{}{}
			return
		case <-tick:
			go g.transaction()
		}
	}
}

// Cancel cancels the load generator. Cancel is non-blocking. Cancel can be called more than once.
func (g *LoadGen) Cancel() {
	select {
	case g.cancel <- struct{}{}:
	default:
	}
}

// Done returns the channel that signals that the load generator has finished processing.
func (g *LoadGen) Done() <-chan struct{} {
	return g.done
}

func (g *LoadGen) updateCounts(passed bool) {
	var total int64
	if passed {
		total = g.Passed.Inc() + g.Failed.Value()
	} else {
		total = g.Failed.Inc() + g.Passed.Value()
	}
	if g.maxTransactions > 0 && total >= int64(g.maxTransactions) {
		g.Cancel()
	}
}

func (g *LoadGen) transaction() {
	if !g.rateLimit.Acquire() {
		return
	}
	defer g.rateLimit.Release()
	err := g.callback()
	g.updateCounts(err == nil)
}
