package safe

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/panjf2000/ants/v2"
	"runtime/debug"
	"sync"

	"github.com/rs/zerolog/log"
)

type routineCtx func(ctx context.Context)

// Pool is a pool of go routines.
type Pool struct {
	antsPool  *ants.Pool
	waitGroup sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewPool creates a Pool.
func NewPool(parentCtx context.Context, size int) *Pool {
	if size == 0 {
		size = 10000
	}

	ctx, cancel := context.WithCancel(parentCtx)
	p, err := ants.NewPool(size)
	if err != nil {
		panic(err)
	}
	return &Pool{
		antsPool: p,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// GoCtx starts a recoverable goroutine with a context.
func (p *Pool) GoCtx(goroutine routineCtx) {
	p.waitGroup.Add(1)
	p.Go(func() {
		defer p.waitGroup.Done()
		goroutine(p.ctx)
	})
}

// Stop stops all started routines, waiting for their termination.
func (p *Pool) Stop() {
	p.cancel()
	p.waitGroup.Wait()
	p.antsPool.Release()
}

// Go starts a recoverable goroutine.
func (p *Pool) Go(goroutine func()) {
	p.GoWithRecover(goroutine, defaultRecoverGoroutine)
}

// GoWithRecover starts a recoverable goroutine using given customRecover() function.
func (p *Pool) GoWithRecover(goroutine func(), customRecover func(err interface{})) {
	err := p.antsPool.Submit(func() {
		defer func() {
			if err := recover(); err != nil {
				customRecover(err)
			}
		}()
		goroutine()
	})
	if err != nil {
		log.Error().Interface("error", err).Msg("Error Ants Pool Submit")
		return
	}
}

func defaultRecoverGoroutine(err interface{}) {
	log.Error().Interface("error", err).Msg("Error in Go routine")
	log.Error().Msgf("Stack: %s", debug.Stack())
}

// OperationWithRecover wrap a backoff operation in a Recover.
func OperationWithRecover(operation backoff.Operation) backoff.Operation {
	return func() (err error) {
		defer func() {
			if res := recover(); res != nil {
				defaultRecoverGoroutine(res)
				err = fmt.Errorf("panic in operation: %w", err)
			}
		}()
		return operation()
	}
}
