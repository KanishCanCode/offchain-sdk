package job

import (
	"context"
	"time"

	jobtypes "github.com/berachain/offchain-sdk/job/types"
	"github.com/berachain/offchain-sdk/worker"
	"github.com/ethereum/go-ethereum"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

// ============================================
// Polling Jobs
// ============================================

// Polling represents a polling job.
type Polling interface {
	Basic
	IntervalTime(ctx context.Context) time.Duration
}

// WrapPolling wraps a polling job into a conditional job, this is possible since,
// polling jobs are simply conditional jobs where `Condition()` always returns true.
func WrapPolling(c Polling) HasProducer {
	return &conditional{&polling{c}}
}

// Remember, polling is just a conditional job where the condition is always true.
var _ Conditional = (*polling)(nil)

// polling is a wrapper for a polling job.
type polling struct {
	Polling
}

// Condition returns true.
func (p *polling) Condition(context.Context) bool {
	return true
}

// ============================================
// Conditional Jobs
// ============================================

// Conditional represents a conditional job.
type Conditional interface {
	Polling
	Condition(ctx context.Context) bool
}

// Wrap Conditional, wraps a conditional job proivded by a user into a conditional job that
// can be intialized by the job manager.
func WrapConditional(c Conditional) HasProducer {
	return &conditional{c}
}

// conditional is a wrapper for a conditional job.
type conditional struct {
	Conditional
}

// ConditionalProducer produces a job when the condition is met.
func (cj *conditional) Producer(ctx context.Context, pool worker.Pool) error {
	for {
		select {
		// If the context is cancelled, return.
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Sleep for a period of time.
			time.Sleep(cj.IntervalTime(ctx))

			// Check if the condition is true.
			if cj.Condition(ctx) {
				// If true add a job
				pool.AddJob(jobtypes.NewPayload(ctx, cj, nil))
			}
		}
	}
}

// Subscribable represents a subscribable job.
type Subscribable interface {
	Basic
	Subscribe(ctx context.Context) chan any
}

// EthSubscribable represents a subscription to an ethereum event.
type EthSubscribable interface {
	Basic
	Subscribe(ctx context.Context) (ethereum.Subscription, chan coretypes.Log)
	Unsubscribe(ctx context.Context)
}