package loadgen

import "time"

// DefaultTps is the default Tps configuration.
const DefaultTps = 1.0

// DefaultTransactionFunc is the default transaction function. It is a no-op.
var DefaultTransactionFunc = func() error { return nil }

// Config represents load generator configuration.
type Config struct {
	// Tps configures per second of the load generator. Can be omitted.
	Tps float64
	// MaxTransactions configures the maximum number of transactions the load generator will execute.
	// If it is zero or omitted, the load generator will not have a transaction limit.
	MaxTransactions int
	// MaxDurations configures the maximum runtime duration of the load generator.
	// If it is zero or omitted, the load generator will not have a duration limit.
	MaxDuration time.Duration
	// MaxConcurrency configures the maximum number of concurrent transactions for the load generator.
	// If it is zero or omitted, the load generator will not have a concurrency limit.
	MaxConcurrency int
	// TransactionFunc is the actual function that the load generator will invoke.
	// Load generator checks the error value to determine whether the transaction passed or failed>
	// If it is nil or omitted, the load generator will use a no-op function.
	TransactionFunc func() error
}
